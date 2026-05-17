// Package encodings provides character encoding types for PDF fonts.
package encodings

import "errors"

// Difference maps a character code to a glyph name override.
type Difference struct {
	Code int
	Name string
}

// DifferenceBasedEncoding represents an encoding created by combining a base
// encoding with a set of differences that override or add entries.
type DifferenceBasedEncoding struct {
	*Encoding
	name string
}

// NewDifferenceBasedEncoding creates a new DifferenceBasedEncoding from a base
// encoding and a list of (code, name) difference pairs. Differences are applied
// first, then any base entries not overridden by a difference are copied in.
func NewDifferenceBasedEncoding(base *Encoding, differences []Difference) (*DifferenceBasedEncoding, error) {
	if base == nil {
		return nil, errors.New("base encoding must not be nil")
	}

	e := &DifferenceBasedEncoding{
		Encoding: NewEncoding(),
		name:     "Difference " + base.EncodingName(),
	}

	for _, d := range differences {
		e.Add(d.Code, d.Name)
	}

	differenceCodes := make(map[int]bool, len(differences))
	for _, d := range differences {
		differenceCodes[d.Code] = true
	}

	for code, name := range base.CodeToName {
		if !differenceCodes[code] {
			e.Add(code, name)
		}
	}

	return e, nil
}

// EncodingName returns the name of this encoding.
func (e *DifferenceBasedEncoding) EncodingName() string {
	return e.name
}
