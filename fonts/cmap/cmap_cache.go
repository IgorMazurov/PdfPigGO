package cmap

import (
	"strings"
	"sync"

	"github.com/uglytoad/pdfpig/go/core"
)

var (
	cacheMu sync.Mutex
	cache   = make(map[string]*CMap, 0)

	cmapParser cmapParserIface
)

// cmapParserIface defines the methods needed for parsing CMaps.
// Implemented by fonts/parser.CMapParser but kept here to avoid import cycles.
type cmapParserIface interface {
	TryParseExternal(name string) (*CMap, bool)
	Parse(bytes core.InputBytes) (*CMap, error)
}

// SetCMapParser sets the global CMap parser instance.
// Call once at application startup with a fully constructed parser.
func SetCMapParser(p cmapParserIface) {
	cacheMu.Lock()
	defer cacheMu.Unlock()
	cmapParser = p
}

// TryGet attempts to retrieve a predefined CMap by name from the global cache.
// If not cached, it delegates to the CMapParser to parse an external CMap.
// Returns true and the CMap if found or successfully parsed, false otherwise.
func TryGet(name string) (*CMap, bool) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	key := strings.ToLower(name)

	if result, ok := cache[key]; ok {
		return result, true
	}

	if cmapParser == nil {
		return nil, false
	}

	result, ok := cmapParser.TryParseExternal(name)
	if ok {
		cache[key] = result
		return result, true
	}

	return nil, false
}

// Parse parses CMap data from the given input bytes using the global parser.
func Parse(bytes core.InputBytes) (*CMap, error) {
	if bytes == nil {
		return nil, core.NewPdfDocumentFormatException("bytes must not be null")
	}

	if cmapParser == nil {
		return nil, core.NewPdfDocumentFormatException("CMap parser has not been initialized; call SetCMapParser first")
	}

	return cmapParser.Parse(bytes)
}
