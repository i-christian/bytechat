package server

import (
	"context"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"bytechat/internal/database"

	"github.com/google/uuid"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

type Server struct {
	queries   *database.Queries
	conn      *pgxpool.Pool
	SecretKey []byte
	port      int
}

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

// Setup database migrations and closes database connection afterwards
func setUpMigration() {
	db, err := sql.Open(os.Getenv("GOOSE_DRIVER"), os.Getenv("DB_URL"))
	if err != nil {
		slog.Error("Failed to open database for migration", "msg", err.Error())
		return
	}

	defer db.Close()
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("postgres"); err != nil {
		slog.Error("Failed to select postgres database", "msg", err.Error())
	}

	if err := goose.Up(db, os.Getenv("GOOSE_MIGRATION_DIR")); err != nil {
		slog.Error("Unable to run migrations:\n", "error", err.Error())
	}
}

// Checks if required env vars are all set during server startup
func validateEnvVars() {
	requiredVars := []string{"DB_URL", "PORT", "RANDOM_HEX", "DOMAIN", "RANDOM_HEX", "PROJECT_NAME", "GOOSE_DRIVER", "GOOSE_MIGRATION_DIR", "SUPERUSER_ROLE", "SUPERUSER_EMAIL", "SUPERUSER_PASSWORD", "ENV"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			slog.Error(fmt.Sprintf("Environment variable %s is required", v))
			os.Exit(1)
		}
	}
}

func NewServer() (*Server, *http.Server) {
	validateEnvVars()
	setUpMigration()

	SecretKey, err := hex.DecodeString(os.Getenv("RANDOM_HEX"))
	if err != nil {
		slog.Error(err.Error())
	}

	ctx := context.Background()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	conn, err := pgxpool.New(ctx, os.Getenv("DB_URL"))
	if err != nil {
		slog.Error("Unable to connect to database: \n", "error", err.Error())
		os.Exit(1)
	}

	generatedQeries := database.New(conn)
	createSuperUser(ctx, generatedQeries)

	AppServer := &Server{
		port:      port,
		conn:      conn,
		queries:   generatedQeries,
		SecretKey: SecretKey,
	}

	// Declare Server config
	httpserver := &http.Server{
		Addr:         fmt.Sprintf(":%d", AppServer.port),
		Handler:      AppServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return AppServer, httpserver
}

func createSuperUser(ctx context.Context, queries *database.Queries) {
	role := os.Getenv("SUPERUSER_ROLE")
	email := os.Getenv("SUPERUSER_EMAIL")
	password := os.Getenv("SUPERUSER_PASSWORD")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Failed to hash password")
		os.Exit(1)
	}

	adminUser, err := queries.GetUserByEmail(ctx, email)

	if adminUser.UserID != uuid.Nil {
		slog.Info("Superuser already exists")
		return
	}

	user := database.CreateUserParams{
		FirstName: "Admin",
		LastName:  "Admin",
		Email:     email,
		Gender:    "M",
		Password:  string(hashedPassword),
		Name:      role,
	}

	_, err = queries.CreateUser(ctx, user)
	if err != nil {
		slog.Error("Failed to create superuser:", "error", err.Error())
	} else {
		slog.Info("Superuser created successfully")
	}
}
