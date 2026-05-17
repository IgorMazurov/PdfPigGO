package names

// TrueTypeWindowsEncodingIdentifier identifies the meaning of the platform-specific
// encoding when the platform identifier is TrueTypePlatformIdentifier.Windows.
type TrueTypeWindowsEncodingIdentifier uint16

const (
	// Symbol is the symbol encoding.
	Symbol TrueTypeWindowsEncodingIdentifier = 0

	// UnicodeBmp is the Unicode Basic Multilingual Plane encoding.
	UnicodeBmp TrueTypeWindowsEncodingIdentifier = 1

	// ShiftJis is the Shift JIS encoding for Japanese text.
	ShiftJis TrueTypeWindowsEncodingIdentifier = 2

	// Prc is the PRC (People's Republic of China) GBK encoding.
	Prc TrueTypeWindowsEncodingIdentifier = 3

	// Big5 is the Big5 encoding for Traditional Chinese text.
	Big5 TrueTypeWindowsEncodingIdentifier = 4

	// Wansung is the Wansung encoding for Korean text.
	Wansung TrueTypeWindowsEncodingIdentifier = 5

	// Johab is the JOHAB encoding for Korean text.
	Johab TrueTypeWindowsEncodingIdentifier = 6

	// Reserved7 is a reserved encoding identifier.
	Reserved7 TrueTypeWindowsEncodingIdentifier = 7

	// Reserved8 is a reserved encoding identifier.
	Reserved8 TrueTypeWindowsEncodingIdentifier = 8

	// Reserved9 is a reserved encoding identifier.
	Reserved9 TrueTypeWindowsEncodingIdentifier = 9

	// WindowsFullUnicode is the Unicode full repertoire encoding.
	WindowsFullUnicode TrueTypeWindowsEncodingIdentifier = 10
)
