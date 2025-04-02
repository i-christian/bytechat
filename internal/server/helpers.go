package server

import (
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"net/http"

	"github.com/a-h/templ"
)

// renderDashboardComponent renders a component either as a full dashboard page
// (when not an HTMX request) or just the component (when it's an HTMX request).
func (s *Server) renderComponent(w http.ResponseWriter, r *http.Request, children templ.Component) {
	if r.Header.Get("HX-Request") == "true" {
		if err := children.Render(r.Context(), w); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			slog.Error("Failed to render dashboard component", "error", err)
		}
	}
}

// Generate a random 6-digit numeric password
func generateNumericPassword() (string, error) {
	const passwordLength = 6
	password := ""

	for i := 0; i < passwordLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		password += fmt.Sprintf("%d", num.Int64())
	}

	return password, nil
}
