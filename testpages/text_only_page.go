package testpages

import (
	"strings"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/geometry"
	"github.com/uglytoad/pdfpig/go/graphics"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/tokenization"
)

// TextOnlyPage holds only the page number and extracted text.
type TextOnlyPage struct {
	Number int
	Text   string
}

// TextOnlyPageInformationFactory creates TextOnlyPage instances using ContentStreamProcessor
// to collect letters from the page content stream.
type TextOnlyPageInformationFactory struct {
	content.BasePageFactory[TextOnlyPage]
}

func (f *TextOnlyPageInformationFactory) OutputTypeName() string {
	return "*testpages.TextOnlyPage"
}

func NewTextOnlyPageInformationFactory(
	pdfScanner tokenization.PdfTokenScanner,
	resourceStore content.ResourceStore,
	filterProvider content.LookupFilterProvider,
	pageContentParser content.PageContentParser,
	parsingOptions *content.ParsingOptions,
) *TextOnlyPageInformationFactory {
	f := &TextOnlyPageInformationFactory{}
	f.BasePageFactory = *content.NewBasePageFactory[TextOnlyPage](
		pdfScanner,
		resourceStore,
		filterProvider,
		pageContentParser,
		parsingOptions,
	)

	f.ProcessPageHook = func(
		pageNumber int,
		dictionary *tokens.DictionaryToken,
		namedDestinations *destinations.NamedDestinations,
		mediaBox *content.MediaBox,
		cropBox *content.CropBox,
		userSpaceUnit geometry.UserSpaceUnit,
		rotation content.PageRotationDegrees,
		initialMatrix core.TransformationMatrix,
		operations []content.GraphicsStateOperation,
	) TextOnlyPage {
		if len(operations) == 0 {
			return TextOnlyPage{Number: pageNumber, Text: ""}
		}

		context := graphics.NewContentStreamProcessor(
			pageNumber,
			f.ResourceStore,
			f.PdfScanner,
			f.PageContentParser,
			f.FilterProvider,
			cropBox,
			userSpaceUnit,
			rotation,
			initialMatrix,
			f.ParsingOptions,
		)

		pageContent := context.Process(pageNumber, operations)

		var sb strings.Builder
		for _, letter := range pageContent.Letters() {
			sb.WriteString(letter.Value)
		}

		return TextOnlyPage{Number: pageNumber, Text: sb.String()}
	}

	return f
}

func (f *TextOnlyPageInformationFactory) Create(
	number int,
	dictionary *tokens.DictionaryToken,
	pageTreeMembers *content.PageTreeMembers,
	namedDestinations any,
) (any, error) {
	result, err := f.BasePageFactory.Create(number, dictionary, pageTreeMembers, namedDestinations.(*destinations.NamedDestinations))
	if err != nil {
		return nil, err
	}
	return &result, nil
}
