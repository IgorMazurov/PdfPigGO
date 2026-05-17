package fonts

import (
	"github.com/uglytoad/pdfpig/go/tokens"
)

// DescriptorFontFile holds the location and type of the stream containing
// the corresponding font program. This can be a Type 1 font program
// (FontFile), a TrueType font program (FontFile2), or a font program whose
// format is given by the Subtype of the stream dictionary (FontFile3).
// At most only one of these entries is present.
type DescriptorFontFile struct {
	// ObjectKey is the indirect reference to the stream object containing the font program.
	ObjectKey *tokens.IndirectReferenceToken

	// FileType indicates the type of the font program represented by this descriptor.
	FileType FontFileType
}

// NewDescriptorFontFile creates a new DescriptorFontFile with the given
// indirect reference and file type.
func NewDescriptorFontFile(objectKey *tokens.IndirectReferenceToken, fileType FontFileType) *DescriptorFontFile {
	return &DescriptorFontFile{
		ObjectKey: objectKey,
		FileType:  fileType,
	}
}

// FontFileType represents the type of font program in a font descriptor stream.
type FontFileType int

const (
	// Type1 indicates a Type 1 font program (FontFile entry).
	Type1 FontFileType = iota

	// TrueType indicates a TrueType font program (FontFile2 entry).
	TrueType

	// FromSubtype indicates a font program whose format is defined by the
	// Subtype entry of the stream dictionary (FontFile3 entry).
	FromSubtype
)
