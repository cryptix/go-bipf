package bipf

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/ssb-ngi-pointer/go-bipf/internal/varint"
)

var (
	//  dbg = ioutil.Discard
	dbg = os.Stderr
)

// Decoder holds the internal state of the reader portion of the bipf implementation
type Decoder struct {
	input io.ReadSeeker

	currentType Type
	currentLen  uint64
}

var errTODO = fmt.Errorf("bipf: todo - not implemented")

// NewDecoder initializes the decoder
func NewDecoder(rd io.ReadSeeker) *Decoder {

	dec := &Decoder{
		input:       rd,
		currentType: typeUninited,
	}

	return dec
}

// Next advances to the next value
func (d *Decoder) Next() error {
	return errTODO
}

// Skip discards the current value
func (d *Decoder) Skip() error {
	return errTODO
}

// SeekToLabel seeks through the stream until it finds a value with that chain of object key names to it
func (d *Decoder) SeekToLabel(p string) error {
	return errTODO
}

// Type returns the type of the current value
func (d *Decoder) Type() (Type, error) {

	// start with 1 byte and append to it until we get a clean varint
	var (
		tag      uint64
		tagBytes []byte
	)

readTagByte:
	for {
		var singleByte = make([]byte, 1)
		_, err := io.ReadFull(d.input, singleByte)
		if err != nil {
			return typeUninited, err
		}
		tagBytes = append(tagBytes, singleByte[0])

		var byteCount int
		tag, byteCount = varint.ConsumeVarint(tagBytes)
		switch {
		case byteCount == varint.ErrCodeTruncated:
			continue readTagByte
		case byteCount > 0:
			fmt.Fprintln(dbg, "\tvarint byteCount:", byteCount)
			break readTagByte // we got a varint!
		default:
			return typeUninited, fmt.Errorf("bipf: broken varint tag field")
		}
	}

	fmt.Fprintf(dbg, "\tdecoded %x to tag: %d\n", tagBytes, tag)

	// apply mask to get type
	d.currentType = Type(tag & tagMask)
	if d.currentType >= TypeReserved {
		return 0, fmt.Errorf("bipf: invalid type: %s", d.currentType)
	}

	// shift right to get length
	d.currentLen = uint64(tag >> tagSize)

	// drop some debugging info
	fmt.Fprintln(dbg, "\tvalue type:", d.currentType)
	fmt.Fprintln(dbg, "\tvalue length:", d.currentLen)
	fmt.Fprintln(dbg)
	dbg.Sync()

	return d.currentType, nil
}

// Bool returns the value if the current type is a bool
func (d *Decoder) Bool() (bool, error) {
	if want := TypeBool; d.currentType != want {
		return false, ErrUnexpectedType{Want: want, Got: d.currentType}
	}

	if d.currentLen != 1 {
		return false, fmt.Errorf("bipf/bool: expected 1 bytes of value, not %d", d.currentLen)
	}

	var valueByte = make([]byte, 1)
	_, err := io.ReadFull(d.input, valueByte)
	if err != nil {
		return false, fmt.Errorf("bipf/bool: failed to get value byte: %w", err)
	}

	if valueByte[0] == 0 {
		return false, nil
	}

	if v := valueByte[0]; v != 1 {
		return false, fmt.Errorf("bipf: unexpected bool value: %d", v)
	}

	return true, nil
}

// CopyString returns a copy of value if the current type is a string
func (d *Decoder) CopyString() (string, error) {
	if want := TypeString; d.currentType != want {
		return "", ErrUnexpectedType{Want: want, Got: d.currentType}
	}

	// maybe we want some max size?
	// if d.currentLen > ??? { ... }

	var valueBuffer = make([]byte, d.currentLen)
	_, err := io.ReadFull(d.input, valueBuffer)
	if err != nil {
		return "", fmt.Errorf("bipf/string: failed to read value: %w", err)
	}

	return string(valueBuffer), nil
}

// Double returns the floating-point value if the current type is a double
func (d *Decoder) Double() (float64, error) {
	if want := TypeDouble; d.currentType != want {
		return -1, ErrUnexpectedType{Want: want, Got: d.currentType}
	}

	if d.currentLen == 0 {
		return 0, nil
	}

	// just read 4 bytes of input
	limit := io.LimitReader(d.input, 4)

	var iv float64
	err := binary.Read(limit, binary.LittleEndian, &iv)
	if err != nil {
		return -1, fmt.Errorf("bipf/double: failed to decode value: %w", err)
	}

	return iv, nil
}

// Int32 returns a the 32bit integer value if the current type is a integer
func (d *Decoder) Int32() (int32, error) {
	if want := TypeInt32; d.currentType != want {
		return -1, ErrUnexpectedType{Want: want, Got: d.currentType}
	}

	if d.currentLen != 4 {
		return -1, fmt.Errorf("bipf/int32: expected 4 bytes of value, not %d", d.currentLen)
	}

	// just read 4 bytes of input
	limit := io.LimitReader(d.input, 4)

	var iv int32
	err := binary.Read(limit, binary.LittleEndian, &iv)
	if err != nil {
		return -1, fmt.Errorf("bipf/int32: failed to read value: %w", err)
	}

	return iv, nil
}

type ErrUnexpectedType struct {
	Got, Want Type
}

func (err ErrUnexpectedType) Error() string {
	return fmt.Sprintf("bipf: unexpected type %d but wanted %s", err.Got, err.Want)
}
