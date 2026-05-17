package document_layout_analysis

import (
	"os"
	"path/filepath"
)

var dlaFolder, integrationFolder string

func init() {
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}

	// Walk up to find the module root (where testdata/ exists)
	for {
		if _, err := os.Stat(filepath.Join(dir, "testdata", "integration", "Documents")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	dlaFolder = filepath.Join(dir, "testdata", "integration", "Dla", "Documents")
	integrationFolder = filepath.Join(dir, "testdata", "integration", "Documents")
}

// GetDocumentPath returns the full path to a test document.
// If isPdf is true and name does not end with ".pdf", the extension is appended.
// First checks the Dla folder; falls back to the Integration folder if not found.
func GetDocumentPath(name string, isPdf bool) string {
	if isPdf && filepath.Ext(name) != ".pdf" {
		name += ".pdf"
	}

	doc := filepath.Join(dlaFolder, name)
	if _, err := os.Stat(doc); err == nil {
		return doc
	}

	return filepath.Join(integrationFolder, name)
}
