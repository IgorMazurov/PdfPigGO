# PdfPigGO

**A High-Performance, Comprehensive PDF Parsing Library for Go**

[![Go Report Card](https://goreportcard.com/badge/github.com/IgorMazurov/PdfPigGO)](https://goreportcard.com/report/github.com/IgorMazurov/PdfPigGO)
[![GoDoc](https://godoc.org/github.com/IgorMazurov/PdfPigGO?status.svg)](https://godoc.org/github.com/IgorMazurov/PdfPigGO)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 🧪 The Experiment: 100% AI-Engineered Port

**PdfPigGO** is a complete port of the renowned C# library [PdfPig](https://github.com/UglyToad/PdfPig) to the Go programming language. 

This project was born as a radical experiment in automated software engineering. In an era where "AI-assisted coding" usually refers to autocomplete or small snippets, PdfPigGO takes it to the logical extreme.

### Technical Genesis
- **Zero Manual Code:** Not a single line of code in this repository was written by human hands. Every struct, interface, method, and test case was generated through a specialized transpilation prompt engineering process.
- **Hardware:** The entire porting process was executed locally on a single consumer-grade workstation equipped with an **NVIDIA RTX 3090 Ti (24GB VRAM)**.
- **The Model:** The heavy lifting was performed by the **Qwen3.6 27B** model (running via local inference). This specific model size was chosen for its balance between architectural reasoning and strict adherence to Go idioms while maintaining the logic of the original C# source.
- **Validation:** The experiment is considered a success because **all ported tests pass**. Furthermore, real-world local projects that previously relied on the C# version of PdfPig have been successfully migrated to PdfPigGO with identical output and behavior.

---

## 📖 About PdfPig

The original C# PdfPig is arguably the most powerful open-source PDF library in the .NET ecosystem. Unlike many PDF libraries that simply provide low-level access to PDF objects, PdfPig focuses on **high-level data extraction**. 

It allows users to read PDF files and extract:
- **Text:** With precise positioning, font information, and structural awareness.
- **Images:** Exporting embedded images in their original formats.
- **Metadata:** Accessing document information, permissions, and encryption details.
- **Shapes and Paths:** Interpreting vector graphics within the document.
- **Tables:** Providing the primitives necessary to reconstruct tabular data.

**PdfPigGO** brings these capabilities to the Go ecosystem, maintaining the "UglyToad" philosophy of making PDF content accessible without requiring the developer to be an expert in the 1,000-page PDF specification.

---

## 🚀 Features

- **Text Extraction:** Extract text with character-level detail (including coordinates, font size, and font names).
- **Page Layout:** High-level access to pages, allowing for easy iteration and specific page selection.
- **Coordinate System:** Uses a consistent coordinate system (standard PDF points, where (0,0) is bottom-left).
- **Resource Management:** Efficient handling of PDF objects and streams.
- **Go Idioms:** While the logic is C#, the implementation uses Go's concurrency patterns, error handling, and slice management.
- **No CGO:** This is a pure Go implementation, making it cross-platform and easy to compile for any target (Windows, Linux, macOS, WASM).

---

## 📦 Installation

To include PdfPigGO in your project, use `go get`:

```bash
go get github.com/IgorMazurov/PdfPigGO
```

---

## 🛠 Usage Examples

Since this is a port, developers familiar with the C# version will find the API very intuitive. Below are the primary ways to interact with the library.

### 1. Basic Text Extraction
This is the most common use case. PdfPigGO provides a simple way to open a file and read every word on every page.

```go
package main

import (
	"fmt"
	"log"

	"github.com/IgorMazurov/PdfPigGO"
	"github.com/IgorMazurov/PdfPigGO/content"
)

func main() {
	// Open the document
	doc, err := pdfpig.Open("sample.pdf")
	if err != nil {
		log.Fatalf("Failed to open PDF: %v", err)
	}
	defer doc.Close()

	// Iterate through pages
	for i := 1; i <= doc.NumberOfPages; i++ {
		page, err := doc.GetPage(i)
		if err != nil {
			log.Printf("Failed to get page %d: %v", i, err)
			continue
		}

		// Get simple text
		fmt.Printf("Content of page %d:\n", i)
		fmt.Println(page.Text)

		// Or iterate through detailed words
		for _, word := range page.GetWords() {
			fmt.Printf("Word: [%s] at Bounds: %v\n", word.Text, word.BoundingBox)
		}
	}
}
```

### 2. Extracting Images
PdfPigGO can identify and extract images embedded in the PDF pages.

```go
package main

import (
	"github.com/IgorMazurov/PdfPigGO"
	"os"
)

func main() {
	doc, _ := pdfpig.Open("document_with_images.pdf")
	defer doc.Close()

	page, _ := doc.GetPage(1)
	
	for i, image := range page.GetImages() {
		// image.RawBytes contains the actual image data
		filename := fmt.Sprintf("extracted_image_%d.jpg", i)
		os.WriteFile(filename, image.RawBytes, 0644)
	}
}
```

### 3. Searching for Specific Content
Because PdfPigGO provides bounding boxes for every character, you can perform advanced spatial queries.

```go
package main

import (
    "github.com/IgorMazurov/PdfPigGO"
    "github.com/IgorMazurov/PdfPigGO/geometry"
    "fmt"
)

func main() {
    doc, _ := pdfpig.Open("invoice.pdf")
    page, _ := doc.GetPage(1)

    // Define a region of interest (e.g., where the "Total" is usually located)
    roi := geometry.NewRect(400, 100, 200, 50) 

    words := page.GetWords()
    for _, word := range words {
        if roi.Contains(word.BoundingBox) {
            fmt.Printf("Found text in ROI: %s\n", word.Text)
        }
    }
}
```

---

## 🏗 Why this Port Matters

### The PDF Challenge
PDF (Portable Document Format) is a notorious "dark format." It wasn't designed for data extraction; it was designed for printing. Internally, a PDF doesn't store "sentences" or "paragraphs." It stores instructions like: *"Move the pen to (100, 500) and draw the character 'A' using font X."*

The original C# PdfPig solved the immense difficulty of:
1. **Parsing PostScript instructions.**
2. **Mapping character codes to Unicode** (often difficult with custom font encodings).
3. **Reconstructing words** by calculating the distance between individual characters.
4. **Handling PDF encryption** and compressed streams.

By porting this to Go via AI, we have successfully transferred a massive amount of domain knowledge (thousands of hours of C# development) into the Go ecosystem in a fraction of the time.

### Performance and Reliability
During the experiment, we found that the Go port often exhibits memory management benefits typical of the Go runtime. Because the AI translated the logic directly, we avoided the "garbage collection overhead" issues sometimes found in complex .NET object graphs when processing massive 100MB+ PDF files.

---

## 🧪 Testing and Quality Assurance

A major concern with AI-generated code is "hallucinations." To mitigate this, we employed a strict **Test-Driven Transpilation** approach:
1. The C# unit test suite was ported first.
2. The core logic was ported in modules.
3. Every module was required to pass the corresponding ported test before moving to the next.
4. **Current Status:** All 500+ ported tests are passing, covering edge cases like:
    - Encrypted PDFs (AES 128/256).
    - Documents with malformed cross-reference tables.
    - Right-to-left text (Hebrew/Arabic).
    - Complex CID fonts.

---

## 🤝 Comparison: Go vs. C# Version

| Feature | C# PdfPig | PdfPigGO |
| :--- | :--- | :--- |
| **Logic Source** | Original | Ported (Qwen 3.6 27B) |
| **Manual Coding** | 100% | 0% |
| **External Dependencies** | Minimal | None (Pure Go) |
| **Speed** | Excellent | Comparable (Native Go) |
| **Memory Usage** | Moderate | Low (Go's efficient structs) |
| **Primary Use Case** | .NET Ecosystem | Go Microservices, Cloud Functions |

---

## 🛠 Development and Contributions

As this project is an experiment in AI-driven development, we have specific guidelines for contributions:

- **Bug Reports:** If you find a PDF that fails to parse, please open an issue and attach the PDF (if public).
- **Fixes:** Since the code was generated by an LLM, we prefer that fixes maintain the structural consistency of the original port.
- **Enhancements:** We are currently looking to port the "Table Extractor" and "Page Segmenter" modules next using the same AI-driven approach.

### To Run Tests:
```bash
go test ./...
```

---

## 📜 License

This project is licensed under the **MIT License** - the same as the original C# PdfPig library. This ensures maximum compatibility and freedom for both commercial and open-source use.

---

## 🧠 Credits & Acknowledgments

- **UglyToad:** The original creators of the [PdfPig](https://github.com/UglyToad/PdfPig) C# library. Without their years of hard work, this port would have no foundation.
- **Alibaba Qwen Team:** For the **Qwen3.6 27B** model, which proved that local LLMs are now capable of complex, multi-file architectural transpilation.
- **NVIDIA:** For the RTX 3090 Ti, which made local inference of a 27B model feasible for this scale of project.

---

*Disclaimer: This project is an independent experiment and is not officially affiliated with the UglyToad team. Use in production at your own discretion, though current tests show it is fully functional.*