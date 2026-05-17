package annotations

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/tokenization"
)

// HyperlinkFactory creates Hyperlink instances from page annotations.
type HyperlinkFactory struct{}

// GetHyperlinks extracts all hyperlinks from the given page by inspecting link annotations
// with /URI actions and collecting the letters that fall within each link's bounding rectangle.
func (f *HyperlinkFactory) GetHyperlinks(
	pageLetters []*content.Letter,
	scanner tokenization.PdfTokenScanner,
	annotations []*Annotation,
) []*content.Hyperlink {
	linkAnnotations := make([]content.LinkAnnotationIface, 0, len(annotations))
	for _, a := range annotations {
		if a.Type() == Link {
			linkAnnotations = append(linkAnnotations, newLinkAnnotationAdapter(a))
		}
	}

	return content.GetHyperlinksFromPage(pageLetters, scanner, linkAnnotations)
}
