package lzw

import (
	"testing"
)

func TestBitStream_GetNumbers(t *testing.T) {
	data := []byte{
		0b00101001,
		0b11011100,
		0b01000110,
		0b11111011,
		0b00101010,
		0b11010111,
		0b10010001,
		0b11011011,
		0b11110000,
		0b00010111,
		0b10101011,
	}

	bitStream := NewBitStream(data)

	first, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	third, err := bitStream.Get(11)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fourth, err := bitStream.Get(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fifth, err := bitStream.Get(17)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 0b001010011 {
		t.Errorf("first = %d, want %d", first, 0b001010011)
	}
	if second != 0b101110001 {
		t.Errorf("second = %d, want %d", second, 0b101110001)
	}
	if third != 0b00011011111 {
		t.Errorf("third = %d, want %d", third, 0b00011011111)
	}
	if fourth != 0b01100 {
		t.Errorf("fourth = %d, want %d", fourth, 0b01100)
	}
	if fifth != 0b10101011010111100 {
		t.Errorf("fifth = %d, want %d", fifth, 0b10101011010111100)
	}
}

func TestBitStream_GetNumbersCrossingBoundaries(t *testing.T) {
	data := []byte{
		0b00101001,
		0b11011100,
		0b01000110,
		0b11111011,
		0b00101010,
		0b11010111,
		0b10010001,
		0b11011011,
		0b11110000,
		0b00010111,
		0b10101011,
	}

	bitStream := NewBitStream(data)

	first, err := bitStream.Get(13)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := bitStream.Get(15)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	third, err := bitStream.Get(13)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first != 0b0010100111011 {
		t.Errorf("first = %d, want %d", first, 0b0010100111011)
	}
	if second != 0b100010001101111 {
		t.Errorf("second = %d, want %d", second, 0b100010001101111)
	}
	if third != 0b1011001010101 {
		t.Errorf("third = %d, want %d", third, 0b1011001010101)
	}
}

func TestBitStream_GetNumbersUntilOffsetResets(t *testing.T) {
	data := []byte{
		0b00101001,
		0b11011100,
		0b01000110,
		0b11111011,
		0b00101010,
		0b11010111,
		0b10010001,
		0b11011011,
		0b11110000,
		0b00010111,
		0b10101011,
	}

	bitStream := NewBitStream(data)

	first, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	third, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fourth, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fifth, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sixth, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	seventh, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	eighth, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ninth, err := bitStream.Get(9)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	end, err := bitStream.Get(7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expecteds := []int32{
		0b001010011,
		0b101110001,
		0b000110111,
		0b110110010,
		0b101011010,
		0b111100100,
		0b011101101,
		0b111110000,
		0b000101111,
	}

	gots := []int32{first, second, third, fourth, fifth, sixth, seventh, eighth, ninth}

	for i, got := range gots {
		if got != expecteds[i] {
			t.Errorf("value[%d] = %d, want %d", i+1, got, expecteds[i])
		}
	}

	wantEnd := int32(0b0101011)
	if end != wantEnd {
		t.Errorf("end = %d, want %d", end, wantEnd)
	}
}
