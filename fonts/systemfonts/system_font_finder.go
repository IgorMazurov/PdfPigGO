package systemfonts

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	standard14fonts "github.com/uglytoad/pdfpig/go/fonts/standard14_fonts"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
)

// SystemFontFinder finds named fonts from the host operating system.
type SystemFontFinder interface {
	// GetTrueTypeFont returns the TrueType font with the specified name,
	// trying various naming conventions and substitute names.
GetTrueTypeFont(name string) *truetypeparser.TrueTypeFont
}

var (
	nameSubstitutes  map[string][]string
	availableFonts   []SystemFontRecord
	fontsByFirstChar map[rune][]SystemFontRecord
	initOnce         sync.Once
	initErr          error
)

// fontCache caches parsed TrueType fonts by name.
type fontCache struct {
	mu   sync.RWMutex
	data map[string]*truetypeparser.TrueTypeFont
}

func newFontCache() *fontCache {
	return &fontCache{data: make(map[string]*truetypeparser.TrueTypeFont)}
}

func (c *fontCache) get(name string) (*truetypeparser.TrueTypeFont, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	font, ok := c.data[name]
	return font, ok
}

func (c *fontCache) add(name string, font *truetypeparser.TrueTypeFont) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, exists := c.data[name]; !exists {
		c.data[name] = font
	}
}

var cache = newFontCache()

// nameToFileMap caches font name to file path mappings discovered during search.
type nameToFileMap struct {
	mu   sync.RWMutex
	data map[string]string
}

func newNameToFileMap() *nameToFileMap {
	return &nameToFileMap{data: make(map[string]string)}
}

func (m *nameToFileMap) get(name string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[name]
	return v, ok
}

func (m *nameToFileMap) add(name, path string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.data[name]; !exists {
		m.data[name] = path
	}
}

var nameToFileNameMap = newNameToFileMap()

// readFilesSet tracks font file paths already scanned during a search.
type readFilesSet struct {
	mu   sync.Mutex
	data map[string]struct{}
}

func newReadFilesSet() *readFilesSet {
	return &readFilesSet{data: make(map[string]struct{})}
}

func (s *readFilesSet) contains(path string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.data[path]
	return ok
}

func (s *readFilesSet) add(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[path] = struct{}{}
}

var readFiles = newReadFilesSet()

// Instance is the singleton SystemFontFinder instance.
var Instance SystemFontFinder

func init() {
	initOnce.Do(func() {
		nameSubstitutes = buildNameSubstitutes()
		lister := createLister()
		if lister == nil {
			initErr = notSupportedError(runtime.GOOS)
			return
		}
		availableFonts = lister.GetAllFonts()
		fontsByFirstChar = buildFontsByFirstChar(availableFonts)
		Instance = &systemFontFinderImpl{}
	})
}

func notSupportedError(goos string) error {
	return &notSupportedErrorType{goos: goos}
}

type notSupportedErrorType struct{ goos string }

func (e *notSupportedErrorType) Error() string {
	return "unsupported operating system: " + e.goos
}

// Is implements errors.Is for notSupportedError.
func (e *notSupportedErrorType) Is(target error) bool {
	_, ok := target.(*notSupportedErrorType)
	return ok
}

var ErrUnsupportedOS = &notSupportedErrorType{}

func buildNameSubstitutes() map[string][]string {
	dict := map[string][]string{
		"Courier":             {"CourierNew", "CourierNewPSMT", "LiberationMono", "NimbusMonL-Regu"},
		"Courier-Bold":        {"CourierNewPS-BoldMT", "CourierNew-Bold", "LiberationMono-Bold", "NimbusMonL-Bold"},
		"Courier-Oblique":     {"CourierNewPS-ItalicMT", "CourierNew-Italic", "LiberationMono-Italic", "NimbusMonL-ReguObli"},
		"Courier-BoldOblique": {"CourierNewPS-BoldItalicMT", "CourierNew-BoldItalic", "LiberationMono-BoldItalic", "NimbusMonL-BoldObli"},
		"Helvetica":           {"ArialMT", "Arial", "LiberationSans", "NimbusSanL-Regu"},
		"Helvetica-Bold":      {"Arial-BoldMT", "Arial-Bold", "LiberationSans-Bold", "NimbusSanL-Bold"},
		"Helvetica-BoldOblique": {"Arial-BoldItalicMT", "Helvetica-BoldItalic", "LiberationSans-BoldItalic", "NimbusSanL-BoldItal"},
		"Helvetica-Oblique":   {"Arial-ItalicMT", "Arial-Italic", "Helvetica-Italic", "LiberationSans-Italic", "NimbusSanL-ReguItal"},
		"Times-Roman":         {"TimesNewRomanPSMT", "TimesNewRoman", "TimesNewRomanPS", "LiberationSerif", "NimbusRomNo9L-Regu"},
		"Times-Bold":          {"TimesNewRomanPS-BoldMT", "TimesNewRomanPS-Bold", "TimesNewRoman-Bold", "LiberationSerif-Bold", "NimbusRomNo9L-Medi"},
		"Times-Italic":        {"TimesNewRomanPS-ItalicMT", "TimesNewRomanPS-Italic", "TimesNewRoman-Italic", "LiberationSerif-Italic", "NimbusRomNo9L-ReguItal"},
		"Times-BoldItalic":    {"TimesNewRomanPS-BoldItalicMT", "TimesNewRomanPS-BoldItalic", "TimesNewRoman-BoldItalic", "LiberationSerif-BoldItalic", "NimbusRomNo9L-MediItal"},
		"Symbol":              {"SymbolMT", "StandardSymL"},
		"ZapfDingbats":        {"ZapfDingbatsITC", "Dingbats", "MS-Gothic"},
	}

	names := standard14fonts.GetNames()
	for _, name := range names {
		if _, exists := dict[name]; !exists {
			mappedName := standard14fonts.GetMappedFontName(name)
			if subs, ok := dict[mappedName]; ok {
				dict[name] = subs
			} else {
				dict[name] = []string{mappedName}
			}
		}
	}

	return dict
}

func createLister() SystemFontLister {
	switch runtime.GOOS {
	case "windows":
		return NewWindowsSystemFontLister()
	case "darwin":
		return NewMacSystemFontLister()
	case "linux":
		return NewLinuxSystemFontLister()
	default:
		return nil
	}
}

func buildFontsByFirstChar(fonts []SystemFontRecord) map[rune][]SystemFontRecord {
	byChar := make(map[rune][]SystemFontRecord)

	for _, record := range fonts {
		fn := filepath.Base(record.Path())
		if fn == "" {
			continue
		}

		key := toUpperRune(rune(fn[0]))
		byChar[key] = append(byChar[key], record)
	}

	return byChar
}

func toUpperRune(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 'a' + 'A'
	}
	return r
}

// systemFontFinderImpl implements SystemFontFinder.
type systemFontFinderImpl struct{}

// NewSystemFontFinder creates a new SystemFontFinder instance.
func NewSystemFontFinder() SystemFontFinder {
	return &systemFontFinderImpl{}
}

// GetTrueTypeFont returns the TrueType font with the specified name,
// trying various naming conventions and substitute names.
func (f *systemFontFinderImpl) GetTrueTypeFont(name string) *truetypeparser.TrueTypeFont {
	if initErr != nil {
		return nil
	}

	var triedNames []string
	var result *truetypeparser.TrueTypeFont
	var noDash, noComma, sub, regName string

	result = f.getTrueTypeFontNamed(name)
	if result != nil {
		goto done
	}
	triedNames = append(triedNames, name)

	noDash = strings.ReplaceAll(name, "-", "")
	if noDash != name {
		result = f.getTrueTypeFontNamed(noDash)
		if result != nil {
			goto done
		}
		triedNames = append(triedNames, noDash)
	}

	noComma = strings.ReplaceAll(name, ",", "-")
	if noComma != name {
		result = f.getTrueTypeFontNamed(noComma)
		if result != nil {
			goto done
		}
		triedNames = append(triedNames, noComma)
	}

	for _, sub = range getSubstituteNames(name) {
		result = f.getTrueTypeFontNamed(sub)
		if result != nil {
			goto done
		}
		triedNames = append(triedNames, sub)
	}

	regName = name + "-Regular"
	result = f.getTrueTypeFontNamed(regName)

done:
	if result != nil {
		for _, n := range triedNames {
			cache.add(n, result)
		}
	}
	return result
}

func getSubstituteNames(name string) []string {
	name = strings.ReplaceAll(name, " ", "")

	if subs, ok := nameSubstitutes[name]; ok {
		return subs
	}

	return nil
}

func (f *systemFontFinderImpl) getTrueTypeFontNamed(name string) *truetypeparser.TrueTypeFont {
	if name == "" {
		return nil
	}

	if cached, ok := cache.get(name); ok {
		return cached
	}

	if fileName, ok := nameToFileNameMap.get(name); ok {
		if result := f.tryReadFile(fileName, false, name); result != nil {
			return result
		}
		return nil
	}

	firstChar := toUpperRune(rune(name[0]))

	if candidates, ok := fontsByFirstChar[firstChar]; ok {
		for _, record := range candidates {
			if font := f.tryGetTrueTypeFont(name, record); font != nil {
				return font
			}
		}
	}

	for _, record := range availableFonts {
		fn := filepath.Base(record.Path())
		if fn == "" {
			continue
		}

		localFirstChar := toUpperRune(rune(fn[0]))
		if localFirstChar == firstChar {
			continue
		}

		if font := f.tryGetTrueTypeFont(name, record); font != nil {
			return font
		}
	}

	return nil
}

func (f *systemFontFinderImpl) tryGetTrueTypeFont(name string, record SystemFontRecord) *truetypeparser.TrueTypeFont {
	if record.Type() != TrueType {
		return nil
	}

	if readFiles.contains(record.Path()) {
		return nil
	}

	return f.tryReadFile(record.Path(), true, name)
}

func (f *systemFontFinderImpl) tryReadFile(fileName string, readNameFirst bool, fontName string) *truetypeparser.TrueTypeFont {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil
	}

	data := truetypeparser.NewTrueTypeDataBytes(bytes)

	var psName string

	if readNameFirst {
		nameTable := truetypeparser.GetNameTable(data)

		if nameTable == nil {
			readFiles.add(fileName)
			return nil
		}

		psName = nameTable.GetPostscriptName()
		fontNameFromFile := psName
		if fontNameFromFile == "" {
			fontNameFromFile = nameTable.FontName()
		}

		nameToFileNameMap.add(fontNameFromFile, fileName)

		if !strings.EqualFold(fontNameFromFile, fontName) {
			readFiles.add(fileName)
			return nil
		}

		data.Seek(0, 0) // io.SeekStart
	}

	font := truetypeparser.Parse(data)
	if font == nil {
		readFiles.add(fileName)
		return nil
	}

	if psName == "" {
		psName = font.Name()
	}

	cache.add(psName, font)
	readFiles.add(fileName)

	return font
}


