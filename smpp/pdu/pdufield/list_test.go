// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"testing"
)

func TestListDecoder_Fixed(t *testing.T) {
	l := List{DataCoding}
	want := []byte{0x02}
	b := bytes.NewBuffer(want)
	m, err := l.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	f, ok := m[DataCoding]
	if !ok {
		t.Fatalf("missing %q key: %#v", DataCoding, m)
	}
	v, ok := f.(*Fixed)
	if !ok {
		t.Fatalf("field is not type Fixed: %#v", f)
	}
	if !bytes.Equal(want, v.Bytes()) {
		t.Fatalf("unexpected data: want %q, have %q", want, v)
	}
}

func TestListDecoder_Variable(t *testing.T) {
	l := List{SystemID}
	want := []byte{'h', 'e', 'l', 'l', 'o', 0x00}
	b := bytes.NewBuffer(want)
	m, err := l.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	f, ok := m[SystemID]
	if !ok {
		t.Fatalf("missing %q key: %#v", SystemID, m)
	}
	v, ok := f.(*Variable)
	if !ok {
		t.Fatalf("field is not type Variable: %#v", f)
	}
	if !bytes.Equal(want, v.Bytes()) {
		t.Fatalf("unexpected data: want %q, have %q", want, v)
	}
}

func TestListDecoder_SM(t *testing.T) {
	l := List{SMLength, ShortMessage}
	want := []byte{0x05, 'h', 'e', 'l', 'l', 'o', 0x0A, 0x0B}
	b := bytes.NewBuffer(want)
	m, err := l.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	f, ok := m[ShortMessage]
	if !ok {
		t.Fatalf("missing %q key: %#v", ShortMessage, m)
	}
	v, ok := f.(*SM)
	if !ok {
		t.Fatalf("field is not type SM: %#v", f)
	}
	want = []byte("hello")
	if !bytes.Equal(want, v.Bytes()) {
		t.Fatalf("unexpected data: want %q, have %q", want, v)
	}
}
