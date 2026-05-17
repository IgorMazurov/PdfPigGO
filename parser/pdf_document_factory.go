package parser

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/uglytoad/pdfpig/go/acroforms"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/crossreference"
	"github.com/uglytoad/pdfpig/go/document"
	"github.com/uglytoad/pdfpig/go/encryption"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/geometry"
	cmap "github.com/uglytoad/pdfpig/go/fonts/cmap"
	"github.com/uglytoad/pdfpig/go/fonts/cidfonts"
	"github.com/uglytoad/pdfpig/go/fonts/composite"
	"github.com/uglytoad/pdfpig/go/fonts/systemfonts"
	"github.com/uglytoad/pdfpig/go/graphics"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/outline"
	pdffonts "github.com/uglytoad/pdfpig/go/pdf_fonts"
	"github.com/uglytoad/pdfpig/go/parser/filestructure"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/pdf_fonts/parser/handlers"
	pdfparser "github.com/uglytoad/pdfpig/go/fonts/parser"
	pdffontparser "github.com/uglytoad/pdfpig/go/pdf_fonts/parser"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// OpenMemory opens a PDF document from the given byte slice.
func OpenMemory(memory []byte, options *content.ParsingOptions) (*content.PdfDocument, error) {
	inputBytes := core.NewMemoryInputBytes(memory)
	return open(inputBytes, options)
}

// OpenFile opens a PDF document from the file at the given path.
func OpenFile(filename string, options *content.ParsingOptions) (*content.PdfDocument, error) {
	info, err := osStat(filename)
	if err != nil {
		return nil, fmt.Errorf("no file exists at: %s", filename)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory: %s", filename)
	}

	data, err := readFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %w", filename, err)
	}

	return OpenMemory(data, options)
}

// OpenStream opens a PDF document from the given seekable stream.
// If the stream does not support seeking, it will be copied into memory first.
func OpenStream(stream io.ReadSeeker, options *content.ParsingOptions) (*content.PdfDocument, error) {
	if stream == nil {
		return nil, fmt.Errorf("stream cannot be nil")
	}

	initialPosition, err := stream.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("cannot determine stream position: %w", err)
	}

	var inputBytes core.InputBytes

	_, seekErr := stream.Seek(0, io.SeekStart)
	if seekErr != nil {
		buf := new(bytes.Buffer)
		if _, copyErr := buf.ReadFrom(stream); copyErr != nil {
			return nil, fmt.Errorf("failed to copy non-seekable stream: %w", copyErr)
		}
		inputBytes = core.NewMemoryInputBytes(buf.Bytes())
	} else {
		streamInput, streamErr := core.NewStreamInputBytes(stream, false)
		if streamErr != nil {
			return nil, fmt.Errorf("failed to create stream input bytes: %w", streamErr)
		}
		inputBytes = streamInput
	}

	doc, err := open(inputBytes, options)
	if err != nil {
		if initialPosition != 0 {
			return nil, fmt.Errorf("could not parse document due to an error, the input stream was not at position zero when provided to the Open method: %w", err)
		}
		return nil, err
	}

	return doc, nil
}

// open opens a PDF document from the given input bytes.
func open(inputBytes core.InputBytes, options *content.ParsingOptions) (doc *content.PdfDocument, err error) {
	defer func() {
		if r := recover(); r != nil {
			if formatErr, ok := r.(*core.PdfDocumentFormatException); ok {
				doc = nil
				err = formatErr
				return
			}
			panic(r)
		}
	}()
	if options == nil {
		options = &content.ParsingOptions{
			UseLenientParsing: true,
			MaxStackDepth:     256,
		}
	}

	// Apply defaults matching C# behavior: UseLenientParsing defaults to true,
	// MaxStackDepth defaults to 256.
	maxStackDepth := options.MaxStackDepth
	if maxStackDepth <= 0 {
		maxStackDepth = 256
	}

	if options.Logger == nil {
		options.Logger = logging.NoopLog
	}

	stackDepthGuard, err := core.NewStackDepthGuard(maxStackDepth)
	if err != nil {
		return nil, fmt.Errorf("failed to create stack depth guard: %w", err)
	}

	passwords := preparePasswords(options)

	tokenScanner := tokenization.NewCoreTokenScanner(
		inputBytes, true, stackDepthGuard, tokenization.ScannerScopeNone,
		nil, options.UseLenientParsing, false)

	return openDocument(inputBytes, tokenScanner, passwords, options, stackDepthGuard)
}

// preparePasswords builds the password list from parsing options.
// Following C# logic: collect Password and non-empty entries from Passwords,
// then ensure empty string is always present as fallback.
func preparePasswords(options *content.ParsingOptions) []string {
	passwords := make([]string, 0)

	if options.Password != "" {
		passwords = append(passwords, options.Password)
	}

	for _, p := range options.Passwords {
		if p != "" {
			passwords = append(passwords, p)
		}
	}

	foundEmpty := false
	for _, p := range passwords {
		if p == "" {
			foundEmpty = true
			break
		}
	}
	if !foundEmpty {
		passwords = append(passwords, "")
	}

	options.Passwords = passwords
	return passwords
}

// openDocument performs the full document parsing pipeline: header parse, first pass,
// trailer resolution, encryption setup, font handler chain construction, resource store,
// page factory, catalog creation, and final PdfDocument assembly.
func openDocument(
	inputBytes core.InputBytes,
	scanner tokenization.SeekableTokenScanner,
	passwords []string,
	parsingOptions *content.ParsingOptions,
	stackDepthGuard *core.StackDepthGuard,
) (*content.PdfDocument, error) {
	filterProvider := createFilterProvider(parsingOptions)

	version, err := filestructure.Parse(scanner, inputBytes, parsingOptions.UseLenientParsing, parsingOptions.Logger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file header: %w", err)
	}

	fileHeaderOffset := filestructure.NewFileHeaderOffset(int(version.OffsetInFile))

	initialParse, err := filestructure.FirstPassParse(fileHeaderOffset, inputBytes, scanner, parsingOptions.Logger)
	if err != nil {
		return nil, fmt.Errorf("first pass parse failed: %w", err)
	}

	if initialParse.Trailer == nil {
		return nil, core.NewPdfDocumentFormatException(
			"Could not find an xref trailer or stream dictionary in the input file.")
	}

	trailer, err := crossreference.NewTrailerDictionary(initialParse.Trailer, parsingOptions.UseLenientParsing)
	if err != nil {
		return nil, fmt.Errorf("failed to create trailer dictionary: %w", err)
	}

	locationProvider := tokenization.NewObjectLocationProviderImpl(
		initialParse.XrefOffsets,
		initialParse.BruteForceOffsets,
		inputBytes,
	)

	streamDecoder := newStreamDecoder(filterProvider)

	pdfScanner := tokenization.NewPdfTokenScanner(
		inputBytes, locationProvider, encryption.NoOpInstance, streamDecoder,
		int(version.OffsetInFile),
		&tokenization.ScannerParsingOptions{UseLenientParsing: parsingOptions.UseLenientParsing},
		stackDepthGuard)

	rootRef, rootDictionary, encryptionDictionary, err := parseTrailer(
		trailer, parsingOptions.UseLenientParsing, pdfScanner)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trailer: %w", err)
	}

	encHandler, encDict, err := resolveEncryption(encryptionDictionary, trailer, passwords, pdfScanner, parsingOptions.Logger)
	if err != nil {
		return nil, fmt.Errorf("resolve encryption: %w", err)
	}

	pdfScanner.UpdateEncryptionHandler(encHandler)

	cidFontFactory := pdfparser.NewCidFontFactory(parsingOptions.Logger, pdfScanner, filterProvider)

	encodingReader := pdffontparser.NewEncodingReader(pdfScanner)

	cmap.SetCMapParser(pdfparser.NewCMapParser())

	cmapCache := cmap.NewCMapLocalCache(filterProvider, pdfScanner)

	type0Handler := handlers.NewType0FontHandler(
		newCidFontAdapterFactory(cidFontFactory),
		pdfScanner,
		newCMapCacheAdapter(cmapCache),
		createCMapCacheLookup(),
		newParsingOptionsAdapter(parsingOptions),
		createType0FontConstructor(),
	)

	type1Handler := handlers.NewType1FontHandler(
		pdfScanner, filterProvider, encodingReader, newCMapCacheAdapter(cmapCache), stackDepthGuard, parsingOptions.UseLenientParsing, parsingOptions.Logger)

	trueTypeHandler := handlers.NewTrueTypeFontHandler(
		parsingOptions.Logger, pdfScanner, filterProvider, encodingReader, newCMapCacheAdapter(cmapCache),
		systemfonts.Instance, type1Handler,
		pdfparser.GetWidths,
		pdfparser.GetFontDescriptor,
		pdfparser.GetName)

	fontFactory := pdffonts.NewFontFactoryImpl(
		parsingOptions.Logger, type0Handler, trueTypeHandler, type1Handler,
		handlers.NewType3FontHandler(pdfScanner, encodingReader, newCMapCacheAdapter(cmapCache)))

	lookupAdapter := newLookupFilterAdapter(filterProvider)

	resourceContainer := content.NewResourceStore(
		pdfScanner, fontFactory, lookupAdapter, parsingOptions)

	information, err := CreateDocumentInformation(pdfScanner, trailer, parsingOptions.UseLenientParsing)
	if err != nil {
		return nil, fmt.Errorf("failed to create document information: %w", err)
	}

	pageContentParser := NewPageContentParser(
		graphicGetReflectionFactory(), stackDepthGuard, parsingOptions.UseLenientParsing)

	pageFactory := NewPageFactory(
		pdfScanner, resourceContainer, lookupAdapter, pageContentParser, parsingOptions)

	catalog, err := CreateCatalog[*content.Page](
		rootRef, rootDictionary, pdfScanner, pageFactory,
		parsingOptions.Logger, parsingOptions.UseLenientParsing)
	if err != nil {
		return nil, fmt.Errorf("failed to create catalog: %w", err)
	}

	objectOffsets := initialParse.XrefOffsets
	if initialParse.BruteForceOffsets != nil && len(initialParse.BruteForceOffsets) > 0 {
		objectOffsets = initialParse.BruteForceOffsets
	}

	acroFormFactory, acroErr := acroforms.NewAcroFormFactory(
		pdfScanner, filterProvider, objectOffsets)
	if acroErr != nil {
		parsingOptions.Logger.Warn(fmt.Sprintf("Failed to create AcroForm factory: %v", acroErr))
	}

	bookmarksProviderInstance := outline.NewBookmarksProvider(parsingOptions.Logger, pdfScanner)
	bookmarksRetriever := &bookmarksRetrieverAdapter{provider: bookmarksProviderInstance}

	advancedAccess, advErr := document.NewAdvancedPdfDocumentAccess(
		pdfScanner, lookupAdapter, catalog)
	if advErr != nil {
		parsingOptions.Logger.Warn(fmt.Sprintf("Failed to create AdvancedPdfDocumentAccess: %v", advErr))
	}

	document, docErr := content.NewPdfDocument(
		inputBytes, version, catalog, information, encDict,
		pdfScanner, lookupAdapter, acroFormFactory,
		bookmarksRetriever, parsingOptions, advancedAccess)
	if docErr != nil {
		return nil, fmt.Errorf("failed to create PDF document: %w", docErr)
	}

	return document, nil
}

// resolveEncryption creates the encryption handler and dictionary from the trailer's
// encryption token. Returns NoOpInstance when no encryption is present.
func resolveEncryption(
	encryptionDictionary *tokens.DictionaryToken,
	trailer *crossreference.TrailerDictionary,
	passwords []string,
	scanner tokenization.PdfTokenScanner,
	log logging.Log,
) (encryption.EncryptionHandler, *encryption.EncryptionDictionary, error) {
	if encryptionDictionary == nil {
		return encryption.NoOpInstance, nil, nil
	}

	encDict, err := encryption.ReadEncryptionDictionary(encryptionDictionary, scanner)
	if err != nil || encDict == nil {
		log.Warn(fmt.Sprintf("Failed to read encryption dictionary: %v", err))
		return encryption.NoOpInstance, nil, fmt.Errorf("read encryption dictionary: %w", err)
	}

	handler, handlerErr := encryption.NewEncryptionHandler(encDict, trailer, passwords)
	if handlerErr != nil || handler == nil {
		log.Warn(fmt.Sprintf("Failed to create encryption handler: %v", handlerErr))
		return encryption.NoOpInstance, encDict, fmt.Errorf("create encryption handler: %w", handlerErr)
	}

	return handler, encDict, nil
}

// parseTrailer resolves the root dictionary from the trailer and optionally extracts
// the encryption dictionary. In lenient mode, a missing /Type entry is auto-corrected.
func parseTrailer(
	trailer *crossreference.TrailerDictionary,
	isLenientParsing bool,
	pdfTokenScanner tokenization.PdfTokenScanner,
) (core.IndirectReference, *tokens.DictionaryToken, *tokens.DictionaryToken, error) {
	encryptionDictionary, _ := getEncryptionDictionary(trailer, pdfTokenScanner)

	rootRef := trailer.Root()
	rootDictionary, err := parts.GetByRef[*tokens.DictionaryToken](rootRef, pdfTokenScanner)
	if err != nil {
		return core.IndirectReference{}, nil, nil, err
	}
	if rootDictionary == nil {
		return core.IndirectReference{}, nil, nil, core.NewPdfDocumentFormatException(
			"The root object in the trailer did not resolve to a readable dictionary.")
	}

	if _, hasType := rootDictionary.TryGet(tokens.Type); !hasType && isLenientParsing {
		rootDictionary = rootDictionary.With(tokens.Type, tokens.Catalog)
	}

	return trailer.Root(), rootDictionary, encryptionDictionary, nil
}

// getEncryptionDictionary extracts the encryption dictionary token from the trailer if present.
// Returns (nil, false) when there is no encryption or the token resolves to null.
func getEncryptionDictionary(
	trailer *crossreference.TrailerDictionary,
	pdfTokenScanner tokenization.PdfTokenScanner,
) (*tokens.DictionaryToken, bool) {
	encryptionToken := trailer.EncryptionToken()
	if encryptionToken == nil {
		return nil, false
	}

	encDict, ok := parts.TryGet[*tokens.DictionaryToken](encryptionToken, pdfTokenScanner)
	if !ok {
		if _, isNull := parts.TryGet[*tokens.NullToken](encryptionToken, pdfTokenScanner); isNull {
			return nil, false
		}
		return nil, false
	}

	return encDict, true
}

// createFilterProvider creates the appropriate filter provider from options.
func createFilterProvider(options *content.ParsingOptions) *filters.FilterProviderWithLookup {
	var baseProvider filters.FilterProvider = filters.Instance
	if options.FilterProvider != nil {
		if fp, ok := options.FilterProvider.(filters.FilterProvider); ok {
			baseProvider = fp
		}
	}
	return filters.NewFilterProviderWithLookup(baseProvider)
}

// graphicGetReflectionFactory returns the reflection-based graphics operation factory.
func graphicGetReflectionFactory() graphics.GraphicsStateOperationFactory {
	return graphics.GetReflectionFactory()
}

// createType0FontConstructor returns a Type0FontConstructor that delegates to composite.NewType0Font.
// CombinedCidFont carries all methods needed by NewType0Font (via its internal cidFontRequirements).
// CMapProvider values from the cmap cache are *cmap.CMap instances, so the type assertions succeed.
func createType0FontConstructor() handlers.Type0FontConstructor {
	return func(
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

		cmapValConcrete, okCMap1 := cmapVal.(*cmap.CMap)
		if !okCMap1 {
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
}

// lookupFilterAdapter wraps FilterProviderWithLookup to satisfy content.LookupFilterProvider.
type lookupFilterAdapter struct {
	provider *filters.FilterProviderWithLookup
}

func newLookupFilterAdapter(p *filters.FilterProviderWithLookup) *lookupFilterAdapter {
	return &lookupFilterAdapter{provider: p}
}

// DecodeStream decodes the given stream token using the wrapped filter provider.
func (a *lookupFilterAdapter) DecodeStream(stream *tokens.StreamToken, scanner tokenization.PdfTokenScanner) []byte {
	filtersList, err := a.provider.GetFiltersWithScanner(stream.StreamDictionary, scanner)
	if err != nil || len(filtersList) == 0 {
		return stream.Data()
	}
	data := stream.Data()
	for i, f := range filtersList {
		result, decodeErr := f.Decode(data, stream.StreamDictionary, a.provider, i)
		if decodeErr != nil {
			return data
		}
		data = result
	}
	return data
}

func (a *lookupFilterAdapter) GetFilters(dictionary *tokens.DictionaryToken) ([]filters.Filter, error) {
	return a.provider.GetFilters(dictionary)
}

func (a *lookupFilterAdapter) GetNamedFilters(names []*tokens.NameToken) ([]filters.Filter, error) {
	return a.provider.GetNamedFilters(names)
}

func (a *lookupFilterAdapter) GetAllFilters() []filters.Filter {
	return a.provider.GetAllFilters()
}

func (a *lookupFilterAdapter) GetFiltersWithScanner(dictionary *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner) ([]filters.Filter, error) {
	return a.provider.GetFiltersWithScanner(dictionary, scanner)
}

var _ content.LookupFilterProvider = (*lookupFilterAdapter)(nil)

// streamDecoderWrapper implements tokenization.StreamDecoder by wrapping a filter provider.
type streamDecoderWrapper struct {
	provider *filters.FilterProviderWithLookup
}

func newStreamDecoder(p any) *streamDecoderWrapper {
	if fp, ok := p.(*filters.FilterProviderWithLookup); ok {
		return &streamDecoderWrapper{provider: fp}
	}
	return nil
}

// DecodeStream decodes the stream data using the filter provider.
func (d *streamDecoderWrapper) DecodeStream(stream *tokens.StreamToken) []byte {
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

var _ tokenization.StreamDecoder = (*streamDecoderWrapper)(nil)

// bookmarksRetrieverAdapter adapts outline.BookmarksProvider to content.BookmarksRetriever.
type bookmarksRetrieverAdapter struct {
	provider *outline.BookmarksProvider
}

func (a *bookmarksRetrieverAdapter) GetBookmarks(catalog *content.Catalog, allowContainerNode bool) (any, error) {
	bookmarks, err := a.provider.GetBookmarks(catalog, allowContainerNode)
	if bookmarks == nil {
		return nil, err
	}
	return bookmarks, err
}

// createCMapCacheLookup returns a CMapCacheLookup callback that delegates to
// the global cmap cache. This bridges the import cycle between handlers -> cmap -> pdf_fonts -> handlers.
func createCMapCacheLookup() handlers.CMapCacheLookup {
	return func(name string) (fonts.CMapProvider, handlers.CidFontSystemInfo, bool) {
		result, ok := cmap.TryGet(name)
		if !ok || result == nil {
			return nil, nil, false
		}
		return result, result.Info(), true
	}
}

// cidFontAdapterFactory adapts pdfparser.CidFontFactory.Generate (returns error)
// to handlers.CidFontFactoryProvider.Generate (no error return).
type cidFontAdapterFactory struct {
	factory *pdfparser.CidFontFactory
}

func newCidFontAdapterFactory(factory *pdfparser.CidFontFactory) *cidFontAdapterFactory {
	return &cidFontAdapterFactory{factory: factory}
}

func (a *cidFontAdapterFactory) Generate(dictionary *tokens.DictionaryToken) handlers.CombinedCidFont {
	result, err := a.factory.Generate(dictionary)
	if err != nil || result == nil {
		return nil
	}
	// The concrete CID font type satisfies both cidfonts.CidFont and handlers.CombinedCidFont
	if combined, ok := any(result).(handlers.CombinedCidFont); ok {
		return combined
	}
	// If direct assertion fails (due to CidFontSystemInfo interface mismatch), wrap it.
	return newCidFontWrapper(result)
}

// cidFontWrapper wraps a CidFont and adapts SystemInfo() return type from
// pdffonts.CidFontSystemInfo to handlers.CidFontSystemInfo (same methods, different package).
type cidFontWrapper struct {
	cid cidfonts.CidFont
}

func newCidFontWrapper(cid cidfonts.CidFont) *cidFontWrapper {
	return &cidFontWrapper{cid: cid}
}

func (w *cidFontWrapper) SystemInfo() handlers.CidFontSystemInfo {
	return w.cid.SystemInfo()
}

func (w *cidFontWrapper) Details() fonts.FontDetails {
	return w.cid.Details()
}

func (w *cidFontWrapper) FontMatrix() core.TransformationMatrix {
	return w.cid.FontMatrix()
}

func (w *cidFontWrapper) GetDescent() float64 {
	return w.cid.GetDescent()
}

func (w *cidFontWrapper) GetAscent() float64 {
	return w.cid.GetAscent()
}

func (w *cidFontWrapper) GetWidthFromDictionary(cid int) float64 {
	return w.cid.GetWidthFromDictionary(cid)
}

func (w *cidFontWrapper) GetWidthFromFont(characterIdentifier int) float64 {
	return w.cid.GetWidthFromFont(characterIdentifier)
}

func (w *cidFontWrapper) GetBoundingBox(characterIdentifier int) (core.PdfRectangle, error) {
	return w.cid.GetBoundingBox(characterIdentifier)
}

func (w *cidFontWrapper) GetPositionVector(characterIdentifier int) geometry.PdfVector {
	return w.cid.GetPositionVector(characterIdentifier)
}

func (w *cidFontWrapper) GetDisplacementVector(characterIdentifier int) geometry.PdfVector {
	return w.cid.GetDisplacementVector(characterIdentifier)
}

func (w *cidFontWrapper) GetFontMatrix(characterIdentifier int) core.TransformationMatrix {
	return w.cid.GetFontMatrix(characterIdentifier)
}

func (w *cidFontWrapper) TryGetPath(characterCode int) ([]core.PdfSubpath, bool) {
	return w.cid.TryGetPath(characterCode)
}

func (w *cidFontWrapper) TryGetNormalisedPath(characterCode int) ([]core.PdfSubpath, bool) {
	return w.cid.TryGetNormalisedPath(characterCode)
}

// cmapCacheAdapter adapts cmap.CMapLocalCache return types (*cmap.CMap) to
// handlers.CMapLocalCacheProvider expected types (fonts.CMapProvider).
type cmapCacheAdapter struct {
	cache *cmap.CMapLocalCache
}

func newCMapCacheAdapter(cache *cmap.CMapLocalCache) *cmapCacheAdapter {
	return &cmapCacheAdapter{cache: cache}
}

func (a *cmapCacheAdapter) TryGetByName(name string) (fonts.CMapProvider, bool) {
	result, ok := a.cache.TryGetByName(name)
	if !ok || result == nil {
		return nil, false
	}
	return result, true
}

func (a *cmapCacheAdapter) TryGetByStream(streamToken *tokens.StreamToken) (fonts.CMapProvider, bool) {
	result, ok := a.cache.TryGetByStream(streamToken)
	if !ok || result == nil {
		return nil, false
	}
	return result, true
}

var _ handlers.CMapLocalCacheProvider = (*cmapCacheAdapter)(nil)

// parsingOptionsAdapter adapts content.ParsingOptions to handlers.ParsingOptionsProvider.
type parsingOptionsAdapter struct {
	opts *content.ParsingOptions
}

func newParsingOptionsAdapter(opts *content.ParsingOptions) *parsingOptionsAdapter {
	return &parsingOptionsAdapter{opts: opts}
}

func (a *parsingOptionsAdapter) UseLenientParsing() bool {
	return a.opts.UseLenientParsing
}

func (a *parsingOptionsAdapter) Logger() logging.Log {
	return a.opts.Logger
}

// FileInfo represents the minimal file info needed by OpenFile.
type FileInfo interface {
	IsDir() bool
}

// osStat wraps os.Stat for testability.
var osStat = func(path string) (FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return &osFileInfo{info: info}, nil
}

// readFile wraps os.ReadFile for testability.
var readFile = os.ReadFile

// osFileInfo wraps os.FileInfo to satisfy FileInfo.
type osFileInfo struct {
	info os.FileInfo
}

func (f *osFileInfo) IsDir() bool {
	return f.info.IsDir()
}
