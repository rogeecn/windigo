//go:build windows

package utl

import (
	"encoding/binary"
)

// Assembles an uint16 from two uint8.
func Make16(lo, hi uint8) uint16 {
	return (uint16(lo) & 0xff) | ((uint16(hi) & 0xff) << 8)
}

// Assembles an uint32 from two uint16.
func Make32(lo, hi uint16) uint32 {
	return (uint32(lo) & 0xffff) | ((uint32(hi) & 0xffff) << 16)
}

// Assembles an uint64 from two uint32.
func Make64(lo, hi uint32) uint64 {
	return (uint64(lo) & 0xffff_ffff) | ((uint64(hi) & 0xffff_ffff) << 32)
}

// Breaks an uint16 into low and high uint8.
func Break16(val uint16) (lo, hi uint8) {
	return uint8(val & 0xff), uint8(val >> 8 & 0xff)
}

// Breaks an uint32 into low and high uint16.
func Break32(val uint32) (lo, hi uint16) {
	return uint16(val & 0xffff), uint16(val >> 16 & 0xffff)
}

// Breaks an uint64 into low and high uint32.
func Break64(val uint64) (lo, hi uint32) {
	return uint32(val & 0xffff_ffff), uint32(val >> 32 & 0xffff_ffff)
}

// Tells whether the number has the nth bit set.
//
// bitPosition must be in the range 0-7.
func BitIsSet(number, bitPosition uint8) bool {
	return (number & (1 << bitPosition)) > 0
}

// Returns a new number with the nth bit set or clear.
//
// bitPosition must be in the range 0-7.
func BitSet(number, bitPosition uint8, doSet bool) uint8 {
	if doSet {
		return number | (1 << bitPosition)
	} else {
		return number &^ (1 << bitPosition)
	}
}

// Reverses the bytes, not the bits.
func ReverseBytes64(n uint64) uint64 {
	var buf64 [8]byte
	binary.LittleEndian.PutUint64(buf64[:], n)
	return binary.BigEndian.Uint64(buf64[:])
}
