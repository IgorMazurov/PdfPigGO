package systemfonts

import (
	"os"
	"path/filepath"
)

// AndroidSystemFontLister lists fonts on Android systems by scanning /system/fonts.
type AndroidSystemFontLister struct{}

// NewAndroidSystemFontLister creates a new AndroidSystemFontLister.
func NewAndroidSystemFontLister() AndroidSystemFontLister {
	return AndroidSystemFontLister{}
}

// GetAllFonts returns all system font records found under /system/fonts.
func (l AndroidSystemFontLister) GetAllFonts() []SystemFontRecord {
	directories := []string{
		"/system/fonts",
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
