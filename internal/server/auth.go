package server

import (
	"log/slog"
	"net/http"
	"time"

	"bytechat/internal/cookies"

	"bytechat/cmd/web"
	"bytechat/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type LoginUser struct {
	Password string
	UserID   uuid.UUID
}

// LoginHandler authenticates the user and creates a session.
func (s *Server) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, "bad request")
		return
	}

	identifier := r.FormValue("identifier")
	password := r.FormValue("password")

	user, err := s.queries.GetUserByEmail(r.Context(), identifier)
	if err != nil {
		slog.Error("login request denied", "user name", identifier, "password", password, "error", err.Error())
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Create a new sessionID and expiry.
	sessionID := uuid.New()
	expiry := pgtype.Timestamptz{Time: time.Now().Add(2 * 7 * 24 * time.Hour), Valid: true}

	sessionParams := database.CreateSessionParams{
		SessionID: sessionID,
		UserID:    user.UserID,
		Expires:   expiry,
	}

	returnedSession, err := s.queries.CreateSession(r.Context(), sessionParams)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		slog.Error("Failed to create session", "error", err.Error())
		return
	}

	cookie := createSessionCookie(sessionID)

	if err := cookies.WriteEncrypted(w, cookie, s.SecretKey); err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	currentUserRole, err := s.queries.GetRedirectPath(r.Context(), returnedSession)

	var redirectPath string
	if currentUserRole == "admin" {
		redirectPath = "/admin"
	} else {
		redirectPath = "/"
	}

	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("HX-Redirect", redirectPath)
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, redirectPath, http.StatusFound)
}

// LogoutHandler to log users out
func (s *Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(User)
	if !ok {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
	}

	if err := s.queries.DeleteSession(r.Context(), user.UserID); err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionid",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (s *Server) LogoutConfirmHandler(w http.ResponseWriter, r *http.Request) {
	s.renderComponent(w, r, web.LogoutConfirmHandler())
}

func (s *Server) LogoutCancelHandler(w http.ResponseWriter, r *http.Request) {
	redirectPath := r.Referer()
	http.Redirect(w, r, redirectPath, http.StatusFound)
}
