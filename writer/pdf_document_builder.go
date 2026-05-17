// Package writer provides types for building PDF documents programmatically.
package writer

import (
	"bytes"
	"fmt"
	"io"
	"math"

	"github.com/uglytoad/pdfpig/go/actions"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	standard14fonts "github.com/uglytoad/pdfpig/go/fonts/standard14_fonts"
	truetypeparser "github.com/uglytoad/pdfpig/go/fonts/truetype/parser"
	"github.com/uglytoad/pdfpig/go/graphics"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/outline"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/parser"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
	writerfonts "github.com/uglytoad/pdfpig/go/writer/fonts"
)

// PdfDocumentBuilder provides methods to construct new PDF documents.
type PdfDocumentBuilder struct {
	context             PdfStreamWriter
	pages               map[int]*PdfPageBuilder
	fonts               map[string]*FontStored
	completed           bool
	fontId              int
	version             float64
	archiveStandard     PdfAStandard
	includeDocInfo      bool
	documentInformation *DocumentInformationBuilder
	bookmarks           *outline.Bookmarks
	xmpMetadata         *string
	existingCopies      map[int]map[core.IndirectReference]*tokens.IndirectReferenceToken
	existingTrees       map[string]map[int]*pageInfo
	docCounter          int
}

// pageInfo holds a page dictionary and its parent dictionaries in the page tree.
type pageInfo struct {
	page    *tokens.DictionaryToken
	parents []*tokens.DictionaryToken
}

// DefaultProcSet is the default procedure set for PDF content streams.
var DefaultProcSet = tokens.NewArrayToken([]tokens.Token{
	tokens.MustCreate("PDF"),
	tokens.Text,
	tokens.ImageB,
	tokens.ImageC,
	tokens.ImageI,
})

// NewPdfDocumentBuilder creates a document builder keeping resources in memory with default PDF version 1.7.
func NewPdfDocumentBuilder() *PdfDocumentBuilder {
	return newPdfDocumentBuilder(1.7)
}

// NewPdfDocumentBuilderWithVersion creates a document builder keeping resources in memory with the specified PDF version.
func NewPdfDocumentBuilderWithVersion(version float64) *PdfDocumentBuilder {
	return newPdfDocumentBuilder(version)
}

func newPdfDocumentBuilder(version float64) *PdfDocumentBuilder {
	buf := &bytes.Buffer{}
	ctx := NewPdfStreamWriter(buf, true, nil, func(v float64) {})
	ctx.InitializePdf(version)

	return &PdfDocumentBuilder{
		context:             ctx,
		pages:               make(map[int]*PdfPageBuilder),
		fonts:               make(map[string]*FontStored),
		version:             version,
		includeDocInfo:      true,
		documentInformation: NewDocumentInformationBuilder(),
		existingCopies:      make(map[int]map[core.IndirectReference]*tokens.IndirectReferenceToken),
		existingTrees:       make(map[string]map[int]*pageInfo),
	}
}

// NewPdfDocumentBuilderWithStream creates a document builder using the supplied stream and writer type.
func NewPdfDocumentBuilderWithStream(
	stream io.WriteSeeker,
	disposeStream bool,
	writerType PdfWriterType,
	version float64,
	tokenWriter TokenWriter,
) *PdfDocumentBuilder {
	var ctx PdfStreamWriter

	switch writerType {
	case PdfWriterObjectInMemoryDedup:
		ctx = NewPdfDedupStreamWriter(stream, disposeStream, tokenWriter, func(v float64) {})
	default:
		ctx = NewPdfStreamWriter(stream, disposeStream, tokenWriter, func(v float64) {})
	}

	if ctx != nil {
		ctx.InitializePdf(version)
	}

	return &PdfDocumentBuilder{
		context:             ctx,
		pages:               make(map[int]*PdfPageBuilder),
		fonts:               make(map[string]*FontStored),
		version:             version,
		includeDocInfo:      true,
		documentInformation: NewDocumentInformationBuilder(),
		existingCopies:      make(map[int]map[core.IndirectReference]*tokens.IndirectReferenceToken),
		existingTrees:       make(map[string]map[int]*pageInfo),
	}
}

// ArchiveStandard returns the PDF/A compliance standard of the generated document.
func (b *PdfDocumentBuilder) ArchiveStandard() PdfAStandard {
	return b.archiveStandard
}

// SetArchiveStandard sets the PDF/A compliance standard for the generated document.
func (b *PdfDocumentBuilder) SetArchiveStandard(standard PdfAStandard) {
	b.archiveStandard = standard
}

// IncludeDocumentInformation returns whether to include the document information dictionary.
func (b *PdfDocumentBuilder) IncludeDocumentInformation() bool {
	return b.includeDocInfo
}

// SetIncludeDocumentInformation sets whether to include the document information dictionary.
func (b *PdfDocumentBuilder) SetIncludeDocumentInformation(v bool) {
	b.includeDocInfo = v
}

// DocumentInformation returns the document information builder for setting metadata values.
func (b *PdfDocumentBuilder) DocumentInformation() *DocumentInformationBuilder {
	return b.documentInformation
}

// SetDocumentInformation sets the document information builder.
func (b *PdfDocumentBuilder) SetDocumentInformation(info *DocumentInformationBuilder) {
	b.documentInformation = info
}

// Bookmarks returns the bookmark nodes for the document outline dictionary.
func (b *PdfDocumentBuilder) Bookmarks() *outline.Bookmarks {
	return b.bookmarks
}

// SetBookmarks sets the bookmark nodes for the document outline.
func (b *PdfDocumentBuilder) SetBookmarks(bookmarks *outline.Bookmarks) {
	b.bookmarks = bookmarks
}

// XmpMetadata returns the document level metadata in XMP format.
func (b *PdfDocumentBuilder) XmpMetadata() *string {
	return b.xmpMetadata
}

// SetXmpMetadata sets the document level XMP metadata.
func (b *PdfDocumentBuilder) SetXmpMetadata(metadata *string) {
	b.xmpMetadata = metadata
}

// Pages returns the current page builders keyed by 1-indexed page number.
func (b *PdfDocumentBuilder) Pages() map[int]*PdfPageBuilder {
	return b.pages
}

// Fonts returns the fonts available in the document builder keyed by id.
func (b *PdfDocumentBuilder) Fonts() map[string]*FontStored {
	return b.fonts
}

// CanUseTrueTypeFont determines whether the bytes of a TrueType font file can be used in a PDF document.
func (b *PdfDocumentBuilder) CanUseTrueTypeFont(fontFileBytes []byte) (bool, []string) {
	reasons := []string{}

	if len(fontFileBytes) == 0 {
		reasons = append(reasons, "Provided bytes were empty.")
		return false, reasons
	}

	font, err := truetypeparser.ParseFull(truetypeparser.NewTrueTypeDataBytes(fontFileBytes))
	if err != nil {
		reasons = append(reasons, err.Error())
		return false, reasons
	}

	tr := font.TableRegister()
	if tr.CMapTable == nil {
		reasons = append(reasons, "The provided font did not contain a cmap table, used to map character codes to glyph codes.")
		return false, reasons
	}

	if tr.Os2Table.Tag() == "" {
		reasons = append(reasons, "The provided font did not contain an OS/2 table, used to fill in the font descriptor dictionary.")
		return false, reasons
	}

	if tr.PostScriptTable.Tag() == "" {
		reasons = append(reasons, "The provided font did not contain a post PostScript table, used to map character codes to glyph codes.")
		return false, reasons
	}

	return true, reasons
}

// AddTrueTypeFont adds a TrueType font to the builder so that pages can use it.
func (b *PdfDocumentBuilder) AddTrueTypeFont(fontFileBytes []byte) (*AddedFont, error) {
	font, err := truetypeparser.ParseFull(truetypeparser.NewTrueTypeDataBytes(fontFileBytes))
	if err != nil {
		return nil, fmt.Errorf("writing only supports TrueType fonts, please provide a valid TrueType font: %w", err)
	}

	id := fmt.Sprintf("font-%d", b.fontId)
	b.fontId++
	ref := b.context.ReserveObjectNumber()
	added := NewAddedFont(id, ref)
	writingFont := writerfonts.NewTrueTypeWritingFont(font, fontFileBytes)
	b.fonts[id] = &FontStored{
		FontKey:     added,
		FontProgram: writingFont,
	}
	return added, nil
}

// AddStandard14Font adds one of the Standard 14 fonts included by default in PDF programs.
func (b *PdfDocumentBuilder) AddStandard14Font(fontType standard14fonts.Standard14Font) (*AddedFont, error) {
	if b.archiveStandard != PdfANone {
		return nil, fmt.Errorf("PDF/A %d requires the font to be embedded in the file, only AddTrueTypeFont is supported", b.archiveStandard)
	}

	id := fmt.Sprintf("font-%d", b.fontId)
	b.fontId++
	ref := b.context.ReserveObjectNumber()
	added := NewAddedFont(id, ref)

	metrics := standard14fonts.GetAdobeFontMetricsByType(fontType)
	writingFont := writerfonts.NewStandard14WritingFont(metrics)
	b.fonts[id] = &FontStored{
		FontKey:     added,
		FontProgram: writingFont,
	}
	return added, nil
}

// AddImage writes an image stream token and returns its indirect reference.
func (b *PdfDocumentBuilder) AddImage(dictionary *tokens.DictionaryToken, bytes []byte) *tokens.IndirectReferenceToken {
	streamToken, _ := tokens.NewStreamToken(dictionary, bytes)
	return b.context.WriteToken(streamToken)
}

// CopyToken copies a token from the source document into this builder's stream,
// resolving all indirect references and tracking them to avoid duplication.
func (b *PdfDocumentBuilder) CopyToken(scanner tokenScannerForCopy, token tokens.Token) tokens.Token {
	scannerKey := getScannerCopyKey(b, scanner)

	refs, ok := b.existingCopies[scannerKey]
	if !ok {
		refs = make(map[core.IndirectReference]*tokens.IndirectReferenceToken)
		b.existingCopies[scannerKey] = refs
	}
	return CopyToken(b.context, token, scanner, refs, nil)
}

// getScannerCopyKey returns a unique integer key for the given scanner to track copied references.
func getScannerCopyKey(b *PdfDocumentBuilder, scanner tokenScannerForCopy) int {
	b.docCounter++
	return b.docCounter
}

// PdfStreamWriterContext returns the underlying PDF stream writer context.
func (b *PdfDocumentBuilder) PdfStreamWriterContext() PdfStreamWriter {
	return b.context
}

// GetOrCreateCopyRefs creates or retrieves a reference tracking map for the given scanner key.
func (b *PdfDocumentBuilder) GetOrCreateCopyRefs(scannerKey int) map[core.IndirectReference]*tokens.IndirectReferenceToken {
	refs, ok := b.existingCopies[scannerKey]
	if !ok {
		refs = make(map[core.IndirectReference]*tokens.IndirectReferenceToken)
		b.existingCopies[scannerKey] = refs
	}
	return refs
}

// NewScannerCopyKey generates a new unique key for tracking copied references from a source document.
func (b *PdfDocumentBuilder) NewScannerCopyKey() int {
	b.docCounter++
	return b.docCounter
}

// AddPageWithSize adds a new page with the specified width and height in points.
func (b *PdfDocumentBuilder) AddPageWithSize(width, height float64) (*PdfPageBuilder, error) {
	if width < 0 {
		return nil, fmt.Errorf("width cannot be negative, got: %f", width)
	}
	if height < 0 {
		return nil, fmt.Errorf("height cannot be negative, got: %f", height)
	}

	var builder *PdfPageBuilder
	for i := 0; i < len(b.pages); i++ {
		if _, ok := b.pages[i+1]; !ok {
			builder = NewPdfPageBuilder(i+1, b)
			break
		}
	}

	if builder == nil {
		builder = NewPdfPageBuilder(len(b.pages)+1, b)
	}

	builder.SetPageSize(content.NewMediaBox(core.NewPdfRectangleFloat(0, 0, width, height)))
	b.pages[builder.PageNumber()] = builder
	return builder, nil
}

// documentStructureAdapter wraps PdfDocument.Structure to implement tokenScannerForCopy.
type documentStructureAdapter struct {
	document *content.PdfDocument
}

func (a *documentStructureAdapter) Get(reference core.IndirectReference) *tokens.ObjectToken {
	obj, _ := a.document.Structure.GetObject(reference)
	return obj
}

// AddPageWithOptions adds a page copied from an existing document with options.
func (b *PdfDocumentBuilder) AddPageWithOptions(
	document *content.PdfDocument,
	pageNumber int,
	options *AddPageOptions,
) (*PdfPageBuilder, error) {
	if document == nil {
		return nil, fmt.Errorf("document must not be nil")
	}

	treeKey := fmt.Sprintf("%p", document)

	b.docCounter++
	scannerKey := b.docCounter

	refs, ok := b.existingCopies[scannerKey]
	if !ok {
		refs = make(map[core.IndirectReference]*tokens.IndirectReferenceToken)
		b.existingCopies[scannerKey] = refs
	}

	pagesInfos, ok := b.existingTrees[treeKey]
	if !ok {
		pagesInfos = make(map[int]*pageInfo)

		pageTreeRoot := document.Structure.Catalog().Pages().PageTree()
		walkResults := WalkTree(pageTreeRoot)

		for i, result := range walkResults {
			pagesInfos[i+1] = &pageInfo{
				page:    result.Node,
				parents: result.Parents,
			}
		}

		b.existingTrees[treeKey] = pagesInfos
	}

	srcPageInfo, ok := pagesInfos[pageNumber]
	if !ok || srcPageInfo.page == nil {
		return nil, fmt.Errorf("page %d was not found in the source document", pageNumber)
	}

	scanner := &documentStructureAdapter{document: document}

	copiedPageDict := make(map[*tokens.NameToken]tokens.Token)
	links := make([]CopiedLink, 0)
	resources := make(map[string]tokens.Token)

	for _, dict := range srcPageInfo.parents {
		if dict == nil {
			continue
		}
		if resTok, hasRes := dict.TryGet(tokens.Resources); hasRes && resTok != nil {
			b.copyResourceDict(resTok, resources, scanner, refs, document)
		}
		if mbTok, hasMb := dict.TryGet(tokens.MediaBox); hasMb && mbTok != nil {
			copiedPageDict[tokens.MediaBox] = CopyToken(b.context, mbTok, scanner, refs, nil)
		}
		if cbTok, hasCb := dict.TryGet(tokens.CropBox); hasCb && cbTok != nil {
			copiedPageDict[tokens.CropBox] = CopyToken(b.context, cbTok, scanner, refs, nil)
		}
		if rtTok, hasRt := dict.TryGet(tokens.Rotate); hasRt && rtTok != nil {
			copiedPageDict[tokens.Rotate] = CopyToken(b.context, rtTok, scanner, refs, nil)
		}
	}

	for nameStr, val := range srcPageInfo.page.Data() {
		keyName := tokens.MustCreate(nameStr)

		if keyName.Equals(tokens.Contents) || keyName.Equals(tokens.Parent) || keyName.Equals(tokens.Type) {
			continue
		}

		if keyName.Equals(tokens.Resources) {
			b.copyResourceDict(val, resources, scanner, refs, document)
			continue
		}

		if keyName.Equals(tokens.Annots) {
			if options == nil || !options.KeepAnnotations {
				continue
			}

			copiedTokens := b.copyAnnotationsFromPageSource(
				val, scanner, refs, pageNumber, document,
				func(a actions.Action) actions.Action {
					if options != nil && options.CopyLinkFunc != nil {
						return options.CopyLinkFunc(a)
					}
					return a
				},
				func(l CopiedLink) { links = append(links, l) },
			)
			copiedPageDict[tokens.Annots] = tokens.NewArrayToken(copiedTokens)
			continue
		}

		copiedPageDict[keyName] = CopyToken(b.context, val, scanner, refs, nil)
	}

	resDict, _ := tokens.WithMap(resources)
	copiedPageDict[tokens.Resources] = resDict

	var copiedStreams []PageContentStream
	if contentsTok, hasContents := srcPageInfo.page.TryGet(tokens.Contents); hasContents && contentsTok != nil {
		prev := b.context.AttemptDeduplication()
		b.context.SetAttemptDeduplication(false)
		b.context.SetWritingPageContents(true)

		contentRefs := make([]*tokens.IndirectReferenceToken, 0)
		if arr, isArray := contentsTok.(*tokens.ArrayToken); isArray {
			for _, item := range arr.Data() {
				if ir, isRef := item.(*tokens.IndirectReferenceToken); isRef {
					contentRefs = append(contentRefs, ir)
				}
			}
		} else if ir, isRef := contentsTok.(*tokens.IndirectReferenceToken); isRef {
			contentRefs = append(contentRefs, ir)
		}

		for _, ref := range contentRefs {
			globalTransform := GetGlobalTransformFromStream(ref, scanner, document)
			updatedIr := CopyToken(b.context, ref, scanner, refs, nil).(*tokens.IndirectReferenceToken)
			copiedStreams = append(copiedStreams, &copiedContentStream{reference: updatedIr, globalTransform: globalTransform})
		}

		b.context.SetAttemptDeduplication(prev)
		b.context.SetWritingPageContents(false)
	}

	if copiedStreams == nil {
		copiedStreams = make([]PageContentStream, 0)
	}

	builderNum := len(b.pages) + 1
	builder := NewPdfPageBuilderFromCopied(builderNum, b, copiedStreams, copiedPageDict, links)
	b.pages[builder.PageNumber()] = builder
	return builder, nil
}

func (b *PdfDocumentBuilder) copyResourceDict(
	token tokens.Token,
	dest map[string]tokens.Token,
	scanner tokenScannerForCopy,
	refs map[core.IndirectReference]*tokens.IndirectReferenceToken,
	document *content.PdfDocument,
) {
	dict := b.getRemoteDict(token, document)
	if dict == nil {
		return
	}

	for keyStr, val := range dict.Data() {
		if _, exists := dest[keyStr]; !exists {
			if ir, isIR := val.(*tokens.IndirectReferenceToken); isIR {
				obj, err := document.Structure.GetObject(ir.Data())
				if err != nil || obj == nil {
					dest[keyStr] = CopyToken(b.context, val, scanner, refs, nil)
					continue
				}
				if _, isStream := obj.Data().(*tokens.StreamToken); isStream {
					dest[keyStr] = CopyToken(b.context, val, scanner, refs, nil)
				} else {
					dest[keyStr] = CopyToken(b.context, obj.Data(), scanner, refs, nil)
				}
			} else {
				dest[keyStr] = CopyToken(b.context, val, scanner, refs, nil)
			}
			continue
		}

		subDict := b.getRemoteDict(val, document)
		destSubDict, ok := dest[keyStr].(*tokens.DictionaryToken)
		if !ok || destSubDict == nil || subDict == nil {
			if ir, isIR := val.(*tokens.IndirectReferenceToken); isIR {
				obj, err := document.Structure.GetObject(ir.Data())
				if err != nil || obj == nil {
					dest[keyStr] = CopyToken(b.context, val, scanner, refs, nil)
				} else {
					dest[keyStr] = CopyToken(b.context, obj.Data(), scanner, refs, nil)
				}
			} else {
				dest[keyStr] = CopyToken(b.context, val, scanner, refs, nil)
			}
			continue
		}

		mutableSubDict := make(map[string]tokens.Token)
		for k, v := range destSubDict.Data() {
			mutableSubDict[k] = v
		}

		for subKeyStr, subVal := range subDict.Data() {
			mutableSubDict[subKeyStr] = CopyToken(b.context, subVal, scanner, refs, nil)
		}

		newDict, _ := tokens.WithMap(mutableSubDict)
		dest[keyStr] = newDict
	}
}

func (b *PdfDocumentBuilder) getRemoteDict(token tokens.Token, document *content.PdfDocument) *tokens.DictionaryToken {
	if ir, isIR := token.(*tokens.IndirectReferenceToken); isIR {
		obj, err := document.Structure.GetObject(ir.Data())
		if err != nil || obj == nil {
			return nil
		}
		if dt, ok := obj.Data().(*tokens.DictionaryToken); ok {
			return dt
		}
	}
	if dt, ok := token.(*tokens.DictionaryToken); ok {
		return dt
	}
	return nil
}

func (b *PdfDocumentBuilder) copyAnnotationsFromPageSource(
	val tokens.Token,
	scanner tokenScannerForCopy,
	refs map[core.IndirectReference]*tokens.IndirectReferenceToken,
	pageNumber int,
	document *content.PdfDocument,
	linkCopyFunc func(actions.Action) actions.Action,
	deferredActionUpdate func(CopiedLink),
) []tokens.Token {
	permittedLinkActionTypes := map[*tokens.NameToken]bool{
		tokens.Uri:    true,
		tokens.GoToR:  true,
		tokens.Launch: true,
	}

	arr := b.resolveToArray(val, document)
	if arr == nil {
		return []tokens.Token{}
	}

	copiedAnnotations := make([]tokens.Token, 0)

	for _, annotEntry := range arr.Data() {
		annotDict := b.getRemoteDict(annotEntry, document)
		if annotDict == nil {
			continue
		}

		removedKeys := make([]*tokens.NameToken, 0)

		if pTok, hasP := annotDict.TryGet(tokens.P); hasP && pTok != nil {
			removedKeys = append(removedKeys, tokens.P)
		}

		if spTok, hasSp := annotDict.TryGet(tokens.StructParent); hasSp && spTok != nil {
			removedKeys = append(removedKeys, tokens.StructParent)
		}

		subtypeTok, hasSubtype := annotDict.TryGet(tokens.Subtype)
		isLink := false
		if hasSubtype {
			if nt, ok := subtypeTok.(*tokens.NameToken); ok {
				isLink = nt.Equals(tokens.Link)
			}
		}

		if !isLink {
			copiedRef := CopyToken(b.context, b.copyWithSkippedKeys(annotDict, removedKeys), scanner, refs, nil)
			copiedAnnotations = append(copiedAnnotations, copiedRef)
			continue
		}

		docScanner := document.Structure.Scanner()
		pages := document.Structure.Catalog().Pages()

		// Match C# AnnotationProvider.GetAction: first check /Dest, then /A
		action := getAnnotationAction(annotDict, docScanner, pages)

		if action != nil && linkCopyFunc != nil && deferredActionUpdate != nil {
			copiedLink := linkCopyFunc(action)
			if copiedLink != action && copiedLink != nil {
				copiedToken := CopyToken(b.context, annotDict, scanner, refs, nil)
				if dictTok, ok := copiedToken.(*tokens.DictionaryToken); ok {
					deferredActionUpdate(CopiedLink{token: dictTok, action: copiedLink})
				}
				continue
			}
		}

		hasA := false
		actionTypePermitted := true
		if aTok, hasATok := annotDict.TryGet(tokens.A); hasATok && aTok != nil {
			hasA = true
			actionDict := b.getRemoteDict(aTok, document)
			if actionDict != nil {
				sTok, hasS := actionDict.TryGet(tokens.S)
				if !hasS || sTok == nil {
					continue
				}
				if nameTok, ok := sTok.(*tokens.NameToken); ok {
					actionTypePermitted = permittedLinkActionTypes[nameTok]
				} else {
					actionTypePermitted = false
				}
			}
		}

		if hasA && !actionTypePermitted {
			continue
		}

		if destTok, hasDest := annotDict.TryGet(tokens.Dest); hasDest && destTok != nil {
			continue
		}

		finalCopiedRef := CopyToken(b.context, b.copyWithSkippedKeys(annotDict, removedKeys), scanner, refs, nil)
		copiedAnnotations = append(copiedAnnotations, finalCopiedRef)
	}

	return copiedAnnotations
}

// getAnnotationAction replicates C# AnnotationProvider.GetAction: first checks /Dest, then /A.
func getAnnotationAction(annotDict *tokens.DictionaryToken, scanner tokenization.PdfTokenScanner, pages destinations.Pages) actions.Action {
	namedDests := destinations.NewNamedDestinations(nil, pages)

	destProvider := destinations.DestinationProvider{}
	if dest, ok := destProvider.TryGetDestination(annotDict, tokens.Dest, namedDests, scanner, nil, false); ok {
		return actions.NewGoToAction(dest)
	}

	actionProvider := actions.ActionProvider{}
	action, hasAction, _ := actionProvider.TryGetAction(annotDict, namedDests, scanner, nil)
	if hasAction && action != nil {
		return action
	}

	return nil
}

// unwrapIndirectRef resolves an IndirectReferenceToken through the scanner.
func unwrapIndirectRef(token tokens.Token, scanner tokenization.PdfTokenScanner) tokens.Token {
	refToken, ok := token.(*tokens.IndirectReferenceToken)
	if !ok || scanner == nil {
		return token
	}
	obj := scanner.Get(refToken.Data())
	if obj == nil {
		return token
	}
	return obj.Data()
}

func (b *PdfDocumentBuilder) resolveToArray(token tokens.Token, document *content.PdfDocument) *tokens.ArrayToken {
	if ir, isIR := token.(*tokens.IndirectReferenceToken); isIR {
		obj, err := document.Structure.GetObject(ir.Data())
		if err != nil || obj == nil {
			return nil
		}
		if arr, ok := obj.Data().(*tokens.ArrayToken); ok {
			return arr
		}
	}
	if arr, ok := token.(*tokens.ArrayToken); ok {
		return arr
	}
	return nil
}

func (b *PdfDocumentBuilder) copyWithSkippedKeys(source *tokens.DictionaryToken, skipped []*tokens.NameToken) *tokens.DictionaryToken {
	dict := make(map[string]tokens.Token)
	for keyStr, val := range source.Data() {
		name := tokens.MustCreate(keyStr)
		ignore := false
		for _, skipName := range skipped {
			if skipName.Equals(name) {
				ignore = true
				break
			}
		}
		if !ignore {
			dict[keyStr] = val
		}
	}
	result, _ := tokens.WithMap(dict)
	return result
}

// Build completes the document and returns the PDF bytes.
func (b *PdfDocumentBuilder) Build() ([]byte, error) {
	if err := b.completeDocument(); err != nil {
		return nil, fmt.Errorf("completing document: %w", err)
	}

	stream := b.context.Stream()
	if buf, ok := stream.(*bytes.Buffer); ok {
		return buf.Bytes(), nil
	}

	if ptw, ok := stream.(*positionTrackingWriter); ok {
		if buf, ok := ptw.writer.(*bytes.Buffer); ok {
			return buf.Bytes(), nil
		}
	}

	if seeker, ok := stream.(interface {
		io.Reader
		io.Seeker
	}); ok {
		seeker.Seek(0, io.SeekStart)
		var result bytes.Buffer
		io.Copy(&result, seeker)
		return result.Bytes(), nil
	}

	return nil, fmt.Errorf("build with external non-seekable stream")
}

// Close completes the document if not already done and disposes the underlying stream.
func (b *PdfDocumentBuilder) Close() error {
	if !b.completed {
		if err := b.completeDocument(); err != nil {
			_ = b.context.Close()
			return fmt.Errorf("completing document: %w", err)
		}
	}
	return b.context.Close()
}

func (b *PdfDocumentBuilder) completeDocument() error {
	ds := b.context.(*defaultPdfStreamWriter)
	ds.stream.DebugPosition("completeDocument start")
	for _, font := range b.fonts {
		font.FontProgram.WriteFont(b.context, font.FontKey.Reference)
	}
	ds.stream.DebugPosition("after fonts")

	const desiredLeafSize = 25
	numLeaves := int(math.Ceil(float64(len(b.pages)) / float64(desiredLeafSize)))

	leafRefs := make([]*tokens.IndirectReferenceToken, numLeaves)
	leafChildren := make([][]*tokens.IndirectReferenceToken, numLeaves)
	leaves := make([]map[*tokens.NameToken]tokens.Token, numLeaves)

	for i := 0; i < numLeaves; i++ {
		leaves[i] = map[*tokens.NameToken]tokens.Token{
			tokens.Type: tokens.Pages,
		}
		leafChildren[i] = make([]*tokens.IndirectReferenceToken, 0)
		leafRefs[i] = b.context.ReserveObjectNumber()
	}

	leafNum := 0
	pageReferences := make(map[int]*tokens.IndirectReferenceToken)
	for pageNum := range b.pages {
		pageReferences[pageNum] = b.context.ReserveObjectNumber()
	}

	sortedPageNumbers := make([]int, 0, len(b.pages))
	for pn := range b.pages {
		sortedPageNumbers = append(sortedPageNumbers, pn)
	}
	for i := 0; i < len(sortedPageNumbers); i++ {
		for j := i + 1; j < len(sortedPageNumbers); j++ {
			if sortedPageNumbers[i] > sortedPageNumbers[j] {
				sortedPageNumbers[i], sortedPageNumbers[j] = sortedPageNumbers[j], sortedPageNumbers[i]
			}
		}
	}

	for _, pageNum := range sortedPageNumbers {
		page := b.pages[pageNum]
		pageDict := page.PageDictionary()
		pageDict[tokens.Type] = tokens.Page
		pageDict[tokens.Parent] = leafRefs[leafNum]
		pageDict[tokens.ProcSet] = DefaultProcSet

		if _, ok := pageDict[tokens.MediaBox]; !ok {
			pageDict[tokens.MediaBox] = rectangleToArray(page.PageSize())
		}

		if rotation := page.Rotation(); rotation != nil {
			pageDict[tokens.Rotate] = tokens.NewNumericTokenFromInt(*rotation)
		}

		prev := b.context.AttemptDeduplication()
		b.context.SetAttemptDeduplication(false)

		contentStreams := page.ContentStreams()
		hasContent := make([]PageContentStream, 0)
		for _, cs := range contentStreams {
			if cs.HasContent() {
				hasContent = append(hasContent, cs)
			}
		}

		if len(hasContent) == 0 {
			defaultStream := &defaultContentStream{operations: make([]content.GraphicsStateOperation, 0)}
			pageDict[tokens.Contents] = defaultStream.Write(b.context)
		} else if len(hasContent) == 1 {
			pageDict[tokens.Contents] = hasContent[0].Write(b.context)
		} else {
			var streamTokens []tokens.Token
			for _, cs := range hasContent {
				streamTokens = append(streamTokens, cs.Write(b.context))
			}
			pageDict[tokens.Contents] = tokens.NewArrayToken(streamTokens)
		}
		b.context.SetAttemptDeduplication(prev)

		if links := page.Links(); len(links) > 0 {
			var annots []tokens.Token
			if existingAnnots, ok := pageDict[tokens.Annots]; ok {
				if arr, ok := existingAnnots.(*tokens.ArrayToken); ok {
					for _, item := range arr.Data() {
						annots = append(annots, item)
					}
				}
			}

			for _, link := range links {
				linkToken, err := createLinkAnnotationToken(link.Token(), link.Action(), pageReferences)
				if err != nil {
					return fmt.Errorf("creating link annotation: %w", err)
				}
				annots = append(annots, linkToken)
			}
			pageDict[tokens.Annots] = tokens.NewArrayToken(annots)
		}

		dictToken, _ := tokens.NewDictionary(pageDict)
		leafChildren[leafNum] = append(leafChildren[leafNum],
			b.context.WriteTokenAt(dictToken, pageReferences[page.PageNumber()]))

		if len(leafChildren[leafNum]) >= desiredLeafSize {
			leafNum++
		}
	}

	dummyName := tokens.MustCreate("#/ObjId")
	for i := 0; i < numLeaves; i++ {
		leaves[i][tokens.Kids] = tokens.NewArrayToken(tokenSliceFromRefs(leafChildren[i]))
		leaves[i][tokens.Count] = tokens.NewNumericTokenFromInt(len(leafChildren[i]))
		leaves[i][dummyName] = leafRefs[i]
	}

	catalogDict := map[*tokens.NameToken]tokens.Token{
		tokens.Type: tokens.Catalog,
	}

	if numLeaves == 1 {
		leaf := leaves[0]
		id := leaf[dummyName].(*tokens.IndirectReferenceToken)
		delete(leaf, dummyName)
		leafDict, _ := tokens.NewDictionary(leaf)
		catalogDict[tokens.Pages] = b.context.WriteTokenAt(leafDict, id)
	} else {
		rootPageInfo := b.createPageTree(leaves, nil, leafRefs, dummyName)
		catalogDict[tokens.Pages] = rootPageInfo.ref
	}

	if b.bookmarks != nil && len(b.bookmarks.Roots()) > 0 {
		bookmarks, err := b.createBookmarkTree(b.bookmarks.Roots(), pageReferences, nil)
		if err != nil {
			return fmt.Errorf("creating bookmark tree: %w", err)
		}
		outlineDict := map[*tokens.NameToken]tokens.Token{
			tokens.Type:  tokens.Outlines,
			tokens.Count: tokens.NewNumericTokenFromInt(len(b.bookmarks.Roots())),
			tokens.First: bookmarks[0],
			tokens.Last:  bookmarks[len(bookmarks)-1],
		}
		outlineToken, _ := tokens.NewDictionary(outlineDict)
		catalogDict[tokens.Outlines] = b.context.WriteToken(outlineToken)
	}

	if b.archiveStandard != PdfANone {
		objectWriter := func(t tokens.Token) *tokens.IndirectReferenceToken {
			return b.context.WriteToken(t)
		}

		pdfABaselineRuleBuilder{}.Obey(catalogDict, objectWriter,
			b.documentInformation.ToContentDocumentInformation(), b.archiveStandard, b.version, b.xmpMetadata)

		switch b.archiveStandard {
		case PdfA1A, PdfA2A, PdfA3A:
			Obey(catalogDict)
		}
	}

	catalogToken, _ := tokens.NewDictionary(catalogDict)
	catalogRef := b.context.WriteToken(catalogToken)

	var infoRef *tokens.IndirectReferenceToken
	if b.includeDocInfo && b.documentInformation != nil {
		infoDict := b.documentInformation.ToDictionaryTokens()
		if len(infoDict) > 0 {
			infoToken, _ := tokens.NewDictionary(infoDict)
			infoRef = b.context.WriteToken(infoToken)
		}
	}

	b.context.CompletePdf(catalogRef, infoRef)
	b.completed = true
	return nil
}

type pageTreeNode struct {
	count int
	ref   *tokens.IndirectReferenceToken
}

func (b *PdfDocumentBuilder) createPageTree(
	pagesNodes []map[*tokens.NameToken]tokens.Token,
	parent *tokens.IndirectReferenceToken,
	leafRefs []*tokens.IndirectReferenceToken,
	dummyName *tokens.NameToken,
) *pageTreeNode {
	const desiredLeafSize = 25

	count := 0
	thisObj := b.context.ReserveObjectNumber()
	children := make([]*tokens.IndirectReferenceToken, 0)

	if len(pagesNodes) > desiredLeafSize {
		currentTreeDepth := int(math.Ceil(math.Log2(float64(len(pagesNodes))) / math.Log2(float64(desiredLeafSize))))
		if currentTreeDepth < 1 {
			currentTreeDepth = 1
		}
		perBranch := int(math.Pow(float64(desiredLeafSize), float64(currentTreeDepth-1)))
		if perBranch < 1 {
			perBranch = 1
		}
		branches := int(math.Ceil(float64(len(pagesNodes)) / float64(perBranch)))
		for i := 0; i < branches; i++ {
			start := i * perBranch
			end := start + perBranch
			if end > len(pagesNodes) {
				end = len(pagesNodes)
			}
			part := pagesNodes[start:end]
			result := b.createPageTree(part, thisObj, leafRefs, dummyName)
			count += result.count
			children = append(children, result.ref)
		}
	} else {
		for _, page := range pagesNodes {
			page[tokens.Parent] = thisObj
			id := page[dummyName].(*tokens.IndirectReferenceToken)
			delete(page, dummyName)
			if cntTok, ok := page[tokens.Count]; ok {
				if nt, ok := cntTok.(*tokens.NumericToken); ok {
					count += int(nt.Data())
				}
			}
			pageDict, _ := tokens.NewDictionary(page)
			children = append(children, b.context.WriteTokenAt(pageDict, id))
		}
	}

	node := map[*tokens.NameToken]tokens.Token{
		tokens.Type:  tokens.Pages,
		tokens.Kids:  tokens.NewArrayToken(tokenSliceFromRefs(children)),
		tokens.Count: tokens.NewNumericTokenFromInt(count),
	}
	if parent != nil {
		node[tokens.Parent] = parent
	}
	nodeDict, _ := tokens.NewDictionary(node)
	return &pageTreeNode{count: count, ref: b.context.WriteTokenAt(nodeDict, thisObj)}
}

func (b *PdfDocumentBuilder) createBookmarkTree(
	nodes []any,
	pageReferences map[int]*tokens.IndirectReferenceToken,
	parent *tokens.IndirectReferenceToken,
) ([]*tokens.IndirectReferenceToken, error) {
	childRefs := make([]*tokens.IndirectReferenceToken, len(nodes))
	for i := 0; i < len(nodes); i++ {
		childRefs[i] = b.context.ReserveObjectNumber()
	}

	for i := 0; i < len(nodes); i++ {
		node := nodes[i]
		objNum := childRefs[i]

		title := ""
		children := []any(nil)
		var docDest destinations.ExplicitDestination
		var uriVal string
		isDocBookmark := false
		isUriBookmark := false

		switch n := node.(type) {
		case *outline.DocumentBookmarkNode:
			title = n.Title
			children = extractBookmarkChildren(n.Children())
			docDest = n.Destination
			isDocBookmark = true
		case *outline.UriBookmarkNode:
			title = n.Title
			children = extractBookmarkChildren(n.Children())
			uriVal = n.Uri
			isUriBookmark = true
		case *outline.BookmarkNode:
			title = n.Title
			children = extractBookmarkChildren(n.Children())
		default:
			return nil, fmt.Errorf("%T is not a supported bookmark node type", node)
		}

		data := map[*tokens.NameToken]tokens.Token{
			tokens.Title: tokens.NewStringToken(title),
			tokens.Count: tokens.NewNumericTokenFromInt(len(children)),
		}

		if parent != nil {
			data[tokens.Parent] = parent
		}
		if i > 0 {
			data[tokens.Prev] = childRefs[i-1]
		}
		if i < len(childRefs)-1 {
			data[tokens.Next] = childRefs[i+1]
		}

		if len(children) > 0 {
			childrenRefs, err := b.createBookmarkTree(children, pageReferences, objNum)
			if err != nil {
				return nil, fmt.Errorf("creating child bookmark tree: %w", err)
			}
			data[tokens.First] = childrenRefs[0]
			data[tokens.Last] = childrenRefs[len(childrenRefs)-1]
		}

		if isDocBookmark {
			dest, err := createExplicitDestinationToken(docDest, pageReferences)
			if err != nil {
				return nil, fmt.Errorf("creating destination for bookmark: %w", err)
			}
			data[tokens.Dest] = dest
		}

		if isUriBookmark {
			actionDict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
				tokens.S:   tokens.Uri,
				tokens.Uri: tokens.NewStringToken(uriVal),
			})
			data[tokens.A] = actionDict
		}

		dataToken, _ := tokens.NewDictionary(data)
		b.context.WriteTokenAt(dataToken, objNum)
	}

	return childRefs, nil
}

func extractBookmarkChildren(children []*outline.BookmarkNode) []any {
	result := make([]any, len(children))
	for i, c := range children {
		result[i] = c
	}
	return result
}

func createExplicitDestinationToken(
	dest destinations.ExplicitDestination,
	pageReferences map[int]*tokens.IndirectReferenceToken,
) (*tokens.ArrayToken, error) {
	page, ok := pageReferences[dest.PageNumber]
	if !ok {
		return nil, fmt.Errorf("page %d was not found in the source document", dest.PageNumber)
	}

	switch dest.Type {
	case destinations.XyzCoordinates:
		left := float64(0)
		top := float64(0)
		if dest.Coordinates != nil && dest.Coordinates.Left != nil {
			left = *dest.Coordinates.Left
		}
		if dest.Coordinates != nil && dest.Coordinates.Top != nil {
			top = *dest.Coordinates.Top
		}
		return tokens.NewArrayToken([]tokens.Token{
			page,
			tokens.XYZ,
			tokens.NewNumericToken(left),
			tokens.NewNumericToken(top),
			tokens.NewNumericToken(0),
		}), nil

	case destinations.FitPage:
		return tokens.NewArrayToken([]tokens.Token{page, tokens.Fit}), nil

	case destinations.FitHorizontally:
		top := float64(0)
		if dest.Coordinates != nil && dest.Coordinates.Top != nil {
			top = *dest.Coordinates.Top
		}
		return tokens.NewArrayToken([]tokens.Token{page, tokens.FitH, tokens.NewNumericToken(top)}), nil

	case destinations.FitVertically:
		left := float64(0)
		if dest.Coordinates != nil && dest.Coordinates.Left != nil {
			left = *dest.Coordinates.Left
		}
		return tokens.NewArrayToken([]tokens.Token{page, tokens.FitV, tokens.NewNumericToken(left)}), nil

	case destinations.FitRectangle:
		left, top, right, bottom := float64(0), float64(0), float64(0), float64(0)
		if dest.Coordinates != nil {
			if dest.Coordinates.Left != nil {
				left = *dest.Coordinates.Left
			}
			if dest.Coordinates.Top != nil {
				top = *dest.Coordinates.Top
			}
			if dest.Coordinates.Right != nil {
				right = *dest.Coordinates.Right
			}
			if dest.Coordinates.Bottom != nil {
				bottom = *dest.Coordinates.Bottom
			}
		}
		return tokens.NewArrayToken([]tokens.Token{
			page, tokens.FitR,
			tokens.NewNumericToken(left), tokens.NewNumericToken(top),
			tokens.NewNumericToken(right), tokens.NewNumericToken(bottom),
		}), nil

	case destinations.FitBoundingBox:
		return tokens.NewArrayToken([]tokens.Token{page, tokens.FitB}), nil

	case destinations.FitBoundingBoxHorizontally:
		left := float64(0)
		if dest.Coordinates != nil && dest.Coordinates.Left != nil {
			left = *dest.Coordinates.Left
		}
		return tokens.NewArrayToken([]tokens.Token{page, tokens.FitBH, tokens.NewNumericToken(left)}), nil

	case destinations.FitBoundingBoxVertically:
		left := float64(0)
		if dest.Coordinates != nil && dest.Coordinates.Left != nil {
			left = *dest.Coordinates.Left
		}
		return tokens.NewArrayToken([]tokens.Token{page, tokens.FitBV, tokens.NewNumericToken(left)}), nil

	default:
		return nil, fmt.Errorf("%v is not a supported bookmark destination type", dest.Type)
	}
}

func createLinkAnnotationToken(
	token *tokens.DictionaryToken,
	action any,
	pageReferences map[int]*tokens.IndirectReferenceToken,
) (tokens.Token, error) {
	data := make(map[*tokens.NameToken]tokens.Token)

	for nameStr, value := range token.Data() {
		nameTok := tokens.MustCreate(nameStr)
		if nameTok.Equals(tokens.A) || nameTok.Equals(tokens.Dest) {
			continue
		}
		data[nameTok] = value
	}

	actionToken, err := createActionToken(action, pageReferences)
	if err != nil {
		return nil, fmt.Errorf("creating action token: %w", err)
	}
	data[tokens.A] = actionToken
	result, _ := tokens.NewDictionary(data)
	return result, nil
}

func createActionToken(
	action any,
	pageReferences map[int]*tokens.IndirectReferenceToken,
) (tokens.Token, error) {
	switch a := action.(type) {
	case *actions.UriAction:
		dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
			tokens.S:   tokens.Uri,
			tokens.Uri: tokens.NewStringToken(a.Uri),
		})
		return dict, nil

	case *actions.GoToAction:
		d, err := createExplicitDestinationToken(a.Destination, pageReferences)
		if err != nil {
			return nil, fmt.Errorf("creating destination for GoTo action: %w", err)
		}
		dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
			tokens.S: tokens.GoTo,
			tokens.D: d,
		})
		return dict, nil

	case *actions.GoToEAction:
		d, err := createExplicitDestinationToken(a.Destination, pageReferences)
		if err != nil {
			return nil, fmt.Errorf("creating destination for GoToE action: %w", err)
		}
		dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
			tokens.S: tokens.GoToE,
			tokens.F: tokens.NewStringToken(a.FileSpecification),
			tokens.D: d,
		})
		return dict, nil

	case *actions.GoToRAction:
		d, err := createExplicitDestinationToken(a.Destination, pageReferences)
		if err != nil {
			return nil, fmt.Errorf("creating destination for GoToR action: %w", err)
		}
		dict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
			tokens.S: tokens.GoToR,
			tokens.F: tokens.NewStringToken(a.Filename),
			tokens.D: d,
		})
		return dict, nil

	default:
		return nil, fmt.Errorf("%T is not a supported PDF action type", action)
	}
}

func rectangleToArray(rect *content.MediaBox) tokens.Token {
	return tokens.NewArrayToken([]tokens.Token{
		tokens.NewNumericToken(rect.Bounds.BottomLeft.X),
		tokens.NewNumericToken(rect.Bounds.BottomLeft.Y),
		tokens.NewNumericToken(rect.Bounds.TopRight.X),
		tokens.NewNumericToken(rect.Bounds.TopRight.Y),
	})
}

func tokenSliceFromRefs(refs []*tokens.IndirectReferenceToken) []tokens.Token {
	result := make([]tokens.Token, len(refs))
	for i, ref := range refs {
		result[i] = ref
	}
	return result
}

// GetGlobalTransformFromStream tries to parse the content stream and extract global transform.
// Errors during parsing are silently ignored, matching C# behavior where failures are caught.
func GetGlobalTransformFromStream(
	ref *tokens.IndirectReferenceToken,
	scanner tokenScannerForCopy,
	document *content.PdfDocument,
) *core.TransformationMatrix {
	obj, err := document.Structure.GetObject(ref.Data())
	if err != nil || obj == nil {
		return nil
	}

	streamTok, ok := obj.Data().(*tokens.StreamToken)
	if !ok {
		return nil
	}

	// Decode the stream bytes using the default filter provider.
	contentBytes, decodeErr := builderDecodeStream(streamTok)
	if decodeErr != nil {
		return nil
	}

	inputBytes := core.NewMemoryInputBytes(contentBytes)
	pageContentParser := parser.NewPageContentParser(
		graphics.GetReflectionFactory(),
		core.Infinite,
		true,
	)

	ops := pageContentParser.Parse(0, inputBytes, logging.NoopLog)
	return GetGlobalTransform(ops)
}

// builderDecodeStream decodes stream bytes through the filter chain using DefaultFilterProvider.
func builderDecodeStream(stream *tokens.StreamToken) ([]byte, error) {
	data := stream.Data()

	filterName, hasFilter := stream.StreamDictionary.TryGet(tokens.Filter)
	if !hasFilter {
		return data, nil
	}

	var filterNames []*tokens.NameToken
	switch ft := filterName.(type) {
	case *tokens.NameToken:
		filterNames = append(filterNames, ft)
	case *tokens.ArrayToken:
		for _, item := range ft.Data() {
			if nt, ok := item.(*tokens.NameToken); ok {
				filterNames = append(filterNames, nt)
			}
		}
	default:
		return data, nil
	}

	filterInstances, err := filters.Instance.GetNamedFilters(filterNames)
	if err != nil {
		return nil, fmt.Errorf("getting named filters: %w", err)
	}

	for i, f := range filterInstances {
		if !f.IsSupported() {
			continue
		}
		decoded, decodeErr := f.Decode(data, stream.StreamDictionary, filters.Instance, i)
		if decodeErr != nil {
			return nil, fmt.Errorf("decoding with filter %T at index %d: %w", f, i, decodeErr)
		}
		data = decoded
	}

	return data, nil
}

// FontStored holds a font key and its writing program for the document builder.
type FontStored struct {
	FontKey     *AddedFont
	FontProgram WritingFont
}

// AddedFont is a key representing a font available to use on the current document builder.
type AddedFont struct {
	Id        string
	Reference *tokens.IndirectReferenceToken
}

// NewAddedFont creates a new AddedFont instance.
func NewAddedFont(id string, reference *tokens.IndirectReferenceToken) *AddedFont {
	return &AddedFont{
		Id:        id,
		Reference: reference,
	}
}

// AddPageOptions controls how a page is copied when using AddPage with an existing document.
type AddPageOptions struct {
	KeepAnnotations bool
	CopyLinkFunc    func(actions.Action) actions.Action
}

// NewAddPageOptions creates default AddPageOptions with annotations preserved.
func NewAddPageOptions() *AddPageOptions {
	return &AddPageOptions{
		KeepAnnotations: true,
	}
}

// CopiedLink holds a copied annotation token and its associated action.
type CopiedLink struct {
	token  *tokens.DictionaryToken
	action any // stores full concrete action type (UriAction, GoToAction, etc.)
}

// Token returns the dictionary token for this link.
func (cl *CopiedLink) Token() *tokens.DictionaryToken {
	return cl.token
}

// Action returns the PDF action associated with this link.
func (cl *CopiedLink) Action() any {
	return cl.action
}

// DocumentInformationBuilder sets metadata values for the document being created.
type DocumentInformationBuilder struct {
	CustomMetadata map[string]string
	Title          *string
	Author         *string
	Subject        *string
	Keywords       *string
	Creator        *string
	Producer       string
	CreationDate   *string
	ModifiedDate   *string
}

// NewDocumentInformationBuilder creates a new DocumentInformationBuilder with default values.
func NewDocumentInformationBuilder() *DocumentInformationBuilder {
	return &DocumentInformationBuilder{
		CustomMetadata: make(map[string]string),
		Producer:       "PdfPig",
	}
}

// ToDictionaryTokens converts this builder to a map of NameToken-keyed token entries.
func (b *DocumentInformationBuilder) ToDictionaryTokens() map[*tokens.NameToken]tokens.Token {
	result := make(map[*tokens.NameToken]tokens.Token)

	for k, v := range b.CustomMetadata {
		if k != "" && v != "" {
			result[tokens.MustCreate(k)] = tokens.NewStringToken(v)
		}
	}

	if b.Title != nil && *b.Title != "" {
		result[tokens.Title] = tokens.NewStringToken(*b.Title)
	}
	if b.Author != nil && *b.Author != "" {
		result[tokens.Author] = tokens.NewStringToken(*b.Author)
	}
	if b.Subject != nil && *b.Subject != "" {
		result[tokens.Subject] = tokens.NewStringToken(*b.Subject)
	}
	if b.Keywords != nil && *b.Keywords != "" {
		result[tokens.Keywords] = tokens.NewStringToken(*b.Keywords)
	}
	if b.Creator != nil && *b.Creator != "" {
		result[tokens.Creator] = tokens.NewStringToken(*b.Creator)
	}
	if b.Producer != "" {
		result[tokens.Producer] = tokens.NewStringToken(b.Producer)
	}
	if b.CreationDate != nil && *b.CreationDate != "" {
		result[tokens.CreationDate] = tokens.NewStringToken(*b.CreationDate)
	}
	if b.ModifiedDate != nil && *b.ModifiedDate != "" {
		result[tokens.ModDate] = tokens.NewStringToken(*b.ModifiedDate)
	}

	return result
}

// ToContentDocumentInformation converts this builder to a content.DocumentInformation.
func (b *DocumentInformationBuilder) ToContentDocumentInformation() *content.DocumentInformation {
	if b == nil {
		return nil
	}
	info := &content.DocumentInformation{}
	if b.Title != nil {
		info.Title = *b.Title
	}
	if b.Author != nil {
		info.Author = *b.Author
	}
	if b.Subject != nil {
		info.Subject = *b.Subject
	}
	if b.Keywords != nil {
		info.Keywords = *b.Keywords
	}
	if b.Creator != nil {
		info.Creator = *b.Creator
	}
	info.Producer = b.Producer
	if b.CreationDate != nil {
		info.CreationDate = *b.CreationDate
	}
	if b.ModifiedDate != nil {
		info.ModifiedDate = *b.ModifiedDate
	}
	return info
}
