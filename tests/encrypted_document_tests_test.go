//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"errors"
	"testing"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/encryption"
	"github.com/uglytoad/pdfpig/go/testutil"
)

const encryptedFileName = "encrypted-password-is-password.pdf"
const encryptedPassword = "password"

func getEncryptedPath() string {
	return testutil.GetSpecificTestDocumentPath(encryptedFileName, true)
}

// TestNoPasswordThrows verifies that opening an encrypted document without a
// password returns a PdfDocumentEncryptedException. This matches C#
// EncryptedDocumentTests.NoPasswordThrows.
func TestNoPasswordThrows(t *testing.T) {
	path := getEncryptedPath()

	opts := &content.ParsingOptions{}

	_, err := pdfpig.OpenFile(path, opts)
	if err == nil {
		t.Fatal("expected error when opening encrypted document without password")
	}

	var encErr *encryption.PdfDocumentEncryptedException
	if !errors.As(err, &encErr) {
		t.Fatalf("expected PdfDocumentEncryptedException, got %T: %v", err, err)
	}
}

// TestCanOpenDocumentAndGetPage verifies that an encrypted document can be opened
// with the correct password and all pages have non-empty text. This matches C#
// EncryptedDocumentTests.CanOpenDocumentAndGetPage.
func TestCanOpenDocumentAndGetPage(t *testing.T) {
	path := getEncryptedPath()

	opts := &content.ParsingOptions{
		Password: encryptedPassword,
	}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for i := 1; i <= doc.NumberOfPages(); i++ {
		pageAny, err := doc.GetPage(i)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", i, err)
		}

		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page, got %T", i, pageAny)
		}

		if page.Text() == "" {
			t.Errorf("page %d: text is empty", i)
		}
	}
}

// TestCanProvideMultiplePasswords verifies that an encrypted document can be opened
// when the correct password is among a list of candidate passwords. This matches C#
// EncryptedDocumentTests.CanProvideMultiplePasswords.
func TestCanProvideMultiplePasswords(t *testing.T) {
	path := getEncryptedPath()

	opts := &content.ParsingOptions{
		Passwords: []string{"pangolin", "harpsichord", encryptedPassword},
	}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for i := 1; i <= doc.NumberOfPages(); i++ {
		pageAny, err := doc.GetPage(i)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", i, err)
		}

		page, ok := pageAny.(*content.Page)
		if !ok {
			t.Fatalf("GetPage(%d): expected *content.Page, got %T", i, pageAny)
		}

		if page.Text() == "" {
			t.Errorf("page %d: text is empty", i)
		}
	}
}

// TestCanReadDocumentWithUEAsString verifies that a document with the /U entry as
// a string (rather than hex) can be opened and its metadata read. This matches C#
// EncryptedDocumentTests.CanReadDocumentWithUEAsString.
func TestCanReadDocumentWithUEAsString(t *testing.T) {
	path := testutil.GetSpecificTestDocumentPath("string_encryption_key.pdf", true)

	opts := &content.ParsingOptions{}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if doc.Information == nil || doc.Information.Producer == "" {
		t.Error("expected non-empty Producer in document information")
	}
}

// TestCanReadDocumentWithEmptyStringEncryptedWithAESEncryptionAndOnlyIV verifies that
// a document containing an empty string encrypted with AES (only IV present) can be
// opened and its metadata read. This matches C#
// EncryptedDocumentTests.CanReadDocumentWithEmptyStringEncryptedWithAESEncryptionAndOnlyIV.
func TestCanReadDocumentWithEmptyStringEncryptedWithAESEncryptionAndOnlyIV(t *testing.T) {
	path := testutil.GetSpecificTestDocumentPath("r4_aes_empty_string.pdf", true)

	opts := &content.ParsingOptions{}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if doc.Information == nil || doc.Information.Producer != "" {
		t.Error("expected empty Producer in document information")
	}
}

// TestCanReadDocumentWithNoKeyLengthAndRevision4 verifies that a revision 4 AESv2
// encrypted document with no key length specified can be opened and its metadata read.
// This matches C# EncryptedDocumentTests.CanReadDocumentWithNoKeyLengthAndRevision4.
func TestCanReadDocumentWithNoKeyLengthAndRevision4(t *testing.T) {
	path := testutil.GetSpecificTestDocumentPath("r4_aesv2_no_length.pdf", true)

	opts := &content.ParsingOptions{}

	doc, err := pdfpig.OpenFile(path, opts)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	if doc.Information == nil || doc.Information.Producer != "" {
		t.Error("expected empty Producer in document information")
	}
}
