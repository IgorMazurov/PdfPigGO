package graphics

import (
	"github.com/uglytoad/pdfpig/go/content"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// GraphicsStateOperationFactory creates graphics state operations from operator tokens and operands.
type GraphicsStateOperationFactory interface {
	Create(op *tokens.OperatorToken, operands []tokens.Token) (content.GraphicsStateOperation, error)
}
