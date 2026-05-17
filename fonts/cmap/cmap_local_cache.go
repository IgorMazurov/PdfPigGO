package cmap

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/util"
)

var cmapNameTag = []byte("/CMapName ")

// CMapLocalCache provides a local (per document) cache for CMap objects, allowing
// efficient retrieval and storage of CMap instances based on their names and unique
// identifiers. This class is designed to cache CMap objects to improve performance
// by avoiding redundant parsing of CMap data. It uses a combination of CMap names
// and GUIDs derived from the CMap data to uniquely identify and store CMap instances.
type CMapLocalCache struct {
	cache          map[string]map[string]*CMap
	mu             sync.Mutex
	filterProvider filters.LookupFilterProvider
	scanner        tokenization.PdfTokenScanner
}

// NewCMapLocalCache creates a new CMapLocalCache with the given filter provider and scanner.
func NewCMapLocalCache(
	filterProvider filters.LookupFilterProvider,
	scanner tokenization.PdfTokenScanner,
) *CMapLocalCache {
	return &CMapLocalCache{
		cache:          make(map[string]map[string]*CMap),
		filterProvider: filterProvider,
		scanner:        scanner,
	}
}

// TryGetByName attempts to retrieve a CMap by name from the global cache.
// Returns true and the CMap if found, false otherwise.
func (c *CMapLocalCache) TryGetByName(name string) (*CMap, bool) {
	return TryGet(name)
}

// TryGetByStream attempts to retrieve or parse a CMap from the given stream token.
// If the stream data is empty, returns false and nil. Otherwise it decodes the stream,
// checks the local cache using a hash-based identifier, and falls back to parsing if
// not found. Returns true and the CMap (either cached or newly parsed).
func (c *CMapLocalCache) TryGetByStream(token *tokens.StreamToken) (*CMap, bool) {
	if token == nil || len(token.Data()) == 0 {
		return nil, false
	}

	decoded, err := decodeStreamWithScanner(token, c.filterProvider, c.scanner)
	if err != nil {
		return nil, false
	}

	cmapName, found := tryGetNameFast(decoded)
	if !found {
		result, _ := Parse(core.NewMemoryInputBytes(decoded))
		return result, true
	}

	guid := getGuid(decoded)

	c.mu.Lock()
	defer c.mu.Unlock()

	cMaps, ok := c.cache[cmapName]
	if !ok {
		cMaps = make(map[string]*CMap)
		c.cache[cmapName] = cMaps
	}

	result, ok := cMaps[guid]
	if ok {
		return result, true
	}

	result, _ = Parse(core.NewMemoryInputBytes(decoded))
	cMaps[guid] = result

	return result, true
}

// decodeStreamWithScanner decodes the stream data using a scanner-aware filter provider.
func decodeStreamWithScanner(
	stream *tokens.StreamToken,
	provider filters.LookupFilterProvider,
	scanner tokenization.PdfTokenScanner,
) ([]byte, error) {
	fl, err := provider.GetFiltersWithScanner(stream.StreamDictionary, scanner)
	if err != nil {
		return nil, err
	}

	totalMaxEstSize := float64(len(stream.Data())) * 100
	transform := stream.Data()

	for i, filter := range fl {
		totalMaxEstSize *= getEstimatedSizeMultiplier(filter)

		transform, err = filter.Decode(transform, stream.StreamDictionary, provider, i)
		if err != nil {
			return nil, err
		}

		if i < len(fl)-1 && float64(len(transform)) > totalMaxEstSize {
			return nil, core.NewPdfDocumentFormatException(
				fmt.Sprintf("decoded stream size exceeds the estimated maximum size. Current decoded stream length: %d, %d filters applied out of %d",
					len(transform), i+1, len(fl)))
		}
	}

	return transform, nil
}

func getEstimatedSizeMultiplier(filter filters.Filter) float64 {
	switch filter.(type) {
	case *filters.Ascii85Filter, *filters.AsciiHexDecodeFilter:
		return 1.3
	case *filters.CcittFaxDecodeFilter, *filters.DctDecodeFilter,
		*filters.Jbig2DecodeFilter, *filters.JpxDecodeFilter:
		return 1
	default:
		return 10
	}
}

func getGuid(b []byte) string {
	return string(util.ComputeX64_128(b))
}

func tryGetNameFast(data []byte) (string, bool) {
	nameIndex := bytes.Index(data, cmapNameTag)
	if nameIndex < 0 {
		return "", false
	}

	nameIndex += len(cmapNameTag)

	defIndex := bytes.Index(data[nameIndex:], []byte("def"))
	if defIndex < 0 {
		return "", false
	}

	name := string(data[nameIndex : nameIndex+defIndex-1])
	return name, true
}


