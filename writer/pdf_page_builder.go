// Package writer provides types for building PDF documents programmatically.
package writer

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/uglytoad/pdfpig/go/actions"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/fonts"
	"github.com/uglytoad/pdfpig/go/graphics"
	"github.com/uglytoad/pdfpig/go/graphics/colors"
	"github.com/uglytoad/pdfpig/go/images"
	png "github.com/uglytoad/pdfpig/go/images/png"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// nameConflictSolver generates unique names for XObjects, graphics states, etc.
type nameConflictSolver struct {
	prefix          string
	key             int
	xobjectNamesUsed map[string]bool
}

func newNameConflictSolver(prefix string) *nameConflictSolver {
	return &nameConflictSolver{
		prefix:           prefix,
		xobjectNamesUsed: make(map[string]bool),
	}
}

func (s *nameConflictSolver) extractPrefix(name string) string {
	if name == "" {
		return s.prefix
	}

	i := 0
	for i < len(name) && (name[i] < '0' || name[i] > '9') {
		i++
	}

	if i != 0 {
		return name[:i]
	}
	return s.prefix
}

func (s *nameConflictSolver) newName(originalName string) string {
	newPrefix := s.extractPrefix(originalName)

	name := fmt.Sprintf("%s%d", newPrefix, s.key)

	for s.xobjectNamesUsed[name] {
		s.key++
		name = fmt.Sprintf("%s%d", newPrefix, s.key)
	}

	s.xobjectNamesUsed[name] = true
	return name
}

func (s *nameConflictSolver) fixName(name string) string {
	if s.xobjectNamesUsed[name] {
		return s.newName(name)
	}

	s.xobjectNamesUsed[name] = true
	return name
}

// PageContentStream represents a content stream that can be written to a PDF page.
type PageContentStream interface {
	HasContent() bool
	Write(writer PdfStreamWriter) *tokens.IndirectReferenceToken
	ReadOnly() bool
	GlobalTransform() *core.TransformationMatrix
	Add(operation content.GraphicsStateOperation)
	Operations() []content.GraphicsStateOperation
}

// PdfPageBuilder builds the content and resources for a single PDF page.
type PdfPageBuilder struct {
	documentBuilder *PdfDocumentBuilder
	pageNumber      int
	pageSize        *content.MediaBox
	pageDictionary  map[*tokens.NameToken]tokens.Token
	contentStreams  []PageContentStream
	currentStream   PageContentStream
	links           []CopiedLink
	documentFonts   map[string]*tokens.NameToken
	nextFontId      int
	textSequence    int
	xobjectsNames   *nameConflictSolver
	gStateNames     *nameConflictSolver
	rotation        *int
}

// NewPdfPageBuilder creates a new PdfPageBuilder for the given page number.
func NewPdfPageBuilder(pageNumber int, documentBuilder *PdfDocumentBuilder) *PdfPageBuilder {
	if documentBuilder == nil {
		panic("documentBuilder must not be nil")
	}

	currentStream := &defaultContentStream{operations: make([]content.GraphicsStateOperation, 0)}

	return &PdfPageBuilder{
		documentBuilder: documentBuilder,
		pageNumber:      pageNumber,
		pageDictionary:  make(map[*tokens.NameToken]tokens.Token),
		contentStreams:  []PageContentStream{currentStream},
		currentStream:   currentStream,
		links:           make([]CopiedLink, 0),
		documentFonts:   make(map[string]*tokens.NameToken),
		nextFontId:      1,
		xobjectsNames:   newNameConflictSolver("I"),
		gStateNames:     newNameConflictSolver("GS"),
	}
}

// NewPdfPageBuilderFromCopied creates a PdfPageBuilder from copied content streams and page dictionary.
func NewPdfPageBuilderFromCopied(
	pageNumber int,
	documentBuilder *PdfDocumentBuilder,
	copied []PageContentStream,
	pageDict map[*tokens.NameToken]tokens.Token,
	links []CopiedLink,
) *PdfPageBuilder {
	if documentBuilder == nil {
		panic("documentBuilder must not be nil")
	}

	writeableContentStream := &defaultContentStream{operations: make([]content.GraphicsStateOperation, 0)}

	if len(copied) > 0 {
		for i := len(copied) - 1; i >= 0; i-- {
			gt := copied[i].GlobalTransform()
			if gt != nil {
				inverse := gt.Inverse()
				writeableContentStream.Add(&graphics.ModifyCurrentTransformationMatrix{
			Value: [6]float64{inverse.A, inverse.B, inverse.C, inverse.D, inverse.E, inverse.F},
				})
				break
			}
		}
	}

	builder := &PdfPageBuilder{
		documentBuilder: documentBuilder,
		pageNumber:      pageNumber,
		pageDictionary:  pageDict,
		contentStreams:  append(copied, writeableContentStream),
		currentStream:   writeableContentStream,
		links:           links,
		documentFonts:   make(map[string]*tokens.NameToken),
		nextFontId:      1,
		xobjectsNames:   newNameConflictSolver("I"),
		gStateNames:     newNameConflictSolver("GS"),
	}

	return builder
}

// PageNumber returns the 1-indexed page number.
func (b *PdfPageBuilder) PageNumber() int {
	return b.pageNumber
}

// PageSize returns the current size of the page.
func (b *PdfPageBuilder) PageSize() *content.MediaBox {
	return b.pageSize
}

// SetPageSize sets the page size rectangle.
func (b *PdfPageBuilder) SetPageSize(size *content.MediaBox) {
	b.pageSize = size
}

// CurrentStream returns the currently active content stream.
func (b *PdfPageBuilder) CurrentStream() PageContentStream {
	return b.currentStream
}

// ContentStreams returns all content streams for this page.
func (b *PdfPageBuilder) ContentStreams() []PageContentStream {
	return b.contentStreams
}

// Resources returns the resources dictionary, creating it if necessary.
func (b *PdfPageBuilder) Resources() map[string]tokens.Token {
	resources, ok := b.pageDictionary[tokens.Resources].(*tokens.DictionaryToken)
	if !ok || resources == nil {
		dict, _ := tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
		b.pageDictionary[tokens.Resources] = dict
		resources = dict
	}
	return resources.Data()
}

// Rotation returns the optional page rotation angle.
func (b *PdfPageBuilder) Rotation() *int {
	return b.rotation
}

// PageDictionary returns the internal page dictionary being built.
func (b *PdfPageBuilder) PageDictionary() map[*tokens.NameToken]tokens.Token {
	return b.pageDictionary
}

// Links returns the list of link annotations for this page.
func (b *PdfPageBuilder) Links() []CopiedLink {
	return b.links
}

// NewContentStreamBefore inserts a new content stream before the current one and selects it.
func (b *PdfPageBuilder) NewContentStreamBefore() {
	index := 0
	for i, cs := range b.contentStreams {
		if cs == b.currentStream {
			index = i - 1
			break
		}
	}
	if index < 0 {
		index = 0
	}

	b.currentStream = &defaultContentStream{operations: make([]content.GraphicsStateOperation, 0)}

	newStreams := make([]PageContentStream, 0, len(b.contentStreams)+1)
	newStreams = append(newStreams, b.contentStreams[:index]...)
	newStreams = append(newStreams, b.currentStream)
	newStreams = append(newStreams, b.contentStreams[index:]...)
	b.contentStreams = newStreams
}

// NewContentStreamAfter inserts a new content stream after the current one and selects it.
func (b *PdfPageBuilder) NewContentStreamAfter() {
	index := len(b.contentStreams)
	for i, cs := range b.contentStreams {
		if cs == b.currentStream {
			index = i + 1
			break
		}
	}
	if index > len(b.contentStreams) {
		index = len(b.contentStreams)
	}

	b.currentStream = &defaultContentStream{operations: make([]content.GraphicsStateOperation, 0)}

	newStreams := make([]PageContentStream, 0, len(b.contentStreams)+1)
	newStreams = append(newStreams, b.contentStreams[:index]...)
	newStreams = append(newStreams, b.currentStream)
	newStreams = append(newStreams, b.contentStreams[index:]...)
	b.contentStreams = newStreams
}

// SelectContentStream selects a content stream by index.
func (b *PdfPageBuilder) SelectContentStream(index int) error {
	if index < 0 || index >= len(b.contentStreams) {
		return fmt.Errorf("content stream index %d out of range [0, %d)", index, len(b.contentStreams))
	}
	b.currentStream = b.contentStreams[index]
	return nil
}

// DrawLine draws a line between two points with the specified line width.
func (b *PdfPageBuilder) DrawLine(from, to core.PdfPoint, lineWidth float64) *PdfPageBuilder {
	if lineWidth != 1 {
		b.currentStream.Add(&graphics.SetLineWidth{Width: lineWidth})
	}

	b.currentStream.Add(&graphics.BeginNewSubpath{X: from.X, Y: from.Y})
	b.currentStream.Add(&graphics.AppendStraightLineSegment{X: to.X, Y: to.Y})
	b.currentStream.Add(graphics.InstanceStrokePath)

	if lineWidth != 1 {
		b.currentStream.Add(&graphics.SetLineWidth{Width: 1})
	}

	return b
}

// DrawRectangle draws a rectangle at the specified position with given dimensions.
func (b *PdfPageBuilder) DrawRectangle(position core.PdfPoint, width, height, lineWidth float64, fill bool) *PdfPageBuilder {
	if lineWidth != 1 {
		b.currentStream.Add(&graphics.SetLineWidth{Width: lineWidth})
	}

	b.currentStream.Add(&graphics.AppendRectangle{LowerLeftX: position.X, LowerLeftY: position.Y, Width: width, Height: height})

	if fill {
		b.currentStream.Add(graphics.InstanceFillPathEvenOddRuleAndStroke)
	} else {
		b.currentStream.Add(graphics.InstanceStrokePath)
	}

	if lineWidth != 1 {
		b.currentStream.Add(&graphics.SetLineWidth{Width: 1})
	}

	return b
}

// SetRotation sets the number of degrees by which the page is rotated clockwise.
func (b *PdfPageBuilder) SetRotation(degrees content.PageRotationDegrees) *PdfPageBuilder {
	val := int(degrees.Value)
	b.rotation = &val
	return b
}

// DrawTriangle draws a triangle with three specified points.
func (b *PdfPageBuilder) DrawTriangle(p1, p2, p3 core.PdfPoint, lineWidth float64, fill bool) *PdfPageBuilder {
	if lineWidth != 1 {
		b.currentStream.Add(&graphics.SetLineWidth{Width: lineWidth})
	}

	b.currentStream.Add(&graphics.BeginNewSubpath{X: p1.X, Y: p1.Y})
	b.currentStream.Add(&graphics.AppendStraightLineSegment{X: p2.X, Y: p2.Y})
	b.currentStream.Add(&graphics.AppendStraightLineSegment{X: p3.X, Y: p3.Y})
	b.currentStream.Add(&graphics.AppendStraightLineSegment{X: p1.X, Y: p1.Y})

	if fill {
		b.currentStream.Add(graphics.InstanceFillPathEvenOddRuleAndStroke)
	} else {
		b.currentStream.Add(graphics.InstanceStrokePath)
	}

	if lineWidth != 1 {
		b.currentStream.Add(&graphics.SetLineWidth{Width: 1})
	}

	return b
}

// DrawCircle draws a circle centered at the given point with the specified diameter.
func (b *PdfPageBuilder) DrawCircle(center core.PdfPoint, diameter, lineWidth float64, fill bool) *PdfPageBuilder {
	return b.DrawEllipsis(center, diameter, diameter, lineWidth, fill)
}

// DrawEllipsis draws an ellipse centered at the given point with the specified width and height.
func (b *PdfPageBuilder) DrawEllipsis(center core.PdfPoint, width, height, lineWidth float64, fill bool) *PdfPageBuilder {
	width /= 2
	height /= 2

	cc := 0.55228474983079

	if lineWidth != 1 {
		b.currentStream.Add(&graphics.SetLineWidth{Width: lineWidth})
	}

	b.currentStream.Add(&graphics.BeginNewSubpath{X: center.X - width, Y: center.Y})
	b.currentStream.Add(&graphics.AppendDualControlPointBezierCurve{
		X1: center.X - width, Y1: center.Y + height*cc,
		X2: center.X - width*cc, Y2: center.Y + height,
		X3: center.X, Y3: center.Y + height,
	})
	b.currentStream.Add(&graphics.AppendDualControlPointBezierCurve{
		X1: center.X + width*cc, Y1: center.Y + height,
		X2: center.X + width, Y2: center.Y + height*cc,
		X3: center.X + width, Y3: center.Y,
	})
	b.currentStream.Add(&graphics.AppendDualControlPointBezierCurve{
		X1: center.X + width, Y1: center.Y - height*cc,
		X2: center.X + width*cc, Y2: center.Y - height,
		X3: center.X, Y3: center.Y - height,
	})
	b.currentStream.Add(&graphics.AppendDualControlPointBezierCurve{
		X1: center.X - width*cc, Y1: center.Y - height,
		X2: center.X - width, Y2: center.Y - height*cc,
		X3: center.X - width, Y3: center.Y,
	})

	if fill {
		b.currentStream.Add(graphics.InstanceFillPathEvenOddRuleAndStroke)
	} else {
		b.currentStream.Add(graphics.InstanceStrokePath)
	}

	if lineWidth != 1 {
		b.currentStream.Add(&graphics.SetLineWidth{Width: 1})
	}

	return b
}

// SetStrokeColor sets the stroke color to an RGB value (0-255 per channel).
func (b *PdfPageBuilder) SetStrokeColor(r, g, bv byte) *PdfPageBuilder {
	b.currentStream.Add(graphics.InstancePush)
	b.currentStream.Add(&graphics.SetStrokeColorDeviceRgb{R: rgbToDouble(r), G: rgbToDouble(g), B: rgbToDouble(bv)})
	return b
}

// SetStrokeColorExact sets the stroke color with exact values between 0 and 1.
func (b *PdfPageBuilder) SetStrokeColorExact(r, g, bv float64) (*PdfPageBuilder, error) {
	cr, err := checkRgbDouble(r, "r")
	if err != nil {
		return nil, err
	}
	cg, err := checkRgbDouble(g, "g")
	if err != nil {
		return nil, err
	}
	cb, err := checkRgbDouble(bv, "b")
	if err != nil {
		return nil, err
	}

	b.currentStream.Add(graphics.InstancePush)
	b.currentStream.Add(&graphics.SetStrokeColorDeviceRgb{R: cr, G: cg, B: cb})
	return b, nil
}

// SetTextAndFillColor sets the fill and text color to an RGB value (0-255 per channel).
func (b *PdfPageBuilder) SetTextAndFillColor(r, g, bv byte) *PdfPageBuilder {
	b.currentStream.Add(graphics.InstancePush)
	b.currentStream.Add(&graphics.SetNonStrokeColorDeviceRgb{R: rgbToDouble(r), G: rgbToDouble(g), B: rgbToDouble(bv)})
	return b
}

// ResetColor restores the stroke, text and fill color to default (black).
func (b *PdfPageBuilder) ResetColor() *PdfPageBuilder {
	b.currentStream.Add(graphics.InstancePop)
	return b
}

// MeasureText calculates letter sizes and positions without modifying page state.
func (b *PdfPageBuilder) MeasureText(text string, fontSize float64, position core.PdfPoint, font *AddedFont) ([]*content.Letter, error) {
	if font == nil {
		return nil, errors.New("font must not be nil")
	}
	if text == "" {
		return nil, errors.New("text must not be empty")
	}

	fontStore, ok := b.documentBuilder.Fonts()[font.Id]
	if !ok || fontStore == nil {
		return nil, fmt.Errorf("no font has been added to the PdfDocumentBuilder with Id: %s", font.Id)
	}
	if fontSize <= 0 {
		return nil, errors.New("fontSize must be greater than 0")
	}

	fontProgram := fontStore.FontProgram
	fm := fontProgram.GetFontMatrix()
	textMatrix := core.FromValues(1, 0, 0, 1, position.X, position.Y)

	letters, err := b.drawLetters(nil, text, fontProgram, fm, fontSize, textMatrix)
	if err != nil {
		return nil, fmt.Errorf("measureText failed: %w", err)
	}
	return letters, nil
}

// AddText draws text on the page at the specified position and returns the drawn letters.
func (b *PdfPageBuilder) AddText(text string, fontSize float64, position core.PdfPoint, font *AddedFont) ([]*content.Letter, error) {
	if font == nil {
		return nil, errors.New("font must not be nil")
	}
	if text == "" {
		return nil, errors.New("text must not be empty")
	}

	fontStore, ok := b.documentBuilder.Fonts()[font.Id]
	if !ok || fontStore == nil {
		return nil, fmt.Errorf("no font has been added to the PdfDocumentBuilder with Id: %s", font.Id)
	}
	if fontSize <= 0 {
		return nil, errors.New("fontSize must be greater than 0")
	}

	fontName := b.getAddedFont(font)
	fontProgram := fontStore.FontProgram
	fm := fontProgram.GetFontMatrix()
	textMatrix := core.FromValues(1, 0, 0, 1, position.X, position.Y)

	letters, err := b.drawLetters(fontName, text, fontProgram, fm, fontSize, textMatrix)
	if err != nil {
		return nil, fmt.Errorf("addText failed: %w", err)
	}

	b.currentStream.Add(graphics.InstanceBeginText)
	b.currentStream.Add(&graphics.SetFontAndSize{Font: fontName, Size: fontSize})
	b.currentStream.Add(&graphics.MoveToNextLineWithOffset{Tx: position.X, Ty: position.Y})

	bytesPerShow := make([]byte, 0)
	for _, letter := range text {
		if letter == ' ' || letter == '\t' || letter == '\n' || letter == '\r' {
			if len(bytesPerShow) > 0 {
				copied := make([]byte, len(bytesPerShow))
				copy(copied, bytesPerShow)
				b.currentStream.Add(&graphics.ShowText{Bytes: copied})
				bytesPerShow = bytesPerShow[:0]
			}
		}
		byteVal := fontProgram.GetValueForCharacter(letter)
		bytesPerShow = append(bytesPerShow, byteVal)
	}

	if len(bytesPerShow) > 0 {
		copied := make([]byte, len(bytesPerShow))
		copy(copied, bytesPerShow)
		b.currentStream.Add(&graphics.ShowText{Bytes: copied})
	}

	b.currentStream.Add(graphics.InstanceEndText)

	return letters, nil
}

// SetTextRenderingMode sets the text rendering mode for future AddText calls.
func (b *PdfPageBuilder) SetTextRenderingMode(mode core.TextRenderingMode) *PdfPageBuilder {
	b.currentStream.Add(&graphics.SetTextRenderingMode{Mode: mode})
	return b
}

func (b *PdfPageBuilder) getAddedFont(font *AddedFont) *tokens.NameToken {
	if existing, ok := b.documentFonts[font.Id]; ok {
		return existing
	}

	value := tokens.Create(fmt.Sprintf("F%d", b.nextFontId))
	b.nextFontId++

	resourcesData := b.Resources()
	fontsDict, ok := resourcesData["Font"]
	if !ok || fontsDict == nil {
		newFonts, _ := tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
		resourcesData["Font"] = newFonts
		fontsDict = newFonts
	}

	fontsMap := fontsDict.(*tokens.DictionaryToken).Data()
	for fontsMap[value.Data()] != nil {
		value = tokens.Create(fmt.Sprintf("F%d", b.nextFontId))
		b.nextFontId++
	}

	b.documentFonts[font.Id] = value
	fontsMap[value.Data()] = font.Reference

	return value
}

// AddedImage represents an image that has been added to a PDF document.
type AddedImage struct {
	Id        string
	Reference *tokens.IndirectReferenceToken
	Width     int
	Height    int
}

func newAddedImage(reference *tokens.IndirectReferenceToken, width, height int) *AddedImage {
	return &AddedImage{
		Id:        fmt.Sprintf("%p", reference),
		Reference: reference,
		Width:     width,
		Height:    height,
	}
}

// AddJpeg adds a JPEG image from bytes at the specified location.
func (b *PdfPageBuilder) AddJpeg(fileBytes []byte, placementRectangle core.PdfRectangle) (*AddedImage, error) {
	stream := bytes.NewReader(fileBytes)
	return b.AddJpegStream(stream, placementRectangle)
}

// AddJpegStream adds a JPEG image from an io.ReadSeeker at the specified location.
func (b *PdfPageBuilder) AddJpegStream(fileStream jpegReader, placementRectangle core.PdfRectangle) (*AddedImage, error) {
	startFrom, err := fileStream.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream position: %w", err)
	}

	info, err := images.GetInformation(fileStream)
	if err != nil {
		return nil, fmt.Errorf("failed to read JPEG information: %w", err)
	}

	if placementRectangle == (core.PdfRectangle{}) {
		placementRectangle = core.PdfRectangle{
			BottomLeft: core.PdfPoint{},
			TopRight:   core.PdfPoint{X: float64(info.Width), Y: float64(info.Height)},
		}
	}

	var data bytes.Buffer
	fileStream.Seek(startFrom, io.SeekStart)
	if _, err := io.Copy(&data, fileStream); err != nil {
		return nil, fmt.Errorf("failed to read JPEG data: %w", err)
	}

	colorSpace := tokens.Devicergb
	if info.NumberOfComponents == 1 {
		colorSpace = tokens.Devicegray
	} else if info.NumberOfComponents == 4 {
		colorSpace = tokens.Devicecmyk
	}

	imgDict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
		tokens.Type:             tokens.Xobject,
		tokens.Subtype:          tokens.Image,
		tokens.Width:           tokens.NewNumericTokenFromInt(info.Width),
		tokens.Height:          tokens.NewNumericTokenFromInt(info.Height),
		tokens.BitsPerComponent: tokens.NewNumericTokenFromInt(info.BitsPerComponent),
		tokens.ColorSpace:      colorSpace,
		tokens.Filter:          tokens.DctDecode,
		tokens.Length:          tokens.NewNumericTokenFromInt(data.Len()),
	})

	reference := b.documentBuilder.AddImage(imgDict, data.Bytes())
	resourcesData := b.Resources()
	xObjects := getOrCreateXObjects(resourcesData)

	key := tokens.Create(b.xobjectsNames.newName(""))
	xObjects[key.Data()] = reference

	b.currentStream.Add(graphics.InstancePush)
	b.currentStream.Add(&graphics.ModifyCurrentTransformationMatrix{
		Value: [6]float64{
			placementRectangle.Width, 0,
			0, placementRectangle.Height,
			placementRectangle.BottomLeft.X, placementRectangle.BottomLeft.Y,
		},
	})
	b.currentStream.Add(&graphics.InvokeNamedXObject{Name: key})
	b.currentStream.Add(graphics.InstancePop)

	return newAddedImage(reference, info.Width, info.Height), nil
}

// AddJpegReference reuses a previously added JPEG image to avoid duplication.
func (b *PdfPageBuilder) AddJpegReference(image *AddedImage, placementRectangle core.PdfRectangle) {
	b.AddImageRef(image, placementRectangle)
}

// AddImageRef adds an already-added image at the specified location without duplicating data.
func (b *PdfPageBuilder) AddImageRef(image *AddedImage, placementRectangle core.PdfRectangle) {
	resourcesData := b.Resources()
	xObjects := getOrCreateXObjects(resourcesData)

	key := tokens.Create(b.xobjectsNames.newName(""))
	xObjects[key.Data()] = image.Reference

	b.currentStream.Add(graphics.InstancePush)
	b.currentStream.Add(&graphics.ModifyCurrentTransformationMatrix{
		Value: [6]float64{
			placementRectangle.Width, 0,
			0, placementRectangle.Height,
			placementRectangle.BottomLeft.X, placementRectangle.BottomLeft.Y,
		},
	})
	b.currentStream.Add(&graphics.InvokeNamedXObject{Name: key})
	b.currentStream.Add(graphics.InstancePop)
}

// AddPng adds a PNG image from bytes at the specified location.
func (b *PdfPageBuilder) AddPng(pngBytes []byte, placementRectangle core.PdfRectangle) (*AddedImage, error) {
	stream := bytes.NewReader(pngBytes)
	return b.AddPngStream(stream, placementRectangle)
}

// AddPngStream adds a PNG image from an io.Reader at the specified location.
func (b *PdfPageBuilder) AddPngStream(pngStream io.Reader, placementRectangle core.PdfRectangle) (*AddedImage, error) {
	pngImg, err := png.Open(pngStream, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open PNG: %w", err)
	}

	if placementRectangle == (core.PdfRectangle{}) {
		placementRectangle = core.PdfRectangle{
			BottomLeft: core.PdfPoint{},
			TopRight:   core.PdfPoint{X: float64(pngImg.Width()), Y: float64(pngImg.Height())},
		}
	}

	width := pngImg.Width()
	height := pngImg.Height()
	data := make([]byte, width*height*3)
	pixelIndex := 0

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			pixel := pngImg.GetPixel(col, row)
			data[pixelIndex] = pixel.R
			data[pixelIndex+1] = pixel.G
			data[pixelIndex+2] = pixel.B
			pixelIndex += 3
		}
	}

	widthToken := tokens.NewNumericTokenFromInt(width)
	heightToken := tokens.NewNumericTokenFromInt(height)

	var smaskReference *tokens.IndirectReferenceToken

	if pngImg.HasAlphaChannel() && b.documentBuilder.ArchiveStandard() != PdfA1B && b.documentBuilder.ArchiveStandard() != PdfA1A {
		smaskData := make([]byte, len(data)/3)
		for row := 0; row < height; row++ {
			for col := 0; col < width; col++ {
				pixel := pngImg.GetPixel(col, row)
				index := row*width + col
				smaskData[index] = pixel.A
			}
		}

		compressedSmask := CompressBytes(smaskData)

		smaskDict, _ := tokens.NewDictionary(map[*tokens.NameToken]tokens.Token{
			tokens.Type:             tokens.Xobject,
			tokens.Subtype:          tokens.Image,
			tokens.Width:           widthToken,
			tokens.Height:          heightToken,
			tokens.ColorSpace:      tokens.Devicegray,
			tokens.BitsPerComponent: tokens.NewNumericTokenFromInt(8),
			tokens.Decode:          tokens.NewArrayToken([]tokens.Token{tokens.NewNumericToken(0), tokens.NewNumericToken(1)}),
			tokens.Length:          tokens.NewNumericTokenFromInt(len(compressedSmask)),
			tokens.Filter:          tokens.FlateDecode,
		})

		smaskReference = b.documentBuilder.AddImage(smaskDict, compressedSmask)
	}

	compressed := CompressBytes(data)

	imgDictEntries := map[*tokens.NameToken]tokens.Token{
		tokens.Type:             tokens.Xobject,
		tokens.Subtype:          tokens.Image,
		tokens.Width:           widthToken,
		tokens.Height:          heightToken,
		tokens.BitsPerComponent: tokens.NewNumericTokenFromInt(8),
		tokens.ColorSpace:      tokens.Devicergb,
		tokens.Filter:          tokens.FlateDecode,
		tokens.Length:          tokens.NewNumericTokenFromInt(len(compressed)),
	}

	if smaskReference != nil {
		imgDictEntries[tokens.Smask] = smaskReference
	}

	imgDict, _ := tokens.NewDictionary(imgDictEntries)
	reference := b.documentBuilder.AddImage(imgDict, compressed)

	resourcesData := b.Resources()
	xObjects := getOrCreateXObjects(resourcesData)

	key := tokens.Create(b.xobjectsNames.newName(""))
	xObjects[key.Data()] = reference

	b.currentStream.Add(graphics.InstancePush)
	b.currentStream.Add(&graphics.ModifyCurrentTransformationMatrix{
		Value: [6]float64{
			placementRectangle.Width, 0,
			0, placementRectangle.Height,
			placementRectangle.BottomLeft.X, placementRectangle.BottomLeft.Y,
		},
	})
	b.currentStream.Add(&graphics.InvokeNamedXObject{Name: key})
	b.currentStream.Add(graphics.InstancePop)

	return newAddedImage(reference, width, height), nil
}

// AddLink adds a URL link annotation at the specified rectangle area.
func (b *PdfPageBuilder) AddLink(url string, linkArea core.PdfRectangle) *PdfPageBuilder {
	action := &actions.UriAction{Uri: url}
	link := NewLinkAnnotation(action, linkArea, nil, nil, nil, nil)
	return b.AddLinkAnnotation(link)
}

// AddLinkDestination adds an internal document link annotation.
func (b *PdfPageBuilder) AddLinkDestination(destination destinations.ExplicitDestination, linkArea core.PdfRectangle) *PdfPageBuilder {
	action := actions.NewGoToAction(destination)
	link := NewLinkAnnotation(action, linkArea, nil, nil, nil, nil)
	return b.AddLinkAnnotation(link)
}

// AddLinkAnnotation adds a link annotation to the page.
func (b *PdfPageBuilder) AddLinkAnnotation(link *LinkAnnotation) *PdfPageBuilder {
	b.links = append(b.links, CopiedLink{
		token:  link.ToToken(),
		action: link.Action(),
	})
	return b
}

// CopyFrom copies page content from an existing page into this builder.
func (b *PdfPageBuilder) CopyFrom(srcPage *content.Page) (*PdfPageBuilder, error) {
	if srcPage == nil {
		return nil, errors.New("srcPage must not be nil")
	}

	currentOps := b.currentStream.Operations()
	if len(currentOps) > 0 {
		b.NewContentStreamAfter()
	}

	destinationStream := b.currentStream

	srcDict := srcPage.Dictionary()
	if srcDict == nil {
		srcOps := srcPage.Operations()
		for _, op := range srcOps {
			b.currentStream.Add(op)
		}
		return b, nil
	}

	var srcResourceDictionary *tokens.DictionaryToken
	resourcesToken, ok := srcDict.TryGet(tokens.Resources)
	if !ok || resourcesToken == nil {
		srcOps := srcPage.Operations()
		for _, op := range srcOps {
			b.currentStream.Add(op)
		}
		return b, nil
	}

	srcResourceDictionary = resolveToDictionary(resourcesToken, srcPage.PdfScanner())
	if srcResourceDictionary == nil {
		srcOps := srcPage.Operations()
		for _, op := range srcOps {
			b.currentStream.Add(op)
		}
		return b, nil
	}

	scanner := srcPage.PdfScanner()
	operations := make([]content.GraphicsStateOperation, len(srcPage.Operations()))
	copy(operations, srcPage.Operations())

	resourcesData := b.Resources()

	for keyStr, value := range srcResourceDictionary.Data() {
		nameToken := tokens.Create(keyStr)
		if nameToken.Equals(tokens.Font) || nameToken.Equals(tokens.Xobject) {
			continue
		}

		if _, exists := resourcesData[keyStr]; !exists {
			copiedValue := b.documentBuilder.CopyToken(scanner, value)
			resourcesData[keyStr] = copiedValue
			continue
		}
	}

	if fontsToken, ok := srcResourceDictionary.TryGet(tokens.Font); ok && fontsToken != nil {
		fontsDict := resolveToDictionary(fontsToken, scanner)
		if fontsDict == nil {
			return b, nil
		}

		pageFontsData := getOrCreateFonts(resourcesData)

		for fontKeyStr, fontValue := range fontsDict.Data() {
			fontName := tokens.Create(fontKeyStr)
			if pageFontsData[fontName.Data()] != nil {
				newName := tokens.Create(fmt.Sprintf("F%d", b.nextFontId))
				b.nextFontId++
				for pageFontsData[newName.Data()] != nil {
					newName = tokens.Create(fmt.Sprintf("F%d", b.nextFontId))
					b.nextFontId++
				}

				oldName := fontName.Data()
				var newNT *tokens.NameToken = newName
				operations = renameSetFontOperations(operations, oldName, newNT)

				fontName = newName
			}

			if refToken, ok := fontValue.(*tokens.IndirectReferenceToken); ok {
				copiedRef := b.documentBuilder.CopyToken(scanner, refToken)
				pageFontsData[fontName.Data()] = copiedRef
			} else {
				return b, &core.PdfDocumentFormatException{
					Message: fmt.Sprintf("expected an IndirectReferenceToken for the font, got a %T", fontValue),
				}
			}
		}
	}

	if xobjectsToken, ok := srcResourceDictionary.TryGet(tokens.Xobject); ok && xobjectsToken != nil {
		xobjectsDict := resolveToDictionary(xobjectsToken, scanner)
		if xobjectsDict == nil {
			return b, nil
		}

		pageXObjectsData := getOrCreateXObjects(resourcesData)

		for xobjectKeyStr, xobjectValue := range xobjectsDict.Data() {
			xobjectName := xobjectKeyStr
			newName := b.xobjectsNames.fixName(xobjectName)

			if xobjectName != newName {
				oldName := xobjectKeyStr
				var newNT *tokens.NameToken = tokens.Create(newName)
				operations = renameInvokeXObjectOperations(operations, oldName, newNT)

				xobjectName = newName
			}

			if refToken, ok := xobjectValue.(*tokens.IndirectReferenceToken); ok {
				copiedRef := b.documentBuilder.CopyToken(scanner, refToken)
				pageXObjectsData[xobjectName] = copiedRef
			} else {
				return b, &core.PdfDocumentFormatException{
					Message: fmt.Sprintf("expected an IndirectReferenceToken for the XObject, got a %T", xobjectValue),
				}
			}
		}
	}

	if gsToken, ok := srcResourceDictionary.TryGet(tokens.ExtGState); ok && gsToken != nil {
		gsDict := resolveToDictionary(gsToken, scanner)
		if gsDict == nil {
			return b, nil
		}

		pageGstateData := getOrCreateExtGState(resourcesData)

		for gstateKeyStr, gstateValue := range gsDict.Data() {
			gstateName := gstateKeyStr
			newName := b.gStateNames.fixName(gstateName)

			if newName != gstateName {
				oldName := gstateKeyStr
				var newNT *tokens.NameToken = tokens.Create(newName)
				operations = renameSetGStateOperations(operations, oldName, newNT)

				gstateName = newName
			}

			copiedValue := b.documentBuilder.CopyToken(scanner, gstateValue)
			pageGstateData[gstateName] = copiedValue
		}
	}

	globalTransform := GetGlobalTransform(operations)
	if globalTransform != nil {
		inverse := globalTransform.Inverse()
		operations = append(operations, &graphics.ModifyCurrentTransformationMatrix{
			Value: [6]float64{inverse.A, inverse.B, inverse.C, inverse.D, inverse.E, inverse.F},
		})
	}

	for _, op := range operations {
		destinationStream.Add(op)
	}

	return b, nil
}

func (b *PdfPageBuilder) drawLetters(name *tokens.NameToken, text string, font WritingFont, fontMatrix core.TransformationMatrix, fontSize float64, textMatrix core.TransformationMatrix) ([]*content.Letter, error) {
	horizontalScaling := 1.0
	rise := 0.0

	letters := make([]*content.Letter, 0)

	renderingMatrix := core.FromValues(fontSize*horizontalScaling, 0, 0, fontSize, 0, rise)

	width := 0.0
	b.textSequence++

	fontName := ""
	if name != nil {
		fontName = name.Data()
	}

	for _, c := range text {
		rect, ok := font.TryGetBoundingBox(c)
		if !ok {
			return nil, fmt.Errorf("the font does not contain a character: '%c' (0x%X)", c, c)
		}

		charWidth, ok := font.TryGetAdvanceWidth(c)
		if !ok {
			return nil, fmt.Errorf("the font does not contain advance width for character: '%c'", c)
		}

		advanceRect := core.NewPdfRectangleFloat(0, 0, charWidth, 0)
		advanceRect = textMatrix.TransformRect(renderingMatrix.TransformRect(fontMatrix.TransformRect(advanceRect)))

		documentSpace := textMatrix.TransformRect(renderingMatrix.TransformRect(fontMatrix.TransformRect(*rect)))

		fontDetails := fonts.GetDefault(fontName)

		letter := content.NewLetterWithDetails(
			string(c),
			documentSpace,
			documentSpace,
			advanceRect.BottomLeft,
			advanceRect.BottomRight,
			width,
			fontSize,
			fontDetails,
			core.FillText,
			colors.GrayBlack,
			colors.GrayBlack,
			fontSize,
			b.textSequence,
		)

		letters = append(letters, letter)

		tx := advanceRect.Width * horizontalScaling
		ty := 0.0

		translate := core.GetTranslationMatrix(tx, ty)
		width += tx
		textMatrix = translate.Multiply(textMatrix)
	}

	return letters, nil
}

func rgbToDouble(value byte) float64 {
	res := math.Max(0, float64(value)/255.0)
	res = math.Round(math.Min(1, res)*10000) / 10000
	return res
}

func checkRgbDouble(value float64, argument string) (float64, error) {
	if value < 0 {
		return 0, fmt.Errorf("provided double for RGB color %s was less than zero: %f", argument, value)
	}
	if value > 1 {
		return 0, fmt.Errorf("provided double for RGB color %s was greater than one: %f", argument, value)
	}
	return value, nil
}

// defaultContentStream is a writable content stream that accumulates graphics operations.
type defaultContentStream struct {
	operations []content.GraphicsStateOperation
}

func (d *defaultContentStream) HasContent() bool {
	return len(d.operations) > 0
}

func (d *defaultContentStream) ReadOnly() bool {
	return false
}

func (d *defaultContentStream) GlobalTransform() *core.TransformationMatrix {
	return nil
}

func (d *defaultContentStream) Add(operation content.GraphicsStateOperation) {
	d.operations = append(d.operations, operation)
}

func (d *defaultContentStream) Operations() []content.GraphicsStateOperation {
	return d.operations
}

type writeOperation interface {
	Write(w io.Writer) error
}

func (d *defaultContentStream) Write(writer PdfStreamWriter) *tokens.IndirectReferenceToken {
	var buf bytes.Buffer
	for _, op := range d.operations {
		if wo, ok := any(op).(writeOperation); ok {
			_ = wo.Write(&buf)
		}
	}

	stream := CompressToStream(buf.Bytes())
	ref := writer.WriteToken(stream)
	return ref
}

// copiedContentStream represents a content stream copied from another document.
type copiedContentStream struct {
	reference       *tokens.IndirectReferenceToken
	globalTransform *core.TransformationMatrix
}

func (c *copiedContentStream) HasContent() bool {
	return c.reference != nil
}

func (c *copiedContentStream) ReadOnly() bool {
	return true
}

func (c *copiedContentStream) GlobalTransform() *core.TransformationMatrix {
	return c.globalTransform
}

func (c *copiedContentStream) Add(_ content.GraphicsStateOperation) {
	panic("writing to a copied content stream is not supported")
}

func (c *copiedContentStream) Operations() []content.GraphicsStateOperation {
	panic("reading raw operations from a copied content stream is not supported")
}

func (c *copiedContentStream) Write(_ PdfStreamWriter) *tokens.IndirectReferenceToken {
	return c.reference
}

// Helper functions for resource dictionary manipulation.

func getOrCreateXObjects(resourcesData map[string]tokens.Token) map[string]tokens.Token {
	xObjectsRaw, ok := resourcesData["XObject"]
	if !ok || xObjectsRaw == nil {
		newDict, _ := tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
		resourcesData["XObject"] = newDict
		xObjectsRaw = newDict
	}
	return xObjectsRaw.(*tokens.DictionaryToken).Data()
}

func getOrCreateFonts(resourcesData map[string]tokens.Token) map[string]tokens.Token {
	fontsRaw, ok := resourcesData["Font"]
	if !ok || fontsRaw == nil {
		newDict, _ := tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
		resourcesData["Font"] = newDict
		fontsRaw = newDict
	}
	return fontsRaw.(*tokens.DictionaryToken).Data()
}

func getOrCreateExtGState(resourcesData map[string]tokens.Token) map[string]tokens.Token {
	gsRaw, ok := resourcesData["ExtGState"]
	if !ok || gsRaw == nil {
		newDict, _ := tokens.NewDictionary(make(map[*tokens.NameToken]tokens.Token))
		resourcesData["ExtGState"] = newDict
		gsRaw = newDict
	}
	return gsRaw.(*tokens.DictionaryToken).Data()
}

// Operation renaming helpers for CopyFrom.

func renameSetFontOperations(operations []content.GraphicsStateOperation, oldName string, newName *tokens.NameToken) []content.GraphicsStateOperation {
	result := make([]content.GraphicsStateOperation, len(operations))
	copy(result, operations)
	for i, op := range result {
		if setFont, ok := op.(*graphics.SetFontAndSize); ok && setFont.Font.Data() == oldName {
			result[i] = &graphics.SetFontAndSize{Font: newName, Size: setFont.Size}
		}
	}
	return result
}

func renameInvokeXObjectOperations(operations []content.GraphicsStateOperation, oldName string, newName *tokens.NameToken) []content.GraphicsStateOperation {
	result := make([]content.GraphicsStateOperation, len(operations))
	copy(result, operations)
	for i, op := range result {
		if invoke, ok := op.(*graphics.InvokeNamedXObject); ok && invoke.Name.Data() == oldName {
			result[i] = &graphics.InvokeNamedXObject{Name: newName}
		}
	}
	return result
}

func renameSetGStateOperations(operations []content.GraphicsStateOperation, oldName string, newName *tokens.NameToken) []content.GraphicsStateOperation {
	result := make([]content.GraphicsStateOperation, len(operations))
	copy(result, operations)
	for i, op := range result {
		if gs, ok := op.(*graphics.SetGraphicsStateParametersFromDictionary); ok && gs.Name.Data() == oldName {
			result[i] = &graphics.SetGraphicsStateParametersFromDictionary{Name: newName}
		}
	}
	return result
}

// jpegReader combines read, seek capabilities needed for JPEG parsing.
type jpegReader interface {
	io.Reader
	io.Seeker
	io.ByteReader
}

func resolveToDictionary(token tokens.Token, scanner tokenization.PdfTokenScanner) *tokens.DictionaryToken {
	if dict, ok := token.(*tokens.DictionaryToken); ok {
		return dict
	}
	if ref, ok := token.(*tokens.IndirectReferenceToken); ok {
		obj := scanner.Get(ref.Data())
		if obj != nil {
			if dict, ok := obj.Data().(*tokens.DictionaryToken); ok {
				return dict
			}
		}
	}
	return nil
}
