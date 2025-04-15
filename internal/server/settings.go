package server

import (
	"log/slog"
	"net/http"
	"strings"

	"bytechat/cmd/web/settings"
	"bytechat/internal/database"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) ShowUserSettings(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(User)
	if !ok {
		writeError(w, http.StatusUnauthorized, "not logged in")
		slog.Error("user not logged in, failed to read userID from userContextKey")
		return
	}

	userDetails, err := s.queries.GetUserDetails(r.Context(), user.UserID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch user")
		slog.Error("failed to fetch user", "UserID", user.UserID, "error", err.Error())
		return
	}

	s.renderComponent(w, r, settings.UserSettings(userDetails))
}

// EditUserProfile updates user information
// expects form data with user information from
func (s *Server) EditUserProfile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userContextKey).(User)
	if !ok {
		writeError(w, http.StatusUnauthorized, "not logged in")
		slog.Error("user not logged in, failed to read userID from userContextKey")
		return
	}

	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	gender := r.FormValue("gender")
	currentPassword := r.FormValue("current_password")
	newPassword := r.FormValue("new_password")
	confirmPassword := r.FormValue("confirm_password")

	if len(strings.TrimSpace(currentPassword)) > 0 {
		if confirmPassword != newPassword {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`
					<div id="popover" class="custom-popover show" style="background-color: #dc2626;">
						<span>❌ New Password does not match confirmed password</span>
					</div>
					<script>
						setTimeout(() => {
							document.getElementById('popover').classList.add('hide');
							setTimeout(() => document.getElementById('popover').remove(), 500);
						}, 3000);
					</script>
				`))
			return

		}
		userDetails, err := s.queries.GetUserDetails(r.Context(), user.UserID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to get user")
			slog.Error("failed to get user", "error", err.Error())
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(currentPassword)); err != nil {
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte(`
					<div id="popover" class="custom-popover show" style="background-color: #dc2626;">
						<span>❌ Incorrect Current Password</span>
					</div>
					<script>
						setTimeout(() => {
							document.getElementById('popover').classList.add('hide');
							setTimeout(() => document.getElementById('popover').remove(), 500);
						}, 3000);
					</script>
				`))
			return
		}

	}

	tx, err := s.conn.Begin(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	defer tx.Rollback(r.Context())
	qtx := s.queries.WithTx(tx)

	updateInfo := database.EditMyProfileParams{
		UserID:    user.UserID,
		FirstName: firstName,
		LastName:  lastName,
		Gender:    gender,
		Email:     email,
	}

	err = qtx.EditMyProfile(r.Context(), updateInfo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if len(strings.TrimSpace(newPassword)) > 0 {
		hashedPassword, err := hashPassword(newPassword)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		changePasswdParams := database.EditMyPasswordParams{
			UserID:   user.UserID,
			Password: string(hashedPassword),
		}

		err = qtx.EditMyPassword(r.Context(), changePasswdParams)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			slog.Error("failed to change password", ":", err.Error())
			return
		}
	}

	tx.Commit(r.Context())

	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`
			<div id="popover" class="custom-popover show" style="background-color: #16a34a;">
				<span>✅ Profile updated successfully</span>
			</div>
			<script>
				setTimeout(() => {
					document.getElementById('popover').classList.add('hide');
					setTimeout(() => document.getElementById('popover').remove(), 500);
				}, 3000);
				window.location.reload()
			</script>
		`))
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
