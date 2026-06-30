package engine

import (
	"io/fs"
	"os"
)

// EmbeddedFS holds the combined embedded filesystem from the main package.
// When set, all file reads go through it instead of the OS filesystem.
var EmbeddedFS fs.FS

// SetEmbedFS sets the embedded filesystem for the engine.
func SetEmbedFS(f fs.FS) {
	EmbeddedFS = f
}

// ReadEmbedFile reads a file from the embedded FS if available, falls back to os.ReadFile.
func ReadEmbedFile(path string) ([]byte, error) {
	if EmbeddedFS != nil {
		return fs.ReadFile(EmbeddedFS, path)
	}
	// Fallback to OS — for CLI tools and development
	return os.ReadFile(path)
}

// ReadEmbedDir reads a directory from the embedded FS if available, falls back to os.ReadDir.
func ReadEmbedDir(path string) ([]fs.DirEntry, error) {
	if EmbeddedFS != nil {
		return fs.ReadDir(EmbeddedFS, path)
	}
	return os.ReadDir(path)
}
