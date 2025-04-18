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

// createPrivateRoom method creates a private room for two users.
// It accepts {friend_id}
func (s *Server) createPrivateRoom(w http.ResponseWriter, r *http.Request) {
	friendID, err := uuid.Parse(r.PathValue("friend_id"))
	if err != nil {
		slog.Error("Invalid room ID format", "roomID", r.PathValue("friend_id"), "error", err)
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	roomName := r.URL.Query().Get("room_name")

	user, ok := r.Context().Value(userContextKey).(User)
	if !ok || user.UserID == uuid.Nil {
		slog.Warn("Chat page access attempt without authentication")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	params := database.InitiatePrivateRoomParams{
		UserA: user.UserID,
		UserB: friendID,
	}

	roomID, err := s.queries.InitiatePrivateRoom(r.Context(), params)
	if err != nil {
		slog.Error("Failed to Create Private Room", "UserA", user.UserID, "UserB", friendID)
		return
	}

	if !roomID.Valid {
		slog.Error("internal server error", "dbRoomID", roomID)
		http.Redirect(w, r, r.Referer(), http.StatusFound)
	}

	http.Redirect(w, r, fmt.Sprintf("/chat/%v?room_name=%s", roomID, roomName), http.StatusFound)
}

// roomMembers displays the latest members in a given chat room
func (s *Server) roomMembers(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(r.PathValue("room_id"))
	if err != nil {
		slog.Error("Invalid room ID format", "roomID", r.PathValue("room_id"), "error", err)
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	roomMembersParams := database.GetUsersInRoomParams{
		RoomID: roomID,
		Limit:  20,
	}
	roomMembers, err := s.queries.GetUsersInRoom(r.Context(), roomMembersParams)
	if err != nil {
		slog.Error("Could not fetch room details", "roomID", roomID, "error", err.Error())
	}

	user, ok := r.Context().Value(userContextKey).(User)
	if !ok || user.UserID == uuid.Nil {
		slog.Warn("Chat page access attempt without authentication")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	s.renderComponent(w, r, chat.RoomMembers(user.UserID, roomMembers))
}

// showSpecificChatPage renders the chat page for a given room ID.
func (s *Server) showSpecificChatPage(w http.ResponseWriter, r *http.Request) {
	roomIDStr := r.PathValue("room_id")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		slog.Error("Invalid room ID format", "roomID", roomIDStr, "error", err)
		writeError(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	roomName := r.URL.Query().Get("room_name")

	user, ok := r.Context().Value(userContextKey).(User)
	if !ok || user.UserID == uuid.Nil {
		slog.Warn("Chat page access attempt without authentication")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	joinParams := database.JoinRoomParams{
		UserID: user.UserID,
		RoomID: roomID,
	}

	err = s.queries.JoinRoom(r.Context(), joinParams)
	if err != nil {
		slog.Error("Failed to join room", "user id", user.UserID, "room id", roomID, "error", err.Error())
	}

	roomDetails, err := s.queries.GetRoomDetails(r.Context(), roomID)
	if err == nil && roomDetails.RoomType == "public" {
		roomName = roomDetails.Name
	} else {
		slog.Warn("Could not fetch room details", "roomID", roomID, "error", err)
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
		RoomType:        roomDetails.RoomType,
		InitialMessages: initialMessages,
		CurrentUserID:   user.UserID,
	}

	s.renderComponent(w, r, chat.ChatPage(pageData))
}

// showPrivateChats method handles direct messages
func (s *Server) showPrivateChats(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(User)
	if !ok || user.UserID == uuid.Nil {
		slog.Warn("attempt to retrieve private chats")
		writeError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		return
	}

	privateRooms, err := s.queries.GetPrivateRooms(r.Context(), user.UserID)
	if err != nil {
		slog.Error("Failed to get private rooms", "error", err.Error())
		writeError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}

	s.renderComponent(w, r, chat.Messages(user.UserID, privateRooms))
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s.addSubscriber(ctx, roomID, sub)
	defer s.deleteSubscriber(ctx, roomID, sub)

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
				dberr := s.queries.UserOfflineStatus(ctx, user.UserID)
				if dberr != nil {
					slog.Warn("User Status Update failed", "warning", dberr.Error())
				}

				if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
					websocket.CloseStatus(err) == websocket.StatusGoingAway ||
					errors.Is(err, context.Canceled) ||
					errors.Is(err, net.ErrClosed) {
					slog.Info("WebSocket reader closed gracefully", "userID", user.UserID, "roomID", roomID)
				} else {
					slog.Error("WebSocket read error", "userID", user.UserID, "roomID", roomID, "error", err.Error())
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
