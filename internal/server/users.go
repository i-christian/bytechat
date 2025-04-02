package server

import (
	"log/slog"
	"net/http"

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
	phoneNumber := r.FormValue("phone_number")
	email := r.FormValue("email")
	gender := r.FormValue("gender")
	role := r.FormValue("role") // Role name

	if firstName == "" || lastName == "" || phoneNumber == "" || gender == "" || role == "" {
		writeError(w, http.StatusBadRequest, "all fields except email are required")
		return
	}

	// Generate a 6-digit numeric password
	password, err := generateNumericPassword()
	if err != nil {
		slog.Error("Failed to generate password")
		writeError(w, http.StatusInternalServerError, "internal server error")
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
