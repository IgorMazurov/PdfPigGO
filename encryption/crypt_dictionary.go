package encryption

// CryptDictionary represents a crypt filter dictionary used to specify
// how streams and strings are encrypted or decrypted in a PDF document.
type CryptDictionary struct {
	Name       Method
	Event      TriggerEvent
	Length     int
	IsIdentity bool
}

// Identity is the identity crypt dictionary that performs no encryption.
var Identity = &CryptDictionary{
	Name:       MethodNone,
	IsIdentity: true,
}

// NewCryptDictionary creates a new CryptDictionary with the given parameters.
func NewCryptDictionary(name Method, event TriggerEvent, length int) *CryptDictionary {
	return &CryptDictionary{
		Name:       name,
		Event:      event,
		Length:     length,
		IsIdentity: false,
	}
}

// Method represents the method used by the consumer application to decrypt data.
type Method byte

const (
	// MethodNone indicates the application does not decrypt data but directs
	// the input stream to the security handler for decryption.
	MethodNone Method = iota

	// MethodV2 indicates the application asks the security handler for the
	// encryption key and implicitly decrypts data using the RC4 algorithm.
	MethodV2

	// MethodAesV2 (PDF 1.6) indicates the application asks the security handler
	// for the encryption key and implicitly decrypts data using the AES algorithm
	// in Cipher Block Chaining (CBC) mode with a 16-byte block size and an
	// initialization vector that is randomly generated and placed as the first
	// 16 bytes in the stream or string.
	MethodAesV2

	// MethodAesV3 indicates the application asks the security handler for the
	// encryption key and implicitly decrypts data using the AES-256 algorithm
	// in Cipher Block Chaining (CBC) with padding mode with a 16-byte block size
	// and an initialization vector that is randomly generated and placed as the
	// first 16 bytes in the stream or string. The key size shall be 256 bits.
	MethodAesV3
)

// TriggerEvent represents the event used to trigger the authorization required
// to access encryption keys used by a crypt filter.
type TriggerEvent byte

const (
	// DocumentOpen indicates authorization is required when a document is opened.
	DocumentOpen TriggerEvent = iota

	// EmbeddedFileOpen indicates authorization is required when accessing embedded files.
	EmbeddedFileOpen
)
