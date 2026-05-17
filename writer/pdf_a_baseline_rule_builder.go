package writer

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/tokens"
	"github.com/uglytoad/pdfpig/go/writer/colors"
)

// pdfABaselineRuleBuilder mirrors the C# internal static class PdfABaselineRuleBuilder.
// The type is unexported; methods are called via the zero value.
type pdfABaselineRuleBuilder struct{}

// Obey adds the required /OutputIntents and /Metadata entries to the given catalog
// dictionary so that it satisfies the baseline requirements of the specified PDF/A standard.
func (pdfABaselineRuleBuilder) Obey(
	catalog map[*tokens.NameToken]tokens.Token,
	objectWriter func(token tokens.Token) *tokens.IndirectReferenceToken,
	documentInfo *content.DocumentInformation,
	standard PdfAStandard,
	version float64,
	xmpMetadata *string,
) {
	outputIntents := buildOutputIntentsArray(objectWriter)
	catalog[tokens.OutputIntents] = outputIntents

	xmpStream := GenerateXmpStream(documentInfo, version, standard, xmpMetadata)
	xmpRef := objectWriter(xmpStream)
	catalog[tokens.Metadata] = xmpRef
}

// buildOutputIntentsArray creates the /OutputIntents array token required for PDF/A.
func buildOutputIntentsArray(
	objectWriter func(token tokens.Token) *tokens.IndirectReferenceToken,
) tokens.Token {
	return colors.GetOutputIntentsArray(
		func(token tokens.Token) *tokens.IndirectReferenceToken {
			return objectWriter(token)
		},
		DataCompressorCompressBytes,
	)
}


