package fields

import (
	"github.com/uglytoad/pdfpig/go/core"
	"strings"
)

// AcroFieldCommonInformation holds information from the field dictionary which is common across all field types.
// All of this information is optional.
type AcroFieldCommonInformation struct {
	// Parent is the reference to the field which is the parent of this one, if applicable.
	Parent *core.IndirectReference

	// PartialName is the partial field name for this field. The fully qualified field name is the
	// period '.' joined name of all parents' partial names and this field's partial name.
	PartialName string

	// AlternateName is the alternate field name to be used instead of the fully qualified field name where
	// the field is being identified on the user interface or by screen readers.
	AlternateName string

	// MappingName is the mapping name used when exporting form field data from the document.
	MappingName string
}

// NewAcroFieldCommonInformation creates a new AcroFieldCommonInformation.
func NewAcroFieldCommonInformation(parent *core.IndirectReference, partialName, alternateName, mappingName string) *AcroFieldCommonInformation {
	return &AcroFieldCommonInformation{
		Parent:        parent,
		PartialName:   partialName,
		AlternateName: alternateName,
		MappingName:   mappingName,
	}
}

// String returns a human-readable representation of the field common information.
func (info *AcroFieldCommonInformation) String() string {
	var b strings.Builder

	if info.Parent != nil {
		b.WriteString("Parent: ")
		b.WriteString(info.Parent.String())
		b.WriteByte('.')
	}

	appendIfNotEmpty := func(val, label string) {
		if val == "" {
			return
		}
		if b.Len() > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(label)
		b.WriteString(": ")
		b.WriteString(val)
		b.WriteByte('.')
	}

	appendIfNotEmpty(info.PartialName, "Partial Name")
	appendIfNotEmpty(info.AlternateName, "Alternate Name")
	appendIfNotEmpty(info.MappingName, "Mapping Name")

	return b.String()
}
