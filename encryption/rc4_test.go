package encryption

import (
	"encoding/hex"
	"testing"
)

func TestRC4Encrypt(t *testing.T) {
	tests := []struct {
		name         string
		message      string
		keyText      string
		cipherTextHex string
	}{
		{
			name:          "Plaintext with Key",
			message:       "Plaintext",
			keyText:       "Key",
			cipherTextHex: "BBF316E8D940AF0AD3",
		},
		{
			name:          "pedia with Wiki",
			message:       "pedia",
			keyText:       "Wiki",
			cipherTextHex: "1021BF0420",
		},
		{
			name:          "Attack at dawn with Secret",
			message:       "Attack at dawn",
			keyText:       "Secret",
			cipherTextHex: "45A01F645FC35B383552544B9BF5",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.message)
			key := []byte(tc.keyText)

			result := Encrypt(key, data)

			expectedBytes, err := hex.DecodeString(tc.cipherTextHex)
			if err != nil {
				t.Fatalf("hex decode error: %v", err)
			}

			if len(result) != len(expectedBytes) {
				t.Errorf("length mismatch: expected %d, got %d", len(expectedBytes), len(result))
			}
			for i := range result {
				if i < len(expectedBytes) && result[i] != expectedBytes[i] {
					t.Errorf("byte mismatch at index %d: expected 0x%02X, got 0x%02X", i, expectedBytes[i], result[i])
					break
				}
			}

			reversed := Encrypt(key, result)

			if len(reversed) != len(data) {
				t.Errorf("reverse length mismatch: expected %d, got %d", len(data), len(reversed))
			}
			for i := range data {
				if i < len(reversed) && reversed[i] != data[i] {
					t.Errorf("reverse byte mismatch at index %d: expected 0x%02X, got 0x%02X", i, data[i], reversed[i])
					break
				}
			}
		})
	}
}
