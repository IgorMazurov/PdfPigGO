package systemfonts

// IOSSystemFontLister is a placeholder font lister for iOS platforms.
// Very early version, intended to help developing support for iOS.
type IOSSystemFontLister struct{}

// NewIOSSystemFontLister creates a new IOSSystemFontLister.
func NewIOSSystemFontLister() IOSSystemFontLister {
	return IOSSystemFontLister{}
}

// GetAllFonts returns an empty slice as iOS support is not yet implemented.
func (l IOSSystemFontLister) GetAllFonts() []SystemFontRecord {
	return nil
}
