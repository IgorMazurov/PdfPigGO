package systemfonts

import (
	"os"
	"path/filepath"
)

// LinuxSystemFontLister lists fonts on Linux systems by scanning standard font directories.
type LinuxSystemFontLister struct{}

// NewLinuxSystemFontLister creates a new LinuxSystemFontLister.
func NewLinuxSystemFontLister() LinuxSystemFontLister {
	return LinuxSystemFontLister{}
}

// GetAllFonts returns all system font records found under standard Linux font directories.
func (l LinuxSystemFontLister) GetAllFonts() []SystemFontRecord {
	directories := []string{
		"/usr/local/fonts",
		"/usr/local/share/fonts",
		"/usr/share/fonts",
		"/usr/X11R6/lib/X11/fonts",
	}

	if home := os.Getenv("HOME"); home != "" {
		directories = append(directories, filepath.Join(home, ".fonts"))
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
