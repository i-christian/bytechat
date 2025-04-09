package server

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"bytechat/cmd/web/chat"
	"bytechat/internal/database"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// showDashboardHome renders a default view for the authenticated root path
func (s *Server) showDashboardHome(w http.ResponseWriter, r *http.Request) {
	s.renderComponent(w, r, chat.HomePlaceholder())
}

// showChatRooms renders a list of available public chatrooms
func (s *Server) showChatRooms(w http.ResponseWriter, r *http.Request) {
}

// showSpecificChatPage renders the chat page for a given room ID.
func (s *Server) showSpecificChatPage(w http.ResponseWriter, r *http.Request) {
	roomIDStr := r.PathValue("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		slog.Error("Invalid room ID format", "roomID", roomIDStr, "error", err)
		http.Error(w, "Invalid Room ID", http.StatusBadRequest)
		return
	}

	user, ok := r.Context().Value(userContextKey).(User)
	if !ok || user.UserID == uuid.Nil {
		slog.Warn("Chat page access attempt without authentication")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	roomDetails, err := s.queries.GetRoomDetails(r.Context(), roomID)
	roomName := "Chat Room"
	if err == nil {
		roomName = roomDetails
	} else {
		slog.Warn("Could not fetch room details", "roomID", roomID, "error", err)
		return
	}

	initialMessages, err := s.queries.ListMessagesByRoom(r.Context(), database.ListMessagesByRoomParams{
		RoomID: roomID,
		Limit:  50,
	})
	if err != nil {
		slog.Error("Failed to fetch initial messages", "roomID", roomID, "error", err)
		initialMessages = []database.ListMessagesByRoomRow{}
	}

	pageData := chat.ChatPageData{
		RoomID:          roomID,
		RoomName:        roomName,
		InitialMessages: initialMessages,
		CurrentUserID:   user.UserID,
	}

	s.renderComponent(w, r, chat.ChatPage(pageData))
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

// publish broadcasts a message HTML bytes to all subscribers in a specific room.
func (s *Server) publish(roomID uuid.UUID, payload broadcastPayload) {
	s.subscribersMu.Lock()
	subsCopy := make([]*subscriber, 0)
	if subs, ok := s.subscribersByRoom[roomID]; ok {
		for sub := range subs {
			subsCopy = append(subsCopy, sub)
		}
	}
	s.subscribersMu.Unlock()

	if len(subsCopy) == 0 {
		slog.Warn("Attempted to publish data to room with no subscribers", "roomID", roomID)
		return
	}

	for _, sub := range subsCopy {
		select {
		case sub.msgs <- payload:
			// Payload sent successfully
		default:
			// Subscriber's buffer is full
			go sub.closeSlow()
		}
	}
}

// Helper function for writing messages with timeout
func writeTimeout(ctx context.Context, timeout time.Duration, c *websocket.Conn, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.Write(ctx, websocket.MessageText, msg)
}

// handleWebSocket is the HTTP handler responsible for upgrading a connection
// to the WebSocket protocol and managing the full lifecycle of that connection for chat.
// It requires the user to be authenticated and expects a room ID in the URL path.
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
		msgs:   make(chan broadcastPayload, s.subscriberMessageBuffer),
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errc := make(chan error, 1)

	go func() {
		defer close(errc)
		for {
			var incomingMsg struct {
				Headers interface{} `json:"HEADERS"`
				Text    string      `json:"text"`
			}

			err = wsjson.Read(ctx, c, &incomingMsg)
			if err != nil {
				if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
					websocket.CloseStatus(err) == websocket.StatusGoingAway ||
					errors.Is(err, context.Canceled) ||
					errors.Is(err, net.ErrClosed) {
					slog.Info("WebSocket reader closed gracefully", "userID", user.UserID, "roomID", roomID, "error", err)
				} else {
					slog.Error("WebSocket read error", "userID", user.UserID, "roomID", roomID, "error", err)
				}
				return
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

			payload := broadcastPayload{
				dbMessage: dbMsg,
				sender:    user,
			}

			s.publish(sub.roomID, payload)

		}
	}()

	// Goroutine to write outgoing messages to the client
	for {
		select {
		case payload := <-sub.msgs:
			messageData := chat.ChatMessageData{
				Text: payload.dbMessage.Text.String,
				Sender: chat.SenderInfo{
					ID:        payload.sender.UserID,
					FirstName: payload.sender.FirstName,
					LastName:  payload.sender.LastName,
				},
				Timestamp:     payload.dbMessage.CreatedAt.Time,
				IsCurrentUser: payload.sender.UserID == sub.userID,
			}

			// Render the templ component to HTML
			var buf bytes.Buffer
			err := chat.ChatMessage(messageData).Render(context.Background(), &buf)
			if err != nil {
				slog.Error("Failed to render ChatMessage component", "userID", sub.userID, "roomID", sub.roomID, "error", err)
				continue
			}

			singleMessageHTML := buf.String()

			// fixing the HTMX web socket swapping problem
			// We have to wrap the response with an id which we are targeting. And also add hx-swap-oob=beforeend attribute
			oobWrapperHTML := fmt.Sprintf(
				`<div id="chat-messages" hx-swap-oob="beforeend">%s</div>`,
				singleMessageHTML,
			)

			htmlBytes := []byte(oobWrapperHTML)

			err = writeTimeout(ctx, time.Second*5, c, htmlBytes)
			if err != nil {
				slog.Error("WebSocket write error", "userID", sub.userID, "roomID", sub.roomID, "error", err)
				select {
				case readErr := <-errc:
					slog.Error("Write loop terminating due to read error", "readError", readErr)
				default:
				}
				return
			} else {
				slog.Info("HTML message sent successfully", "userID", sub.userID, "roomID", sub.roomID, "bytes", len(htmlBytes))
			}

		case <-ctx.Done(): // Connection closed by client or server
			slog.Error("WebSocket context done (write loop)", "userID", sub.userID, "roomID", sub.roomID, "error", ctx.Err())
			select {
			case readErr := <-errc:
				if readErr != nil && !errors.Is(readErr, context.Canceled) && !errors.Is(readErr, net.ErrClosed) && websocket.CloseStatus(readErr) < 0 {
					slog.Error("WebSocket closed due to read error", "userID", sub.userID, "roomID", sub.roomID, "error", readErr)
				} else {
					slog.Warn("Read loop closed gracefully or context canceled", "userID", sub.userID, "roomID", sub.roomID, "error", readErr)
				}
			default:
				slog.Info("Write loop context done without prior read error signal")
			}
			return

		}
	}
}
