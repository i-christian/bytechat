package server

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"bytechat/cmd/web/chat"
	"bytechat/internal/database"

	"github.com/coder/websocket"
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

	// Optional: Fetch initial messages for the room
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

// publish broadcasts a message (now HTML bytes) to all subscribers in a specific room.
func (s *Server) publish(roomID uuid.UUID, msgHTML []byte) {
	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	subs, ok := s.subscribersByRoom[roomID]
	if !ok {
		slog.Warn("Attempted to publish to room with no subscribers", "roomID", roomID)
		return
	}

	slog.Info("Publishing message HTML", "roomID", roomID, "subscribers", len(subs), "html_len", len(msgHTML))
	for sub := range subs {
		htmlToSend := make([]byte, len(msgHTML))
		copy(htmlToSend, msgHTML)

		select {
		case sub.msgs <- htmlToSend:
		default:
			slog.Warn("Subscriber channel full, closing slow connection", "userID", sub.userID, "roomID", roomID)
			go sub.closeSlow()
		}
	}
}

// publishHandler reads the request body with a limit o 8192 and then publishes the received message
func (s *Server) publishHandler(w http.ResponseWriter, r *http.Request) {
	roomIDStr := r.PathValue("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		slog.Error("Invalid room ID format", "roomID", roomIDStr, "error", err)
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	user, ok := r.Context().Value(userContextKey).(User)
	if !ok || user.UserID == uuid.Nil {
		slog.Warn("Publish attempt without authentication")
		writeError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	err = r.ParseForm()
	if err != nil {
		slog.Error("Failed to parse form", "error", err)
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	msgText := r.FormValue("message")
	if msgText == "" {
		slog.Warn("Attempted to send empty message", "userID", user.UserID, "roomID", roomID)
		writeError(w, http.StatusBadRequest, "Message cannot be empty")
		return
	}

	messageParams := database.CreateMessageParams{
		UserID: user.UserID,
		RoomID: roomID,
		Text:   pgtype.Text{String: msgText, Valid: true},
	}
	dbMsg, err := s.queries.CreateMessage(r.Context(), messageParams)
	if err != nil {
		slog.Error("Failed to save message to DB", "userID", user.UserID, "roomID", roomID, "error", err)
		writeError(w, http.StatusInternalServerError, "Could not save message")
		return
	}
	slog.Info("Message saved to DB", "messageID", dbMsg.MessageID)

	messageTimestamp := dbMsg.CreatedAt.Time
	if messageTimestamp.IsZero() {
		messageTimestamp = time.Now().UTC()
	}

	chatMessageData := chat.ChatMessageData{
		Text: msgText,
		Sender: chat.SenderInfo{
			ID:        user.UserID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
		Timestamp:     messageTimestamp,
		IsCurrentUser: false,
	}

	var buf bytes.Buffer
	s.renderComponent(w, r, chat.ChatMessage(chatMessageData))
	htmlBytes := buf.Bytes()

	s.publish(roomID, htmlBytes)

	w.WriteHeader(http.StatusOK)
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
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
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

	ctx, cancel := context.WithTimeout(r.Context(), time.Minute*10)
	defer cancel()

	ctx = c.CloseRead(ctx)
	errc := make(chan error, 1)
	defer close(errc)

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
