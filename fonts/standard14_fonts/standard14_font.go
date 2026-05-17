package standard14fonts

// Standard14Font represents the 14 standard Type 1 fonts included by default in PDF readers.
type Standard14Font int

const (
	// TimesRoman is Times New Roman regular.
	TimesRoman Standard14Font = iota

	// TimesBold is Times New Roman bold.
	TimesBold

	// TimesItalic is Times New Roman italic.
	TimesItalic

	// TimesBoldItalic is Times New Roman bold and italic.
	TimesBoldItalic

	// Helvetica is Helvetica regular.
	Helvetica

	// HelveticaBold is Helvetica bold.
	HelveticaBold

	// HelveticaOblique is Helvetica oblique (italic without different font shapes).
	HelveticaOblique

	// HelveticaBoldOblique is Helvetica bold and oblique.
	HelveticaBoldOblique

	// Courier is Courier regular.
	Courier

	// CourierBold is Courier bold.
	CourierBold

	// CourierOblique is Courier oblique.
	CourierOblique

	// CourierBoldOblique is Courier bold and oblique.
	CourierBoldOblique

	// Symbol is the Symbol font (mathematical symbols, Greek letters).
	Symbol

	// ZapfDingbats is the Zapf Dingbats font (ornamental symbols).
	ZapfDingbats
)
