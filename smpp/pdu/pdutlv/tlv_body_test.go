// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutlv

import (
	"bytes"
	"testing"
)

func TestNewTLV_Raw(t *testing.T) {
	want := []byte("hello")
	wb := []byte{0x05, 0x01, 0x00, 0x05}
	wb = append(wb, want...)
	d := NewTLV(0x0501, want)
	f, ok := d.(*Field)
	if !ok {
		t.Fatalf("unexpected field type: want SM, have %#v", f)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(wb, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", wb, v)
	}
	if !bytes.Equal(want, f.Bytes()) {
		t.Fatalf("unexpected field data: want %q, have %q", want, f.Bytes())
	}
}

func TestNewTLV_Nil(t *testing.T) {
	want := make([]byte, 0)
	wb := []byte{0x05, 0x01, 0x00, 0x00}
	d := NewTLV(0x0501, nil)
	f, ok := d.(*Field)
	if !ok {
		t.Fatalf("unexpected field type: want SM, have %#v", f)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(wb, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", wb, v)
	}
	if !bytes.Equal(want, f.Bytes()) {
		t.Fatalf("unexpected field data: want %q, have %q", want, f.Bytes())
	}

	d = NewTLV(0x0501, want)
	f, ok = d.(*Field)
	if !ok {
		t.Fatalf("unexpected field type: want SM, have %#v", f)
	}
	b.Reset()
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(wb, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", wb, v)
	}
	if !bytes.Equal(want, f.Bytes()) {
		t.Fatalf("unexpected field data: want %q, have %q", want, f.Bytes())
	}
}
