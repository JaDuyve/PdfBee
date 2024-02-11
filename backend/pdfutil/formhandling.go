package pdfutil

import (
	"errors"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/create"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/form"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

type Form struct {
	TextFields []TextField `json:"textFields,omitempty"`
}

type TextField struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
	Multiline bool   `json:"multiline"`
}

var (
	ErrFormNotFound         = errors.New("file doesn't contain a form")
	ErrNoFormFieldsAffected = errors.New("no form fields affected")
)

func (f Form) String() string {
	builder := strings.Builder{}

	builder.WriteString("TextFields filled in:\n")
	for _, field := range f.TextFields {
		if field.Value != "" {
			builder.WriteString(fmt.Sprintf("%s: %s\n", field.Name, field.Value))
		}
	}

	builder.WriteString("TextFields not filled in:\n")
	for _, field := range f.TextFields {
		if field.Value == "" {
			builder.WriteString(fmt.Sprintf("%s: %s\n", field.Name, field.Value))
		}
	}

	return builder.String()
}

func GetFormFields(filePath string) (*Form, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	formGroup, err := api.ExportForm(file, filePath, nil)
	if err != nil {
		return nil, err
	}

	if len(formGroup.Forms) == 0 {
		return nil, ErrFormNotFound
	}

	exportTextFields := make([]TextField, len(formGroup.Forms[0].TextFields))

	for textFieldIndex, textField := range formGroup.Forms[0].TextFields {
		exportTextFields[textFieldIndex] = TextField{
			Id:        textField.ID,
			Name:      textField.Name,
			Value:     textField.Value,
			Multiline: textField.Multiline,
		}
	}

	sort.Slice(exportTextFields, func(i, j int) bool {
		return strings.Compare(exportTextFields[i].Name, exportTextFields[j].Name) == -1
	})

	return &Form{
		TextFields: exportTextFields,
	}, nil
}

func FillForm(inputForm Form, filePath string, w io.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.EXPORTFORMFIELDS

	ctx, _, _, _, err := api.ReadValidateAndOptimize(file, conf, time.Now())
	if err != nil {
		return err
	}

	if err := ctx.EnsurePageCount(); err != nil {
		return err
	}

	formGroup, ok, err := form.ExportForm(ctx.XRefTable, filePath)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNoFormFieldsAffected
	}

	if len(formGroup.Forms) == 0 {
		return ErrFormNotFound
	}

	for _, field := range formGroup.Forms[0].TextFields {
		inputField, _ := inputForm.findTextFieldByName(field.Name)
		if inputField != nil {
			fmt.Printf("field: %s, locked: %v, page: %v", field.Name, field.Locked, field.Pages)
			field.Value = inputField.Value
		} else {
			log.Printf("field not found: %s", inputField.Name)
		}
	}

	conf.Cmd = model.FILLFORMFIELDS
	ctx.RemoveSignature()

	f := formGroup.Forms[0]

	ok, pp, err := form.FillForm(ctx, form.FillDetails(&f, nil), f.Pages, form.JSON)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNoFormFieldsAffected
	}

	if _, _, err := create.UpdatePageTree(ctx, pp, nil); err != nil {
		return err
	}

	if conf.ValidationMode != model.ValidationNone {
		if err = api.ValidateContext(ctx); err != nil {
			return err
		}
	}

	return api.WriteContext(ctx, w)
}

func (f Form) findTextFieldByName(name string) (*TextField, error) {
	for _, field := range f.TextFields {
		if field.Name == name {
			return &field, nil
		}
	}

	return nil, fmt.Errorf("field with name '%s' not found inside form", name)
}

func FillFormWithFieldName(filePath string, w io.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	conf := model.NewDefaultConfiguration()
	conf.Cmd = model.EXPORTFORMFIELDS

	ctx, _, _, _, err := api.ReadValidateAndOptimize(file, conf, time.Now())
	if err != nil {
		return err
	}

	if err := ctx.EnsurePageCount(); err != nil {
		return err
	}

	formGroup, ok, err := form.ExportForm(ctx.XRefTable, filePath)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNoFormFieldsAffected
	}

	if len(formGroup.Forms) == 0 {
		return ErrFormNotFound
	}

	for _, field := range formGroup.Forms[0].TextFields {
		field.Value = field.Name
		//field.Locked = false
	}

	conf.Cmd = model.FILLFORMFIELDS
	ctx.RemoveSignature()

	f := formGroup.Forms[0]

	ok, pp, err := form.FillForm(ctx, form.FillDetails(&f, nil), f.Pages, form.JSON)
	if err != nil {
		return err
	}
	if !ok {
		return ErrNoFormFieldsAffected
	}

	if _, _, err := create.UpdatePageTree(ctx, pp, nil); err != nil {
		return err
	}

	if conf.ValidationMode != model.ValidationNone {
		if err = api.ValidateContext(ctx); err != nil {
			return err
		}
	}

	return api.WriteContext(ctx, w)
}
