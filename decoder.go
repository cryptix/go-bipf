package bipf

import (
	"fmt"
	"io"
)

// Decoder holds the internal state of the reader portion of the bipf implementation
type Decoder struct {
	input io.ReadSeeker
}

var errTODO = fmt.Errorf("bipf: todo - not implemented")

// NewDecoder initializes the decoder
func NewDecoder(rd io.ReadSeeker) *Decoder {
	return &Decoder{input: rd}
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
	return 0, errTODO
}

// Bool returns the value if the current type is a bool
func (d *Decoder) Bool() (bool, error) {
	return false, errTODO
}

// CopyString returns a copy of value if the current type is a string
func (d *Decoder) CopyString() (string, error) {
	return "false", errTODO
}

// Double returns the floating-point value if the current type is a double
func (d *Decoder) Double() (float64, error) {
	return -1, errTODO
}

// Int32 returns a the 32bit integer value if the current type is a integer
func (d *Decoder) Int32() (int32, error) {
	return -1, errTODO
}
