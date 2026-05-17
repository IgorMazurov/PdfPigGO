package util

import (
	"errors"
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/parser/parts"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

var (
	errNilDictionary = errors.New("dictionary cannot be nil")
	errNilArray      = errors.New("array cannot be nil")
)

// GetObjectOrDefault retrieves the entry with a given name from the dictionary.
// Returns nil if the key is not found.
func GetObjectOrDefault(dict *tokens.DictionaryToken, name *tokens.NameToken) tokens.Token {
	if token, ok := dict.TryGet(name); ok {
		return token
	}
	return nil
}

// GetObjectOrDefaultTwoKeys retrieves the entry with a given name from the dictionary,
// falling back to a second key if the first is not found. Returns nil if neither
// key exists.
func GetObjectOrDefaultTwoKeys(dict *tokens.DictionaryToken, first, second *tokens.NameToken) tokens.Token {
	if token, ok := dict.TryGet(first); ok {
		return token
	}
	if token, ok := dict.TryGet(second); ok {
		return token
	}
	return nil
}

// GetInt retrieves an integer value from the dictionary by key name.
func GetInt(dict *tokens.DictionaryToken, name *tokens.NameToken) (int, error) {
	if dict == nil {
		return 0, fmt.Errorf("get int: %w", errNilDictionary)
	}

	numeric, ok := GetObjectOrDefault(dict, name).(*tokens.NumericToken)
	if !ok || numeric == nil {
		return 0, core.NewPdfDocumentFormatException(fmt.Sprintf("the dictionary did not contain a number with the key %v. Dictionary way: %v", name, dict))
	}

	return numeric.IntVal(), nil
}

// GetIntOrDefault retrieves an integer value from the dictionary by key name,
// returning defaultValue if the key is not found or is not a numeric token.
func GetIntOrDefault(dict *tokens.DictionaryToken, name *tokens.NameToken, defaultValue int) (int, error) {
	if dict == nil {
		return 0, fmt.Errorf("get int or default: %w", errNilDictionary)
	}

	if numeric, ok := GetObjectOrDefault(dict, name).(*tokens.NumericToken); ok && numeric != nil {
		return numeric.IntVal(), nil
	}

	return defaultValue, nil
}

// GetIntOrDefaultTwoKeys retrieves an integer value from the dictionary using two
// candidate keys, returning defaultValue if neither yields a numeric token.
func GetIntOrDefaultTwoKeys(dict *tokens.DictionaryToken, first, second *tokens.NameToken, defaultValue int) (int, error) {
	if dict == nil {
		return 0, fmt.Errorf("get int or default two keys: %w", errNilDictionary)
	}

	if numeric, ok := GetObjectOrDefaultTwoKeys(dict, first, second).(*tokens.NumericToken); ok && numeric != nil {
		return numeric.IntVal(), nil
	}

	return defaultValue, nil
}

// GetLongOrDefault retrieves a long value from the dictionary by key name.
// Returns (value, true) if found and is numeric, or (0, false) otherwise.
func GetLongOrDefault(dict *tokens.DictionaryToken, name *tokens.NameToken) (int64, bool, error) {
	if dict == nil {
		return 0, false, fmt.Errorf("get long or default: %w", errNilDictionary)
	}

	if numeric, ok := GetObjectOrDefault(dict, name).(*tokens.NumericToken); ok && numeric != nil {
		return numeric.LongVal(), true, nil
	}

	return 0, false, nil
}

// GetBooleanOrDefault retrieves a boolean value from the dictionary by key name,
// returning defaultValue if the key is not found or is not a boolean token.
func GetBooleanOrDefault(dict *tokens.DictionaryToken, name *tokens.NameToken, defaultValue bool) (bool, error) {
	if dict == nil {
		return false, fmt.Errorf("get boolean or default: %w", errNilDictionary)
	}

	if boolean, ok := GetObjectOrDefault(dict, name).(*tokens.BooleanToken); ok && boolean != nil {
		return boolean.Data(), nil
	}

	return defaultValue, nil
}

// GetNameOrDefault retrieves a name token from the dictionary by key name.
// Returns nil if not found or not a name token.
func GetNameOrDefault(dict *tokens.DictionaryToken, name *tokens.NameToken) (*tokens.NameToken, error) {
	if dict == nil {
		return nil, fmt.Errorf("get name or default: %w", errNilDictionary)
	}

	nameToken, _ := GetObjectOrDefault(dict, name).(*tokens.NameToken)
	return nameToken, nil
}

// TryGetOptionalTokenDirect attempts to get a token of type T from the dictionary
// by resolving indirect references through the scanner. Returns true if found.
func TryGetOptionalTokenDirect[T tokens.Token](dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (T, bool) {
	var zero T

	token, ok := dict.TryGet(name)
	if !ok {
		return zero, false
	}

	result, found := parts.TryGet[T](token, scanner)
	if found {
		return result, true
	}

	return zero, false
}

// TryGetOptionalStringDirect attempts to get a string value from the dictionary by
// resolving indirect references through the scanner. Tries both StringToken and
// HexToken formats. Returns true if found.
func TryGetOptionalStringDirect(dict *tokens.DictionaryToken, name *tokens.NameToken, scanner tokenization.PdfTokenScanner) (string, bool) {
	if stringToken, ok := TryGetOptionalTokenDirect[*tokens.StringToken](dict, name, scanner); ok {
		return stringToken.Data(), true
	}

	if hexToken, ok := TryGetOptionalTokenDirect[*tokens.HexToken](dict, name, scanner); ok {
		return hexToken.Data(), true
	}

	return "", false
}

// GetNumeric retrieves the numeric token at the given index from an array.
func GetNumeric(arr *tokens.ArrayToken, index int) (*tokens.NumericToken, error) {
	if arr == nil {
		return nil, fmt.Errorf("get numeric: %w", errNilArray)
	}

	data := arr.Data()
	if index < 0 || index >= len(data) {
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("cannot index into array at index %d. Array was: %v.", index, arr))
	}

	if numeric, ok := data[index].(*tokens.NumericToken); ok {
		return numeric, nil
	}

	return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("the array did not contain a number at index %d. Array was: %v.", index, arr))
}

// ToRectangle converts an array token into a PdfRectangle using floating-point
// coordinates. The array must have at least 4 elements; extra elements are ignored.
func ToRectangle(arr *tokens.ArrayToken, scanner tokenization.PdfTokenScanner) (*core.PdfRectangle, error) {
	if arr == nil {
		return nil, fmt.Errorf("to rectangle: %w", errNilArray)
	}

	data := arr.Data()
	if len(data) < 4 {
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("cannot convert array to rectangle, expected 4 values instead got: %d.", len(data)))
	}

	n0, err := parts.GetByToken[*tokens.NumericToken](data[0], scanner)
	if err != nil {
		return nil, err
	}
	n1, err := parts.GetByToken[*tokens.NumericToken](data[1], scanner)
	if err != nil {
		return nil, err
	}
	n2, err := parts.GetByToken[*tokens.NumericToken](data[2], scanner)
	if err != nil {
		return nil, err
	}
	n3, err := parts.GetByToken[*tokens.NumericToken](data[3], scanner)
	if err != nil {
		return nil, err
	}

	rect := core.NewPdfRectangleFloat(n0.DoubleVal(), n1.DoubleVal(), n2.DoubleVal(), n3.DoubleVal())
	return &rect, nil
}

// ToIntRectangle converts an array token into a PdfRectangle using integer
// coordinates. The array must have exactly 4 elements.
func ToIntRectangle(arr *tokens.ArrayToken, scanner tokenization.PdfTokenScanner) (*core.PdfRectangle, error) {
	if arr == nil {
		return nil, fmt.Errorf("to int rectangle: %w", errNilArray)
	}

	data := arr.Data()
	if len(data) != 4 {
		return nil, core.NewPdfDocumentFormatException(fmt.Sprintf("cannot convert array to rectangle, expected 4 values instead got: %d.", len(data)))
	}

	n0, err := parts.GetByToken[*tokens.NumericToken](data[0], scanner)
	if err != nil {
		return nil, err
	}
	n1, err := parts.GetByToken[*tokens.NumericToken](data[1], scanner)
	if err != nil {
		return nil, err
	}
	n2, err := parts.GetByToken[*tokens.NumericToken](data[2], scanner)
	if err != nil {
		return nil, err
	}
	n3, err := parts.GetByToken[*tokens.NumericToken](data[3], scanner)
	if err != nil {
		return nil, err
	}

	rect := core.NewPdfRectangleFromInt(n0.IntVal(), n1.IntVal(), n2.IntVal(), n3.IntVal())
	return &rect, nil
}
