//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"os"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/outline"
	"github.com/uglytoad/pdfpig/go/testutil"
)

// TestCanReadAccentedBookmarksCorrectly verifies that bookmarks containing
// accented characters are read correctly from a PDF document. This matches C#
// AccentedCharactersInBookmarksTests.CanReadAccentedBookmarksCorrectly.
func TestCanReadAccentedBookmarksCorrectly(t *testing.T) {
	path := testutil.GetDocumentPath("bookmarks-with-accented-characters.pdf", true)

	if _, err := os.Stat(path); err != nil {
		t.Skipf("test document not found: %s", path)
	}

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	resultAny, isSuccess, err := doc.TryGetBookmarks(false)
	if err != nil {
		t.Fatalf("TryGetBookmarks: %v", err)
	}

	if !isSuccess {
		t.Fatal("expected bookmarks to be present")
	}

	bookmarks, ok := resultAny.(*outline.Bookmarks)
	if !ok {
		t.Fatalf("expected *outline.Bookmarks, got %T", resultAny)
	}

	nodes := bookmarks.GetNodes()
	expected := []string{
		"\u017E",
		"\u017E\u010D",
		"\u017E\u0111",
		"\u017E\u0107",
		"\u017E\u0161",
		"\u017E ajklyghvbnmxcseqwuioprtzdf",
		"\u0161",
		"\u0161\u010D",
		"\u0161\u0111",
		"\u0161\u0107",
		"\u0161\u017E",
		"\u0161 ajklyghvbnmxcseqwuioprtzdf",
	}

	if len(nodes) != len(expected) {
		t.Fatalf("expected %d bookmark nodes, got %d", len(expected), len(nodes))
	}

	for i, node := range nodes {
		if node.Title != expected[i] {
			t.Errorf("bookmark[%d]: expected title %q, got %q", i, expected[i], node.Title)
		}
	}
}

// TestCanReadContainerBookmarksCorrectly verifies that bookmarks with and without
// container nodes are returned correctly. This matches C#
// AccentedCharactersInBookmarksTests.CanReadContainerBookmarksCorrectly.
func TestCanReadContainerBookmarksCorrectly(t *testing.T) {
	path := testutil.GetDocumentPath("dotnet-ai.pdf", true)

	if _, err := os.Stat(path); err != nil {
		t.Skipf("test document not found: %s", path)
	}

	doc, err := pdfpig.OpenFile(path, &content.ParsingOptions{})
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	resultAny, isSuccess, err := doc.TryGetBookmarks(false)
	if err != nil {
		t.Fatalf("TryGetBookmarks(false): %v", err)
	}

	if !isSuccess {
		t.Fatal("expected bookmarks to be present")
	}

	bookmarks, ok := resultAny.(*outline.Bookmarks)
	if !ok {
		t.Fatalf("expected *outline.Bookmarks, got %T", resultAny)
	}

	if len(bookmarks.Roots()) != 3 {
		t.Errorf("expected 3 root bookmarks without container nodes, got %d", len(bookmarks.Roots()))
	}

	resultAny2, isSuccess2, err := doc.TryGetBookmarks(true)
	if err != nil {
		t.Fatalf("TryGetBookmarks(true): %v", err)
	}

	if !isSuccess2 {
		t.Fatal("expected bookmarks to be present with container nodes")
	}

	bookmarksWithContainers, ok := resultAny2.(*outline.Bookmarks)
	if !ok {
		t.Fatalf("expected *outline.Bookmarks, got %T", resultAny2)
	}

	if len(bookmarksWithContainers.Roots()) <= 3 {
		t.Errorf("expected more than 3 root bookmarks with container nodes, got %d", len(bookmarksWithContainers.Roots()))
	}
}
