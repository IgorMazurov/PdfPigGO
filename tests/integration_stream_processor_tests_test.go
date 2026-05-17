//go:build integration

package pdfpig_test

import (
	"strings"
	"testing"

	pdfpig "github.com/uglytoad/pdfpig/go"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/fonts"
	cmap "github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/fonts/cidfonts"
	"github.com/uglytoad/pdfpig/go/fonts/composite"
	pdfparser "github.com/uglytoad/pdfpig/go/fonts/parser"
	"github.com/uglytoad/pdfpig/go/fonts/systemfonts"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/graphics"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/parser"
	pdffontparser "github.com/uglytoad/pdfpig/go/pdf_fonts/parser"
	handlers "github.com/uglytoad/pdfpig/go/pdf_fonts/parser/handlers"
	pdffonts "github.com/uglytoad/pdfpig/go/pdf_fonts"
	"github.com/uglytoad/pdfpig/go/testpages"
	"github.com/uglytoad/pdfpig/go/testutil"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
)

func TestStreamProcessorTextOnly(t *testing.T) {
	path := testutil.GetDocumentPath("cat-genetics", true)

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pageAny, err := doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}

	page, ok := pageAny.(*content.Page)
	if !ok {
		t.Fatalf("GetPage(1): expected *content.Page, got %T", pageAny)
	}

	letters := page.Letters()
	if len(letters) == 0 {
		t.Fatal("page has no letters")
	}

	expected := buildExpectedText(letters)

	parsingOpts := &content.ParsingOptions{Logger: logging.NoopLog, UseLenientParsing: true}
	pdfScanner, resourceStore, filterProvider, pageContentParser, err := createDependencies(doc, path, parsingOpts)
	if err != nil {
		t.Fatalf("createDependencies: %v", err)
	}

	factory := testpages.NewTextOnlyPageInformationFactory(
		pdfScanner, resourceStore, filterProvider, pageContentParser, parsingOpts,
	)
	doc.AddPageFactory(factory)

	textOnlyPagePtr, err := content.GetTypedPage[testpages.TextOnlyPage](doc, 1)
	if err != nil {
		t.Logf("GetTypedPage error: %v", err)
		t.Fatalf("GetTypedPage[testpages.TextOnlyPage](1): %v", err)
	}

	textOnlyPage := *textOnlyPagePtr

	if textOnlyPage.Text != expected {
		t.Fatalf("Text mismatch:\nExpected length: %d\nActual length:   %d\nExpected[:100]: %q\nActual[:100]:    %q",
			len(expected), len(textOnlyPage.Text),
			truncate(expected, 100), truncate(textOnlyPage.Text, 100),
		)
	}

	t.Logf("Text match: %d chars from %d letters", len(expected), len(letters))
}

// createDependencies constructs all dependencies needed for TextOnlyPageInformationFactory.
// Uses the document's scanner to ensure consistent object resolution.
func createDependencies(doc *content.PdfDocument, path string, parsingOptions *content.ParsingOptions) (
	tokenization.PdfTokenScanner,
	content.ResourceStore,
	content.LookupFilterProvider,
	content.PageContentParser,
	error,
) {
	// Use the document's existing scanner - it has cached objects and correct xref data
	pdfScanner := doc.Structure.Scanner()

	maxStackDepth := 256
	stackDepthGuard, err := core.NewStackDepthGuard(maxStackDepth)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	filterProviderWithLookup := filters.NewFilterProviderWithLookup(filters.Instance)

	fontFactory := createFontFactory(pdfScanner, filterProviderWithLookup, stackDepthGuard, parsingOptions)

	lookupAdapter := &testLookupFilterAdapter{provider: filterProviderWithLookup}

	resourceStore := content.NewResourceStore(
		pdfScanner, fontFactory, lookupAdapter, parsingOptions,
	)

	pageContentParser := parser.NewPageContentParser(
		graphics.GetReflectionFactory(), stackDepthGuard, parsingOptions.UseLenientParsing,
	)

	return pdfScanner, resourceStore, lookupAdapter, pageContentParser, nil
}

// testStreamDecoder implements tokenization.StreamDecoder.
type testStreamDecoder struct {
	provider *filters.FilterProviderWithLookup
}

func (d *testStreamDecoder) DecodeStream(stream *tokens.StreamToken) []byte {
	if d.provider == nil {
		return stream.Data()
	}
	filtersList, err := d.provider.GetFilters(stream.StreamDictionary)
	if err != nil || len(filtersList) == 0 {
		return stream.Data()
	}
	data := stream.Data()
	for _, f := range filtersList {
		result, decodeErr := f.Decode(data, stream.StreamDictionary, d.provider, 0)
		if decodeErr != nil {
			return data
		}
		data = result
	}
	return data
}

// testEncryptionNoOp implements the encryptionHandler interface.
type testEncryptionNoOp struct{}

func (testEncryptionNoOp) Decrypt(reference core.IndirectReference, token tokens.Token) tokens.Token {
	return token
}

// testLookupFilterAdapter adapts FilterProviderWithLookup to util.PatternParserProvider.
type testLookupFilterAdapter struct {
	provider *filters.FilterProviderWithLookup
}

func (a *testLookupFilterAdapter) DecodeStream(stream *tokens.StreamToken, scanner tokenization.PdfTokenScanner) []byte {
	filtersList, err := a.provider.GetFilters(stream.StreamDictionary)
	if err != nil || len(filtersList) == 0 {
		return stream.Data()
	}
	data := stream.Data()
	for _, f := range filtersList {
		result, decodeErr := f.Decode(data, stream.StreamDictionary, a.provider, 0)
		if decodeErr != nil {
			return data
		}
		data = result
	}
	return data
}

func (a *testLookupFilterAdapter) GetFilters(dictionary *tokens.DictionaryToken) ([]filters.Filter, error) {
	return a.provider.GetFilters(dictionary)
}

func (a *testLookupFilterAdapter) GetNamedFilters(names []*tokens.NameToken) ([]filters.Filter, error) {
	return a.provider.GetNamedFilters(names)
}

func (a *testLookupFilterAdapter) GetAllFilters() []filters.Filter {
	return a.provider.GetAllFilters()
}

func (a *testLookupFilterAdapter) GetFiltersWithScanner(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) ([]filters.Filter, error) {
	return a.provider.GetFiltersWithScanner(dictionary, scanner)
}

// createFontFactory creates a FontFactory with all required handlers.
func createFontFactory(
	pdfScanner tokenization.PdfTokenScanner,
	filterProvider *filters.FilterProviderWithLookup,
	stackDepthGuard *core.StackDepthGuard,
	parsingOptions *content.ParsingOptions,
) pdffonts.FontFactory {
	cidFontFactory := pdfparser.NewCidFontFactory(logging.NoopLog, pdfScanner, filterProvider)
	encodingReader := pdffontparser.NewEncodingReader(pdfScanner)

	cmap.SetCMapParser(pdfparser.NewCMapParser())
	cmapCache := cmap.NewCMapLocalCache(filterProvider, pdfScanner)

	type0Handler := handlers.NewType0FontHandler(
		testCidFontAdapterFactory{factory: cidFontFactory},
		pdfScanner,
		testCMapCacheAdapter{cache: cmapCache},
		testCMapCacheLookup,
		testParsingOptionsAdapter{options: parsingOptions},
		testType0FontConstructor,
	)

	type1Handler := handlers.NewType1FontHandler(
		pdfScanner, filterProvider, encodingReader, testCMapCacheAdapter{cache: cmapCache}, stackDepthGuard, parsingOptions.UseLenientParsing, logging.NoopLog,
	)

	trueTypeHandler := handlers.NewTrueTypeFontHandler(
		logging.NoopLog, pdfScanner, filterProvider, encodingReader, testCMapCacheAdapter{cache: cmapCache},
		systemfonts.Instance, type1Handler,
		pdfparser.GetWidths,
		pdfparser.GetFontDescriptor,
		pdfparser.GetName,
	)

	fontFactory := pdffonts.NewFontFactoryImpl(
		logging.NoopLog, type0Handler, trueTypeHandler, type1Handler,
		handlers.NewType3FontHandler(pdfScanner, encodingReader, testCMapCacheAdapter{cache: cmapCache}),
	)

	return fontFactory
}

// Adapter types bridging between packages

type testCidFontAdapterFactory struct {
	factory *pdfparser.CidFontFactory
}

func (a testCidFontAdapterFactory) Generate(dictionary *tokens.DictionaryToken) handlers.CombinedCidFont {
	result, err := a.factory.Generate(dictionary)
	if err != nil || result == nil {
		return nil
	}
	if combined, ok := any(result).(handlers.CombinedCidFont); ok {
		return combined
	}
	return testCidFontWrapper{cid: result}
}

type testCidFontWrapper struct {
	cid cidfonts.CidFont
}

func (w testCidFontWrapper) SystemInfo() handlers.CidFontSystemInfo {
	return w.cid.SystemInfo()
}

func (w testCidFontWrapper) Details() fonts.FontDetails {
	return w.cid.Details()
}

func (w testCidFontWrapper) FontMatrix() core.TransformationMatrix {
	return w.cid.FontMatrix()
}

func (w testCidFontWrapper) GetDescent() float64 {
	if desc := w.cid.Descriptor(); desc != nil {
		return desc.Descent
	}
	return 0
}

func (w testCidFontWrapper) GetAscent() float64 {
	if desc := w.cid.Descriptor(); desc != nil {
		return desc.Ascent
	}
	return 0
}

func (w testCidFontWrapper) GetWidthFromDictionary(cid int) float64 {
	return w.cid.GetWidthFromDictionary(cid)
}

func (w testCidFontWrapper) GetWidthFromFont(characterIdentifier int) float64 {
	return w.cid.GetWidthFromFont(characterIdentifier)
}

func (w testCidFontWrapper) GetBoundingBox(characterIdentifier int) (core.PdfRectangle, error) {
	return w.cid.GetBoundingBox(characterIdentifier)
}

func (w testCidFontWrapper) GetPositionVector(characterIdentifier int) geometry.PdfVector {
	return w.cid.GetPositionVector(characterIdentifier)
}

func (w testCidFontWrapper) GetDisplacementVector(characterIdentifier int) geometry.PdfVector {
	return w.cid.GetDisplacementVector(characterIdentifier)
}

func (w testCidFontWrapper) GetFontMatrix(characterIdentifier int) core.TransformationMatrix {
	return w.cid.FontMatrix()
}

func (w testCidFontWrapper) TryGetPath(characterCode int) ([]core.PdfSubpath, bool) {
	return nil, false
}

func (w testCidFontWrapper) TryGetNormalisedPath(characterCode int) ([]core.PdfSubpath, bool) {
	return nil, false
}

type testCMapCacheAdapter struct {
	cache *cmap.CMapLocalCache
}

func (a testCMapCacheAdapter) Get(name string) fonts.CMapProvider {
	cmap, _ := a.cache.TryGetByName(name)
	if cmap == nil {
		return nil
	}
	return cmap
}

func (a testCMapCacheAdapter) TryGetByName(name string) (fonts.CMapProvider, bool) {
	cmap, ok := a.cache.TryGetByName(name)
	if !ok || cmap == nil {
		return nil, false
	}
	return cmap, true
}

func (a testCMapCacheAdapter) TryGetByStream(streamToken *tokens.StreamToken) (fonts.CMapProvider, bool) {
	cmap, ok := a.cache.TryGetByStream(streamToken)
	if !ok || cmap == nil {
		return nil, false
	}
	return cmap, true
}

var testCMapCacheLookup = func(name string) (fonts.CMapProvider, handlers.CidFontSystemInfo, bool) {
	return nil, nil, false
}

type testParsingOptionsAdapter struct {
	options *content.ParsingOptions
}

func (a testParsingOptionsAdapter) UseLenientParsing() bool {
	return a.options.UseLenientParsing
}

func (a testParsingOptionsAdapter) Logger() logging.Log {
	if a.options != nil && a.options.Logger != nil {
		return a.options.Logger
	}
	return logging.NoopLog
}

var testType0FontConstructor = func(
	baseFont *tokens.NameToken,
	cidFont handlers.CombinedCidFont,
	cmapVal fonts.CMapProvider,
	toUnicodeCMap fonts.CMapProvider,
	ucs2CMap fonts.CMapProvider,
	useLenientParsing bool,
	isChineseJapaneseOrKorean bool,
) fonts.Font {
	if cidFont == nil || cmapVal == nil {
		return nil
	}

	cmapValConcrete, ok1 := cmapVal.(*cmap.CMap)
	if !ok1 {
		return nil
	}

	toUnicodeConcrete := (*cmap.CMap)(nil)
	if toUnicodeCMap != nil {
		if tc, ok := toUnicodeCMap.(*cmap.CMap); ok {
			toUnicodeConcrete = tc
		}
	}

	ucs2Concrete := (*cmap.CMap)(nil)
	if ucs2CMap != nil {
		if uc, ok := ucs2CMap.(*cmap.CMap); ok {
			ucs2Concrete = uc
		}
	}

	font, err := composite.NewType0Font(
		baseFont, cidFont, cmapValConcrete, toUnicodeConcrete, ucs2Concrete,
		useLenientParsing, isChineseJapaneseOrKorean)
	if err != nil {
		return nil
	}
	return font
}

func buildExpectedText(letters []*content.Letter) string {
	var sb strings.Builder
	for _, l := range letters {
		sb.WriteString(l.Value)
	}
	return sb.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
