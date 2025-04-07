package server

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"bytechat/cmd/web/chat"
	"bytechat/internal/database"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// showChatPage renders a chat page
func (s *Server) showChatPage(w http.ResponseWriter, r *http.Request) {
	s.renderComponent(w, r, chat.ChatPage())
}

// showChatRooms renders a list of available public chatrooms
func (s *Server) showChatRooms(w http.ResponseWriter, r *http.Request) {
}

// addSubscriber registers a subscriber for a given room.
func (s *Server) addSubscriber(roomID uuid.UUID, sub *subscriber) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	if _, ok := s.subscribersByRoom[roomID]; !ok {
		s.subscribersByRoom[roomID] = make(map[*subscriber]struct{})
	}
	s.subscribersByRoom[roomID][sub] = struct{}{}
	slog.Info("Subscriber added", "userID", sub.userID, "roomID", roomID)
}

// deleteSubscriber removes a subscriber from a given room.
func (s *Server) deleteSubscriber(roomID uuid.UUID, sub *subscriber) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	if subs, ok := s.subscribersByRoom[roomID]; ok {
		delete(subs, sub)
		if len(subs) == 0 {
			delete(s.subscribersByRoom, roomID)
		}
		slog.Info("Subscriber deleted", "userID", sub.userID, "roomID", roomID)
	}
}

// publish broadcasts a message to all subscribers in a specific room.
func (s *Server) publish(roomID uuid.UUID, msg []byte) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	subs, ok := s.subscribersByRoom[roomID]
	if !ok {
		slog.Warn("Attempted to publish to room with no subscribers", "roomID", roomID)
		return
	}

	slog.Info("Publishing message", "roomID", roomID, "subscribers", len(subs))
	for sub := range subs {
		select {
		case sub.msgs <- msg:
			// Message sent successfully
		default:
			// Subscriber's buffer is full, they are too slow.
			go sub.closeSlow()
		}
	}
}

// Helper function for writing with timeout
func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, msg)
}

// handleWebSocket upgrades the connection to WebSocket and manages the lifecycle.
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	roomIDStr := r.PathValue("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		slog.Error("Invalid room ID format", "roomID", roomIDStr, "error", err)
		http.Error(w, "Invalid Room ID", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(userContextKey).(User)
	if !ok || user.UserID == uuid.Nil {
		slog.Warn("WebSocket connection attempt without authentication")
		writeError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	slog.Info("Attempting WebSocket upgrade", "userID", user.UserID, "roomID", roomID)

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		slog.Error("WebSocket upgrade failed", "error", err)
		return
	}

	defer c.CloseNow()
	slog.Info("WebSocket connection established", "userID", user.UserID, "roomID", roomID)

	var mu sync.Mutex
	var closed bool

	sub := &subscriber{
		msgs:   make(chan []byte, s.subscriberMessageBuffer),
		userID: user.UserID,
		roomID: roomID,
		closeSlow: func() {
			mu.Lock()
			defer mu.Unlock()
			if !closed {
				closed = true
				slog.Warn("Closing slow connection", "userID", user.UserID, "roomID", roomID)
				c.Close(websocket.StatusPolicyViolation, "connection too slow")
			}
		},
	}
	s.addSubscriber(roomID, sub)
	defer s.deleteSubscriber(roomID, sub)

	ctx := c.CloseRead(context.Background())
	errc := make(chan error, 1)

	go func() {
		defer close(errc)
		for {
			msgType, r, err := c.Reader(ctx)
			if err != nil {
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
					websocket.CloseStatus(err) == websocket.StatusGoingAway ||
					errors.Is(err, context.Canceled) ||
					errors.Is(err, net.ErrClosed) {
					slog.Info("WebSocket reader closed gracefully", "userID", sub.userID, "roomID", sub.roomID, "error", err)
				} else {
					slog.Error("WebSocket read error", "userID", sub.userID, "roomID", sub.roomID, "error", err)
					errc <- err
				}
				return
			}

			const maxMessageSize = 8192
			lr := io.LimitedReader{R: r, N: maxMessageSize + 1} // +1 to detect overflow
			msgBytes, err := io.ReadAll(&lr)
			if err != nil {
				slog.Error("Error reading message bytes", "userID", sub.userID, "roomID", sub.roomID, "error", err)
				continue
			}
			if lr.N <= 0 {
				slog.Warn("Received message exceeds size limit", "userID", sub.userID, "roomID", sub.roomID, "limit", maxMessageSize)
				c.Close(websocket.StatusMessageTooBig, "message too large")
				errc <- errors.New("message too large")
				return
			}

			if msgType == websocket.MessageText {
				slog.Info("Received message", "userID", sub.userID, "roomID", sub.roomID, "size", len(msgBytes))

				var incomingMsg struct {
					Text string `json:"text"`
				}
				if err := json.Unmarshal(msgBytes, &incomingMsg); err != nil || incomingMsg.Text == "" {
					slog.Warn("Received invalid message format or empty text", "userID", sub.userID, "roomID", sub.roomID, "data", string(msgBytes))
					continue
				}

				messageParams := database.CreateMessageParams{
					UserID: sub.userID,
					RoomID: sub.roomID,
					Text:   pgtype.Text{String: incomingMsg.Text, Valid: true},
				}
				dbMsg, err := s.queries.CreateMessage(ctx, messageParams)
				if err != nil {
					slog.Error("Failed to save message to DB", "userID", sub.userID, "roomID", sub.roomID, "error", err)
					continue
				}
				slog.Info("Message saved to DB", "messageID", dbMsg.MessageID)

				broadcastMsg := map[string]any{
					"type": "message",
					"text": incomingMsg.Text,
					"sender": map[string]string{
						"id":        user.UserID.String(),
						"firstName": user.FirstName,
						"lastName":  user.LastName,
					},
					"timestamp": time.Now().UTC().Format(time.RFC3339),
				}
				broadcastBytes, _ := json.Marshal(broadcastMsg)

				s.publish(sub.roomID, broadcastBytes)

			} else {
				slog.Warn("Received non-text message type", "type", msgType, "userID", sub.userID, "roomID", sub.roomID)
			}
		}
	}()

	// Goroutine to write outgoing messages to the client
	for {
		select {
		case msg := <-sub.msgs:
			err := writeTimeout(ctx, time.Second*5, c, msg)
			if err != nil {
				slog.Error("WebSocket write error", "userID", sub.userID, "roomID", sub.roomID, "error", err)
				return
			}
		case <-ctx.Done(): // Connection closed by client or server (CloseRead)
			slog.Info("WebSocket context done (write loop)", "userID", sub.userID, "roomID", sub.roomID, "error", ctx.Err())
			readErr := <-errc

			if readErr != nil && !errors.Is(readErr, context.Canceled) {
				slog.Error("WebSocket closed due to read error", "userID", sub.userID, "roomID", sub.roomID, "error", readErr)
			}
			return
		}
	}
}
