// SPDX-License-Identifier: MIT

// Package varint is a copy of the protobuf varint encoding.
//
// https://github.com/protocolbuffers/protobuf-go/blob/fb30439f551a7e79e413e7b4f5f4dfb58e117d73/encoding/protowire/wire.go#L180
package varint

import (
	"io"
	"math/bits"
)

// WriteVarint writes v to w as a varint-encoded uint64.
func WriteVarint(w io.Writer, v uint64) error {
	var b []byte
	switch {

	case v < 1<<7:
		b = []byte{
			byte(v),
		}
	case v < 1<<14:
		b = []byte{
			byte((v>>0)&0x7f | 0x80),
			byte(v >> 7),
		}
	case v < 1<<21:
		b = []byte{
			byte((v>>0)&0x7f | 0x80),
			byte((v>>7)&0x7f | 0x80),
			byte(v >> 14),
		}
	case v < 1<<28:
		b = []byte{
			byte((v>>0)&0x7f | 0x80),
			byte((v>>7)&0x7f | 0x80),
			byte((v>>14)&0x7f | 0x80),
			byte(v >> 21),
		}
	case v < 1<<35:
		b = []byte{
			byte((v>>0)&0x7f | 0x80),
			byte((v>>7)&0x7f | 0x80),
			byte((v>>14)&0x7f | 0x80),
			byte((v>>21)&0x7f | 0x80),
			byte(v >> 28),
		}
	case v < 1<<42:
		b = []byte{
			byte((v>>0)&0x7f | 0x80),
			byte((v>>7)&0x7f | 0x80),
			byte((v>>14)&0x7f | 0x80),
			byte((v>>21)&0x7f | 0x80),
			byte((v>>28)&0x7f | 0x80),
			byte(v >> 35),
		}
	case v < 1<<49:
		b = []byte{
			byte((v>>0)&0x7f | 0x80),
			byte((v>>7)&0x7f | 0x80),
			byte((v>>14)&0x7f | 0x80),
			byte((v>>21)&0x7f | 0x80),
			byte((v>>28)&0x7f | 0x80),
			byte((v>>35)&0x7f | 0x80),
			byte(v >> 42),
		}
	case v < 1<<56:
		b = []byte{
			byte((v>>0)&0x7f | 0x80),
			byte((v>>7)&0x7f | 0x80),
			byte((v>>14)&0x7f | 0x80),
			byte((v>>21)&0x7f | 0x80),
			byte((v>>28)&0x7f | 0x80),
			byte((v>>35)&0x7f | 0x80),
			byte((v>>42)&0x7f | 0x80),
			byte(v >> 49),
		}
	case v < 1<<63:
		b = []byte{
			byte((v>>0)&0x7f | 0x80),
			byte((v>>7)&0x7f | 0x80),
			byte((v>>14)&0x7f | 0x80),
			byte((v>>21)&0x7f | 0x80),
			byte((v>>28)&0x7f | 0x80),
			byte((v>>35)&0x7f | 0x80),
			byte((v>>42)&0x7f | 0x80),
			byte((v>>49)&0x7f | 0x80),
			byte(v >> 56),
		}
	default:
		b = []byte{
			byte((v>>0)&0x7f | 0x80),
			byte((v>>7)&0x7f | 0x80),
			byte((v>>14)&0x7f | 0x80),
			byte((v>>21)&0x7f | 0x80),
			byte((v>>28)&0x7f | 0x80),
			byte((v>>35)&0x7f | 0x80),
			byte((v>>42)&0x7f | 0x80),
			byte((v>>49)&0x7f | 0x80),
			byte((v>>56)&0x7f | 0x80),
			1}

	}
	_, err := w.Write(b)
	return err
}

const (
	_ = -iota
	errCodeTruncated
	errCodeFieldNumber
	errCodeOverflow
)

// ConsumeVarint parses b as a varint-encoded uint64, reporting its length.
// This returns a negative length upon an error (see ParseError).
func ConsumeVarint(b []byte) (v uint64, n int) {
	var y uint64
	if len(b) <= 0 {
		return 0, errCodeTruncated
	}
	v = uint64(b[0])
	if v < 0x80 {
		return v, 1
	}
	v -= 0x80

	if len(b) <= 1 {
		return 0, errCodeTruncated
	}
	y = uint64(b[1])
	v += y << 7
	if y < 0x80 {
		return v, 2
	}
	v -= 0x80 << 7

	if len(b) <= 2 {
		return 0, errCodeTruncated
	}
	y = uint64(b[2])
	v += y << 14
	if y < 0x80 {
		return v, 3
	}
	v -= 0x80 << 14

	if len(b) <= 3 {
		return 0, errCodeTruncated
	}
	y = uint64(b[3])
	v += y << 21
	if y < 0x80 {
		return v, 4
	}
	v -= 0x80 << 21

	if len(b) <= 4 {
		return 0, errCodeTruncated
	}
	y = uint64(b[4])
	v += y << 28
	if y < 0x80 {
		return v, 5
	}
	v -= 0x80 << 28

	if len(b) <= 5 {
		return 0, errCodeTruncated
	}
	y = uint64(b[5])
	v += y << 35
	if y < 0x80 {
		return v, 6
	}
	v -= 0x80 << 35

	if len(b) <= 6 {
		return 0, errCodeTruncated
	}
	y = uint64(b[6])
	v += y << 42
	if y < 0x80 {
		return v, 7
	}
	v -= 0x80 << 42

	if len(b) <= 7 {
		return 0, errCodeTruncated
	}
	y = uint64(b[7])
	v += y << 49
	if y < 0x80 {
		return v, 8
	}
	v -= 0x80 << 49

	if len(b) <= 8 {
		return 0, errCodeTruncated
	}
	y = uint64(b[8])
	v += y << 56
	if y < 0x80 {
		return v, 9
	}
	v -= 0x80 << 56

	if len(b) <= 9 {
		return 0, errCodeTruncated
	}
	y = uint64(b[9])
	v += y << 63
	if y < 2 {
		return v, 10
	}
	return 0, errCodeOverflow
}

// SizeVarint returns the encoded size of a varint.
// The size is guaranteed to be within 1 and 10, inclusive.
func SizeVarint(v uint64) int {
	// This computes 1 + (bits.Len64(v)-1)/7.
	// 9/64 is a good enough approximation of 1/7
	return int(9*uint32(bits.Len64(v))+64) / 64
}
