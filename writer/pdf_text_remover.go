// Package writer provides types for writing PDF documents.

package writer

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/parser"
)

// RemoveTextBytes reads a PDF from the given file path, removes all text operations
// from page content streams, and returns the result as a byte slice.
// If pages is nil, all pages are processed; otherwise only the listed 1-based page numbers are emitted.
func RemoveTextBytes(filePath string, pages []int) ([]byte, error) {
	out := &bytes.Buffer{}
	if err := RemoveTextToFile(out, filePath, pages); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// RemoveTextToFile reads a PDF from the given file path, removes all text operations
// from page content streams, and writes the result to the output writer.
// The caller is responsible for managing the output writer's lifecycle.
// If pages is nil, all pages are processed; otherwise only the listed 1-based page numbers are emitted.
func RemoveTextToFile(out io.Writer, filePath string, pages []int) error {
	if out == nil {
		return fmt.Errorf("output writer must not be nil")
	}
	if filePath == "" {
		return fmt.Errorf("file path must not be empty")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer f.Close()

	return RemoveTextFromStream(f, out, pages)
}

// RemoveTextMemory removes text operations from a PDF supplied as a byte slice
// and returns the result as a new byte slice.
// If pages is nil, all pages are processed; otherwise only the listed 1-based page numbers are emitted.
func RemoveTextMemory(file []byte, pages []int) ([]byte, error) {
	if file == nil {
		return nil, fmt.Errorf("file bytes must not be nil")
	}

	doc, err := parser.OpenMemory(file, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open pdf from bytes: %w", err)
	}

	out := &seekableBuffer{}
	if err := RemoveTextFromDocument(doc, out, pages); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

// RemoveTextFromStream reads a PDF from the input stream, removes all text operations
// from page content streams, and writes the result to the output writer.
// The caller is responsible for managing both streams' lifecycle.
// If pages is nil, all pages are processed; otherwise only the listed 1-based page numbers are emitted.
func RemoveTextFromStream(stream io.Reader, out io.Writer, pages []int) error {
	if stream == nil {
		return fmt.Errorf("input stream must not be nil")
	}
	if out == nil {
		return fmt.Errorf("output writer must not be nil")
	}

	data, err := io.ReadAll(stream)
	if err != nil {
		return fmt.Errorf("failed to read stream: %w", err)
	}

	doc, err := parser.OpenMemory(data, nil)
	if err != nil {
		return fmt.Errorf("failed to open pdf from stream: %w", err)
	}

	buf := &seekableBuffer{}
	if err := RemoveTextFromDocument(doc, buf, pages); err != nil {
		return err
	}

	_, err = io.Copy(out, buf)
	return err
}

// RemoveTextFromDocument removes text operations from the given PdfDocument
// and writes the result to the output writer.
// The caller is responsible for managing the output writer's lifecycle.
// If pages is nil, all pages are processed; otherwise only the listed 1-based page numbers are emitted.
func RemoveTextFromDocument(doc *content.PdfDocument, out io.WriteSeeker, pages []int) error {
	if doc == nil {
		return fmt.Errorf("pdf document must not be nil")
	}
	if out == nil {
		return fmt.Errorf("output writer must not be nil")
	}

	tokenWriter := NewNoTextTokenWriter()

	builder := NewPdfDocumentBuilderWithStream(out, false, PdfWriterDefault, doc.Version(), tokenWriter)
	if builder == nil {
		return fmt.Errorf("failed to create pdf document builder")
	}
	defer builder.Close()

	pageList := determinePages(doc.NumberOfPages(), pages)

	for _, pageNum := range pageList {
		tokenWriter.SetPage(pageNum)
		if _, err := builder.AddPageWithOptions(doc, pageNum, NewAddPageOptions()); err != nil {
			return fmt.Errorf("failed to add page %d: %w", pageNum, err)
		}
	}

	return nil
}

// determinePages returns the list of 1-based page numbers to process.
// If pages is nil or empty, all pages from 1 to total are returned.
func determinePages(total int, pages []int) []int {
	if len(pages) == 0 {
		result := make([]int, total)
		for i := 0; i < total; i++ {
			result[i] = i + 1
		}
		return result
	}
	copied := make([]int, len(pages))
	copy(copied, pages)
	return copied
}
