package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"bytechat/cmd/web"

	"github.com/a-h/templ"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

// renderDashboardComponent renders a component either as a full dashboard page
// (when not an HTMX request) or just the component (when it's an HTMX request).
func (s *Server) renderComponent(w http.ResponseWriter, r *http.Request, children templ.Component) {
	if r.Header.Get("HX-Request") == "true" {
		if err := children.Render(r.Context(), w); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			slog.Error("Failed to render dashboard component", "error", err)
		}
	} else {
		userRole, ok := r.Context().Value(userContextKey).(User)
		if !ok {
			writeError(w, http.StatusUnauthorized, "unauthorised")
			return
		}
		user := web.DashboardUserRole{
			Role: userRole.Role,
		}

		rooms, err := s.queries.ListPublicRooms(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "no rooms found")
		}

		ctx := templ.WithChildren(r.Context(), children)
		if err := web.Dashboard(user, rooms).Render(ctx, w); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			slog.Error("Failed to render dashboard layout", "error", err)
		}
	}
}

// addSubscriber registers a subscriber for a given room.
func (s *Server) addSubscriber(ctx context.Context, roomID uuid.UUID, sub *subscriber) {
	err := s.queries.UserOnlineStatus(ctx, sub.userID)
	if err != nil {
		slog.Warn("User Status Update failed", "warning", err.Error())
	}

	s.subscribersMu.Lock()
	defer s.subscribersMu.Unlock()

	if _, ok := s.subscribersByRoom[roomID]; !ok {
		s.subscribersByRoom[roomID] = make(map[*subscriber]struct{})
	}
	s.subscribersByRoom[roomID][sub] = struct{}{}
	slog.Info("Subscriber added", "userID", sub.userID, "roomID", roomID)
}

// deleteSubscriber removes a subscriber from a given room.
func (s *Server) deleteSubscriber(ctx context.Context, roomID uuid.UUID, sub *subscriber) {
	err := s.queries.UserOfflineStatus(ctx, sub.userID)
	if err != nil {
		slog.Warn("User Status Update failed", "warning", err.Error())
	}

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
