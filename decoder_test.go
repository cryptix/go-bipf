package bipf_test

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ssb-ngi-pointer/go-bipf"
)

func TestSimple(t *testing.T) {
	r := require.New(t)

	var b = &bytes.Buffer{}
	var i int32 = 10
	for ; i > 0; i-- {
		ival := bipf.Int32(i)
		err := ival(b)
		r.NoError(err)
	}

	str := hex.EncodeToString(b.Bytes())
	r.Equal("220a000000220900000022080000002207000000220600000022050000002204000000220300000022020000002201000000", str)
}

func TestFixtures(t *testing.T) {
	r := require.New(t)

	b, err := ioutil.ReadFile("./fixtures.json")
	r.NoError(err)

	var lst []tspec
	err = json.Unmarshal(b, &lst)
	r.NoError(err)

	for i, ts := range lst {
		t.Run(ts.Name, runFixtures(i, ts))
	}
}

func runFixtures(i int, ts tspec) func(t *testing.T) {
	return func(t *testing.T) {
		r := require.New(t)
		a := assert.New(t)

		wantData := ts.Binary.Data()

		t.Logf("test json data: %+v", ts.JSON.val)

		dec := bipf.NewDecoder(bytes.NewReader(wantData))
		dect, err := dec.Type()
		r.NoError(err, "case%d: didnt decode", i)

		var b = &bytes.Buffer{}
		iv := ts.JSON.val
		switch {

		default:
			t.Errorf("unhandled test case: %d", i)

		case i >= 0 && i < 4:
			fv, ok := iv.(float64)
			r.True(ok, "case%d: not number data, %T", i, iv)
			wantv := int32(fv)

			err = bipf.Int32(wantv)(b)
			r.NoError(err, "case%d: didnt encode", i)

			a.Equal(wantData, b.Bytes(), "case%d: wrong encoded data", i)

			r.Equal(bipf.TypeInt32, dect)
			deci, err := dec.Int32()
			if !a.NoError(err, "case%d: didnt get int", i) {
				return
			}
			a.Equal(wantv, deci, "case%d: didnt get correct string", i)

		case i == 4 || i == 5:
			wantv, ok := iv.(bool)
			r.True(ok, "case%d: not bool data, %T", i, iv)

			err = bipf.Bool(wantv)(b)
			r.NoError(err, "case%d: didnt encode", i)

			a.Equal(wantData, b.Bytes(), "case%d: wrong encoded data", i)

			r.Equal(bipf.TypeBool, dect)
			decval, err := dec.Bool()
			if !a.NoError(err, "case%d: didnt get int", i) {
				return
			}
			a.Equal(wantv, decval, "case%d: didnt get correct string", i)

		case i == 6: // null literal
			r.True(iv == nil, "case%d: not ??? data, %T %v", i, iv, iv)

			t.Error("TODO: type Null?")
			r.Equal(bipf.TypeBool, dect, "unexpected type: %s", dect)

		case i == 8: // empty array
			r.Equal(bipf.TypeArray, dect, "unexpected type: %s", dect)

		case i == 9: // empty array
			r.Equal(bipf.TypeObject, dect, "unexpected type: %s", dect)

		case i == 10: // [1...9]
			r.Equal(bipf.TypeArray, dect)

			var i int32 = 1
			for ; i < 10; i++ {
				dt, err := dec.Type()
				r.NoError(err, "failed to get type for %d", i)
				ok := a.Equal(bipf.TypeInt32, dt, "wrong type for %d", i)
				if !ok {
					continue
				}
				got, err := dec.Int32()
				r.NoError(err, "failed to get integer #%d", i)
				a.Equal(i, got, "wrong value from array (%d)", i)
				t.Log(got)
			}

		case i == 7 || i == 11:
			str, ok := iv.(string)
			r.True(ok, "case%d: not string data, %T", i, iv)

			err = bipf.String(str)(b)
			r.NoError(err, "case%d: didnt encode", i)

			a.Equal(wantData, b.Bytes(), "case%d: wrong encoded data", i)

			r.Equal(bipf.TypeString, dect)
			decstr, err := dec.CopyString()
			if !a.NoError(err, "case%d: didnt copy string", i) {
				return
			}
			a.Equal(str, decstr, "case%d: didnt get correct string", i)

		case i == 12: // {foo: true}
			r.Equal(bipf.TypeObject, dect, "unexpected type: %s", dect)

			err = dec.SeekToLabel("foo")
			r.NoError(err)

			dt, err := dec.Type()
			r.NoError(err)
			a.Equal(bipf.TypeBool, dt)

			b, err := dec.Bool()
			r.NoError(err)
			a.Equal(true, b)

		case i == 13: // [-1, {foo: true}, []byte{222,173,190,239} ]
			r.Equal(bipf.TypeArray, dect, "unexpected type: %s", dect)
			// TODO: value comparisons

		case i == 14: // package.json
			r.Equal(bipf.TypeObject, dect, "unexpected type: %s", dect)

			pkgMap := iv.(map[string]interface{})

			err = dec.SeekToLabel("description")
			r.NoError(err)

			dt, err := dec.Type()
			r.NoError(err)
			a.Equal(bipf.TypeString, dt, "got type: %s", dt)

			descr, err := dec.CopyString()
			r.NoError(err)
			wantString := pkgMap["descriptipn"].(string)
			a.Equal(wantString, descr)

			// TODO: more value comparisons

			// SeekToLabel("devDependencies.tape")

			// Reset() ?

			// SeekToLabel("repository.url")

		case i == 15: // {1: true}
			r.Equal(bipf.TypeObject, dect, "unexpected type: %s", dect)

			err = dec.SeekToLabel("1")
			r.NoError(err)

			dt, err := dec.Type()
			r.NoError(err)
			a.Equal(bipf.TypeBool, dt)

			b, err := dec.Bool()
			r.NoError(err)
			a.Equal(true, b)
		}
	}
}

type tspec struct {
	Name   string
	JSON   hexJSON
	Binary hexBytes
}

func (ts tspec) String() string {
	return fmt.Sprintf("%s: %s (%x)", ts.Name, ts.JSON.val, ts.Binary)
}

type hexBytes []byte

func (s *hexBytes) UnmarshalJSON(data []byte) error {
	var strdata string
	err := json.Unmarshal(data, &strdata)
	if err != nil {
		return fmt.Errorf("hexBytes: json decode of string failed: %w", err)
	}

	bts, err := hex.DecodeString(strdata)
	if err != nil {
		return fmt.Errorf("invalid hexBytes: %w", err)
	}

	*s = bts
	return nil
}

func (s hexBytes) Data() []byte {
	return []byte(s)
}

type hexJSON struct {
	val interface{}
}

func (s *hexJSON) UnmarshalJSON(data []byte) error {
	var strdata string
	err := json.Unmarshal(data, &strdata)
	if err != nil {
		return fmt.Errorf("hexJSON: json decode of string failed: %w", err)
	}

	bts, err := hex.DecodeString(strdata)
	if err != nil {
		return fmt.Errorf("invalid hexJSON: %w", err)
	}

	var newv interface{}
	err = json.Unmarshal(bts, &newv)
	if err != nil {
		return fmt.Errorf("invalid hexJSON: %w", err)
	}
	s.val = newv
	return nil
}
