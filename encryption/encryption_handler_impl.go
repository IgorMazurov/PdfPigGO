package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	"github.com/uglytoad/pdfpig/go/core"
	"github.com/uglytoad/pdfpig/go/crossreference"
	"github.com/uglytoad/pdfpig/go/tokens"
)

var paddingBytes = []byte{
	0x28, 0xBF, 0x4E, 0x5E,
	0x4E, 0x75, 0x8A, 0x41,
	0x64, 0x00, 0x4E, 0x56,
	0xFF, 0xFA, 0x01, 0x08,
	0x2E, 0x2E, 0x00, 0xB6,
	0xD0, 0x68, 0x3E, 0x80,
	0x2F, 0x0C, 0xA9, 0xFE,
	0x64, 0x53, 0x69, 0x7A,
}

type encryptionHandler struct {
	previouslyDecrypted map[core.IndirectReference]bool
	encryptionDictionary *EncryptionDictionary
	cryptHandler         *CryptHandler
	encryptionKey        []byte
	useAes               bool
}

var _ EncryptionHandler = (*encryptionHandler)(nil)

func NewEncryptionHandler(
	encDict *EncryptionDictionary,
	trailer *crossreference.TrailerDictionary,
	passwords []string,
) (EncryptionHandler, error) {
	if passwords == nil || len(passwords) == 0 {
		passwords = []string{""}
	}

	hasEmpty := false
	for _, p := range passwords {
		if p == "" {
			hasEmpty = true
			break
		}
	}
	if !hasEmpty {
		passwords = append([]string{""}, passwords...)
	}

	var documentIdBytes []byte
	if trailer != nil {
		identifier := trailer.Identifier()
		if len(identifier) == 2 {
			token := identifier[0]
			switch t := token.(type) {
			case *tokens.HexToken:
				documentIdBytes = t.Bytes()
			case *tokens.StringToken:
				documentIdBytes = t.GetBytes()
			default:
				dataStr := ""
				if getter, ok := any(token).(interface{ Data() string }); ok {
					dataStr = getter.Data()
				}
				documentIdBytes = core.StringAsLatin1Bytes(dataStr)
			}
		}
	}

	if encDict == nil {
		return &encryptionHandler{
			previouslyDecrypted: make(map[core.IndirectReference]bool),
		}, nil
	}

	useAes := false
	var cryptHandler *CryptHandler

	if encDict.EncryptionAlgorithmCode() == SecurityHandlerInDocument ||
		encDict.EncryptionAlgorithmCode() == SecurityHandlerInDocument256 {
		var err error
		cryptHandler, err = encDict.TryGetCryptHandler()
		if err != nil {
			return nil, NewPdfDocumentEncryptedException(
				"Document encrypted with security handler in document but no crypt dictionary found.")
		}
		if cryptHandler == nil {
			return nil, NewPdfDocumentEncryptedException(
				"Document encrypted with security handler in document but no crypt dictionary found.")
		}

		sd := cryptHandler.StreamDictionary()
		useAes = sd.Name == MethodAesV2 || sd.Name == MethodAesV3
	}

	charsetGetBytes := func(s string) []byte {
		return core.StringAsLatin1Bytes(s)
	}

	revision := encDict.Revision()
	if revision == 5 || revision == 6 {
		charsetGetBytes = func(s string) []byte {
			return []byte(s)
		}
	}

	length := calcKeyLength(encDict, cryptHandler)

	var encryptionKey []byte
	foundPassword := false

	for _, password := range passwords {
		passwordBytes := charsetGetBytes(password)

		isOwner, userPassBytes := isOwnerPassword(passwordBytes, encDict, length, documentIdBytes)
		if isOwner {
			var decryptionPasswordBytes []byte
			if revision == 5 || revision == 6 {
				decryptionPasswordBytes = passwordBytes
			} else {
				decryptionPasswordBytes = userPassBytes
			}

			isUser := false
			encryptionKey = calculateEncryptionKey(decryptionPasswordBytes, encDict, length, documentIdBytes, isUser)
			foundPassword = true
			break
		}

		isUser := isUserPassword(passwordBytes, encDict, length, documentIdBytes)
		if isUser {
			encryptionKey = calculateEncryptionKey(passwordBytes, encDict, length, documentIdBytes, true)
			foundPassword = true
			break
		}
	}

	if !foundPassword {
		return nil, NewPdfDocumentEncryptedException(
			"The document was encrypted and none of the provided passwords were the user or owner password.")
	}

	return &encryptionHandler{
		previouslyDecrypted: make(map[core.IndirectReference]bool),
		encryptionDictionary: encDict,
		cryptHandler:         cryptHandler,
		encryptionKey:        encryptionKey,
		useAes:               useAes,
	}, nil
}

func calcKeyLength(encDict *EncryptionDictionary, cryptHandler *CryptHandler) int {
	code := encDict.EncryptionAlgorithmCode()
	if code == Rc4OrAes40BitKey {
		return 5
	}

	keyLen := encDict.KeyLength()
	switch {
	case keyLen == nil && code == Rc4OrAesGreaterThan40BitKey:
		return 40 / 8
	case keyLen == nil && code == UnpublishedAlgorithm40To128BitKey:
		return 40 / 8
	case keyLen == nil && cryptHandler != nil && cryptHandler.StreamDictionary().Name == MethodAesV2:
		return 128 / 8
	case keyLen == nil && cryptHandler != nil && cryptHandler.StreamDictionary().Name == MethodAesV3:
		return 256 / 8
	default:
		if keyLen != nil {
			return *keyLen / 8
		}
		return 40 / 8
	}
}

func (h *encryptionHandler) Decrypt(reference core.IndirectReference, token tokens.Token) tokens.Token {
	if token == nil {
		panic("token cannot be nil")
	}

	decrypted, err := h.decryptInternal(reference, token)
	if err != nil {
		panic(fmt.Sprintf("The document was encrypted and decryption of a token failed. Token was: %v. Error: %v", token, err))
	}

	h.previouslyDecrypted[reference] = true
	return decrypted
}

func (h *encryptionHandler) decryptInternal(reference core.IndirectReference, token tokens.Token) (tokens.Token, error) {
	switch t := token.(type) {
	case *tokens.StreamToken:
		if h.cryptHandler != nil {
			sd := h.cryptHandler.StreamDictionary()
			if sd.IsIdentity || sd.Name == MethodNone {
				return token, nil
			}
		}

		if typeNameVal, found := t.StreamDictionary.TryGet(tokens.Type); found {
			if typeName, ok := typeNameVal.(*tokens.NameToken); ok {
				if typeName.Equals(tokens.Xref) {
					return token, nil
				}
				if !h.encryptionDictionary.EncryptMetadata() && typeName.Equals(tokens.Metadata) {
					return token, nil
				}
			}
		}

		decryptedDict, err := h.decryptInternal(reference, t.StreamDictionary)
		if err != nil {
			return nil, err
		}
		decryptedDictToken, ok := decryptedDict.(*tokens.DictionaryToken)
		if !ok {
			return nil, fmt.Errorf("expected DictionaryToken after decrypting stream dictionary")
		}

		decryptedData := h.decryptData(t.Data(), reference)
		result, err := tokens.NewStreamToken(decryptedDictToken, decryptedData)
		if err != nil {
			return nil, err
		}
		return result, nil

	case *tokens.StringToken:
		if h.cryptHandler != nil {
			strd := h.cryptHandler.StringDictionary()
			if strd.IsIdentity || strd.Name == MethodNone {
				return token, nil
			}
		}

		data := t.GetBytes()
		decrypted := h.decryptData(data, reference)
		return getStringTokenFromDecryptedData(decrypted), nil

	case *tokens.HexToken:
		data := t.Bytes()
		decrypted := h.decryptData(data, reference)
		hexStr := toHexString(decrypted)
		runes := []rune(hexStr)
		return tokens.NewHexToken(runes), nil

	case *tokens.DictionaryToken:
		if _, found := t.TryGet(tokens.Cf); found {
			return token, nil
		}

		isSigDict := false
		if typeNameVal, found := t.TryGet(tokens.Type); found {
			if typeName, ok := typeNameVal.(*tokens.NameToken); ok {
				if typeName.Equals(tokens.Sig) || typeName.Equals(tokens.DocTimeStamp) {
					isSigDict = true
				}
			}
		}

		dict := t
		for keyStr, val := range dict.Data() {
			keyName := tokens.Create(keyStr)
			if isSigDict && keyName.Equals(tokens.Contents) {
				continue
			}

			switch val.(type) {
			case *tokens.StringToken, *tokens.ArrayToken, *tokens.DictionaryToken, *tokens.HexToken:
				inner, err := h.decryptInternal(reference, val)
				if err != nil {
					return nil, err
				}
				dict = dict.With(keyName, inner)
			}
		}

		return dict, nil

	case *tokens.ArrayToken:
		result := make([]tokens.Token, len(t.Data()))
		for i, elem := range t.Data() {
			decrypted, err := h.decryptInternal(reference, elem)
			if err != nil {
				return nil, err
			}
			result[i] = decrypted
		}
		return tokens.NewArrayToken(result), nil
	}

	return token, nil
}

func getStringTokenFromDecryptedData(data []byte) *tokens.StringToken {
	if len(data) >= 2 {
		if data[0] == 0xFE && data[1] == 0xFF {
			runes := make([]rune, 0, (len(data)-2)/2)
			for i := 2; i+1 < len(data); i += 2 {
				runes = append(runes, rune(uint16(data[i])<<8|uint16(data[i+1])))
			}
			return tokens.NewStringToken(string(runes), tokens.Utf16BE)
		}
		if data[0] == 0xFF && data[1] == 0xFE {
			codeUnits := make([]uint16, 0, (len(data)-2)/2)
			for i := 2; i+1 < len(data); i += 2 {
				codeUnits = append(codeUnits, uint16(data[i])|uint16(data[i+1])<<8)
			}
			runes := make([]rune, 0)
			for j := 0; j < len(codeUnits); j++ {
				c := codeUnits[j]
				if c >= 0xD800 && c <= 0xDBFF && j+1 < len(codeUnits) {
					high := int(c-0xD800) * 1024
					low := int(codeUnits[j+1]) - 0xDC00
					runes = append(runes, rune(0x10000+high+low))
					j++
				} else {
					runes = append(runes, rune(c))
				}
			}
			return tokens.NewStringToken(string(runes), tokens.Utf16)
		}
	}
	return tokens.NewStringToken(core.BytesAsLatin1String(data), tokens.Iso88591)
}

func (h *encryptionHandler) decryptData(data []byte, reference core.IndirectReference) []byte {
	if h.useAes && len(h.encryptionKey) == 32 {
		decrypted, err := AesDecrypt(data, h.encryptionKey)
		if err != nil {
			return data
		}
		return decrypted
	}

	finalKey := h.getObjectKey(reference)

	if h.useAes {
		decrypted, err := AesDecrypt(data, finalKey)
		if err != nil {
			return data
		}
		return decrypted
	}

	return Encrypt(finalKey, data)
}

func (h *encryptionHandler) getObjectKey(reference core.IndirectReference) []byte {
	aesExtra := 0
	if h.useAes {
		aesExtra = 4
	}

	finalKey := make([]byte, len(h.encryptionKey)+5+aesExtra)
	copy(finalKey, h.encryptionKey)

	objNum := reference.ObjectNumber()
	gen := reference.Generation()

	finalKey[len(h.encryptionKey)] = byte(objNum)
	finalKey[len(h.encryptionKey)+1] = byte(objNum >> 8)
	finalKey[len(h.encryptionKey)+2] = byte(objNum >> 16)

	finalKey[len(h.encryptionKey)+3] = byte(gen)
	finalKey[len(h.encryptionKey)+4] = byte(gen >> 8)

	if h.useAes {
		finalKey[len(h.encryptionKey)+5] = 's'
		finalKey[len(h.encryptionKey)+6] = 'A'
		finalKey[len(h.encryptionKey)+7] = 'l'
		finalKey[len(h.encryptionKey)+8] = 'T'
	}

	hasher := md5.New()
	hasher.Write(finalKey)
	md5Hash := hasher.Sum(nil)

	length := len(h.encryptionKey) + 5
	if length > 16 {
		length = 16
	}

	result := make([]byte, length)
	copy(result, md5Hash[:length])

	return result
}

func calculateEncryptionKey(password []byte, encDict *EncryptionDictionary, length int, documentId []byte, isUserPassword bool) []byte {
	revision := encDict.Revision()
	if revision >= 2 && revision <= 4 {
		return calculateKeyRevisions2To4(password, encDict, length, documentId)
	}
	if revision <= 6 {
		return calculateKeyRevisions5And6(password, encDict, isUserPassword)
	}
	panic(fmt.Sprintf("PDF encrypted with unrecognized revision: %v", encDict))
}

func calculateKeyRevisions2To4(password []byte, encDict *EncryptionDictionary, length int, documentId []byte) []byte {
	passwordFull := getPaddedPassword(password)
	revision := encDict.Revision()

	hasher := md5.New()
	hasher.Write(passwordFull)
	hasher.Write(encDict.OwnerBytes())

	unsigned := uint32(encDict.UserAccessPermissions())
	hasher.Write([]byte{byte(unsigned), byte(unsigned >> 8), byte(unsigned >> 16), byte(unsigned >> 24)})
	hasher.Write(documentId)

	if revision >= 4 && !encDict.EncryptMetadata() {
		hasher.Write([]byte{0xFF, 0xFF, 0xFF, 0xFF})
	}

	if revision == 3 || revision == 4 {
		hasher.Write(nil)
		input := hasher.Sum(nil)

		for i := 0; i < 50; i++ {
				h2 := md5.New()
				h2.Write(input[:length])
				input = h2.Sum(nil)
		}

		result := make([]byte, length)
		copy(result, input[:length])
		return result
	}

	hasher.Write(nil)
	result := make([]byte, length)
	copy(result, hasher.Sum(nil)[:length])
	return result
}

func calculateKeyRevisions5And6(password []byte, encDict *EncryptionDictionary, isUserPassword bool) []byte {
	password = truncatePasswordTo127Bytes(password)

	var intermediateKey []byte
	var encryptedFileKey []byte

	if !isUserPassword {
		ownerKeySalt := make([]byte, 8)
		ownerBytes := encDict.OwnerBytes()
		copy(ownerKeySalt, ownerBytes[40:48])

		if encDict.Revision() == 6 {
			var err error
			intermediateKey, err = computeStupidIsoHash(password, ownerKeySalt, encDict.UserBytes())
			if err != nil {
				return nil
			}
		} else {
			intermediateKey = computeSha256Hash(password, ownerKeySalt, encDict.UserBytes())
		}

		encryptedFileKey = encDict.OwnerEncryptionBytes()
	} else {
		userKeySalt := make([]byte, 8)
		userBytes := encDict.UserBytes()
		copy(userKeySalt, userBytes[40:48])

		if encDict.Revision() == 6 {
			var err error
			intermediateKey, err = computeStupidIsoHash(password, userKeySalt, nil)
			if err != nil {
				return nil
			}
		} else {
			intermediateKey = computeSha256Hash(password, userKeySalt)
		}

		encryptedFileKey = encDict.UserEncryptionBytes()
	}

	iv := make([]byte, 16)

	block, err := aes.NewCipher(intermediateKey)
	if err != nil {
		panic(fmt.Sprintf("cannot create AES cipher: %v", err))
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(encryptedFileKey))
	cbc.CryptBlocks(plaintext, encryptedFileKey)

	return plaintext
}

func isUserPassword(passwordBytes []byte, encDict *EncryptionDictionary, length int, documentIdBytes []byte) bool {
	if encDict.Revision() == 5 || encDict.Revision() == 6 {
		return isUserPasswordRevision5And6(passwordBytes, encDict)
	}

	calculatedEncryptionKey := calculateKeyRevisions2To4(passwordBytes, encDict, length, documentIdBytes)

	var output []byte

	if encDict.Revision() >= 3 {
		hasher := md5.New()
		hasher.Write(paddingBytes)
		hasher.Write(documentIdBytes)
		hasher.Write(nil)
		result := hasher.Sum(nil)

		temp := Encrypt(calculatedEncryptionKey, result)

		for i := byte(1); i <= 19; i++ {
			key := make([]byte, len(calculatedEncryptionKey))
			for j, b := range calculatedEncryptionKey {
				key[j] = b ^ i
			}
			temp = Encrypt(key, temp)
		}

		output = temp
	} else {
		output = Encrypt(calculatedEncryptionKey, paddingBytes)
	}

	if encDict.Revision() >= 3 {
		userBytes := encDict.UserBytes()
		return bytes.Equal(userBytes[:16], output[:16])
	}

	return bytes.Equal(encDict.UserBytes(), output)
}

func isUserPasswordRevision5And6(passwordBytes []byte, encDict *EncryptionDictionary) bool {
	truncatedPassword := truncatePasswordTo127Bytes(passwordBytes)

	userPasswordHash := make([]byte, 32)
	userValidationSalt := make([]byte, 8)
	userBytes := encDict.UserBytes()
	copy(userPasswordHash, userBytes[:32])
	copy(userValidationSalt, userBytes[32:40])

	var result []byte
	if encDict.Revision() == 6 {
		var err error
		result, err = computeStupidIsoHash(truncatedPassword, userValidationSalt, nil)
		if err != nil {
			return false
		}
	} else {
		result = computeSha256Hash(truncatedPassword, userValidationSalt)
	}

	return bytes.Equal(result, userPasswordHash)
}

func isOwnerPassword(passwordBytes []byte, encDict *EncryptionDictionary, length int, documentIdBytes []byte) (bool, []byte) {
	if encDict.Revision() == 5 || encDict.Revision() == 6 {
		return isOwnerPasswordRevision5And6(passwordBytes, encDict), nil
	}

	paddedPassword := getPaddedPassword(passwordBytes)

	hasher := md5.New()
	hasher.Write(paddedPassword)
	hashVal := hasher.Sum(nil)

	if encDict.Revision() >= 3 {
		for i := 0; i < 50; i++ {
			h2 := md5.New()
			h2.Write(hashVal)
			hashVal = h2.Sum(nil)
		}
	}

	key := make([]byte, length)
	copy(key, hashVal[:length])

	var userPassword []byte

	if encDict.Revision() == 2 {
		userPassword = Encrypt(key, encDict.OwnerBytes())
	} else {
		output := make([]byte, len(encDict.OwnerBytes()))
		copy(output, encDict.OwnerBytes())

		for i := 0; i < 20; i++ {
			keyIter := make([]byte, len(key))
			for j, b := range key {
				keyIter[j] = b ^ byte(19-i)
			}

			if i == 0 {
				copy(output, encDict.OwnerBytes())
			}

			output = Encrypt(keyIter, output)
		}

		userPassword = output
	}

	result := isUserPassword(userPassword, encDict, length, documentIdBytes)
	return result, userPassword
}

func isOwnerPasswordRevision5And6(passwordBytes []byte, encDict *EncryptionDictionary) bool {
	truncatedPassword := truncatePasswordTo127Bytes(passwordBytes)

	ownerHash := make([]byte, 32)
	validationSalt := make([]byte, 8)
	ownerBytes := encDict.OwnerBytes()
	copy(ownerHash, ownerBytes[:32])
	copy(validationSalt, ownerBytes[32:40])

	var result []byte
	if encDict.Revision() == 6 {
		var err error
		result, err = computeStupidIsoHash(truncatedPassword, validationSalt, encDict.UserBytes())
		if err != nil {
			return false
		}
	} else {
		result = computeSha256Hash(truncatedPassword, validationSalt, encDict.UserBytes())
	}

	return bytes.Equal(result, ownerHash)
}

func computeSha256Hash(input1 []byte, input2 []byte, input3 ...[]byte) []byte {
	hasher := sha256.New()
	hasher.Write(input1)
	hasher.Write(input2)
	if len(input3) > 0 && input3[0] != nil {
		hasher.Write(input3[0])
	}
	return hasher.Sum(nil)
}

func getPaddedPassword(password []byte) []byte {
	if password == nil || len(password) == 0 {
		result := make([]byte, len(paddingBytes))
		copy(result, paddingBytes)
		return result
	}

	result := make([]byte, 32)
	passwordLen := len(password)
	if passwordLen > 32 {
		passwordLen = 32
	}

	copy(result, password[:passwordLen])

	paddingBytesNeeded := 32 - passwordLen
	if paddingBytesNeeded > 0 {
		copy(result[passwordLen:], paddingBytes[:paddingBytesNeeded])
	}

	return result
}

func truncatePasswordTo127Bytes(password []byte) []byte {
	if len(password) <= 127 {
		result := make([]byte, len(password))
		copy(result, password)
		return result
	}

	result := make([]byte, 127)
	copy(result, password[:127])
	return result
}

func computeStupidIsoHash(password []byte, salt []byte, vector []byte) ([]byte, error) {
	if vector == nil {
		vector = []byte{}
	} else if len(vector) > 0 && len(vector) < 48 {
		return nil, NewPdfDocumentEncryptedException(
			fmt.Sprintf("Vector for revision 6 owner password check (/U) is the wrong length, expected 48 bytes got %d bytes.", len(vector)))
	} else if len(vector) > 48 {
		temp := make([]byte, 48)
		copy(temp, vector[:48])
		vector = temp
	}

	password = truncatePasswordTo127Bytes(password)

	hasher := sha256.New()
	hasher.Write(password)
	hasher.Write(salt)
	hasher.Write(vector)
	input := hasher.Sum(nil)

	key := make([]byte, 16)
	copy(key, input[:16])

	iv := make([]byte, 16)
	copy(iv, input[16:32])

	var x []byte
	i := 0
	for i < 64 || (len(x) > 0 && i < int(x[len(x)-1])+32) {
		roundResult := make([]byte, 0, 64*(len(password)+len(input)+len(vector)))

		for j := 0; j < 64; j++ {
			roundResult = append(roundResult, password...)
			roundResult = append(roundResult, input...)
			if len(vector) > 0 {
				roundResult = append(roundResult, vector...)
			}
		}

		block, err := aes.NewCipher(key)
		if err != nil {
			panic(fmt.Sprintf("cannot create AES cipher: %v", err))
		}

		encryptor := cipher.NewCBCEncrypter(block, iv)
		x = make([]byte, len(roundResult))
		copy(x, roundResult)
		encryptor.CryptBlocks(x, x)

		sum := int64(0)
		for _, b := range x[:16] {
			sum += int64(b)
		}

		mod3 := sum % 3
		var nextHash hash.Hash
		switch mod3 {
		case 0:
			nextHash = sha256.New()
		case 1:
			nextHash = sha512.New384()
		case 2:
			nextHash = sha512.New()
		default:
			panic("Invalid remainder from summing first sixteen bytes of this round's hash.")
		}
		nextHash.Write(x)
		input = nextHash.Sum(nil)

		copy(key, input[:16])
		copy(iv, input[16:32])

		i++
	}

	if len(input) > 32 {
		result := make([]byte, 32)
		copy(result, input[:32])
		return result, nil
	}

	result := make([]byte, len(input))
	copy(result, input)
	return result, nil
}

func toHexString(data []byte) string {
	result := make([]byte, 0, len(data)*2)
	for _, b := range data {
		result = append(result, hexChars[b>>4], hexChars[b&0x0F])
	}
	return string(result)
}

var hexChars = []byte("0123456789ABCDEF")
