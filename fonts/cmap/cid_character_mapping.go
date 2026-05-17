package cmap

import "fmt"

// CidCharacterMapping maps from a single character code to its CID.
type CidCharacterMapping struct {
	SourceCharacterCode int
	DestinationCid      int
}

// NewCidCharacterMapping creates a new mapping from a character code to a CID.
func NewCidCharacterMapping(sourceCharacterCode, destinationCid int) CidCharacterMapping {
	return CidCharacterMapping{
		SourceCharacterCode: sourceCharacterCode,
		DestinationCid:      destinationCid,
	}
}

// String returns a string representation of the mapping.
func (m CidCharacterMapping) String() string {
	return fmt.Sprintf("Code %d -> CID %d", m.SourceCharacterCode, m.DestinationCid)
}
