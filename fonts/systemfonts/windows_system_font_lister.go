package systemfonts

import (
	"os"
	"path/filepath"
)

// WindowsSystemFontLister lists fonts on Windows by scanning the system Fonts and PSFonts directories.
type WindowsSystemFontLister struct{}

// NewWindowsSystemFontLister creates a new WindowsSystemFontLister.
func NewWindowsSystemFontLister() WindowsSystemFontLister {
	return WindowsSystemFontLister{}
}

// GetAllFonts returns all system font records found under the Windows Fonts and PSFonts directories.
func (l WindowsSystemFontLister) GetAllFonts() []SystemFontRecord {
	winDir := os.Getenv("WINDIR")
	if winDir == "" {
		winDir = "C:\\Windows"
	}

	directories := []string{
		filepath.Join(winDir, "Fonts"),
		filepath.Join(winDir, "PSFonts"),
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
