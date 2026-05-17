// Package graphics provides types for processing PDF graphics content streams.
package graphics

import (
	"errors"
	"fmt"

	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/functions"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/util"
)

// SoftMaskType defines the subtype of a soft mask.
type SoftMaskType byte

const (
	// Alpha indicates the alpha channel of the mask image is used as the mask.
	Alpha SoftMaskType = 0
	// Luminosity indicates the luminosity of the mask image pixels is used to compute the mask.
	Luminosity SoftMaskType = 1
)

// SoftMask represents a soft-mask dictionary used for transparency operations in PDF.
type SoftMask struct {
	subtype            SoftMaskType
	transparencyGroup  *tokens.StreamToken
	bc                 []float64
	transferFunction   functions.PdfFunction
}

// Subtype returns the soft mask subtype (Alpha or Luminosity).
func (s *SoftMask) Subtype() SoftMaskType {
	if s == nil {
		return Alpha
	}
	return s.subtype
}

// TransparencyGroup returns the stream token representing the transparency group.
func (s *SoftMask) TransparencyGroup() *tokens.StreamToken {
	if s == nil {
		return nil
	}
	return s.transparencyGroup
}

// BC returns the backdrop color values, or nil if not specified.
func (s *SoftMask) BC() []float64 {
	if s == nil {
		return nil
	}
	return s.bc
}

// TransferFunction returns the transfer function, or nil if not set.
func (s *SoftMask) TransferFunction() functions.PdfFunction {
	if s == nil {
		return nil
	}
	return s.transferFunction
}

// ParseSoftMask parses a soft-mask dictionary token into a SoftMask struct.
// The dictionary must contain /S (subtype) and /G (group stream) entries.
// Corresponds to C# SoftMask.Parse(DictionaryToken, IPdfTokenScanner, ILookupFilterProvider).
func ParseSoftMask(
	dictionaryToken *tokens.DictionaryToken,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.LookupFilterProvider,
) (*SoftMask, error) {
	if dictionaryToken == nil {
		return nil, errors.New("dictionaryToken cannot be nil")
	}

	softMask := &SoftMask{}

	sToken, found := dictionaryToken.TryGet(tokens.S)
	if !found {
		return nil, fmt.Errorf("missing soft-mask dictionary '%s' entry", tokens.S)
	}

	sName, ok := sToken.(*tokens.NameToken)
	if !ok {
		return nil, fmt.Errorf("soft-mask dictionary '%s' entry is not a name token", tokens.S)
	}

	if sName.Equals(tokens.Luminosity) {
		softMask.subtype = Luminosity
	} else if sName.Equals(tokens.Alpha) {
		softMask.subtype = Alpha
	} else {
		return nil, fmt.Errorf("invalid soft-mask Subtype '%s' entry", sName)
	}

	gToken, found := dictionaryToken.TryGet(tokens.G)
	if !found {
		return nil, fmt.Errorf("missing soft-mask dictionary '%s' entry", tokens.G)
	}

	gStream, ok := gToken.(*tokens.StreamToken)
	if !ok {
		return nil, fmt.Errorf("soft-mask dictionary '%s' entry is not a stream token", tokens.G)
	}
	softMask.transparencyGroup = gStream

	bcToken, found := dictionaryToken.TryGet(tokens.Bc)
	if found {
		if bcArray, ok := bcToken.(*tokens.ArrayToken); ok {
			for _, t := range bcArray.Data() {
				if nt, ok := t.(*tokens.NumericToken); ok {
					softMask.bc = append(softMask.bc, nt.DoubleVal())
				}
			}
		}
	}

	trToken, found := dictionaryToken.TryGet(tokens.Tr)
	if found {
		switch tr := trToken.(type) {
		case *tokens.NameToken:
			if !tr.Equals(tokens.Identity) {
				return nil, fmt.Errorf("invalid transfer function name '%s' entry, should be '%s'", tr, tokens.Identity)
			}
		default:
			pdfFunction, err := util.Create(tr, scanner, filterProvider)
			if err != nil {
				return nil, fmt.Errorf("failed to parse transfer function: %w", err)
			}
			softMask.transferFunction = pdfFunction
		}
	}

	return softMask, nil
}
