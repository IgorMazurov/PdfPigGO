//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"errors"
	"path/filepath"
	"testing"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokens"
)

func getSpecificTestDocPath(name string) string {
	return filepath.Join(specificTestDocRoot, name)
}

// TestCanReadDocumentInformation verifies that document metadata fields and custom
// dictionary entries can be read from a PDF with custom properties.
// This matches C# DocumentInformationTests.CanReadDocumentInformation.
func TestCanReadDocumentInformation(t *testing.T) {
	path := getSpecificTestDocPath("custom-properties.pdf")

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	info := doc.Information

	if info.Creator != "Writer" {
		t.Errorf("Creator = %q, want %q", info.Creator, "Writer")
	}
	if info.Keywords != "MoreKeywords" {
		t.Errorf("Keywords = %q, want %q", info.Keywords, "MoreKeywords")
	}
	if info.Producer != "LibreOffice 6.1" {
		t.Errorf("Producer = %q, want %q", info.Producer, "LibreOffice 6.1")
	}
	if info.Subject != "TestSubject" {
		t.Errorf("Subject = %q, want %q", info.Subject, "TestSubject")
	}
	if info.Title != "TestTitle" {
		t.Errorf("Title = %q, want %q", info.Title, "TestTitle")
	}

	infoDict := info.DocumentInformationDictionary

	nameToken := tokens.Create("CustomProperty1")
	valueToken, ok := infoDict.TryGet(nameToken)
	if !ok {
		t.Fatal("first custom property must be present")
	}
	stringToken, ok := valueToken.(*tokens.StringToken)
	if !ok {
		t.Fatalf("expected *tokens.StringToken, got %T", valueToken)
	}
	if stringToken.Data() != "Property Value" {
		t.Errorf("CustomProperty1 = %q, want %q", stringToken.Data(), "Property Value")
	}

	nameToken = tokens.Create("CustomProperty2")
	valueToken2, ok := infoDict.TryGet(nameToken)
	if !ok {
		t.Fatal("second custom property must be present")
	}
	stringToken2, ok := valueToken2.(*tokens.StringToken)
	if !ok {
		t.Fatalf("expected *tokens.StringToken, got %T", valueToken2)
	}
	if stringToken2.Data() != "Another Property Value" {
		t.Errorf("CustomProperty2 = %q, want %q", stringToken2.Data(), "Another Property Value")
	}
}

// TestCanReadInvalidDocumentInformation verifies that a PDF with invalid structure
// in the info dictionary can be read in lenient mode but fails in strict mode.
// This matches C# DocumentInformationTests.CanReadInvalidDocumentInformation.
func TestCanReadInvalidDocumentInformation(t *testing.T) {
	path := getSpecificTestDocPath("invalid-pdf-structure-pdfminer-entire-doc.pdf")

	// Lenient parsing on -> can process
	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile (lenient): %v", err)
	}
	defer doc.Close()

	info := doc.Information

	if info.Creator != "LaTeX with hyperref" {
		t.Errorf("Creator = %q, want %q", info.Creator, "LaTeX with hyperref")
	}
	if info.Keywords != "" {
		t.Errorf("Keywords = %q, want empty", info.Keywords)
	}
	if info.Producer != "pdfTeX-1.40.21" {
		t.Errorf("Producer = %q, want %q", info.Producer, "pdfTeX-1.40.21")
	}
	if info.Subject != "" {
		t.Errorf("Subject = %q, want empty", info.Subject)
	}
	if info.Title != "" {
		t.Errorf("Title = %q, want empty", info.Title)
	}
	if info.Author != "" {
		t.Errorf("Author = %q, want empty", info.Author)
	}
	if info.CreationDate != "D:20230418010134Z" {
		t.Errorf("CreationDate = %q, want %q", info.CreationDate, "D:20230418010134Z")
	}
	if info.ModifiedDate != "D:20230418010134Z" {
		t.Errorf("ModifiedDate = %q, want %q", info.ModifiedDate, "D:20230418010134Z")
	}

	infoDict := info.DocumentInformationDictionary

	nameToken := tokens.Create("Trapped")
	valueToken, ok := infoDict.TryGet(nameToken)
	if !ok {
		t.Fatal("Trapped entry must be present")
	}
	trappedNameToken, ok := valueToken.(*tokens.NameToken)
	if !ok {
		t.Fatalf("expected *tokens.NameToken for Trapped, got %T", valueToken)
	}
	if trappedNameToken.Data() != "False" {
		t.Errorf("Trapped = %q, want %q", trappedNameToken.Data(), "False")
	}

	nameToken = tokens.Create("PTEX.Fullbanner")
	valueToken2, ok := infoDict.TryGet(nameToken)
	if !ok {
		t.Fatal("PTEX.Fullbanner entry must be present")
	}
	stringToken2, ok := valueToken2.(*tokens.StringToken)
	if !ok {
		t.Fatalf("expected *tokens.StringToken for PTEX.Fullbanner, got %T", valueToken2)
	}
	expectedBanner := "This is pdfTeX, Version 3.14159265-2.6-1.40.21 (TeX Live 2020) kpathsea version 6.3.2"
	if stringToken2.Data() != expectedBanner {
		t.Errorf("PTEX.Fullbanner = %q, want %q", stringToken2.Data(), expectedBanner)
	}

	// Lenient parsing off -> throws
	_, err = pdfpig.OpenFile(path, &pdfpig.LenientParsingOff)
	if err == nil {
		t.Fatal("expected error when opening with lenient parsing off")
	}

	var formatErr *core.PdfDocumentFormatException
	if !errors.As(err, &formatErr) {
		t.Fatalf("expected *core.PdfDocumentFormatException, got %T: %v", err, err)
	}
	expectedMsg := "Expected name as dictionary key, instead got: Collaborative"
	if formatErr.Message != expectedMsg {
		t.Errorf("expected message %q, got %q", expectedMsg, formatErr.Message)
	}
}

// TestCanReadDocumentInformationIndirectRef verifies that document information stored
// via an indirect reference can be read correctly. (Issue 706)
// This matches C# DocumentInformationTests.CanReadDocumentInformationIndirectRef.
func TestCanReadDocumentInformationIndirectRef(t *testing.T) {
	path := getSpecificTestDocPath("EBOOK-DIETETYKA-SPORTOWA_copy_1.pdf")

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer doc.Close()

	info := doc.Information

	if info.Title != "EBOOK" {
		t.Errorf("Title = %q, want %q", info.Title, "EBOOK")
	}
	if info.Creator != "Pages" {
		t.Errorf("Creator = %q, want %q", info.Creator, "Pages")
	}
	if info.CreationDate != "D:20190306232856Z00'00'" {
		t.Errorf("CreationDate = %q, want %q", info.CreationDate, "D:20190306232856Z00'00'")
	}
}

// TestCanReadDocumentInformationDirectory verifies that a document with an info
// dictionary in the trailer can be read in lenient mode but fails in strict mode.
// (Issue 884) This matches C# DocumentInformationTests.CanReadDocumentInfromationDirectory.
func TestCanReadDocumentInformationDirectory(t *testing.T) {
	path := getSpecificTestDocPath("info_dictionary.pdf")

	// Lenient parsing on -> can process
	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("OpenFile (lenient): %v", err)
	}
	defer doc.Close()

	info := doc.Information

	if info.Producer != "SumatraPDF 3.2" {
		t.Errorf("Producer = %q, want %q", info.Producer, "SumatraPDF 3.2")
	}

	// Lenient parsing off -> throws
	_, err = pdfpig.OpenFile(path, &pdfpig.LenientParsingOff)
	if err == nil {
		t.Fatal("expected error when opening with lenient parsing off")
	}

	var formatErr *core.PdfDocumentFormatException
	if !errors.As(err, &formatErr) {
		t.Fatalf("expected *core.PdfDocumentFormatException, got %T: %v", err, err)
	}
	expectedMsg := "The info token in the trailer dictionary should only contain indirect references, instead got: <Producer, (SumatraPDF 3.2)>."
	if formatErr.Message != expectedMsg {
		t.Errorf("expected message %q, got %q", expectedMsg, formatErr.Message)
	}
}
