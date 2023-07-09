package main

import (
	"bufio"
	"changeme/backend/pdfutil"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

// App struct
type App struct {
	ctx                  context.Context
	CurrentFile          string
	CurrentGeneratedFile string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetCurrentFile() string {
	file, err := os.Open(a.CurrentFile)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	var size = stats.Size()
	bytes := make([]byte, size)

	buff := bufio.NewReader(file)
	_, err = buff.Read(bytes)

	return base64.StdEncoding.EncodeToString(bytes)
}

func (a *App) OpenFileDialog() []pdfutil.Form {
	options := runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				DisplayName: "pdf",
				Pattern:     "*.pdf",
			},
		},
	}
	filePath, err := runtime.OpenFileDialog(a.ctx, options)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	runtime.LogInfo(a.ctx, filePath)

	form, err := pdfutil.GetFormFields(filePath)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	a.CurrentFile = filePath
	runtime.EventsEmit(a.ctx, "current_file_changed")

	return form
}
