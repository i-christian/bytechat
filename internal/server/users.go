package server

import (
	"log/slog"
	"net/http"
	"strings"

	"bytechat/cmd/web/components"
	"bytechat/cmd/web/users"
	"bytechat/internal/database"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// hashPassword accepts a string and returns a hashed password
func hashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

// An endpoint to create a new user account
func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusUnprocessableEntity, "failed to parse form")
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	gender := r.FormValue("gender")
	role := "user"

	if firstName == "" || lastName == "" || email == "" || gender == "" || password == "" {
		writeError(w, http.StatusBadRequest, "all fields are required")
		return
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	caser := cases.Title(language.English)
	user := database.CreateUserParams{
		FirstName: caser.String(firstName),
		LastName:  caser.String(lastName),
		Email:     email,
		Gender:    gender,
		Password:  string(hashedPassword),
		Name:      role,
	}

	_, err = s.queries.CreateUser(r.Context(), user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		slog.Info("Failed to create user", "message:", err.Error())
		return
	}

	w.Header().Set("Content-Type", "text/html")

	http.Redirect(w, r, "/login", http.StatusFound)
}

// userProfile handler method returns user current logged in user details
func (s *Server) userProfile(w http.ResponseWriter, r *http.Request) {
	contextUser, ok := r.Context().Value(userContextKey).(User)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorised")
		return
	}

	dbUser, err := s.queries.GetUserDetails(r.Context(), contextUser.UserID)
	if err != nil {
		var status int
		var errorMessage string
		if strings.Contains(err.Error(), "unauthorized") {
			status = http.StatusUnauthorized
			errorMessage = "user not authorised"
		} else {
			status = http.StatusInternalServerError
			errorMessage = "internal server error"
		}
		writeError(w, status, errorMessage)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	user := components.User{
		UserID:    dbUser.UserID,
		FirstName: dbUser.FirstName,
		LastName:  dbUser.LastName,
		Gender:    dbUser.Gender,
		Email:     dbUser.Email,
	}
	component := components.UserDetails(user)
	s.renderComponent(w, r, component)
}

// ListUsers handler retrieves all users from the database.
func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	userList, err := s.queries.ListUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	contents := users.UsersList(userList)
	s.renderComponent(w, r, contents)
}

// ShowEditUserForm fetches the user by id and renders the edit modal.
func (s *Server) ShowEditUserForm(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "invalid user id")
		return
	}

	user, err := s.queries.GetUserDetails(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		slog.Error("User not found", "message:", err.Error())
		return
	}

	s.renderComponent(w, r, users.EditUserModal(user))
}

// EditUser handler
// Update user information
// expects form data with user information from
func (s *Server) EditUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "invalid user id")
	}

	if err := r.ParseForm(); err != nil {
		writeError(w, http.StatusBadRequest, "Bad Request")
		return
	}

	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	gender := r.FormValue("gender")
	password := r.FormValue("password")

	tx, err := s.conn.Begin(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	defer tx.Rollback(r.Context())
	qtx := s.queries.WithTx(tx)

	updateInfo := database.EditUserParams{
		UserID:    userID,
		FirstName: firstName,
		LastName:  lastName,
		Gender:    gender,
		Email:     email,
	}

	err = qtx.EditUser(r.Context(), updateInfo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if len(strings.TrimSpace(password)) > 0 {
		hashedPassword, err := hashPassword(password)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		changePasswdParams := database.EditPasswordParams{
			UserID:   userID,
			Password: string(hashedPassword),
		}

		err = qtx.EditPassword(r.Context(), changePasswdParams)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal server error")
			slog.Error("failed to change password", ":", err.Error())
			return
		}
	}

	tx.Commit(r.Context())

	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("HX-Redirect", "/dashboard/userlist")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/dashboard/userlist", http.StatusFound)
}

// ShowDeleteConfirmation renders the delete confirmation modal, passing the user id.
func (s *Server) ShowDeleteConfirmation(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	s.renderComponent(w, r, users.DeleteConfirmationModal(userID))
}

// DeleteUser handler
// Accepts an id parameter
// deletes a user from the database
func (s *Server) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, "failed to parse user id")
		return
	}

	err = s.queries.DeleteUser(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("HX-Redirect", "/dashboard/userlist")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/dashboard/userlist", http.StatusFound)
}
