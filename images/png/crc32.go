package png

// crc32Lookup is the precomputed CRC-32 lookup table using polynomial 0xEDB88320.
var crc32Lookup [256]uint32

func init() {
	polynomial := uint32(0xEDB88320)
	for i := uint32(0); i < 256; i++ {
		value := i
		for j := 0; j < 8; j++ {
			if value&1 != 0 {
				value = (value >> 1) ^ polynomial
			} else {
				value >>= 1
			}
		}
		crc32Lookup[i] = value
	}
}

// Calculate computes the CRC-32 checksum of data.
func Calculate(data []byte) uint32 {
	crc := uint32(0xFFFFFFFF)
	for i := 0; i < len(data); i++ {
		index := (crc ^ uint32(data[i])) & 0xFF
		crc = (crc >> 8) ^ crc32Lookup[index]
	}
	return crc ^ 0xFFFFFFFF
}

// CalculateTwo computes the combined CRC-32 checksum of two byte slices.
func CalculateTwo(data, data2 []byte) uint32 {
	crc := uint32(0xFFFFFFFF)
	for i := 0; i < len(data); i++ {
		index := (crc ^ uint32(data[i])) & 0xFF
		crc = (crc >> 8) ^ crc32Lookup[index]
	}
	for i := 0; i < len(data2); i++ {
		index := (crc ^ uint32(data2[i])) & 0xFF
		crc = (crc >> 8) ^ crc32Lookup[index]
	}
	return crc ^ 0xFFFFFFFF
}

// Crc32State holds the state for incremental CRC-32 computation.
type Crc32State struct {
	state uint32
}

// NewCrc32State creates a new Crc32State initialized to the starting value.
func NewCrc32State() *Crc32State {
	return &Crc32State{state: 0xFFFFFFFF}
}

// Append updates the CRC-32 state with additional data.
func (c *Crc32State) Append(data []byte) {
	for i := 0; i < len(data); i++ {
		index := (c.state ^ uint32(data[i])) & 0xFF
		c.state = (c.state >> 8) ^ crc32Lookup[index]
	}
}

// CurrentHash returns the current CRC-32 hash as a uint32.
func (c *Crc32State) CurrentHash() uint32 {
	return c.state ^ 0xFFFFFFFF
}

// Reset reinitializes the state to the starting value.
func (c *Crc32State) Reset() {
	c.state = 0xFFFFFFFF
}
