package encryption

import (
	"fmt"

	"github.com/uglytoad/pdfpig/go/tokens"
)

// PdfDocumentEncryptedException is returned when a PDF document is encrypted
// and cannot be decrypted with the provided password or without one.
type PdfDocumentEncryptedException struct {
	Message    string
	Inner      error
	Dictionary *EncryptionDictionary
}

// Error implements the error interface.
func (e *PdfDocumentEncryptedException) Error() string {
	if e.Inner != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Inner)
	}
	return e.Message
}

// Unwrap returns the inner error for errors.Is/errors.As support.
func (e *PdfDocumentEncryptedException) Unwrap() error {
	return e.Inner
}

// NewPdfDocumentEncryptedException creates a new PdfDocumentEncryptedException with the given message.
func NewPdfDocumentEncryptedException(message string) *PdfDocumentEncryptedException {
	return &PdfDocumentEncryptedException{Message: message}
}

// NewPdfDocumentEncryptedExceptionWithInner creates a new PdfDocumentEncryptedException
// with a message and an inner error.
func NewPdfDocumentEncryptedExceptionWithInner(message string, inner error) *PdfDocumentEncryptedException {
	return &PdfDocumentEncryptedException{Message: message, Inner: inner}
}

// NewPdfDocumentEncryptedExceptionWithDict creates a new PdfDocumentEncryptedException
// with a message and the encryption dictionary that caused the failure.
func NewPdfDocumentEncryptedExceptionWithDict(message string, dict *EncryptionDictionary) *PdfDocumentEncryptedException {
	return &PdfDocumentEncryptedException{Message: message, Dictionary: dict}
}

// NewPdfDocumentEncryptedExceptionWithDictAndInner creates a new PdfDocumentEncryptedException
// with a message, encryption dictionary, and inner error.
func NewPdfDocumentEncryptedExceptionWithDictAndInner(message string, dict *EncryptionDictionary, inner error) *PdfDocumentEncryptedException {
	return &PdfDocumentEncryptedException{Message: message, Dictionary: dict, Inner: inner}
}

// CryptHandler manages crypt filter dictionaries for stream and string decryption.
type CryptHandler struct {
	cryptDict        *tokens.DictionaryToken
	streamDictionary *CryptDictionary
	stringDictionary *CryptDictionary
}

// NewCryptHandler creates a new CryptHandler from the given crypt dictionary
// and name tokens for stream and string filters.
func NewCryptHandler(cryptDict *tokens.DictionaryToken, streamName, stringName *tokens.NameToken) (*CryptHandler, error) {
	if streamName == nil {
		return nil, fmt.Errorf("streamName must not be nil")
	}
	if stringName == nil {
		return nil, fmt.Errorf("stringName must not be nil")
	}
	if cryptDict == nil {
		return nil, fmt.Errorf("cryptDictionary must not be nil")
	}

	sd, err := parseCryptDictionary(cryptDict, streamName)
	if err != nil {
		return nil, err
	}
	strd, err := parseCryptDictionary(cryptDict, stringName)
	if err != nil {
		return nil, err
	}

	return &CryptHandler{
		cryptDict:        cryptDict,
		streamDictionary: sd,
		stringDictionary: strd,
	}, nil
}

// StreamDictionary returns the crypt dictionary for stream decryption.
func (h *CryptHandler) StreamDictionary() *CryptDictionary {
	return h.streamDictionary
}

// StringDictionary returns the crypt dictionary for string decryption.
func (h *CryptHandler) StringDictionary() *CryptDictionary {
	return h.stringDictionary
}

// GetNamedCryptDictionary returns the crypt dictionary identified by the given name token.
func (h *CryptHandler) GetNamedCryptDictionary(name *tokens.NameToken) (*CryptDictionary, error) {
	if name == nil {
		return nil, fmt.Errorf("name must not be nil")
	}
	return parseCryptDictionary(h.cryptDict, name)
}

func parseCryptDictionary(cryptDict *tokens.DictionaryToken, name *tokens.NameToken) (*CryptDictionary, error) {
	if name.Equals(tokens.Identity) {
		return Identity, nil
	}

	cryptDictTokenVal, ok := cryptDict.TryGet(name)
	if !ok {
		return nil, NewPdfDocumentEncryptedException(fmt.Sprintf("Could not find named crypt filter %s for decryption in crypt dictionary.", name))
	}
	cryptDictToken, ok := cryptDictTokenVal.(*tokens.DictionaryToken)
	if !ok {
		return nil, NewPdfDocumentEncryptedException(fmt.Sprintf("Could not find named crypt filter %s for decryption in crypt dictionary.", name))
	}

	if typeNameVal, found := cryptDictToken.TryGet(tokens.Type); found {
		typeName, ok := typeNameVal.(*tokens.NameToken)
		if ok && !typeName.Equals(tokens.CryptFilter) && !typeName.Equals(tokens.CryptAlgorithm) {
			return nil, NewPdfDocumentEncryptedException(fmt.Sprintf("Invalid crypt dictionary type %s for crypt filter %s: %v.", typeName, name, cryptDictToken))
		}
	}

	var cfmName *tokens.NameToken
	if cfmVal, found := cryptDictToken.TryGet(tokens.Cfm); found {
		if n, ok := cfmVal.(*tokens.NameToken); ok {
			cfmName = n
		}
	}
	if cfmName == nil {
		cfmName = tokens.None
	}

	var method Method
	if cfmName.Equals(tokens.None) {
		method = MethodNone
	} else if cfmName.Equals(tokens.V2) {
		method = MethodV2
	} else if cfmName.Equals(tokens.Aesv2) {
		method = MethodAesV2
	} else if cfmName.Equals(tokens.Aesv3) {
		method = MethodAesV3
	} else {
		return nil, NewPdfDocumentEncryptedException(fmt.Sprintf("Unrecognized CFM option for crypt filter %s: %v.", cfmName, cryptDictToken))
	}

	var eventName *tokens.NameToken
	if authVal, found := cryptDictToken.TryGet(tokens.AuthEvent); found {
		if n, ok := authVal.(*tokens.NameToken); ok {
			eventName = n
		}
	}
	if eventName == nil {
		eventName = tokens.DocOpen
	}

	var event TriggerEvent
	if eventName.Equals(tokens.DocOpen) {
		event = DocumentOpen
	} else if eventName.Equals(tokens.EfOpen) {
		event = EmbeddedFileOpen
	} else {
		return nil, NewPdfDocumentEncryptedException(fmt.Sprintf("Unrecognized AuthEvent option for crypt filter %s: %v.", eventName, cryptDictToken))
	}

	var length int
	if lenVal, found := cryptDictToken.TryGet(tokens.Length); found {
		if nt, ok := lenVal.(*tokens.NumericToken); ok {
			length = nt.IntVal()
		}
	}

	return NewCryptDictionary(method, event, length), nil
}
