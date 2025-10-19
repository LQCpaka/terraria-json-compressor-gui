package handlers

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

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
