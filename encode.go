// SPDX-License-Identifier: MIT

package bipf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ssb-ngi-pointer/go-bipf/internal/varint"
)

//go:generate stringer -type=Type -trimprefix Type

type Type byte

// This block defines the tags
const (
	TypeString   Type = iota // = 0b000
	TypeBuffer               // = 0b001
	TypeInt32                // = 0b010 //32bit int
	TypeDouble               // = 0b011 //use next 8 bytes to encode 64bit float
	TypeArray                // = 0b100
	TypeObject               // = 0b101
	TypeBool                 // = 0b110 // or null
	TypeReserved             // = 0b111

	typeUninited = 255 // sadly 0 is taken
)

const (
	tagSize = 3
	tagMask = 7
)

type Valuer func(io.Writer) error

// String encodes type 1
// TODO: unicode!?
func String(v string) Valuer {
	return func(w io.Writer) error {
		var err error
		if len(v) == 0 {
			_, err = w.Write([]byte{byte(TypeString)})
			return err
		}

		strlen := len(v)
		tag := Type(strlen<<tagSize) | TypeString

		err = varint.WriteVarint(w, uint64(tag))
		if err != nil {
			return err
		}
		n, err := w.Write([]byte(v))
		if err != nil {
			return err
		}
		if n != strlen {
			return fmt.Errorf("short write. %d vs %d", n, strlen)
		}
		return nil
	}
}

func Bytes([]byte) Valuer {
	return func(w io.Writer) error {
		err := fmt.Errorf("TODO: Bytes()")
		return err
	}
}

func Int32(v int32) Valuer {
	return func(w io.Writer) error {
		_, err := w.Write([]byte{0x22})
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, v)
		return err
	}
}
func Double(v float64) Valuer {
	return func(w io.Writer) error {
		_, err := w.Write([]byte{0x43})
		if err != nil {
			return err
		}
		err = binary.Write(w, binary.LittleEndian, v)
		return err
	}
}

// MapOf encodes the passed map as an object.
// If the order of the fields is important, these can be passed as variadic list of strings.
// If it's passed it needs to have the same length as the number of keys in the map.
func MapOf(m map[string]Valuer, order ...string) Valuer {
	return func(w io.Writer) error {
		if len(order) > 0 && len(order) != len(m) {
			return fmt.Errorf("map and orderd field size differ")
		}

		var buf bytes.Buffer

		var keys = make([]string, len(m))
		if len(order) == 0 {
			// no order, just pick any
			i := 0
			for k := range m {
				keys[i] = k
				i++
			}
		} else {
			for i, k := range order {
				if _, has := m[k]; !has {
					return fmt.Errorf("orderd field %q not in map", k)
				}
				keys[i] = k
			}
		}

		for _, k := range keys {
			v := m[k]

			if err := String(k)(&buf); err != nil {
				return fmt.Errorf("map encoding failed during key(%q): %w", k, err)
			}

			if err := v(&buf); err != nil {
				return fmt.Errorf("map encoding failed at value for key(%q): %w", k, err)
			}

		}
		objSize := buf.Len()
		tag := Type(objSize<<tagSize) | TypeObject
		err := varint.WriteVarint(w, uint64(tag))
		if err != nil {
			return err
		}

		n, err := buf.WriteTo(w)
		if err != nil {
			return err
		}

		if int(n) != objSize {
			return fmt.Errorf("short write. %d vs %d", n, objSize)
		}

		return nil
	}
}

func ListOf(items ...Valuer) Valuer {
	return func(w io.Writer) error {
		var buf bytes.Buffer

		for idx, item := range items {
			err := item(&buf)
			if err != nil {
				return fmt.Errorf("list failed to encode item %d: %w", idx, err)
			}
		}

		arrSize := buf.Len()
		tag := Type(arrSize<<tagSize) | TypeArray
		err := varint.WriteVarint(w, uint64(tag))
		if err != nil {
			return err
		}

		n, err := buf.WriteTo(w)
		if err != nil {
			return err
		}

		if int(n) != arrSize {
			return fmt.Errorf("short write. %d vs %d", n, arrSize)
		}

		return nil
	}
}

func Bool(yes bool) Valuer {
	return func(w io.Writer) error {
		var t byte = 0x00
		if yes {
			t = 0x01
		}
		_, err := w.Write([]byte{0x0e, t})
		return err
	}
}
