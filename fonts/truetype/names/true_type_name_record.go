package names

import (
	"fmt"
)

// TrueTypeNameRecord holds a human-readable name for a feature, setting,
// copyright notice, font name or other font-related information in a TrueType font.
type TrueTypeNameRecord struct {
	// PlatformId is the supported platform identifier.
	PlatformId TrueTypePlatformIdentifier

	// PlatformEncodingId is the platform-specific encoding id.
	// Interpretation depends on the value of PlatformId.
	PlatformEncodingId uint16

	// LanguageId uniquely defines the language in which the string is written for this record.
	LanguageId uint16

	// NameId is used to reference this record by other tables in the font.
	NameId uint16

	// Value is the value of this record.
	Value string
}

// NewTrueTypeNameRecord creates a new TrueTypeNameRecord.
func NewTrueTypeNameRecord(platformId TrueTypePlatformIdentifier, platformEncodingId, languageId, nameId uint16, value string) *TrueTypeNameRecord {
	return &TrueTypeNameRecord{
		PlatformId:       platformId,
		PlatformEncodingId: platformEncodingId,
		LanguageId:     languageId,
		NameId:         nameId,
		Value:          value,
	}
}

// String returns a string representation of the TrueTypeNameRecord.
func (r *TrueTypeNameRecord) String() string {
	return fmt.Sprintf("(Platform: %v, Id: %d) - %s", r.PlatformId, r.NameId, r.Value)
}
