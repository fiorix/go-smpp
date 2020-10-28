// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutlv

import (
	"bytes"
	"testing"
)

func TestTag_Hex(t *testing.T) {
	tag := TagDestAddrSubunit
	want := "0005"
	if v := tag.Hex(); v != want {
		t.Fatalf("unexpected hex: want %q have %q", want, v)
	}

	tag = 0x0501
	want = "0501"
	if v := tag.Hex(); v != want {
		t.Fatalf("unexpected hex: want %q have %q", want, v)
	}
}

func TestTLVField(t *testing.T) {
	var want []byte
	want = append(want, []byte{0x13, 0x0C}...) // Tag
	want = append(want, []byte{0x00, 0x06}...) // Length
	want = append(want, []byte("foobar")...)   // Value

	f := &Field{Tag: 0x130C, Data: []byte("foobar")}
	if f.Len() != len(want) {
		t.Fatalf("unexpected len: want %d, have %d", len(want), f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string("foobar") {
		t.Fatalf("unexpected string: want %q have %q", want, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", want, v)
	}
}

func TestTLVField_NullTerminated(t *testing.T) {
	var want []byte
	want = append(want, []byte{0x13, 0x0C}...)   // Tag
	want = append(want, []byte{0x00, 0x07}...)   // Length
	want = append(want, []byte("foobar\x00")...) // Value

	f := &Field{Tag: 0x130C, Data: []byte("foobar\x00")}
	if f.Len() != len(want) {
		t.Fatalf("unexpected len: want %d, have %d", len(want), f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string("foobar") {
		t.Fatalf("unexpected string: want %q have %q", want, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", want, v)
	}
}
