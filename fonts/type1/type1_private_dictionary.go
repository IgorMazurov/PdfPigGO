package type1

import (
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/fonts/type1/charstrings"
)

// Type1PrivateDictionary holds the Private dictionary for a Type 1 font. It contains hints that apply across all characters
// in the font to help preserve properties of character outline shapes when rendered at smaller sizes and lower resolutions.
type Type1PrivateDictionary struct {
	fonts.AdobeStylePrivateDictionary

	// UniqueId optionally uniquely identifies this font.
	UniqueId *int

	// LenIv indicates the number of random bytes used for charstring encryption/decryption. Default: 4.
	LenIv int

	// RoundStemUp is preserved for backwards compatibility. Must be set if LanguageGroup is 1.
	RoundStemUp *bool

	// Password is required for backwards compatibility. Default: 5839.
	Password int

	// MinFeature is required for backwards compatibility. Default: {16, 16}.
	MinFeature MinFeature
}

// Type1PrivateDictionaryBuilder is a mutable builder for constructing a Type1PrivateDictionary. It performs no validation.
type Type1PrivateDictionaryBuilder struct {
	fonts.AdobeStylePrivateDictionaryBuilder

	// Rd holds the temporary storage for the Rd procedure tokens.
	Rd []TokenLike

	// NoAccessPut holds the temporary storage for the No Access Put procedure tokens.
	NoAccessPut []TokenLike

	// NoAccessDef holds the temporary storage for the No Access Def procedure tokens.
	NoAccessDef []TokenLike

	// Subroutines holds the decrypted but raw bytes of the subroutines in this private dictionary.
	Subroutines []*charstrings.Type1CharstringDecryptedBytes

	// UniqueId maps to Type1PrivateDictionary.UniqueId.
	UniqueId *int

	// Password maps to Type1PrivateDictionary.Password.
	Password *int

	// LenIv maps to Type1PrivateDictionary.LenIv.
	LenIv int

	// MinFeature maps to Type1PrivateDictionary.MinFeature.
	MinFeature MinFeature

	// RoundStemUp maps to Type1PrivateDictionary.RoundStemUp.
	RoundStemUp *bool
}

// Build generates a Type1PrivateDictionary from the values in this builder.
func (b *Type1PrivateDictionaryBuilder) Build() *Type1PrivateDictionary {
	base := fonts.NewAdobeStylePrivateDictionary(&b.AdobeStylePrivateDictionaryBuilder)

	pw := 5839
	if b.Password != nil {
		pw = *b.Password
	}

	mf := DefaultMinFeature
	if b.MinFeature.First != 0 || b.MinFeature.Second != 0 {
		mf = b.MinFeature
	}

	return &Type1PrivateDictionary{
		AdobeStylePrivateDictionary: base,
		UniqueId:                    b.UniqueId,
		LenIv:                       b.LenIv,
		RoundStemUp:                 b.RoundStemUp,
		Password:                    pw,
		MinFeature:                  mf,
	}
}
