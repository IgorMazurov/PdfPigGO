package encryption

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// EncryptionDictionary holds the parameters extracted from a PDF encryption
// dictionary entry in the document catalog.
type EncryptionDictionary struct {
	filter                string
	encryptionAlgorithmCode EncryptionAlgorithmCode
	keyLength             *int
	revision              int
	ownerBytes            []byte
	userBytes             []byte
	ownerEncryptionBytes  []byte
	userEncryptionBytes   []byte
	userAccessPermissions UserAccessPermissions
	dictionary            *tokens.DictionaryToken
	encryptMetadata       bool
}

// NewEncryptionDictionary creates a new EncryptionDictionary.
func NewEncryptionDictionary(
	filter string,
	encryptionAlgorithmCode EncryptionAlgorithmCode,
	keyLength *int,
	revision int,
	ownerBytes []byte,
	userBytes []byte,
	ownerEncryptionBytes []byte,
	userEncryptionBytes []byte,
	userAccessPermissions UserAccessPermissions,
	dictionary *tokens.DictionaryToken,
	encryptMetadata bool,
) *EncryptionDictionary {
	return &EncryptionDictionary{
		filter:                filter,
		encryptionAlgorithmCode: encryptionAlgorithmCode,
		keyLength:             keyLength,
		revision:              revision,
		ownerBytes:            ownerBytes,
		userBytes:             userBytes,
		ownerEncryptionBytes:  ownerEncryptionBytes,
		userEncryptionBytes:   userEncryptionBytes,
		userAccessPermissions: userAccessPermissions,
		dictionary:            dictionary,
		encryptMetadata:       encryptMetadata,
	}
}

// Filter returns the encryption filter name.
func (e *EncryptionDictionary) Filter() string {
	return e.filter
}

// EncryptionAlgorithmCode returns the algorithm code used for encryption.
func (e *EncryptionDictionary) EncryptionAlgorithmCode() EncryptionAlgorithmCode {
	return e.encryptionAlgorithmCode
}

// KeyLength returns a pointer to the key length, or nil if not specified.
func (e *EncryptionDictionary) KeyLength() *int {
	return e.keyLength
}

// Revision returns the revision number of the encryption handler.
func (e *EncryptionDictionary) Revision() int {
	return e.revision
}

// OwnerBytes returns the owner password bytes, or nil if not present.
func (e *EncryptionDictionary) OwnerBytes() []byte {
	return e.ownerBytes
}

// UserBytes returns the user password bytes, or nil if not present.
func (e *EncryptionDictionary) UserBytes() []byte {
	return e.userBytes
}

// OwnerEncryptionBytes returns the 32-byte owner encryption key for revision 5+.
func (e *EncryptionDictionary) OwnerEncryptionBytes() []byte {
	return e.ownerEncryptionBytes
}

// UserEncryptionBytes returns the 32-byte user encryption key for revision 5+.
func (e *EncryptionDictionary) UserEncryptionBytes() []byte {
	return e.userEncryptionBytes
}

// UserAccessPermissions returns the user access permission flags.
func (e *EncryptionDictionary) UserAccessPermissions() UserAccessPermissions {
	return e.userAccessPermissions
}

// IsStandardFilter reports whether the filter is "Standard" (case-insensitive).
func (e *EncryptionDictionary) IsStandardFilter() bool {
	return len(e.filter) == 8 &&
		(e.filter[0] == 'S' || e.filter[0] == 's') &&
		(e.filter[1] == 't' || e.filter[1] == 'T') &&
		(e.filter[2] == 'a' || e.filter[2] == 'A') &&
		(e.filter[3] == 'n' || e.filter[3] == 'N') &&
		(e.filter[4] == 'd' || e.filter[4] == 'D') &&
		(e.filter[5] == 'a' || e.filter[5] == 'A') &&
		(e.filter[6] == 'r' || e.filter[6] == 'R') &&
		(e.filter[7] == 'd' || e.filter[7] == 'D')
}

// EncryptMetadata reports whether document metadata should be encrypted.
func (e *EncryptionDictionary) EncryptMetadata() bool {
	return e.encryptMetadata
}

// Dictionary returns the underlying PDF dictionary token.
func (e *EncryptionDictionary) Dictionary() *tokens.DictionaryToken {
	return e.dictionary
}

// TryGetCryptHandler attempts to build a CryptHandler from this encryption
// dictionary. It succeeds only when the algorithm code indicates a security
// handler defined within the document and the required "CF" entry is present.
func (e *EncryptionDictionary) TryGetCryptHandler() (*CryptHandler, error) {
	if e.encryptionAlgorithmCode != SecurityHandlerInDocument &&
		e.encryptionAlgorithmCode != SecurityHandlerInDocument256 {
		return nil, nil
	}

	cryptFilterTok, ok := e.dictionary.TryGet(tokens.Cf)
	if !ok {
		return nil, nil
	}

	cryptFilterDict, ok := cryptFilterTok.(*tokens.DictionaryToken)
	if !ok {
		return nil, nil
	}

	streamFilterName := tokens.Identity
	if tok, found := e.dictionary.TryGet(tokens.StmF); found {
		if n, ok := tok.(*tokens.NameToken); ok {
			streamFilterName = n
		}
	}

	stringFilterName := tokens.Identity
	if tok, found := e.dictionary.TryGet(tokens.StrF); found {
		if n, ok := tok.(*tokens.NameToken); ok {
			stringFilterName = n
		}
	}

	if !streamFilterName.Equals(tokens.Identity) {
		if _, found := cryptFilterDict.TryGet(streamFilterName); !found {
			return nil, NewPdfDocumentEncryptedException(
				fmt.Sprintf("Stream filter %s not found in crypt dictionary: %v.", streamFilterName, cryptFilterDict))
		}
	}

	if !stringFilterName.Equals(tokens.Identity) {
		if _, found := cryptFilterDict.TryGet(stringFilterName); !found {
			return nil, NewPdfDocumentEncryptedException(
				fmt.Sprintf("String filter %s not found in crypt dictionary: %v.", stringFilterName, cryptFilterDict))
		}
	}

	handler, err := NewCryptHandler(cryptFilterDict, streamFilterName, stringFilterName)
	if err != nil {
		return nil, err
	}

	return handler, nil
}
