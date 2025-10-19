package handlers

import (
	"context"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"
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

func OpenFileDialog(ctx context.Context) (string, error) {
	file, err := runtime.OpenFileDialog(ctx, runtime.OpenDialogOptions{
		Title: "Select File - Should be a .csv File",
		Filters: []runtime.FileFilter{
			{DisplayName: "csv", Pattern: "*.csv"},
			{DisplayName: "All Files", Pattern: "*"},
		},
	})

	return file, err
}
