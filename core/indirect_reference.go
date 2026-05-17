package core

import (
	"fmt"
)

// IndirectReference uniquely identifies and refers to objects in a PDF file.
// The lowest 16 bits of the internal value hold the generation number (0-65535),
// the remaining bits hold the object number.
type IndirectReference struct {
	numberAndGeneration int64
}

const (
	// numberOffset is the bit offset for the object number within the combined value.
	numberOffset = 16

	// generationMask extracts the lowest 16 bits (generation number).
	generationMask int64 = (1 << numberOffset) - 1 // 0xFFFF = 65535

	// maxObjectNumber is the maximum allowed object number (~4.7 billion).
	maxObjectNumber int64 = (1<<(64-numberOffset-1)) - 1
)

// NewIndirectReference creates a new IndirectReference with the given object number and generation.
func NewIndirectReference(objectNumber int64, generation int) (IndirectReference, error) {
	if generation < 0 {
		return IndirectReference{}, fmt.Errorf("generation number must not be negative, got %d", generation)
	}

	if objectNumber < -maxObjectNumber || objectNumber > maxObjectNumber {
		return IndirectReference{}, fmt.Errorf("object number must be between %d and %d, got %d", -maxObjectNumber, maxObjectNumber, objectNumber)
	}

	return IndirectReference{
		numberAndGeneration: computeInternalHash(objectNumber, generation),
	}, nil
}

// ObjectNumber returns the positive integer object number.
func (r IndirectReference) ObjectNumber() int64 {
	return r.numberAndGeneration >> numberOffset
}

// Generation returns the non-negative integer generation number (0-65535).
func (r IndirectReference) Generation() int {
	return int(r.numberAndGeneration & generationMask)
}

// computeInternalHash packs object number and generation into a single int64.
func computeInternalHash(num int64, gen int) int64 {
	return num<<numberOffset | int64(gen&int(generationMask))
}

// Equals reports whether r and other are equal.
func (r IndirectReference) Equals(other IndirectReference) bool {
	return r.numberAndGeneration == other.numberAndGeneration
}

// HashCode returns the hash code for this reference.
func (r IndirectReference) HashCode() int64 {
	return r.numberAndGeneration
}

// String returns the string representation in PDF format: "objectNumber generation".
func (r IndirectReference) String() string {
	return fmt.Sprintf("%d %d", r.ObjectNumber(), r.Generation())
}
