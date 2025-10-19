package handlers

import (
	"context"
	"os"
	"path/filepath"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (f *FileHandler) GetFirstFile(ctx context.Context, folderPath string) (string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			return filepath.Join(folderPath, entry.Name()), nil
		}
	}

	return "", nil
}
