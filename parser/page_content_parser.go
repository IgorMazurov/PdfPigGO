package parser

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/graphics"
	"github.com/uglytoad/pdfpig/go/logging"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// PageContentParserInterface parses page content streams into graphics state operations.
type PageContentParser interface {
	// Parse converts the raw input bytes into a list of graphics state operations
	// for the given page number.
	Parse(pageNumber int, inputBytes core.InputBytes, log logging.Log) []content.GraphicsStateOperation
}

// pageContentParserImpl provides functionality to parse the content of a PDF page,
// extracting graphics state operations from the input data. This struct is responsible
// for interpreting the PDF content stream and converting it into a collection of operations.
type pageContentParserImpl struct {
	operationFactory  graphics.GraphicsStateOperationFactory
	stackDepthGuard   *core.StackDepthGuard
	useLenientParsing bool
}

// NewPageContentParser creates a new PageContentParser with the given operation factory,
// stack depth guard, and optional lenient parsing flag.
func NewPageContentParser(
	operationFactory graphics.GraphicsStateOperationFactory,
	stackDepthGuard *core.StackDepthGuard,
	useLenientParsing bool,
) PageContentParser {
	return &pageContentParserImpl{
		operationFactory:  operationFactory,
		stackDepthGuard:   stackDepthGuard,
		useLenientParsing: useLenientParsing,
	}
}

// Parse parses the content of a PDF page and extracts a collection of graphics state operations.
func (p *pageContentParserImpl) Parse(
	pageNumber int,
	inputBytes core.InputBytes,
	log logging.Log,
) []content.GraphicsStateOperation {
	scanner := tokenization.NewCoreTokenScanner(
		inputBytes, false, p.stackDepthGuard,
		tokenization.ScannerScopeNone, nil,
		p.useLenientParsing, true,
	)

	precedingTokens := make([]tokens.Token, 0, 16)
	graphicsStateOperations := make([]content.GraphicsStateOperation, 0, 256)

	var lastEndImageOffset *int64

	for scanner.Advance() {
		token := scanner.Current()

		if inlineImageData, ok := token.(*tokens.InlineImageDataToken); ok {
			dictionary := make(map[*tokens.NameToken]tokens.Token)

			for i := 0; i < len(precedingTokens)-1; i++ {
				t := precedingTokens[i]
				nameToken, ok := t.(*tokens.NameToken)
				if !ok {
					continue
				}

				i++

				dictionary[nameToken] = precedingTokens[i]
			}

			graphicsStateOperations = append(graphicsStateOperations, graphics.InstanceBeginInlineImage)
			graphicsStateOperations = append(graphicsStateOperations, &graphics.BeginInlineImageData{Dictionary: dictionary})
			graphicsStateOperations = append(graphicsStateOperations, graphics.NewEndInlineImage(inlineImageData.Data()))

			offset := scanner.CurrentPosition() - 2
			lastEndImageOffset = &offset

			precedingTokens = precedingTokens[:0]
		} else if op, ok := token.(*tokens.OperatorToken); ok {
			if op.Data() == "EI" {
				var lastOperation content.GraphicsStateOperation
				if len(graphicsStateOperations) > 0 {
					lastOperation = graphicsStateOperations[len(graphicsStateOperations)-1]
				}

				if lastEndImageOffset == nil || lastOperation == nil {
					panic(core.NewPdfDocumentFormatException(
						fmt.Sprintf("Encountered End Image token outside an inline image on page %d at offset in content: %d.",
							pageNumber, scanner.CurrentPosition()),
					))
				}

				lastEndImage, ok := lastOperation.(*graphics.EndInlineImage)
				if !ok {
					panic(core.NewPdfDocumentFormatException(
						fmt.Sprintf("Encountered End Image token outside an inline image on page %d at offset in content: %d.",
							pageNumber, scanner.CurrentPosition()),
					))
				}

				actualEndImageOffset := scanner.CurrentPosition() - 3

				log.Warn(fmt.Sprintf("End inline image (EI) encountered after previous EI, attempting recovery at %d.", actualEndImageOffset))

				gap := int(actualEndImageOffset - *lastEndImageOffset)

				from := inputBytes.CurrentOffset()
				inputBytes.Seek(*lastEndImageOffset, 0) // io.SeekStart

				missingData := make([]byte, gap)
				read, _ := inputBytes.Read(missingData)
				if read != gap {
					panic(fmt.Errorf("Failed to read expected buffer length %d on page %d when reading inline image at offset in content: %d.",
						gap, pageNumber, *lastEndImageOffset))
				}

				fullData := append(lastEndImage.ImageData, missingData...)
				idx := len(graphicsStateOperations) - 1
				graphicsStateOperations[idx] = graphics.NewEndInlineImage(fullData)

				lastEndImageOffset = &actualEndImageOffset

				inputBytes.Seek(from, 0) // io.SeekStart
			} else {
				var operation content.GraphicsStateOperation
				var err error
				func() {
					defer func() {
						if r := recover(); r != nil {
							if e, ok := r.(error); ok {
								err = e
							} else {
								err = fmt.Errorf("%v", r)
							}
						}
					}()
					operation, err = p.operationFactory.Create(op, precedingTokens)
				}()
				if err != nil {
					log.ErrorWithException(
						fmt.Sprintf("Failed reading operation at offset %d for page %d, data: '%s'",
							inputBytes.CurrentOffset(), pageNumber, op.Data()),
						err,
					)
					if _, idx := tryGetLastEndImage(graphicsStateOperations); idx != -1 || p.useLenientParsing {
						operation = nil
					} else {
						panic(err)
					}
				}

				if operation != nil {
					graphicsStateOperations = append(graphicsStateOperations, operation)
				} else if len(graphicsStateOperations) > 0 {
					if prevEndInlineImage, index := tryGetLastEndImage(graphicsStateOperations); prevEndInlineImage != nil && lastEndImageOffset != nil {
						log.Warn(fmt.Sprintf("Operator %s was not understood following end of inline image data at %d, attempting recovery.",
							op.Data(), *lastEndImageOffset))

						nextByteSet := scanner.RecoverFromIncorrectEndImage(*lastEndImageOffset)
						graphicsStateOperations = append(graphicsStateOperations[:index], graphicsStateOperations[index+1:]...)
						newData := append(prevEndInlineImage.ImageData, nextByteSet...)
						graphicsStateOperations = append(graphicsStateOperations, graphics.NewEndInlineImage(newData))
						offset := scanner.CurrentPosition() - 3
						lastEndImageOffset = &offset
					} else if op.Data() == "inf" {
						precedingTokens = append(precedingTokens, tokens.Zero)
						continue
					} else {
						log.Warn(fmt.Sprintf("Operator which was not understood encountered. Values was %s. Ignoring.", op.Data()))
					}
				}
			}

			precedingTokens = precedingTokens[:0]
		} else if _, ok := token.(*tokens.CommentToken); ok {
			// Comments are ignored.
		} else {
			precedingTokens = append(precedingTokens, token)
		}
	}

	return graphicsStateOperations
}

// tryGetLastEndImage searches backwards through the operations list for the most recent
// EndInlineImage operation, stopping at EndText or BeginInlineImageData boundaries.
func tryGetLastEndImage(graphicsStateOperations []content.GraphicsStateOperation) (*graphics.EndInlineImage, int) {
	if len(graphicsStateOperations) == 0 {
		return nil, -1
	}

	for i := len(graphicsStateOperations) - 1; i >= 0; i-- {
		last := graphicsStateOperations[i]

		if ei, ok := last.(*graphics.EndInlineImage); ok {
			return ei, i
		}

		if _, ok := last.(graphics.EndText); ok {
			break
		}
		if _, ok := last.(graphics.BeginInlineImageData); ok {
			break
		}
	}

	return nil, -1
}
