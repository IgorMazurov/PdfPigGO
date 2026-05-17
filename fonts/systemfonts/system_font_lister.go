package systemfonts

// SystemFontLister lists fonts available on the current system.
type SystemFontLister interface {
	// GetAllFonts returns all discovered system font records.
	GetAllFonts() []SystemFontRecord
}
