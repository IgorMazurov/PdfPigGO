package testpages

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// SimplePage is a minimal page representation holding only page number, rotation, and media box.
type SimplePage struct {
	Number   int
	Rotation int
	MediaBox *content.MediaBox
}

// SimplePageFactory creates SimplePage instances from PDF page dictionaries.
type SimplePageFactory struct {
	content.BasePageFactory[any]
}

func (f *SimplePageFactory) OutputTypeName() string {
	return "*testpages.SimplePage"
}

func NewSimplePageFactory(
	pdfScanner any,
	resourceStore content.ResourceStore,
	filterProvider content.LookupFilterProvider,
	pageContentParser content.PageContentParser,
	parsingOptions *content.ParsingOptions,
) *SimplePageFactory {
	return &SimplePageFactory{}
}

func (f *SimplePageFactory) Create(
	number int,
	dictionary *tokens.DictionaryToken,
	pageTreeMembers *content.PageTreeMembers,
	namedDestinations any,
) (any, error) {
	rotation := 0
	if pageTreeMembers.Rotation != nil {
		rotation = *pageTreeMembers.Rotation
	}

	mediaBox := pageTreeMembers.MediaBoxValue
	if mediaBox == nil {
		mediaBox = content.MediaBoxUSLetter
	}

	return &SimplePage{
		Number:   number,
		Rotation: rotation,
		MediaBox: mediaBox,
	}, nil
}

// PageInformation holds page metadata including dimensions and user space unit.
type PageInformation struct {
	Number        int
	Rotation      content.PageRotationDegrees
	Width         float64
	Height        float64
	UserSpaceUnit geometry.UserSpaceUnit
}

// PageInformationFactory creates PageInformation instances from PDF page dictionaries.
type PageInformationFactory struct {
	content.BasePageFactory[any]
}

func (f *PageInformationFactory) OutputTypeName() string {
	return "*testpages.PageInformation"
}

func NewPageInformationFactory(
	pdfScanner any,
	resourceStore content.ResourceStore,
	filterProvider content.LookupFilterProvider,
	pageContentParser content.PageContentParser,
	parsingOptions *content.ParsingOptions,
) *PageInformationFactory {
	return &PageInformationFactory{}
}

func (f *PageInformationFactory) Create(
	number int,
	dictionary *tokens.DictionaryToken,
	pageTreeMembers *content.PageTreeMembers,
	namedDestinations any,
) (any, error) {
	rotation := content.PageRotationDegrees{Value: 0}
	if pageTreeMembers.Rotation != nil {
		rotation.Value = *pageTreeMembers.Rotation
	}

	mediaBox := pageTreeMembers.MediaBoxValue
	if mediaBox == nil {
		mediaBox = content.MediaBoxUSLetter
	}

	cropBox := pageTreeMembers.CropBoxValue
	if cropBox == nil {
		cropBox = content.NewCropBox(mediaBox.Bounds)
	}

	userSpaceUnit := content.GetUserSpaceUnits(dictionary)

	viewBox := geometry.RectangleIntersect(mediaBox.Bounds, cropBox.Bounds)
	if viewBox == nil {
		viewBox = &cropBox.Bounds
	}

	return &PageInformation{
		Number:        number,
		Rotation:      rotation,
		Width:         viewBox.Right() - viewBox.Left(),
		Height:        viewBox.Top() - viewBox.Bottom(),
		UserSpaceUnit: userSpaceUnit,
	}, nil
}

// WrongConstructorFactory is a factory that returns the wrong type for negative testing.
type WrongConstructorFactory struct {
	content.BasePageFactory[any]
}

func (f *WrongConstructorFactory) OutputTypeName() string {
	return "*testpages.PageInformation"
}

func NewWrongConstructorFactory(
	pdfScanner any,
	resourceStore content.ResourceStore,
	filterProvider content.LookupFilterProvider,
	pageContentParser content.PageContentParser,
	parsingOptions *content.ParsingOptions,
) *WrongConstructorFactory {
	return &WrongConstructorFactory{}
}

func (f *WrongConstructorFactory) Create(
	number int,
	dictionary *tokens.DictionaryToken,
	pageTreeMembers *content.PageTreeMembers,
	namedDestinations any,
) (any, error) {
	panic("should never be called")
}
