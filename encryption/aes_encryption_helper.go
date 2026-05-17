package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
)

var errDataTooShort = errors.New("encryption: data too short for AES decryption")

// AesDecrypt decrypts AES-CBC encrypted data with PKCS7 padding.
// The first 16 bytes of data are used as the IV, followed by the ciphertext.
// Returns an empty slice if data length is zero or at most the IV size.
func AesDecrypt(data []byte, finalKey []byte) ([]byte, error) {
	if len(data) == 0 {
		return data, nil
	}

	const ivLen = 16

	if len(data) <= ivLen {
		return []byte{}, nil
	}

	iv := make([]byte, ivLen)
	copy(iv, data[:ivLen])

	block, err := aes.NewCipher(finalKey)
	if err != nil {
		return nil, fmt.Errorf("encryption: cannot create AES cipher: %w", err)
	}

	cbc := cipher.NewCBCDecrypter(block, iv)

	plaintext := make([]byte, len(data)-ivLen)
	cbc.CryptBlocks(plaintext, data[ivLen:])

	return unpadPKCS7(plaintext), nil
}

// unpadPKCS7 removes PKCS7 padding from the decrypted data.
func unpadPKCS7(data []byte) []byte {
	if len(data) == 0 {
		return data
	}

	padding := int(data[len(data)-1])

	if padding > len(data) || padding == 0 {
		return data
	}

	for i := 0; i < padding; i++ {
		if data[len(data)-1-i] != byte(padding) {
			return data
		}
	}

	return data[:len(data)-padding]
}
