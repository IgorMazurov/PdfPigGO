package cmap

import "fmt"

// CidRange associates the beginning and end of a range of character codes
// with the starting CID for the range.
type CidRange struct {
	firstCharacterCode int
	lastCharacterCode  int
	cid                int
}

// NewCidRange creates a new CidRange to associate a range of character codes
// to a range of CIDs. Returns an error if lastCharacterCode is less than
// firstCharacterCode.
func NewCidRange(firstCharacterCode, lastCharacterCode, cid int) (*CidRange, error) {
	if lastCharacterCode < firstCharacterCode {
		return nil, fmt.Errorf("the last character code cannot be lower than the first character code: first: %d, last: %d, cid: %d", firstCharacterCode, lastCharacterCode, cid)
	}

	return &CidRange{
		firstCharacterCode: firstCharacterCode,
		lastCharacterCode:  lastCharacterCode,
		cid:                cid,
	}, nil
}

// Contains determines if this CidRange contains a mapping for the character code.
func (r CidRange) Contains(characterCode int) bool {
	return r.firstCharacterCode <= characterCode && characterCode <= r.lastCharacterCode
}

// TryMap attempts to map the given character code to the corresponding CID in
// this range. Returns true if the character code maps to a CID, false if out
// of range.
func (r CidRange) TryMap(characterCode int) (int, bool) {
	if r.Contains(characterCode) {
		return r.cid + (characterCode - r.firstCharacterCode), true
	}
	return 0, false
}

// String returns a string representation of the range.
func (r CidRange) String() string {
	return fmt.Sprintf("CID %d: Code %d -> %d", r.cid, r.firstCharacterCode, r.lastCharacterCode)
}
