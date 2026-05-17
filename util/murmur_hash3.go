package util

import (
	"encoding/binary"
)

// MurmurHash3 is a non-cryptographic hash function suitable for general hash-based lookup.
// It was created by Austin Appleby in 2008 and released into the public domain.
// The name comes from two basic operations, multiply (MU) and rotate (R), used in its inner loop.
//
// Note: The x86 and x64 variants do not produce the same results, as the
// algorithms are optimized for their respective platforms.

// ComputeX86_128 computes the MurmurHash3 128-bit x86 variant hash of data with seed 0.
// Returns a 16-byte slice.
func ComputeX86_128(data []byte) []byte {
	return computeX86_128(data, len(data), 0)
}

// ComputeX86_128WithSeed computes the MurmurHash3 128-bit x86 variant hash of data with the given seed.
// Returns a 16-byte slice.
func ComputeX86_128WithSeed(data []byte, length int, seed uint32) []byte {
	return computeX86_128(data, length, seed)
}

// ComputeX64_128 computes the MurmurHash3 128-bit x64 variant hash of data with seed 0.
// Returns a 16-byte slice.
func ComputeX64_128(data []byte) []byte {
	return computeX64_128(data, len(data), 0)
}

// ComputeX64_128WithSeed computes the MurmurHash3 128-bit x64 variant hash of data with the given seed.
// Returns a 16-byte slice.
func ComputeX64_128WithSeed(data []byte, length int, seed uint32) []byte {
	return computeX64_128(data, length, seed)
}

func computeX86_128(data []byte, length int, seed uint32) []byte {
	var h [4]uint32
	computeX86_128Core(data, length, seed, &h)

	result := make([]byte, 16)
	binary.BigEndian.PutUint32(result[0:4], h[0])
	binary.BigEndian.PutUint32(result[4:8], h[1])
	binary.BigEndian.PutUint32(result[8:12], h[2])
	binary.BigEndian.PutUint32(result[12:16], h[3])

	return result
}

func computeX64_128(data []byte, length int, seed uint32) []byte {
	var h [2]uint64
	computeX64_128Core(data, length, seed, &h)

	result := make([]byte, 16)
	binary.BigEndian.PutUint64(result[0:8], h[0])
	binary.BigEndian.PutUint64(result[8:16], h[1])

	return result
}

func computeX86_128Core(data []byte, length int, seed uint32, outHash *[4]uint32) {
	const (
		c1 = 0x239b961b
		c2 = 0xab0e9789
		c3 = 0x38b34ae5
		c4 = 0xa1e38b93
	)

	h1 := seed
	h2 := seed
	h3 := seed
	h4 := seed
	nblocks := length / 16

	for i := 0; i < nblocks; i++ {
		offset := i * 16

		k1 := binary.LittleEndian.Uint32(data[offset : offset+4])
		k2 := binary.LittleEndian.Uint32(data[offset+4 : offset+8])
		k3 := binary.LittleEndian.Uint32(data[offset+8 : offset+12])
		k4 := binary.LittleEndian.Uint32(data[offset+12 : offset+16])

		k1 *= c1
		k1 = rotl32(k1, 15)
		k1 *= c2
		h1 ^= k1
		h1 = rotl32(h1, 19)
		h1 += h2
		h1 = h1*5 + 0x561ccd1b

		k2 *= c2
		k2 = rotl32(k2, 16)
		k2 *= c3
		h2 ^= k2
		h2 = rotl32(h2, 17)
		h2 += h3
		h2 = h2*5 + 0x0bcaa747

		k3 *= c3
		k3 = rotl32(k3, 17)
		k3 *= c4
		h3 ^= k3
		h3 = rotl32(h3, 15)
		h3 += h4
		h3 = h3*5 + 0x96cd1c35

		k4 *= c4
		k4 = rotl32(k4, 18)
		k4 *= c1
		h4 ^= k4
		h4 = rotl32(h4, 13)
		h4 += h1
		h4 = h4*5 + 0x32ac3b17
	}

	tailStart := nblocks * 16
	tk1 := uint32(0)
	tk2 := uint32(0)
	tk3 := uint32(0)
	tk4 := uint32(0)

	switch length & 15 {
	case 15:
		tk4 ^= uint32(data[tailStart+14]) << 16
		fallthrough
	case 14:
		tk4 ^= uint32(data[tailStart+13]) << 8
		fallthrough
	case 13:
		tk4 ^= uint32(data[tailStart+12])
		tk4 *= c4
		tk4 = rotl32(tk4, 18)
		tk4 *= c1
		h4 ^= tk4
		fallthrough
	case 12:
		tk3 ^= uint32(data[tailStart+11]) << 24
		fallthrough
	case 11:
		tk3 ^= uint32(data[tailStart+10]) << 16
		fallthrough
	case 10:
		tk3 ^= uint32(data[tailStart+9]) << 8
		fallthrough
	case 9:
		tk3 ^= uint32(data[tailStart+8])
		tk3 *= c3
		tk3 = rotl32(tk3, 17)
		tk3 *= c4
		h3 ^= tk3
		fallthrough
	case 8:
		tk2 ^= uint32(data[tailStart+7]) << 24
		fallthrough
	case 7:
		tk2 ^= uint32(data[tailStart+6]) << 16
		fallthrough
	case 6:
		tk2 ^= uint32(data[tailStart+5]) << 8
		fallthrough
	case 5:
		tk2 ^= uint32(data[tailStart+4])
		tk2 *= c2
		tk2 = rotl32(tk2, 16)
		tk2 *= c3
		h2 ^= tk2
		fallthrough
	case 4:
		tk1 ^= uint32(data[tailStart+3]) << 24
		fallthrough
	case 3:
		tk1 ^= uint32(data[tailStart+2]) << 16
		fallthrough
	case 2:
		tk1 ^= uint32(data[tailStart+1]) << 8
		fallthrough
	case 1:
		tk1 ^= uint32(data[tailStart])
		tk1 *= c1
		tk1 = rotl32(tk1, 15)
		tk1 *= c2
		h1 ^= tk1
	}

	h1 ^= uint32(length)
	h2 ^= uint32(length)
	h3 ^= uint32(length)
	h4 ^= uint32(length)
	h1 += h2
	h1 += h3
	h1 += h4
	h2 += h1
	h3 += h1
	h4 += h1
	h1 = fmix32(h1)
	h2 = fmix32(h2)
	h3 = fmix32(h3)
	h4 = fmix32(h4)
	h1 += h2
	h1 += h3
	h1 += h4
	h2 += h1
	h3 += h1
	h4 += h1

	outHash[0] = h1
	outHash[1] = h2
	outHash[2] = h3
	outHash[3] = h4
}

func computeX64_128Core(data []byte, length int, seed uint32, outHash *[2]uint64) {
	const (
		c1 = 0x87c37b91114253d5
		c2 = 0x4cf5ad432745937f
	)

	h1 := uint64(seed)
	h2 := uint64(seed)
	nblocks := length / 16

	for i := 0; i < nblocks; i++ {
		offset := i * 16
		k1 := binary.LittleEndian.Uint64(data[offset : offset+8])
		k2 := binary.LittleEndian.Uint64(data[offset+8 : offset+16])

		k1 *= c1
		k1 = rotl64(k1, 31)
		k1 *= c2
		h1 ^= k1
		h1 = rotl64(h1, 27)
		h1 += h2
		h1 = h1*5 + 0x52dce729

		k2 *= c2
		k2 = rotl64(k2, 33)
		k2 *= c1
		h2 ^= k2
		h2 = rotl64(h2, 31)
		h2 += h1
		h2 = h2*5 + 0x38495ab5
	}

	tailStart := nblocks * 16
	tk1 := uint64(0)
	tk2 := uint64(0)

	switch length & 15 {
	case 15:
		tk2 ^= uint64(data[tailStart+14]) << 48
		fallthrough
	case 14:
		tk2 ^= uint64(data[tailStart+13]) << 40
		fallthrough
	case 13:
		tk2 ^= uint64(data[tailStart+12]) << 32
		fallthrough
	case 12:
		tk2 ^= uint64(data[tailStart+11]) << 24
		fallthrough
	case 11:
		tk2 ^= uint64(data[tailStart+10]) << 16
		fallthrough
	case 10:
		tk2 ^= uint64(data[tailStart+9]) << 8
		fallthrough
	case 9:
		tk2 ^= uint64(data[tailStart+8])
		tk2 *= c2
		tk2 = rotl64(tk2, 33)
		tk2 *= c1
		h2 ^= tk2
		fallthrough
	case 8:
		tk1 ^= uint64(data[tailStart+7]) << 56
		fallthrough
	case 7:
		tk1 ^= uint64(data[tailStart+6]) << 48
		fallthrough
	case 6:
		tk1 ^= uint64(data[tailStart+5]) << 40
		fallthrough
	case 5:
		tk1 ^= uint64(data[tailStart+4]) << 32
		fallthrough
	case 4:
		tk1 ^= uint64(data[tailStart+3]) << 24
		fallthrough
	case 3:
		tk1 ^= uint64(data[tailStart+2]) << 16
		fallthrough
	case 2:
		tk1 ^= uint64(data[tailStart+1]) << 8
		fallthrough
	case 1:
		tk1 ^= uint64(data[tailStart+0])
		tk1 *= c1
		tk1 = rotl64(tk1, 31)
		tk1 *= c2
		h1 ^= tk1
	}

	h1 ^= uint64(length)
	h2 ^= uint64(length)
	h1 += h2
	h2 += h1
	h1 = fmix64(h1)
	h2 = fmix64(h2)
	h1 += h2
	h2 += h1

	outHash[0] = h1
	outHash[1] = h2
}

func rotl32(x uint32, r int) uint32 {
	return (x << r) | (x >> (32 - r))
}

func rotl64(x uint64, r int) uint64 {
	return (x << r) | (x >> (64 - r))
}

func fmix32(h uint32) uint32 {
	h ^= h >> 16
	h *= 0x85ebca6b
	h ^= h >> 13
	h *= 0xc2b2ae35
	h ^= h >> 16
	return h
}

func fmix64(k uint64) uint64 {
	k ^= k >> 33
	k *= 0xff51afd7ed558ccd
	k ^= k >> 33
	k *= 0xc4ceb9fe1a85ec53
	k ^= k >> 33
	return k
}
