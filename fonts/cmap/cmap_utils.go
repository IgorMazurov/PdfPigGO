package cmap

// ToInt converts a big-endian byte slice to an integer.
// Each byte is shifted into the result, processing from most significant to least significant.
func ToInt(data []byte) int {
	code := 0
	for i := 0; i < len(data); i++ {
		code <<= 8
		code |= int(data[i] & 0xff)
	}
	return code
}

// PutAll copies all key-value pairs from source into target.
func PutAll[K comparable, V any](target map[K]V, source map[K]V) {
	for k, v := range source {
		target[k] = v
	}
}
