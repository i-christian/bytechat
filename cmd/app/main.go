package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bytechat/internal/server"
)

func gracefulShutdown(appServer *server.Server, httpServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	slog.Info("shutting down gracefully, press Ctrl+C again to force")

	// Shutting down database connection
	if err := appServer.CloseDbConn(); err != nil {
		slog.Info("Database connection pool closed successfully")
	}

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		slog.Info("Server forced to shutdown with error, ", "Message", err.Error())
	}

	slog.Info("Server exiting...")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	appServer, httpServer := server.NewServer()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go func() {
		gracefulShutdown(appServer, httpServer, done)
	}()

	log.Printf("The server is starting on: http://%s:%s\n", os.Getenv("DOMAIN"), os.Getenv("PORT"))

	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	slog.Info("Graceful shutdown complete.")
}
