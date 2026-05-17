//go:build gallery

// Package main provides an interactive gallery that lets users pick and run
// individual PdfPig examples by number, analogous to the C# Program.cs entry point.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// example represents a single runnable demonstration with a display name and action.
type example struct {
	name   string
	action func() error
}

// GalleryMain is the entry point for the interactive examples gallery.
func GalleryMain() {
	fmt.Println("Welcome to the PdfPig examples gallery!")

	documentsDir := filepath.Join("..", "..", "test_documents")

	examples := map[int]example{
		1: {
			name: "Extract Words with newline detection (example with algorithm)",
			action: func() error {
				return OpenDocumentAndExtractWords(filepath.Join(documentsDir, "Two Page Text Only - from libre office.pdf"))
			},
		},
		2: {
			name: "Extract Text with newlines (using built-in content extractor)",
			action: func() error {
				return ExtractTextWithNewlines(filepath.Join(documentsDir, "Two Page Text Only - from libre office.pdf"))
			},
		},
		3: {
			name: "Extract images",
			action: func() error {
				return ExtractImages(filepath.Join(documentsDir, "2006_Swedish_Touring_Car_Championship.pdf"))
			},
		},
		4: {
			name: "Merge PDF Documents",
			action: func() error {
				return MergePdfDocuments(
					filepath.Join(documentsDir, "Two Page Text Only - from libre office.pdf"),
					filepath.Join(documentsDir, "2006_Swedish_Touring_Car_Championship.pdf"),
					filepath.Join(documentsDir, "Rotated Text Libre Office.pdf"),
				)
			},
		},
		5: {
			name: "Extract form contents",
			action: func() error {
				return GetFormContents(filepath.Join(documentsDir, "AcroFormsBasicFields.pdf"))
			},
		},
		6: {
			name: "Generate PDF/A-2A compliant file",
			action: func() error {
				return GeneratePdfA2AFile(
					filepath.Join(documentsDir, "..", "..", "Fonts", "TrueType", "Roboto-Regular.ttf"),
					filepath.Join(documentsDir, "smile-250-by-160.jpg"),
				)
			},
		},
		7: {
			name: "Advance text extraction using layout analysis algorithms",
			action: func() error {
				return AdvancedTextExtraction(filepath.Join(documentsDir, "ICML03-081.pdf"))
			},
		},
		8: {
			name: "Extract Words with newline detection (example with algorithm). Issue 512",
			action: func() error {
				return OpenDocumentAndExtractWords(filepath.Join(documentsDir, "OPEN.RABBIT.ENGLISH.LOP.pdf"))
			},
		},
	}

	choices := buildChoices(examples)
	fmt.Println(choices)
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter a number to pick an example (enter 'q' to exit): ")
		if !scanner.Scan() {
			break
		}

		val := strings.TrimSpace(scanner.Text())

		if strings.EqualFold(val, "q") {
			return
		}

		opt, ok := parseChoice(val)
		if !ok {
			fmt.Printf("No option with value: %s.\n", val)
			fmt.Println()
			fmt.Println(choices)
			fmt.Println()
			continue
		}

		ex, exists := examples[opt]
		if !exists {
			fmt.Printf("No option with value: %s.\n", val)
			fmt.Println()
			fmt.Println(choices)
			fmt.Println()
			continue
		}

		if err := ex.action(); err != nil {
			fmt.Fprintf(os.Stderr, "Example error: %v\n", err)
		}

		fmt.Println()
		fmt.Println()
		fmt.Println(choices)
		fmt.Println()
	}
}

// buildChoices returns a formatted string listing all available examples.
func buildChoices(examples map[int]example) string {
	var lines []string
	for i := 1; i <= len(examples); i++ {
		if ex, ok := examples[i]; ok {
			lines = append(lines, fmt.Sprintf("%d: %s", i, ex.name))
		}
	}
	return strings.Join(lines, "\n")
}

// parseChoice converts a string input to an integer option number.
func parseChoice(val string) (int, bool) {
	var opt int
	_, err := fmt.Sscanf(val, "%d", &opt)
	if err != nil {
		return 0, false
	}
	return opt, true
}


func main() {
	GalleryMain()
}
