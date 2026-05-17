package annotations

import "github.com/uglytoad/pdfpig/go/content"

// Hyperlink is an alias to content.Hyperlink for backward compatibility.
// The type has been moved to content package to allow Page.GetHyperlinks() to work
// without creating an import cycle between content and annotations packages.
type Hyperlink = content.Hyperlink

// NewHyperlink is re-exported from content for backward compatibility.
var NewHyperlink = content.NewHyperlink
