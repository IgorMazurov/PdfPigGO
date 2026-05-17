package encryption

// Encrypt encrypts or decrypts data using the RC4 stream cipher algorithm.
// The same function performs both encryption and decryption since RC4 is symmetric:
// calling Encrypt with the ciphertext as input returns the original plaintext.
func Encrypt(key []byte, data []byte) []byte {
	// Key-scheduling algorithm (KSA).
	s := make([]byte, 256)
	for i := 0; i < 256; i++ {
		s[i] = byte(i)
	}

	j := 0
	for i := 0; i < 256; i++ {
		j = (j + int(s[i]) + int(key[i%len(key)])) % 256

		temp := s[i]
		s[i] = s[j]
		s[j] = temp
	}

	result := make([]byte, len(data))

	// Pseudo-random generation algorithm (PRGA).
	j = 0
	i := 0
	for step := 0; step < len(data); step++ {
		i = (i + 1) % 256
		j = (j + int(s[i])) % 256

		temp := s[i]
		s[i] = s[j]
		s[j] = temp

		k := s[(int(s[i])+int(s[j]))%256]
		result[step] = data[step] ^ k
	}

	return result
}
