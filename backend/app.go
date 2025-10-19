package backend

import (
	"context"
	"terraria-json-compressor-gui/backend/handlers"
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
