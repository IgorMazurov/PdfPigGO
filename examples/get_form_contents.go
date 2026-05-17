// Package main demonstrates how to get form contents from a PDF document.
package main

import (
	"fmt"
	"strings"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/acroforms/fields"
)

// GetFormContents gets the form contents from the PDF at filePath and prints field details for page 1.
func GetFormContents(filePath string) error {
	doc, err := pdfpig.OpenFile(filePath, nil)
	if err != nil {
		return fmt.Errorf("opening PDF: %w", err)
	}
	defer doc.Close()

	form, ok, err := doc.TryGetForm()
	if err != nil {
		return fmt.Errorf("getting form: %w", err)
	}
	if !ok {
		fmt.Printf("No form found in file: %s.\n", filePath)
		return nil
	}

	page1Fields, err := form.GetFieldsForPage(1)
	if err != nil {
		return fmt.Errorf("getting fields for page 1: %w", err)
	}

	for _, field := range page1Fields {
		switch f := any(field).(type) {
		case *fields.AcroTextField:
			fmt.Printf("Found text field on page 1 with text: %s.\n", f.Value)
		case *fields.AcroCheckboxesField:
			fmt.Printf("Found checkboxes field on page 1 with %d checkboxes.\n", len(f.Children()))
		case *fields.AcroListBoxField:
			opts := make([]string, len(f.Options))
			for i, opt := range f.Options {
				opts[i] = opt.Name
			}
			fmt.Printf("Found listbox field on page 1 with options: %s.\n", strings.Join(opts, ", "))
		}
	}

	return nil
}


