package config

import (
	"path/filepath"
)

// AppConfig holds all path configurations relative to a base directory.
type AppConfig struct {
	BaseDir      string
	SettingDir   string
	DataDir      string
	ImagesDir    string
	ImagesChrDir string
}

// NewAppConfig creates an AppConfig with paths relative to the given base directory.
func NewAppConfig(baseDir string) *AppConfig {
	return &AppConfig{
		BaseDir:      baseDir,
		SettingDir:   filepath.Join(baseDir, "resources", "setting"),
		DataDir:      filepath.Join(baseDir, "resources", "data"),
		ImagesDir:    filepath.Join(baseDir, "resources", "images"),
		ImagesChrDir: filepath.Join(baseDir, "resources", "images", "chr"),
	}
}

// DefaultBaseDir returns "."
func DefaultBaseDir() string {
	return "."
}
