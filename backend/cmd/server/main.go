package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"sekaitext/backend/internal/api"
	"sekaitext/backend/internal/config"
)

func main() {
	port := flag.Int("port", 9800, "server port")
	dir := flag.String("dir", ".", "base directory for resources")
	flag.Parse()

	// Resolve base directory:
	// - If --dir explicitly provided, use it as-is (relative to CWD).
	// - If default "." and resources aren't found, fall back to the executable's
	//   directory (for Tauri sidecar deployment).
	baseDir := *dir
	if baseDir == "." {
		// Check if CWD has resources/ directory
		if _, err := os.Stat(filepath.Join(".", "resources", "setting")); os.IsNotExist(err) {
			// Fall back to executable directory for sidecar deployment
			exe, err := os.Executable()
			if err == nil {
				baseDir = filepath.Dir(exe)
			}
		}
	}

	cfg := config.NewAppConfig(baseDir)

	// Ensure resource directories exist
	ensureDir(cfg.SettingDir)
	ensureDir(cfg.DataDir)
	ensureDir(cfg.ImagesChrDir)

	router := api.NewRouter(cfg)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("SekaiText server starting on %s", addr)
	log.Printf("Base directory: %s", baseDir)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func ensureDir(path string) {
	if err := os.MkdirAll(path, 0755); err != nil {
		log.Printf("Warning: could not create directory %s: %v", path, err)
	}
}
