package pdfutil

import (
	"errors"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"os"
)

type Form struct {
	TextFields []TextField `json:"textFields,omitempty"`
}

type TextField struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	value string `json:"value"`
}

func GetFormFields(filePath string) ([]Form, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	formGroup, err := api.ExportFormToStruct(file, filePath, nil)
	if err != nil {
		return nil, err
	}

	if len(formGroup.Forms) == 0 {
		return nil, errors.New("file doesn't contain a form")
	}

	exportForms := make([]Form, len(formGroup.Forms))
	for formIndex, form := range formGroup.Forms {
		exportTextFields := make([]TextField, len(form.TextFields))

		for textFieldIndex, textField := range form.TextFields {
			exportTextFields[textFieldIndex] = TextField{
				Id:    textField.ID,
				Name:  textField.Name,
				value: textField.Value,
			}
		}

		exportForm := Form{
			TextFields: exportTextFields,
		}

		exportForms[formIndex] = exportForm
	}

	return exportForms, nil
}
