// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"strconv"
	"testing"
)

func TestFixed(t *testing.T) {
	f := &Fixed{Data: 0x34}
	if f.Len() != 1 {
		t.Fatalf("unexpected len: want 1, have %d", f.Len())
	}
	if v, ok := f.Raw().(uint8); !ok {
		t.Fatalf("unexpected type: want uint8, have %#v", v)
	}
	ws := strconv.Itoa(0x34)
	if v := f.String(); v != string(ws) {
		t.Fatalf("unexpected string: want %q, have %q", ws, v)
	}
	wb := []byte{0x34}
	if v := f.Bytes(); !bytes.Equal(wb, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", wb, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(wb, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", wb, v)
	}
}

func TestVariable(t *testing.T) {
	want := []byte("foobar")
	f := &Variable{Data: want}
	lw := len(want) + 1
	if f.Len() != lw {
		t.Fatalf("unexpected len: want %d, have %d", lw, f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string(want) {
		t.Fatalf("unexpected string: want %q have %q", want, v)
	}
	want = []byte("foobar\x00")
	if v := f.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", want, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", want, v)
	}
}

func TestSM(t *testing.T) {
	want := []byte("foobar")
	f := &SM{Data: want}
	if f.Len() != len(want) {
		t.Fatalf("unexpected len: want %d, have %d", len(want), f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string(want) {
		t.Fatalf("unexpected string: want %q have %q", want, v)
	}
	if v := f.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", want, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", want, v)
	}
}
