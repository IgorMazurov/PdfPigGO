package core

import (
	"strings"
	"testing"
)

func TestIndirectReferenceSetsProperties(t *testing.T) {
	ref, err := NewIndirectReference(129, 45)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ref.ObjectNumber() != 129 {
		t.Errorf("expected ObjectNumber 129, got %d", ref.ObjectNumber())
	}
	if ref.Generation() != 45 {
		t.Errorf("expected Generation 45, got %d", ref.Generation())
	}
}

func TestIndirectReferenceToString(t *testing.T) {
	ref, err := NewIndirectReference(130, 70)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "130 70"
	if ref.String() != expected {
		t.Errorf("expected %q, got %q", expected, ref.String())
	}
}

func TestTwoIndirectReferencesEqual(t *testing.T) {
	ref1, _ := NewIndirectReference(1574, 690)
	ref2, _ := NewIndirectReference(1574, 690)

	if !ref1.Equals(ref2) {
		t.Error("expected references to be equal")
	}
}

func TestTwoIndirectReferencesNotEqual(t *testing.T) {
	ref1, _ := NewIndirectReference(1574, 690)
	ref2, _ := NewIndirectReference(12, 0)

	if ref1.Equals(ref2) {
		t.Error("expected references not to be equal")
	}
}

func TestIndirectReferenceValidation(t *testing.T) {
	tests := []struct {
		name         string
		objectNumber int64
		generation   int
		wantErr      bool
		errContains  string
	}{
		{
			name:         "positive object number",
			objectNumber: 1574,
			generation:   690,
			wantErr:      false,
		},
		{
			name:         "negative object number",
			objectNumber: -1574,
			generation:   690,
			wantErr:      false,
		},
		{
			name:         "large positive object number",
			objectNumber: 58949797283757,
			generation:   16,
			wantErr:      false,
		},
		{
			name:         "max object number with max ushort generation",
			objectNumber: 140737488355327,
			generation:   65535,
			wantErr:      false,
		},
		{
			name:         "min negative object number",
			objectNumber: -140737488355327,
			generation:   65535,
			wantErr:      false,
		},
		{
			name:         "object number exceeds max",
			objectNumber: 140737488355328,
			generation:   0,
			wantErr:      true,
			errContains:  "object number must be between",
		},
		{
			name:         "object number below min",
			objectNumber: -140737488355328,
			generation:   0,
			wantErr:      true,
			errContains:  "object number must be between",
		},
		{
			name:         "negative generation rejected",
			objectNumber: 1574,
			generation:   -1,
			wantErr:      true,
			errContains:  "generation number must not be negative",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := NewIndirectReference(tc.objectNumber, tc.generation)

			if tc.wantErr {
				if err == nil {
					t.Error("expected an error but got none")
					return
				}
		if len(tc.errContains) > 0 && !strings.Contains(err.Error(), tc.errContains) {
				t.Errorf("expected error containing %q, got %q", tc.errContains, err.Error())
			}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if ref.ObjectNumber() != tc.objectNumber {
				t.Errorf("expected ObjectNumber %d, got %d", tc.objectNumber, ref.ObjectNumber())
			}
			if ref.Generation() != tc.generation {
				t.Errorf("expected Generation %d, got %d", tc.generation, ref.Generation())
			}
		})
	}
}

func TestIndirectReferenceObjectNumberPreservedWithLargeGeneration(t *testing.T) {
	tests := []struct {
		name         string
		objectNumber int64
		generation   int
	}{
		{
			name:         "object number preserved with max int generation",
			objectNumber: 1574,
			generation:   2147483647,
		},
		{
			name:         "negative object number preserved above ushort max generation",
			objectNumber: -1574,
			generation:   65545,
		},
		{
			name:         "min negative object number with large generation",
			objectNumber: -140737488355327,
			generation:   65545,
		},
		{
			name:         "max object number with very large generation",
			objectNumber: 140737488355327,
			generation:   655350,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ref, err := NewIndirectReference(tc.objectNumber, tc.generation)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if ref.ObjectNumber() != tc.objectNumber {
				t.Errorf("expected ObjectNumber %d, got %d", tc.objectNumber, ref.ObjectNumber())
			}
		})
	}
}

func TestIndirectReferencesEqualHashConsistency(t *testing.T) {
	ref1, _ := NewIndirectReference(1267775544, 690)
	ref2, _ := NewIndirectReference(1267775544, 690)

	if !ref1.Equals(ref2) {
		t.Error("equal references must report equal")
	}
}

func TestIndirectReferencesDifferentGenerationNotEqual(t *testing.T) {
	ref1, _ := NewIndirectReference(1267775544, 690)
	ref2, _ := NewIndirectReference(1267775544, 12)

	if ref1.Equals(ref2) {
		t.Error("references with different generation must not be equal")
	}
}

func TestTwoIndirectHashCodeEqual(t *testing.T) {
	ref1, _ := NewIndirectReference(1267775544, 690)
	ref2, _ := NewIndirectReference(1267775544, 690)

	if ref1.HashCode() != ref2.HashCode() {
		t.Errorf("expected equal hash codes, got %d and %d", ref1.HashCode(), ref2.HashCode())
	}
}

func TestTwoIndirectHashCodeNotEqual(t *testing.T) {
	ref1, _ := NewIndirectReference(1267775544, 690)
	ref2, _ := NewIndirectReference(1267775544, 12)

	if ref1.HashCode() == ref2.HashCode() {
		t.Errorf("expected different hash codes, got %d and %d", ref1.HashCode(), ref2.HashCode())
	}
}

func TestTwoIndirectHashCodeSimilarValuesNotEqual(t *testing.T) {
	ref1, _ := NewIndirectReference(12, 1)
	ref2, _ := NewIndirectReference(1, 12)

	if ref1.HashCode() == ref2.HashCode() {
		t.Errorf("expected different hash codes, got %d and %d", ref1.HashCode(), ref2.HashCode())
	}
}

func TestIndirectReferencesSwappedValuesNotEqual(t *testing.T) {
	ref1, _ := NewIndirectReference(12, 1)
	ref2, _ := NewIndirectReference(1, 12)

	if ref1.Equals(ref2) {
		t.Error("references with swapped values must not be equal")
	}
}


