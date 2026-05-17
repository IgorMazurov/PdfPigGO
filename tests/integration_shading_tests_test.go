//go:build integration

package pdfpig_test


import (
	pdfpig "github.com/uglytoad/pdfpig/go"
	"path/filepath"
	"testing"
)

func getPDFBOX1869Path() string {
	return filepath.Join(integrationDocRoot, "PDFBOX-1869-4-1.pdf")
}

func getAxialRadial1Path() string {
	return filepath.Join(integrationDocRoot, "68-1990-01_A.pdf")
}

func getAxialRadialTensorProduct1Path() string {
	return filepath.Join(integrationDocRoot, "MOZILLA-3136-0.pdf")
}

func getIronOreQ2Q32013Path() string {
	return filepath.Join(integrationDocRoot, "iron-ore-q2-q3-2013.pdf")
}

// TestIssue702 verifies that a document containing a FunctionBasedShading can be
// opened and its first page retrieved without error. This matches C# ShadingTests.Issue702.
func TestIssue702(t *testing.T) {
	path := getPDFBOX1869Path()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	_, err = doc.GetPage(1)
	if err != nil {
		t.Fatalf("GetPage(1): %v", err)
	}
}

// TestAxialRadial1 verifies that pages with axial and radial shadings can be parsed
// correctly. This matches C# ShadingTests.AxialRadial1.
func TestAxialRadial1(t *testing.T) {
	path := getAxialRadial1Path()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	pages := []int{7, 14, 15, 16, 19}
	for _, pageNum := range pages {
		_, err = doc.GetPage(pageNum)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", pageNum, err)
		}
	}
}

// TestAxialRadialTensorProduct1 verifies that all pages of a document containing axial,
// radial, and tensor product shadings can be parsed without error. This matches
// C# ShadingTests.AxialRadialTensorProduct1.
func TestAxialRadialTensorProduct1(t *testing.T) {
	path := getAxialRadialTensorProduct1Path()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	for i := 0; i < doc.NumberOfPages(); i++ {
		_, err = doc.GetPage(i + 1)
		if err != nil {
			t.Fatalf("GetPage(%d): %v", i+1, err)
		}
	}
}

// TestAxialRadialTensorProductManyFunctions2 verifies that page 8 of a document with
// many shading functions can be parsed without error. This matches C#
// ShadingTests.AxialRadialTensorProductManyFunctions2.
func TestAxialRadialTensorProductManyFunctions2(t *testing.T) {
	path := getIronOreQ2Q32013Path()

	doc, err := pdfpig.OpenFile(path, nil)
	if err != nil {
		t.Fatalf("pdfpig.OpenFile(%q): %v", path, err)
	}
	defer doc.Close()

	_, err = doc.GetPage(8)
	if err != nil {
		t.Fatalf("GetPage(8): %v", err)
	}
}
