package testutil

import (
	"os"
	"path/filepath"
)

func moduleRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

// IntegrationDocumentsRoot is the base path for integration test documents.
var IntegrationDocumentsRoot = func() string {
	if root := moduleRoot(); root != "" {
		return filepath.Join(root, "testdata", "integration", "Documents")
	}
	return "testdata/integration/Documents"
}()

// SpecificTestDocumentsRoot is the base path for specific integration test documents.
var SpecificTestDocumentsRoot = func() string {
	if root := moduleRoot(); root != "" {
		return filepath.Join(root, "testdata", "integration", "SpecificTestDocuments")
	}
	return "testdata/integration/SpecificTestDocuments"
}()

// GetDocumentPath returns the absolute file path to an integration document.
// If isPdf is true and name does not already end with ".pdf", the extension is appended.
func GetDocumentPath(name string, isPdf bool) string {
	if isPdf && filepath.Ext(name) != ".pdf" {
		name += ".pdf"
	}

	return filepath.Join(IntegrationDocumentsRoot, name)
}

// GetSpecificTestDocumentPath returns the absolute file path to a specific integration test document.
// If isPdf is true and name does not already end with ".pdf", the extension is appended.
func GetSpecificTestDocumentPath(name string, isPdf bool) string {
	if isPdf && filepath.Ext(name) != ".pdf" {
		name += ".pdf"
	}

	return filepath.Join(SpecificTestDocumentsRoot, name)
}
