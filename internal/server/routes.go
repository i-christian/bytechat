package server

import (
	"net/http"
	"os"

	"bytechat/cmd/web"
	"bytechat/cmd/web/auth"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Global middlewares
	r.Use(middleware.CleanPath)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(secureHeaders)
	r.Use(middleware.Compress(5, "text/html", "text/css"))
	r.Use(middleware.Recoverer)

	// CORS setup
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("DOMAIN")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Static file server
	fileServer := http.FileServer(http.FS(web.Files))
	r.Handle("/assets/*", fileServer)

	// PUBLIC ROUTES
	r.Group(func(r chi.Router) {
		r.With(s.RedirectIfAuthenticated).Get("/login", templ.Handler(auth.Login()).ServeHTTP)
		r.Get("/register", templ.Handler(auth.CreateUserForm()).ServeHTTP)
		r.Post("/register", s.Register)
		r.Post("/login", s.LoginHandler)
	})

	// AUTHENTICATED USER ROUTES
	r.Group(func(r chi.Router) {
		r.Use(s.AuthMiddleware)

		r.Get("/", s.showDashboardHome)
		r.Get("/logout/confirm", s.LogoutConfirmHandler)
		r.Get("/logout/cancel", s.LogoutCancelHandler)
		r.Post("/logout", s.LogoutHandler)
		r.Get("/profile", s.userProfile)
	})

	// USER MANAGEMENT (ADMIN)
	r.Route("/users", func(r chi.Router) {
		r.Use(s.AuthMiddleware)
		r.Use(s.RequireRoles("admin"))

		r.Get("/{id}/edit", s.ShowEditUserForm)
		r.Put("/{id}", s.EditUser)

		r.Get("/{id}/delete", s.ShowDeleteConfirmation)
		r.Delete("/{id}", s.DeleteUser)
	})

	// ADMIN DASHBOARD (ADMIN)
	r.Route("/admin", func(r chi.Router) {
		r.Use(s.AuthMiddleware)
		r.Use(s.RequireRoles("admin"))

		r.Get("/", s.ListUsers)
	})

	// CHAT PAGE
	r.Route("/chat", func(r chi.Router) {
		r.Use(s.AuthMiddleware)

		r.Get("/{room_id}", s.showSpecificChatPage)
		r.Get("/{room_id}/members", s.roomMembers)
		r.Get("/ws/{room_id}", s.handleWebSocket)
		r.Get("/dm", s.showPrivateChats)
		r.Get("/dm/{friend_id}", s.createPrivateRoom)
	})

	// user settings routes
	r.Route("/settings", func(r chi.Router) {
		r.Use(s.AuthMiddleware)

		r.Get("/", s.ShowUserSettings)
		r.Put("/", s.EditUserProfile)
	})

	return r
}
