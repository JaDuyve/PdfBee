package main

import (
	"bufio"
	"changeme/backend/pdfutil"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
)

var (
	ErrNoFileSelected = errors.New("no file selected")
)

// App struct
type App struct {
	AppName              string
	ctx                  context.Context
	CurrentSelectedFile  string
	CurrentGeneratedFile string
}

// NewApp creates a new App application struct
func NewApp(appName string) *App {
	return &App{
		AppName: appName,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetPreviewFilePath() (string, error) {
	if a.CurrentGeneratedFile != "" {
		return a.CurrentGeneratedFile, nil
	}

	if a.CurrentSelectedFile != "" {
		return a.CurrentSelectedFile, nil
	}

	return "", ErrNoFileSelected
}

func (a *App) OpenFileDialog() {
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

	a.CurrentSelectedFile = filePath
	a.CurrentGeneratedFile = ""
	runtime.EventsEmit(a.ctx, "preview_file_content_updated")
	runtime.EventsEmit(a.ctx, "form_content_updated")
}

func (a *App) GetPreviewContent() string {
	previewFilePath, err := a.GetPreviewFilePath()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return ""
	}

	runtime.LogInfo(a.ctx, previewFilePath)

	file, err := os.Open(previewFilePath)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	defer file.Close()

	stats, statsErr := file.Stat()
	if statsErr != nil {
		runtime.LogError(a.ctx, err.Error())
		return ""
	}

	var size = stats.Size()
	bytes := make([]byte, size)

	buff := bufio.NewReader(file)
	_, err = buff.Read(bytes)

	return base64.StdEncoding.EncodeToString(bytes)
}

func (a *App) GetPdfForm() pdfutil.Form {
	previewFilePath, err := a.GetPreviewFilePath()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return pdfutil.Form{}
	}

	form, err := pdfutil.GetFormFields(previewFilePath)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}

	return *form
}

func (a *App) UpdatePdfForm(form pdfutil.Form) {
	err := a.updatePdfForm(form)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
}

func (a *App) updatePdfForm(form pdfutil.Form) error {
	runtime.LogInfo(a.ctx, form.String())
	runtime.LogInfo(a.ctx, a.CurrentGeneratedFile)

	prevGeneratedFile := a.CurrentGeneratedFile

	pattern := fmt.Sprintf("%s-", a.AppName)
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return err
	}

	defer f.Close()

	err = pdfutil.FillForm(form, a.CurrentSelectedFile, f)
	if err != nil {
		return err
	}

	if prevGeneratedFile != "" {
		err = os.Remove(prevGeneratedFile)
		if err != nil {
			return err
		}
	}

	a.CurrentGeneratedFile = f.Name()

	f.Close()
	runtime.EventsEmit(a.ctx, "preview_file_content_updated")

	return nil
}

func (a *App) UpdatePdfFormWithFieldNames() {
	err := a.updatePdfFormWithFieldNames()
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
	}
}

func (a *App) updatePdfFormWithFieldNames() error {
	runtime.LogInfo(a.ctx, a.CurrentGeneratedFile)

	prevGeneratedFile := a.CurrentGeneratedFile

	pattern := fmt.Sprintf("%s-", a.AppName)
	f, err := os.CreateTemp("", pattern)
	if err != nil {
		return err
	}

	defer f.Close()

	err = pdfutil.FillFormWithFieldName(a.CurrentSelectedFile, f)
	if err != nil {
		return err
	}

	if prevGeneratedFile != "" {
		err = os.Remove(prevGeneratedFile)
		if err != nil {
			return err
		}
	}

	a.CurrentGeneratedFile = f.Name()

	f.Close()
	runtime.EventsEmit(a.ctx, "preview_file_content_updated")
	runtime.EventsEmit(a.ctx, "form_content_updated")

	return nil
}
