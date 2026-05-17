package encryption

// EncryptionAlgorithmCode specifies the algorithm to be used in encrypting
// and decrypting the document.
type EncryptionAlgorithmCode byte

const (
	// Unrecognized indicates an undocumented algorithm that is no longer supported.
	Unrecognized EncryptionAlgorithmCode = iota

	// Rc4OrAes40BitKey indicates RC4 or AES encryption using a key of 40 bits.
	Rc4OrAes40BitKey

	// Rc4OrAesGreaterThan40BitKey indicates RC4 or AES encryption using a key
	// of more than 40 bits.
	Rc4OrAesGreaterThan40BitKey

	// UnpublishedAlgorithm40To128BitKey indicates an unpublished algorithm that
	// permits encryption key lengths ranging from 40 to 128 bits.
	UnpublishedAlgorithm40To128BitKey

	// SecurityHandlerInDocument indicates the security handler defines the use
	// of encryption and decryption in the document with a key length of 128 bits.
	SecurityHandlerInDocument

	// SecurityHandlerInDocument256 indicates the security handler defines the
	// use of encryption and decryption in the document with a key length of
	// 256 bits.
	SecurityHandlerInDocument256

	// UndocumentedDueToIso corresponds to revision 6 whose specification is
	// undocumented because ISO charges for access to the PDF 2 spec.
	UndocumentedDueToIso
)
