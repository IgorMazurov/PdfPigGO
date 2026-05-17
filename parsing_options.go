package pdfpig

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/filters"
	"github.com/uglytoad/pdfpig/go/logging"
)

// ParsingOptions configures options used by the parser when reading PDF documents.
type ParsingOptions struct {
	// ClipPaths indicates whether the parser should apply clipping to paths.
	// Defaults to false. Bezier curves will be transformed into polylines if clipping is set to true.
	ClipPaths bool

	// UseLenientParsing indicates whether the parser should ignore issues where
	// the document does not conform to the PDF specification.
	UseLenientParsing bool

	// Logger is used to record messages raised by the parsing process.
	Logger logging.Log

	// Password is the password to use to open the document if it is encrypted.
	// If you need to supply multiple passwords to test against, use Passwords.
	// The value of Password will be included in the list to test against.
	Password string

	// Passwords are all passwords to try when opening this document.
	// Will include any values set for Password.
	Passwords []string

	// SkipMissingFonts skips extracting content where the font could not be found,
	// which will result in some letters being skipped/missed but will prevent the
	// library throwing where the source PDF has some corrupted text. Also skips
	// XObjects like forms and images when missing.
	SkipMissingFonts bool

	// MaxStackDepth is the maximum allowed stack depth. This property can be used
	// to limit the depth of recursive or nested operations to prevent stack overflows
	// or excessive resource usage.
	MaxStackDepth int

	// FilterProvider is the filter provider to use while parsing the document.
	// The DefaultFilterProvider will be used if set to nil.
	FilterProvider filters.FilterProvider
}

// LenientParsingOff is a default ParsingOptions with UseLenientParsing set to false.
var LenientParsingOff = content.ParsingOptions{
	UseLenientParsing: false,
	Logger:            logging.NoopLog,
	Passwords:         []string{},
	MaxStackDepth:     256,
}
