package names

// TrueTypeUnicodeEncodingIdentifier specifies the meaning of the platform-specific
// encoding when the platform identifier is TrueTypePlatformIdentifier.Unicode.
type TrueTypeUnicodeEncodingIdentifier uint16

const (
	// Default indicates default semantics.
	Default TrueTypeUnicodeEncodingIdentifier = 0

	// Version1Point1 indicates Unicode version 1.1 semantics.
	Version1Point1 TrueTypeUnicodeEncodingIdentifier = 1

	// Iso10646 indicates ISO 10646 1993 semantics (deprecated).
	Iso10646 TrueTypeUnicodeEncodingIdentifier = 2

	// Unicode2BmpOnly indicates Unicode 2.0+ semantics for BMP characters only.
	Unicode2BmpOnly TrueTypeUnicodeEncodingIdentifier = 3

	// Unicode2NonBmpAllowed indicates Unicode 2.0+ semantics including non-BMP characters.
	Unicode2NonBmpAllowed TrueTypeUnicodeEncodingIdentifier = 4

	// UnicodeVariationSequences indicates Unicode Variation Sequences.
	UnicodeVariationSequences TrueTypeUnicodeEncodingIdentifier = 5

	// FullUnicode indicates full Unicode coverage.
	FullUnicode TrueTypeUnicodeEncodingIdentifier = 6
)
