package truetypeparser

import (
	"github.com/uglytoad/pdfpig/go/fonts/truetype"
)

// TrueTypeTableParser defines the contract for parsing a single TrueType table
// from raw font data into a typed table struct.
// The type parameter T is the concrete table type produced by Parse.
type TrueTypeTableParser[T any] interface {
	Parse(header truetype.TrueTypeHeaderTable, data *TrueTypeDataBytes, register *TableRegisterBuilder) (T, error)
}
