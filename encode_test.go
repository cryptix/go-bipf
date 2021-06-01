// SPDX-License-Identifier: MIT
package bipf

import (
	"bytes"
	"testing"
)

func TestMapOfWrongOrderSize(t *testing.T) {
	var b = &bytes.Buffer{}

	v := MapOf(map[string]Valuer{}, "foo")

	err := v(b)
	if err == nil {
		t.Error("expected error")
	}
}

func TestMapOfOrderdFieldNotInMap(t *testing.T) {
	var b = &bytes.Buffer{}

	v := MapOf(map[string]Valuer{"foo": Bool(true)}, "nope")

	err := v(b)
	if err == nil {
		t.Error("expected error")
	}
}
