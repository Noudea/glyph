package core

import (
	"os"
	"path/filepath"
)

// DefaultRootPath returns the default Glyph data root.
func DefaultRootPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".glyph"), nil
}
