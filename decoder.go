package bipf

import (
	"fmt"
	"io"
)

type Decoder struct {
	input io.ReadSeeker
}

var errTODO = fmt.Errorf("bipf: todo - not implemented")

func NewDecoder(rd io.ReadSeeker) *Decoder {
	return &Decoder{input: rd}
}

func (d *Decoder) Bool() (bool, error) {
	return false, errTODO
}
func (d *Decoder) CopyString() (string, error) {
	return "false", errTODO
}
func (d *Decoder) Double() (float64, error) {
	return -1, errTODO
}
func (d *Decoder) Int32() (int32, error) {
	return -1, errTODO
}
func (d *Decoder) Next() error {
	return errTODO
}

func (d *Decoder) SeekToLabel(p string) error {
	return errTODO
}
func (d *Decoder) Type() (Type, error) {
	return 0, errTODO
}
