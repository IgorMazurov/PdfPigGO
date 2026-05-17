package util

import (
	"fmt"
	"testing"
)

func TestMurmurHash3(t *testing.T) {
	cases := []struct {
		name        string
		sentence    string
		expectedX86 string
		expectedX64 string
	}{
		{
			name:        "quick_brown_fox",
			sentence:    "The quick brown fox jumps over the lazy dog",
			expectedX86: "2f1583c3ecee2c675d7bf66ce5e91d2c",
			expectedX64: "e34bbc7bbc071b6c7a433ca9c49a9347",
		},
		{
			name:        "murmurhash_author",
			sentence:    "MurmurHash3 was written by Austin Appleby, and is placed in the public",
			expectedX86: "6d3583489d9d1e5a898493af67e2ad10",
			expectedX64: "a91793d43f82cbabda2fb0c28c24799a",
		},
		{
			name:        "single_zero",
			sentence:    "0",
			expectedX86: "0ab2409ea5eb34f8a5eb34f8a5eb34f8",
			expectedX64: "2ac9debed546a3803a8de9e53c875e09",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(tc.sentence)

			hash := ComputeX86_128(data)
			actual := fmt.Sprintf("%x", hash)
			if actual != tc.expectedX86 {
				t.Errorf("ComputeX86_128: got %s, want %s", actual, tc.expectedX86)
			}

			hash = ComputeX64_128(data)
			actual = fmt.Sprintf("%x", hash)
			if actual != tc.expectedX64 {
				t.Errorf("ComputeX64_128: got %s, want %s", actual, tc.expectedX64)
			}
		})
	}
}
