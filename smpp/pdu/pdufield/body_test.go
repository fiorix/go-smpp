// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"testing"
)

func TestData_SM(t *testing.T) {
	want := []byte("hello")
	d := New(ShortMessage, want)
	f, ok := d.(*SM)
	if !ok {
		t.Fatalf("unexpected field type: want SM, have %#v", f)
	}
	if !bytes.Equal(want, f.Bytes()) {
		t.Fatalf("unexpected field data: want %q, have %q", want, f.Bytes())
	}
	want = []byte("")
	d = New(ShortMessage, nil)
	f, ok = d.(*SM)
	if !ok {
		t.Fatalf("unexpected field type: want SM, have %#v", f)
	}
	if !bytes.Equal(want, f.Bytes()) {
		t.Fatalf("unexpected field data: want %q, have %q", want, f.Bytes())
	}
}

func TestData_Invalid(t *testing.T) {
	d := New("foobar", nil)
	if d != nil {
		t.Fatalf("unexpected field: %#v", d)
	}
}
