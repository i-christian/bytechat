package web

import (
	"embed"
	"os"
)

//go:embed "assets"
var Files embed.FS

// Define the project name
func ProjectName() string {
	return os.Getenv("PROJECT_NAME")
}
