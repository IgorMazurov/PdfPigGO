package filestructure

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/util"
)

var startXRefBytes = []byte("startxref")

// StartXRefLocation holds the result of searching for the startxref token.
type StartXRefLocation struct {
	// StartXRefOperatorToken is the offset in the file where the "startxref" token starts (if found).
	StartXRefOperatorToken *int64

	// StartXRefDeclaredOffset is the offset in the file that the located "startxref" declares
	// the xref table should be at (if found).
	StartXRefDeclaredOffset *int64
}

// IsValidOffset checks whether the declared offset is within valid bounds.
func (s StartXRefLocation) IsValidOffset(bytes core.InputBytes) bool {
	if s.StartXRefDeclaredOffset == nil || *s.StartXRefDeclaredOffset < 0 || *s.StartXRefDeclaredOffset > bytes.Length() {
		return false
	}
	return true
}

// GetFirstCrossReferenceOffset searches backwards from the end of the file
// for the "startxref" token and returns its location along with the declared
// cross-reference offset.
func GetFirstCrossReferenceOffset(bytes core.InputBytes, scanner tokenization.SeekableTokenScanner, log logging.Log) StartXRefLocation {
	fileLength := bytes.Length()

	buffer := util.NewCircularByteBuffer(len(startXRefBytes))

	bytes.Seek(fileLength, 0) // io.SeekStart

	var capturedOffset *int64
	i := 0
	for {
		buffer.AddReverse(bytes.CurrentByte())
		i++

		if i >= len(startXRefBytes) {
			if buffer.IsCurrentlyEqual("startxref") {
				offset := bytes.CurrentOffset() - 1
				capturedOffset = &offset
				break
			}

			if buffer.EndsWith("startref") {
				offset := bytes.CurrentOffset()
				capturedOffset = &offset
				break
			}
		}

		bytes.Seek(bytes.CurrentOffset()-1, 0) // io.SeekStart

		if bytes.CurrentOffset() <= 0 {
			break
		}
	}

	var specifiedXrefOffset *int64
	if capturedOffset != nil {
		scanner.Seek(*capturedOffset, 0) // io.SeekStart

		if scanner.Advance() {
			if opToken, ok := scanner.Current().(*tokens.OperatorToken); ok && (opToken.Data() == "startxref" || opToken.Data() == "startref") {
				specifiedXrefOffset = getNumericTokenFollowingCurrent(scanner)
				log.Debug(fmt.Sprintf("Found startxref at %v", specifiedXrefOffset))
			}
		}
	} else {
		log.Warn("No startxref token found in the document")
	}

	return StartXRefLocation{
		StartXRefOperatorToken:  capturedOffset,
		StartXRefDeclaredOffset: specifiedXrefOffset,
	}
}

func getNumericTokenFollowingCurrent(scanner tokenization.SeekableTokenScanner) *int64 {
	for scanner.Advance() {
		if token, ok := scanner.Current().(*tokens.NumericToken); ok {
			v := token.LongVal()
			return &v
		}
		if _, ok := scanner.Current().(*tokens.CommentToken); !ok {
			break
		}
	}
	return nil
}

// FirstPassResults holds the results of the first-pass PDF parsing.
type FirstPassResults struct {
	Parts           []XrefSection
	BruteForceOffsets map[core.IndirectReference]core.XrefLocation
	XrefOffsets     map[core.IndirectReference]core.XrefLocation
	Trailer         *tokens.DictionaryToken
}

// FirstPassParse performs the first pass over the PDF file to find all xref sections
// and build a mapping of indirect references to their byte offsets.
func FirstPassParse(fileHeaderOffset FileHeaderOffset, bytes core.InputBytes, scanner tokenization.SeekableTokenScanner, log logging.Log) (FirstPassResults, error) {
	if log == nil {
		log = logging.NoopLog
	}

	var bruteForceOffsets map[core.IndirectReference]core.XrefLocation
	var bruteForceTrailer *tokens.DictionaryToken

	startXrefLocation := GetFirstCrossReferenceOffset(bytes, scanner, log)

	streamsAndTables := getXrefPartsDirectly(fileHeaderOffset, bytes, scanner, startXrefLocation, log)

	if len(streamsAndTables) == 0 {
		bruteForce := FindAllXrefsInFileOrder(bytes, scanner, log)

		streamsAndTables = bruteForce.XRefParts
		bruteForceOffsets = bruteForce.ObjectOffsets
		bruteForceTrailer = bruteForce.LastTrailer

		if len(streamsAndTables) == 0 && (bruteForceOffsets == nil || len(bruteForceOffsets) == 0) {
			return FirstPassResults{}, core.NewPdfDocumentFormatException(
				"Could not find any xref tables or streams in this document and could not resolve brute force positions.")
		}
	}

	orderedXrefs := make([]XrefSection, 0, len(streamsAndTables))
	for _, part := range streamsAndTables {
		orderedXrefs = append(orderedXrefs, part)
	}
	sortByOffset(orderedXrefs)

	var lastTrailer *tokens.DictionaryToken
	flattenedOffsets := make(map[core.IndirectReference]core.XrefLocation)

	for _, xrefPart := range orderedXrefs {
		if xrefPart.Dictionary() != nil {
			if xrefPart.Dictionary().ContainsKey(tokens.Root) ||
				lastTrailer == nil ||
				!lastTrailer.ContainsKey(tokens.Root) {
				lastTrailer = xrefPart.Dictionary()
			}
		}

		for k, v := range xrefPart.ObjectOffsets() {
			flattenedOffsets[k] = v
		}
	}

	trailer := lastTrailer
	if trailer == nil {
		trailer = bruteForceTrailer
	}

	return FirstPassResults{
		Parts:             orderedXrefs,
		BruteForceOffsets: bruteForceOffsets,
		XrefOffsets:       flattenedOffsets,
		Trailer:           trailer,
	}, nil
}

func getXrefPartsDirectly(offset FileHeaderOffset, bytes core.InputBytes, scanner tokenization.SeekableTokenScanner, startLocation StartXRefLocation, log logging.Log) []XrefSection {
	if startLocation.StartXRefDeclaredOffset == nil || !startLocation.IsValidOffset(bytes) {
		return nil
	}

	visitedLocations := make(map[int64]bool)
	var results []XrefSection
	nextLocation := *startLocation.StartXRefDeclaredOffset

	for nextLocation != 0 {
		streamOrTable := getXrefStreamOrTable(offset, bytes, scanner, nextLocation, log)

		if visitedLocations[nextLocation] {
			return nil
		}
		visitedLocations[nextLocation] = true

		if streamOrTable == nil {
			return nil
		}

		results = append(results, streamOrTable)

		switch s := streamOrTable.(type) {
		case *XrefTable:
			nextLocationPtr := s.GetPrevious()
			if nextLocationPtr != nil {
				nextLocation = *nextLocationPtr
			} else {
				nextLocation = 0
			}

			xRefStm := s.GetXRefStm()
			if xRefStm != nil {
				stream := getXrefStreamOrTable(offset, bytes, scanner, *xRefStm, log)
				if stream != nil {
					results = append(results, stream)
				}
			}
		case *XrefStream:
			nextLocationPtr := s.GetPrevious()
			if nextLocationPtr != nil {
				nextLocation = *nextLocationPtr
			} else {
				nextLocation = 0
			}
		default:
			nextLocation = 0
		}
	}

	return results
}

func getXrefStreamOrTable(fileHeaderOffset FileHeaderOffset, bytes core.InputBytes, scanner tokenization.SeekableTokenScanner, location int64, log logging.Log) XrefSection {
	table := tryReadTableAtOffset(fileHeaderOffset, location, bytes, scanner, log)
	if table != nil {
		return table
	}

	stream := tryReadStreamAtOffset(fileHeaderOffset, location, bytes, scanner, log)
	if stream == nil {
		return nil
	}
	return stream
}

func sortByOffset(xrefs []XrefSection) {
	for i := 0; i < len(xrefs)-1; i++ {
		for j := i + 1; j < len(xrefs); j++ {
			if xrefs[i].Offset() > xrefs[j].Offset() {
				xrefs[i], xrefs[j] = xrefs[j], xrefs[i]
			}
		}
	}
}

// tryReadTableAtOffset attempts to read an xref table at the given offset.
// Returns nil if no valid table is found.
func tryReadTableAtOffset(fileHeaderOffset FileHeaderOffset, location int64, bytes core.InputBytes, scanner tokenization.SeekableTokenScanner, log logging.Log) *XrefTable {
	return TryReadTableAtOffset(fileHeaderOffset, location, bytes, scanner, log)
}

// tryReadStreamAtOffset attempts to read an xref stream at the given offset.
// Returns nil if no valid stream is found.
func tryReadStreamAtOffset(fileHeaderOffset FileHeaderOffset, location int64, bytes core.InputBytes, scanner tokenization.SeekableTokenScanner, log logging.Log) *XrefStream {
	return TryReadStreamAtOffset(fileHeaderOffset, location, bytes, scanner, log)
}
