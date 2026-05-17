package parts

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// TryGet attempts to get a token of the specified type T.
// If the token is already of type T, it returns it directly.
// If the token is an IndirectReferenceToken, it resolves through the scanner recursively.
func TryGet[T tokens.Token](token tokens.Token, scanner tokenization.PdfTokenScanner) (T, bool) {
	if result, ok := any(token).(T); ok {
		return result, true
	}

	ref, ok := any(token).(*tokens.IndirectReferenceToken)
	if !ok {
		var zero T
		return zero, false
	}

	var obj *tokens.ObjectToken
	func() {
		defer func() { recover() }()
		obj = scanner.Get(ref.Data())
	}()

	if obj == nil {
		var zero T
		return zero, false
	}

	data := obj.Data()
	if data == nil {
		var zero T
		return zero, false
	}

	if result, ok := any(data).(T); ok {
		return result, true
	}

	if nestedRef, ok := any(data).(*tokens.IndirectReferenceToken); ok {
		return TryGet[T](nestedRef, scanner)
	}

	var zero T
	return zero, false
}

// GetByRef gets a token of the specified type T by following an IndirectReference through the scanner.
// It handles null objects, nested indirect references, and single-element arrays.
func GetByRef[T tokens.Token](reference core.IndirectReference, scanner tokenization.PdfTokenScanner) (T, error) {
	var zero T

	guard := scanner.StackDepthGuard()
	if err := guard.Enter(); err != nil {
		return zero, err
	}
	defer guard.Exit()

	obj := scanner.Get(reference)
	if obj == nil {
		return zero, nil
	}

	data := obj.Data()
	if _, isNull := any(data).(*tokens.NullToken); isNull {
		return zero, nil
	}

	if result, ok := any(data).(T); ok {
		return result, nil
	}

	if nestedRef, ok := any(data).(*tokens.IndirectReferenceToken); ok {
		return GetByRef[T](nestedRef.Data(), scanner)
	}

	if array, ok := any(data).(*tokens.ArrayToken); ok && array.Length() == 1 {
		arrayElement := array.Get(0)

		if arrayRef, ok := any(arrayElement).(*tokens.IndirectReferenceToken); ok {
			return GetByRef[T](arrayRef.Data(), scanner)
		}

		if arrayToken, ok := any(arrayElement).(T); ok {
			return arrayToken, nil
		}
	}

	return zero, core.NewPdfDocumentFormatException(fmt.Sprintf("Could not find the object number %s with type %T instead, it was found with type %T.", reference, zero, data))
}

// GetByToken gets a token of the specified type T from any token, using stack depth guard
// to prevent excessive recursion. If the token is an IndirectReferenceToken, it resolves
// through the scanner recursively.
func GetByToken[T tokens.Token](token tokens.Token, scanner tokenization.PdfTokenScanner) (T, error) {
	var zero T

	guard := scanner.StackDepthGuard()
	if err := guard.Enter(); err != nil {
		return zero, err
	}
	defer guard.Exit()

	if result, ok := any(token).(T); ok {
		return result, nil
	}

	if ref, ok := any(token).(*tokens.IndirectReferenceToken); ok {
		return GetByRef[T](ref.Data(), scanner)
	}

	return zero, core.NewPdfDocumentFormatException(fmt.Sprintf("Could not find the object %v with type %T instead, it was found with type %T.", token, zero, token))
}
