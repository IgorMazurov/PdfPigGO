package writer

import (
	"fmt"
	"io"
	"os"

	"github.com/uglytoad/pdfpig/go/actions"
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/outline/destinations"
	"github.com/uglytoad/pdfpig/go/parser"
)

// MergeTwoFiles merges two PDF documents together with the pages from file1 followed by file2.
func MergeTwoFiles(
	file1, file2 string,
	file1Selection, file2Selection []int,
	archiveStandard PdfAStandard,
	docInfo *content.DocumentInformation,
) ([]byte, error) {
	var buf seekableBuffer
	if err := MergeTwoStreams(file1, file2, &buf, file1Selection, file2Selection, archiveStandard, docInfo); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MergeTwoStreams merges two PDF documents from file paths into the output stream.
func MergeTwoStreams(
	file1, file2 string,
	output io.WriteSeeker,
	file1Selection, file2Selection []int,
	archiveStandard PdfAStandard,
	docInfo *content.DocumentInformation,
) error {
	if file1 == "" {
		return fmt.Errorf("file1 must not be empty")
	}
	if file2 == "" {
		return fmt.Errorf("file2 must not be empty")
	}

	stream1, err := os.Open(file1)
	if err != nil {
		return fmt.Errorf("cannot open file1: %w", err)
	}
	defer stream1.Close()

	stream2, err := os.Open(file2)
	if err != nil {
		return fmt.Errorf("cannot open file2: %w", err)
	}
	defer stream2.Close()

	pagesBundle := [][]int{file1Selection, file2Selection}
	return mergeFromStreams([]io.ReadSeeker{stream1, stream2}, output, pagesBundle, archiveStandard, docInfo)
}

// MergeFiles merges multiple PDF documents together with the pages in the order the file paths are provided.
func MergeFiles(filePaths ...string) ([]byte, error) {
	return MergeFilesWithOptions(PdfANone, nil, filePaths...)
}

// MergeFilesWithOptions merges multiple PDF documents with archive standard and document info options.
func MergeFilesWithOptions(
	archiveStandard PdfAStandard,
	docInfo *content.DocumentInformation,
	filePaths ...string,
) ([]byte, error) {
	var buf seekableBuffer
	if err := MergeFilesToStreamWithOptions(&buf, archiveStandard, docInfo, filePaths...); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MergeFilesToStream merges multiple PDF documents into the output stream.
func MergeFilesToStream(output io.WriteSeeker, filePaths ...string) error {
	return MergeFilesToStreamWithOptions(output, PdfANone, nil, filePaths...)
}

// MergeFilesToStreamWithOptions merges multiple PDF documents into the output stream with options.
func MergeFilesToStreamWithOptions(
	output io.WriteSeeker,
	archiveStandard PdfAStandard,
	docInfo *content.DocumentInformation,
	filePaths ...string,
) error {
	streams := make([]io.ReadSeeker, 0, len(filePaths))
	defer func() {
		for _, s := range streams {
			if closer, ok := s.(interface{ Close() error }); ok {
				closer.Close()
			}
		}
	}()

	for i, fp := range filePaths {
		if fp == "" {
			return fmt.Errorf("null filepath at index %d", i)
		}
		f, err := os.Open(fp)
		if err != nil {
			return fmt.Errorf("cannot open file %q: %w", fp, err)
		}
		streams = append(streams, f)
	}

	return mergeFromStreams(streams, output, nil, archiveStandard, docInfo)
}

// MergeBytes merges a set of PDF documents given as byte slices.
func MergeBytes(files [][]byte, pagesBundle [][]int, archiveStandard PdfAStandard, docInfo *content.DocumentInformation) ([]byte, error) {
	if files == nil {
		return nil, fmt.Errorf("files must not be nil")
	}

	var buf seekableBuffer
	if err := MergeBytesToStream(files, &buf, pagesBundle, archiveStandard, docInfo); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MergeBytesToStream merges a set of PDF documents given as byte slices into the output stream.
func MergeBytesToStream(
	files [][]byte,
	output io.WriteSeeker,
	pagesBundle [][]int,
	archiveStandard PdfAStandard,
	docInfo *content.DocumentInformation,
) error {
	docs := make([]*content.PdfDocument, len(files))
	for i, f := range files {
		doc, err := parser.OpenMemory(f, nil)
		if err != nil {
			return fmt.Errorf("cannot open file at index %d: %w", i, err)
		}
		docs[i] = doc
	}
	return mergeFromDocuments(docs, output, pagesBundle, archiveStandard, docInfo)
}

// MergeStreams merges PDF documents from a list of streams into the output stream.
func MergeStreams(
	streams []io.ReadSeeker,
	output io.WriteSeeker,
	pagesBundle [][]int,
	archiveStandard PdfAStandard,
	docInfo *content.DocumentInformation,
) error {
	if streams == nil {
		return fmt.Errorf("streams must not be nil")
	}
	if output == nil {
		return fmt.Errorf("output must not be nil")
	}
	return mergeFromStreams(streams, output, pagesBundle, archiveStandard, docInfo)
}

func mergeFromStreams(
	streams []io.ReadSeeker,
	output io.WriteSeeker,
	pagesBundle [][]int,
	archiveStandard PdfAStandard,
	docInfo *content.DocumentInformation,
) error {
	docs := make([]*content.PdfDocument, len(streams))
	for i, s := range streams {
		data, err := io.ReadAll(s)
		if err != nil {
			return fmt.Errorf("cannot read stream at index %d: %w", i, err)
		}
		doc, err := parser.OpenMemory(data, nil)
		if err != nil {
			return fmt.Errorf("cannot parse pdf from stream at index %d: %w", i, err)
		}
		docs[i] = doc
	}
	return mergeFromDocuments(docs, output, pagesBundle, archiveStandard, docInfo)
}

func mergeFromDocuments(
	files []*content.PdfDocument,
	output io.WriteSeeker,
	pagesBundle [][]int,
	archiveStandard PdfAStandard,
	docInfo *content.DocumentInformation,
) error {
	maxVersion := 0.0
	for _, f := range files {
		if v := f.Version(); v > maxVersion {
			maxVersion = v
		}
	}

	builder := NewPdfDocumentBuilderWithStream(output, false, PdfWriterDefault, maxVersion, nil)
	builder.SetArchiveStandard(archiveStandard)

	if docInfo != nil {
		builder.SetIncludeDocumentInformation(true)
		infoBuilder := convertToDocInfoBuilder(docInfo)
		builder.SetDocumentInformation(infoBuilder)
	}

	for fileIndex := 0; fileIndex < len(files); fileIndex++ {
		existing := files[fileIndex]
		var pages []int
		if pagesBundle != nil && fileIndex < len(pagesBundle) {
			pages = pagesBundle[fileIndex]
		}

		basePageNumber := len(builder.Pages())

		if pages == nil {
			for i := 1; i <= existing.NumberOfPages(); i++ {
				copyLinkFunc := func(action actions.Action) actions.Action {
					return copyLink(action, func(n int) *int {
						v := basePageNumber + n
						return &v
					})
				}
				opts := &AddPageOptions{
					KeepAnnotations: true,
					CopyLinkFunc:    copyLinkFunc,
				}
				builder.AddPageWithOptions(existing, i, opts)
			}
		} else {
			pageNumbers := make(map[int]int)
			for i := 0; i < len(pages); i++ {
				pageNumbers[pages[i]] = basePageNumber + i + 1
			}

			for _, pg := range pages {
				copyLinkFunc := func(action actions.Action) actions.Action {
					return copyLink(action, func(n int) *int {
						v, ok := pageNumbers[n]
						if !ok {
							return nil
						}
						return &v
					})
				}
				opts := &AddPageOptions{
					KeepAnnotations: true,
					CopyLinkFunc:    copyLinkFunc,
				}
				builder.AddPageWithOptions(existing, pg, opts)
			}
		}
	}

	return builder.Close()
}

// copyLink remaps GoTo-type action destinations to new page numbers in the merged document.
// For non-GoTo actions (URI, Launch, etc.) it returns the action unchanged.
// If getPageNumber returns nil for a destination's page number, the link is dropped (returns nil).
func copyLink(action actions.Action, getPageNumber func(int) *int) actions.Action {
	if action == nil {
		return nil
	}

	link, ok := toAbstractGoToAction(action)
	if !ok {
		return action
	}

	newPageNum := getPageNumber(link.Destination.PageNumber)
	if newPageNum == nil {
		return nil
	}

	newDest := destinations.NewExplicitDestination(*newPageNum, link.Destination.Type, link.Destination.Coordinates)

	switch link.Pdf().Type {
	case actions.GoTo:
		return actions.NewGoToAction(newDest)
	case actions.GoToE:
		if goToE, ok := action.(*actions.GoToEAction); ok {
			return actions.NewGoToEAction(newDest, goToE.FileSpecification)
		}
		return actions.NewGoToEAction(newDest, "")
	case actions.GoToR:
		if goToR, ok := action.(*actions.GoToRAction); ok {
			return actions.NewGoToRAction(newDest, goToR.Filename)
		}
		return actions.NewGoToRAction(newDest, "")
	default:
		return action
	}
}

// toAbstractGoToAction type-asserts the action to a concrete GoTo-type action and returns
// its embedded AbstractGoToAction with Destination data intact. Returns (nil, false) for
// non-GoTo action types.
func toAbstractGoToAction(action actions.Action) (*actions.AbstractGoToAction, bool) {
	switch a := action.(type) {
	case *actions.GoToAction:
		return a.AbstractGoToAction, true
	case *actions.GoToEAction:
		return a.AbstractGoToAction, true
	case *actions.GoToRAction:
		return a.AbstractGoToAction, true
	default:
		return nil, false
	}
}

// convertToDocInfoBuilder converts a content.DocumentInformation to DocumentInformationBuilder.
func convertToDocInfoBuilder(info *content.DocumentInformation) *DocumentInformationBuilder {
	builder := NewDocumentInformationBuilder()
	if info.Title != "" {
		builder.Title = &info.Title
	}
	if info.Author != "" {
		builder.Author = &info.Author
	}
	if info.Subject != "" {
		builder.Subject = &info.Subject
	}
	if info.Keywords != "" {
		builder.Keywords = &info.Keywords
	}
	if info.Creator != "" {
		builder.Creator = &info.Creator
	}
	if info.Producer != "" {
		builder.Producer = info.Producer
	}
	if info.CreationDate != "" {
		builder.CreationDate = &info.CreationDate
	}
	if info.ModifiedDate != "" {
		builder.ModifiedDate = &info.ModifiedDate
	}
	return builder
}
