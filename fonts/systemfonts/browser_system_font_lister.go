package systemfonts

// BrowserSystemFontLister is a placeholder font lister for browser environments.
type BrowserSystemFontLister struct{}

// NewBrowserSystemFontLister creates a new BrowserSystemFontLister.
func NewBrowserSystemFontLister() BrowserSystemFontLister {
	return BrowserSystemFontLister{}
}

// GetAllFonts returns an empty slice as no system fonts are available in browser context.
func (l BrowserSystemFontLister) GetAllFonts() []SystemFontRecord {
	return nil
}
