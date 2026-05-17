package encryption

import (
	"errors"
	"fmt"

	"github.com/uglytoad/pdfpig/go/tokenization"
	"github.com/uglytoad/pdfpig/go/tokens"
)

// ReadEncryptionDictionary reads and constructs an EncryptionDictionary from a
// raw PDF dictionary token. The scanner parameter is accepted for API
// compatibility but indirect reference resolution is not yet implemented.
func ReadEncryptionDictionary(
	dictionary *tokens.DictionaryToken,
	scanner tokenization.TokenScanner,
) (*EncryptionDictionary, error) {
	_ = scanner

	if dictionary == nil {
		return nil, errors.New("encryption dictionary cannot be nil")
	}

	filterTok, ok := dictionary.TryGet(tokens.Filter)
	if !ok {
		return nil, fmt.Errorf("encryption dictionary does not contain a Filter entry")
	}

	var filter string
	if nameToken, ok := filterTok.(*tokens.NameToken); ok {
		filter = nameToken.Data()
	} else if refToken, ok := filterTok.(*tokens.IndirectReferenceToken); ok {
		return nil, fmt.Errorf("indirect reference for Filter not yet resolved: %v", refToken)
	} else {
		return nil, fmt.Errorf("Filter entry is not a NameToken, got %T", filterTok)
	}

	code := Unrecognized

	if vNum := tryGetNumeric(dictionary, tokens.V); vNum != nil {
		code = EncryptionAlgorithmCode(vNum.IntVal())
	}

	var length *int
	if lengthToken := tryGetNumeric(dictionary, tokens.Length); lengthToken != nil {
		v := lengthToken.IntVal()
		length = &v
	}

	revision := 0
	if revisionToken := tryGetNumeric(dictionary, tokens.R); revisionToken != nil {
		revision = revisionToken.IntVal()
	}

	ownerBytes := tryGetBytes(dictionary, tokens.O)
	userBytes := tryGetBytes(dictionary, tokens.U)

	var access UserAccessPermissions
	if accessToken := tryGetNumeric(dictionary, tokens.P); accessToken != nil {
		access = UserAccessPermissions(accessToken.LongVal())
	}

	var ownerEncryptionBytes []byte
	var userEncryptionBytes []byte
	if revision >= 5 {
		ownerEncryptionBytes = getEncryptionBytesOrDefault(dictionary, false)
		userEncryptionBytes = getEncryptionBytesOrDefault(dictionary, true)
	}

	encryptMetadata := true
	if encryptMetaToken := tryGetBoolean(dictionary, tokens.EncryptMetaData); encryptMetaToken != nil {
		encryptMetadata = encryptMetaToken.Data()
	}

	return NewEncryptionDictionary(
		filter,
		code,
		length,
		revision,
		ownerBytes,
		userBytes,
		ownerEncryptionBytes,
		userEncryptionBytes,
		access,
		dictionary,
		encryptMetadata,
	), nil
}

func tryGetNumeric(dictionary *tokens.DictionaryToken, name *tokens.NameToken) *tokens.NumericToken {
	tok, ok := dictionary.TryGet(name)
	if !ok {
		return nil
	}
	if numeric, ok := tok.(*tokens.NumericToken); ok {
		return numeric
	}
	return nil
}

func tryGetBoolean(dictionary *tokens.DictionaryToken, name *tokens.NameToken) *tokens.BooleanToken {
	tok, ok := dictionary.TryGet(name)
	if !ok {
		return nil
	}
	if boolean, ok := tok.(*tokens.BooleanToken); ok {
		return boolean
	}
	return nil
}

func tryGetBytes(dictionary *tokens.DictionaryToken, name *tokens.NameToken) []byte {
	tok, ok := dictionary.TryGet(name)
	if !ok {
		return nil
	}

	if stringToken, ok := tok.(*tokens.StringToken); ok {
		return stringToken.GetBytes()
	}

	if hexToken, ok := tok.(*tokens.HexToken); ok {
		b := hexToken.Bytes()
		result := make([]byte, len(b))
		copy(result, b)
		return result
	}

	return nil
}

// getEncryptionBytesOrDefault retrieves the revision 5+ owner (Oe) or user (Ue)
// encryption key bytes from the dictionary. Returns nil if not found.
func getEncryptionBytesOrDefault(dictionary *tokens.DictionaryToken, isUser bool) []byte {
	var name *tokens.NameToken
	if isUser {
		name = tokens.Ue
	} else {
		name = tokens.Oe
	}

	tok, ok := dictionary.TryGet(name)
	if !ok {
		return nil
	}

	if stringToken, ok := tok.(*tokens.StringToken); ok {
		return stringToken.GetBytes()
	}

	if hexToken, ok := tok.(*tokens.HexToken); ok {
		b := hexToken.Bytes()
		result := make([]byte, len(b))
		copy(result, b)
		return result
	}

	return nil
}
