package backend

import (
	"context"
	"fmt"
	"os"
	"terraria-json-compressor-gui/backend/handlers"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx         context.Context
	fileHandler *handlers.FileHandler
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		fileHandler: handlers.NewFileHandler(),
	}
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
}

// Get first file in a folder
func (a *App) GetFirstFile(folderPath string) (string, error) {
	return a.fileHandler.GetFirstFile(a.ctx, folderPath)
}

// SelectFolder opens a folder picker dialog
func (a *App) SelectFile() (string, error) {
	return handlers.OpenFileDialog(a.ctx)
}

func (a *App) PreviewCSVFile(filePath string) ([]handlers.LogEntry, error) {
	return handlers.PreviewCSVFile(a.ctx, filePath)
}

// CompressAndSave compresses the file and prompts user to save
func (a *App) CompressAndSave(inputPath string) []handlers.LogEntry {
	var logs []handlers.LogEntry

	logs = append(logs, handlers.LogEntry{Level: "info", Message: "Starting compression process..."})

	// Compress CSV to JSON (in memory first)
	tempOutputPath := "" // empty = auto-generate temp path
	result, err := handlers.CompressCSVToJSON(
		a.ctx,
		inputPath,
		tempOutputPath,
		false, // includeEmpty - skip empty translations
		true,  // pretty - format JSON with indentation
		true,  // perKeyLog - log each key processed
	)

	// Add logs from compression
	logs = append(logs, result.Logs...)

	if err != nil {
		logs = append(logs, handlers.LogEntry{Level: "error", Message: fmt.Sprintf("Compression failed: %v", err)})
		return logs
	}

	// Check if there were critical errors
	hasError := false
	for _, log := range result.Logs {
		if log.Level == "error" {
			hasError = true
			break
		}
	}

	if hasError {
		logs = append(logs, handlers.LogEntry{Level: "error", Message: "Compression completed with errors. Please fix errors and try again."})
		// Clean up temp file
		if result.OutputPath != "" {
			os.Remove(result.OutputPath)
		}
		return logs
	}

	// Show save dialog
	savePath, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save Compressed Translation File",
		DefaultFilename: "en-US.json",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "JSON Files (*.json)",
				Pattern:     "*.json",
			},
		},
	})

	if err != nil {
		logs = append(logs, handlers.LogEntry{Level: "error", Message: fmt.Sprintf("Save dialog error: %v", err)})
		os.Remove(result.OutputPath)
		return logs
	}

	if savePath == "" {
		logs = append(logs, handlers.LogEntry{Level: "warn", Message: "Save cancelled by user"})
		os.Remove(result.OutputPath)
		return logs
	}

	// Read the temp file and write to user's selected path
	data, err := os.ReadFile(result.OutputPath)
	if err != nil {
		logs = append(logs, handlers.LogEntry{Level: "error", Message: fmt.Sprintf("Failed to read compressed data: %v", err)})
		os.Remove(result.OutputPath)
		return logs
	}

	err = os.WriteFile(savePath, data, 0644)
	if err != nil {
		logs = append(logs, handlers.LogEntry{Level: "error", Message: fmt.Sprintf("Failed to save file: %v", err)})
		os.Remove(result.OutputPath)
		return logs
	}

	// Clean up temp file
	os.Remove(result.OutputPath)

	logs = append(logs, handlers.LogEntry{Level: "ok", Message: fmt.Sprintf("File saved successfully to: %s", savePath)})
	logs = append(logs, handlers.LogEntry{Level: "summary", Message: fmt.Sprintf("Compression complete! %d keys exported", result.KeyCount)})

	return logs
}

// SaveFile opens a save file dialog and returns the selected path
func (a *App) SaveFile() (string, error) {
	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save Compressed File",
		DefaultFilename: "compressed_output.json",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "JSON Files (*.json)",
				Pattern:     "*.json",
			},
			{
				DisplayName: "All Files (*.*)",
				Pattern:     "*.*",
			},
		},
	})

	if err != nil {
		return "", err
	}

	return selection, nil
}
