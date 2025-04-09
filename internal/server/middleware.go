package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"bytechat/internal/cookies"
	"bytechat/internal/database"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type contextKey string

const (
	sessionIDKey   contextKey = "session_id"
	userContextKey contextKey = "user"
)

// refreshSession method updates a near session in the database if its near expiry
func (s *Server) refreshSession(ctx context.Context, session database.GetSessionRow) (uuid.UUID, error) {
	newExpiry := pgtype.Timestamptz{Time: time.Now().Add(2 * 7 * 24 * time.Hour), Valid: true}
	newSessionID := uuid.New()

	refreshParams := database.RefreshSessionParams{
		UserID:    session.UserID,
		Expires:   newExpiry,
		SessionID: newSessionID,
	}

	if err := s.queries.RefreshSession(ctx, refreshParams); err != nil {
		return uuid.Nil, err
	}

	return newSessionID, nil
}

// AuthMiddleware ensures the user is authenticated and loads the user (with role) into the context.
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := cookies.ReadEncrypted(r, "sessionid", s.SecretKey)
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) || errors.Is(err, cookies.ErrInvalidValue) {
				http.Redirect(w, r, "/login", http.StatusFound)
				return
			}
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		parsedSessionID, err := uuid.Parse(sessionID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid session ID")
			return
		}

		session, err := s.queries.GetSession(r.Context(), parsedSessionID)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		if session.Expires.Valid && session.Expires.Time.Before(time.Now()) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		timeLeft := time.Until(session.Expires.Time)
		if timeLeft < 24*time.Hour {
			newSessionID, err := s.refreshSession(r.Context(), session)
			if err != nil {
				writeError(w, http.StatusInternalServerError, "internal server error")
				return
			}

			cookie := http.Cookie{
				Name:     "sessionid",
				Value:    newSessionID.String(),
				Path:     "/",
				MaxAge:   3600 * 24 * 7 * 2, // 2 weeks
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			}

			if err := cookies.WriteEncrypted(w, cookie, s.SecretKey); err != nil {
				writeError(w, http.StatusInternalServerError, "internal server error")
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), sessionIDKey, newSessionID))
		} else {
			r = r.WithContext(context.WithValue(r.Context(), sessionIDKey, session.SessionID))
		}

		user := User{
			UserID:    session.UserID,
			Role:      session.Role,
			LastName:  session.FirstName,
			FirstName: session.FirstName,
		}

		r = r.WithContext(context.WithValue(r.Context(), userContextKey, user))

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// RequireRoles returns a middleware that allows access only if the user's role
// is in the provided list of allowed roles.
func (s *Server) RequireRoles(allowedRoles ...string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{})
	for _, role := range allowedRoles {
		allowed[role] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(userContextKey).(User)
			if !ok {
				writeError(w, http.StatusUnauthorized, "user not authenticated")
				return
			}

			if _, exists := allowed[user.Role]; !exists {
				writeError(w, http.StatusForbidden, "forbidden")
				return
			}

			w.Header().Add("Cache-Control", "no-store")
			next.ServeHTTP(w, r)
		})
	}
}

// RedirectIfAuthenticated checks if a user is already logged in and redirects them to the home page
func (s *Server) RedirectIfAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := cookies.ReadEncrypted(r, "sessionid", s.SecretKey)
		if err == nil {
			parsedSessionID, parseErr := uuid.Parse(sessionID)

			if parseErr == nil {
				session, getSessionErr := s.queries.GetSession(r.Context(), parsedSessionID)

				if getSessionErr == nil && session.Expires.Valid && session.Expires.Time.After(time.Now()) {
					var redirectPath string
					if session.Role == "admin" {
						redirectPath = "/admin"
					} else {
						redirectPath = "/"
					}

					http.Redirect(w, r, redirectPath, http.StatusFound)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}

// secureHeaders middleware automatically adds HTTP security headers to every response inline with current OWASP guidance.
// This middleware acts on every request that is received, and needs it to be executed before a request hits our servemux(router)
func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline' fonts.googleapis.com; font-src 'self' data: fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}
