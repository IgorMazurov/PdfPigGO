package png

// Palette stores color palette data for indexed-color PNG images.
// Each entry is stored as RGBA with alpha defaulting to 255 (fully opaque)
// until SetAlphaValues is called with tRNS chunk data.
type Palette struct {
	hasAlpha bool
	data     []byte
}

// NewPalette creates a new Palette from raw RGB data.
// The input slice must have length divisible by 3, representing consecutive RGB triplets.
// Each triplet is expanded to RGBA with alpha set to 255.
func NewPalette(data []byte) *Palette {
	p := &Palette{
		data: make([]byte, len(data)*4/3),
	}
	di := 0
	for i := 0; i < len(data); i += 3 {
		p.data[di] = data[i]
		p.data[di+1] = data[i+1]
		p.data[di+2] = data[i+2]
		p.data[di+3] = 255
		di += 4
	}
	return p
}

// HasAlphaValues reports whether alpha values have been set via SetAlphaValues.
func (p *Palette) HasAlphaValues() bool {
	return p.hasAlpha
}

// Data returns the internal RGBA byte slice.
func (p *Palette) Data() []byte {
	return p.data
}

// SetAlphaValues sets the alpha channel for each palette entry from tRNS chunk data.
func (p *Palette) SetAlphaValues(bytes []byte) {
	p.hasAlpha = true
	for i := 0; i < len(bytes); i++ {
		p.data[i*4+3] = bytes[i]
	}
}

// GetPixel returns the Pixel at the given palette index.
func (p *Palette) GetPixel(index int) Pixel {
	start := index * 4
	return NewPixel(p.data[start], p.data[start+1], p.data[start+2], p.data[start+3], false)
}
