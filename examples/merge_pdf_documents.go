// Package main demonstrates how to merge multiple PDF documents into a single PDF,
// analogous to the C# MergePdfDocuments example from PdfPig.
package main

import (
	"fmt"
	"os"
	"path/filepath"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/writer"
)

// MergePdfDocuments merges three PDF files given by filePath1, filePath2, and filePath3 into a
// single output PDF written to the current working directory as outputOfMerge.pdf.
func MergePdfDocuments(filePath1, filePath2, filePath3 string) error {
	fileBytes := [][]byte{}
	for _, path := range []string{filePath1, filePath2, filePath3} {
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading file %s: %w", path, err)
		}
		fileBytes = append(fileBytes, data)
	}

	resultFileBytes, err := writer.MergeBytes(fileBytes, nil, writer.PdfANone, nil)
	if err != nil {
		return fmt.Errorf("merging PDFs: %w", err)
	}

	output := filepath.Join(".", "outputOfMerge.pdf")
	if err := os.WriteFile(output, resultFileBytes, 0644); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}
	fmt.Printf("File output to: %s\n", output)

	doc, err := pdfpig.Open(resultFileBytes, nil)
	if err != nil {
		return fmt.Errorf("opening merged PDF: %w", err)
	}
	defer doc.Close()

	fmt.Printf("Generated document with %d pages.\n", doc.NumberOfPages())

	return nil
}


