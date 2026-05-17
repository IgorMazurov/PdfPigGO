package systemfonts

import (
	"os"
	"path/filepath"
)

// MacSystemFontLister lists fonts on macOS by scanning standard font directories.
type MacSystemFontLister struct{}

// NewMacSystemFontLister creates a new MacSystemFontLister.
func NewMacSystemFontLister() MacSystemFontLister {
	return MacSystemFontLister{}
}

// GetAllFonts returns all system font records found under standard macOS font directories.
func (l MacSystemFontLister) GetAllFonts() []SystemFontRecord {
	directories := []string{
		"/Library/Fonts/",
		"/System/Library/Fonts/",
		"/Network/Library/Fonts/",
	}

	if home := os.Getenv("HOME"); home != "" {
		directories = append(directories, filepath.Join(home, "Library", "Fonts"))
	}

	var records []SystemFontRecord

	for _, directory := range directories {
		info, err := os.Stat(directory)
		if err != nil || !info.IsDir() {
			continue
		}

		err = filepath.WalkDir(directory, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				return nil
			}
			if record, ok := TryCreate(path); ok {
				records = append(records, record)
			}
			return nil
		})
		if err != nil {
			continue
		}
	}

	return records
}
