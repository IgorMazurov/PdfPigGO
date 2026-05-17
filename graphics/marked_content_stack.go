package graphics

import (
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/xobjects"
)

// MarkedContentStack handles building MarkedContentElement instances during content stream processing.
type MarkedContentStack struct {
	builderStack []*markedContentBuilder
	number       int
	top          *markedContentBuilder
}

// NewMarkedContentStack creates a new empty MarkedContentStack.
func NewMarkedContentStack() *MarkedContentStack {
	return &MarkedContentStack{
		builderStack: make([]*markedContentBuilder, 0),
		number:       -1,
	}
}

// CanPop reports whether there is a marked content element ready to be popped.
func (m *MarkedContentStack) CanPop() bool {
	return m.top != nil
}

// Push starts a new marked content section with the given name and properties.
func (m *MarkedContentStack) Push(name *tokens.NameToken, properties *tokens.DictionaryToken) {
	if len(m.builderStack) == 0 {
		m.number++
	}

	props := properties
	if props == nil {
		props, _ = tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
	}

	m.top = &markedContentBuilder{
		number:     m.number,
		name:       name,
		properties: props,
		letters:    make([]content.Letter, 0),
		paths:      make([]geometry.PdfPath, 0),
		images:     make([]content.PdfImage, 0),
		children:   make([]content.MarkedContentElement, 0),
	}
	m.builderStack = append(m.builderStack, m.top)
}

// Pop ends the current marked content section and builds the element.
// Returns nil for child elements (they are added to parent's children list).
func (m *MarkedContentStack) Pop(pdfScanner tokenization.PdfTokenScanner) *content.MarkedContentElement {
	if len(m.builderStack) == 0 {
		return nil
	}

	idx := len(m.builderStack) - 1
	builder := m.builderStack[idx]
	m.builderStack = m.builderStack[:idx]

	result := builder.Build(pdfScanner)

	if len(m.builderStack) > 0 {
		m.top = m.builderStack[len(m.builderStack)-1]
		m.top.children = append(m.top.children, *result)
		return nil
	}

	m.top = nil
	return result
}

// AddLetter adds a letter to the current marked content section.
func (m *MarkedContentStack) AddLetter(letter *content.Letter) {
	if m.top != nil {
		m.top.letters = append(m.top.letters, *letter)
	}
}

// AddPath adds a path to the current marked content section.
func (m *MarkedContentStack) AddPath(path *geometry.PdfPath) {
	if m.top != nil {
		m.top.paths = append(m.top.paths, *path)
	}
}

// AddImage adds an image to the current marked content section.
func (m *MarkedContentStack) AddImage(image content.PdfImage) {
	if m.top != nil {
		m.top.images = append(m.top.images, image)
	}
}

// AddXObject adds an XObject image to the current marked content section if applicable.
// Only Image-type XObjects are processed; Form and PostScript types are ignored.
func (m *MarkedContentStack) AddXObject(
	xObject *xobjects.XObjectContentRecord,
	scanner tokenization.PdfTokenScanner,
	filterProvider filters.FilterProvider,
	resourceStore content.ResourceStore,
) {
	if m.top == nil || xObject == nil {
		return
	}

	if xObject.Type != xobjects.Image {
		return
	}

	img, err := content.ReadImage(xObject, scanner, filterProvider, resourceStore)
	if err != nil {
		return
	}

	m.top.images = append(m.top.images, img)
}

// markedContentBuilder collects data for a single marked content element.
type markedContentBuilder struct {
	number     int
	name       *tokens.NameToken
	properties *tokens.DictionaryToken

	letters  []content.Letter
	paths    []geometry.PdfPath
	images   []content.PdfImage
	children []content.MarkedContentElement
}

// Build creates a MarkedContentElement from the collected data.
// For Artifact-tagged elements, it produces an ArtifactMarkedContentElement with
// parsed artifactType, subType, attributeOwners, boundingBox, and attached fields.
func (b *markedContentBuilder) Build(pdfScanner tokenization.PdfTokenScanner) *content.MarkedContentElement {
	mcid := b.extractMcid()

	language := getOptionalString(b.properties, tokens.Lang)
	actualText := getOptionalString(b.properties, tokens.ActualText)
	alternateDescription := getOptionalString(b.properties, tokens.Alternate)
	expandedForm := getOptionalString(b.properties, tokens.E)

	if b.name == nil || b.name.Data() != tokens.Artifact.Data() {
		return b.buildRegular(mcid, language, actualText, alternateDescription, expandedForm)
	}

	return b.buildArtifact(mcid, language, actualText, alternateDescription, expandedForm, pdfScanner)
}

func (b *markedContentBuilder) extractMcid() int {
	if b.properties == nil {
		return -1
	}

	token, ok := b.properties.TryGet(tokens.Mcid)
	if !ok {
		return -1
	}

	if num, ok := token.(*tokens.NumericToken); ok {
		return num.IntVal()
	}

	return -1
}

func (b *markedContentBuilder) buildRegular(
	mcid int,
	language, actualText, alternateDescription, expandedForm *string,
) *content.MarkedContentElement {
	elem, err := content.NewMarkedContentElement(
		mcid,
		b.name,
		b.properties,
		language,
		actualText,
		alternateDescription,
		expandedForm,
		false,
		b.children,
		b.letters,
		b.pathsToInterface(),
		b.images,
		b.number,
	)
	if err != nil {
		return b.fallbackElement(mcid)
	}
	return elem
}

func (b *markedContentBuilder) buildArtifact(
	mcid int,
	language, actualText, alternateDescription, expandedForm *string,
	pdfScanner tokenization.PdfTokenScanner,
) *content.MarkedContentElement {
	_ = pdfScanner

	artifactType := b.parseArtifactType()
	subType := getOptionalString(b.properties, tokens.Subtype)
	boundingBox := b.parseBoundingBox()
	attached := b.parseAttached()

	isTop := false
	isBottom := false
	isLeft := false
	isRight := false
	for _, name := range attached {
		switch name.Data() {
		case "Top":
			isTop = true
		case "Bottom":
			isBottom = true
		case "Left":
			isLeft = true
		case "Right":
			isRight = true
		}
	}

	elem, err := content.NewMarkedContentElement(
		mcid,
		b.name,
		b.properties,
		language,
		actualText,
		alternateDescription,
		expandedForm,
		true,
		b.children,
		b.letters,
		b.pathsToInterface(),
		b.images,
		b.number,
	)
	if err != nil {
		return b.fallbackElement(mcid)
	}

	elem.ArtifactType = artifactType
	elem.SubType = subType
	elem.BoundingBox = boundingBox
	elem.IsTopAttached = isTop
	elem.IsBottomAttached = isBottom
	elem.IsLeftAttached = isLeft
	elem.IsRightAttached = isRight

	return elem
}

func (b *markedContentBuilder) parseArtifactType() content.ArtifactType {
	if b.properties == nil {
		return content.Unknown
	}

	token, ok := b.properties.TryGet(tokens.Type)
	if !ok {
		return content.Unknown
	}

	var data string
	switch t := token.(type) {
	case *tokens.StringToken:
		data = t.Data()
	case *tokens.NameToken:
		data = t.Data()
	default:
		return content.Unknown
	}

	return parseArtifactTypeEnum(data)
}

func parseArtifactTypeEnum(s string) content.ArtifactType {
	switch strings.ToLower(s) {
	case "pagination":
		return content.Pagination
	case "layout":
		return content.Layout
	case "page":
		return content.PageArtifact
	case "background":
		return content.Background
	default:
		return content.Unknown
	}
}

func (b *markedContentBuilder) parseBoundingBox() *core.PdfRectangle {
	if b.properties == nil {
		return nil
	}

	arrayToken, ok := tokens.TryGetTyped[*tokens.ArrayToken](b.properties, tokens.Bbox)
	if !ok {
		return nil
	}

	data := arrayToken.Data()
	var left, bottom, right, top *tokens.NumericToken

	if len(data) == 4 {
		left = asNumeric(data[0])
		bottom = asNumeric(data[1])
		right = asNumeric(data[2])
		top = asNumeric(data[3])
	} else if len(data) == 6 {
		left = asNumeric(data[2])
		bottom = asNumeric(data[3])
		right = asNumeric(data[4])
		top = asNumeric(data[5])
	}

	if left != nil && bottom != nil && right != nil && top != nil {
		rect := core.NewPdfRectangleFloat(
			left.DoubleVal(),
			bottom.DoubleVal(),
			right.DoubleVal(),
			top.DoubleVal(),
		)
		return &rect
	}

	return nil
}

func (b *markedContentBuilder) parseAttached() []*tokens.NameToken {
	if b.properties == nil {
		return make([]*tokens.NameToken, 0)
	}

	arrayToken, ok := tokens.TryGetTyped[*tokens.ArrayToken](b.properties, tokens.Attached)
	if !ok {
		return make([]*tokens.NameToken, 0)
	}

	result := make([]*tokens.NameToken, 0)
	for _, token := range arrayToken.Data() {
		if name, ok := token.(*tokens.NameToken); ok {
			result = append(result, name)
		}
	}

	return result
}

func asNumeric(token tokens.Token) *tokens.NumericToken {
	if num, ok := token.(*tokens.NumericToken); ok {
		return num
	}
	return nil
}

func (b *markedContentBuilder) fallbackElement(mcid int) *content.MarkedContentElement {
	name := tokens.Create("")
	props, _ := tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
	elem, _ := content.NewMarkedContentElement(
		mcid,
		name,
		props,
		nil, nil, nil, nil,
		false,
		b.children,
		b.letters,
		b.pathsToInterface(),
		b.images,
		b.number,
	)
	return elem
}

// pathsToInterface converts the internal geometry.PdfPath slice to []content.PdfPath.
func (b *markedContentBuilder) pathsToInterface() []content.PdfPath {
	result := make([]content.PdfPath, len(b.paths))
	for i := range b.paths {
		result[i] = content.PdfPath{PdfPath: &b.paths[i]}
	}
	return result
}

func getOptionalString(dict *tokens.DictionaryToken, name *tokens.NameToken) *string {
	if dict == nil || name == nil {
		return nil
	}

	token, ok := dict.TryGet(name)
	if !ok {
		return nil
	}

	switch t := token.(type) {
	case *tokens.StringToken:
		s := t.Data()
		return &s
	case *tokens.NameToken:
		s := t.Data()
		return &s
	}

	return nil
}
