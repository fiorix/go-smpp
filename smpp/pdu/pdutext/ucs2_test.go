// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

import (
	"bytes"
	"testing"
)

func TestUCS2Encoder(t *testing.T) {
	want := []byte("\x00O\x00l\x00\xe1\x00 \x00m\x00u\x00n\x00d\x00\xe3\x00o")
	text := []byte("Olá mundão")
	s := UCS2(text)
	if s.Type() != 0x08 {
		t.Fatalf("Unexpected data type; want 0x08, have %d", s.Type())
	}
	have := s.Encode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}

func TestUCS2Decoder(t *testing.T) {
	want := []byte("Olá mundão")
	text := []byte("\x00O\x00l\x00\xe1\x00 \x00m\x00u\x00n\x00d\x00\xe3\x00o")
	s := UCS2(text)
	if s.Type() != 0x08 {
		t.Fatalf("Unexpected data type; want 0x08, have %d", s.Type())
	}
	have := s.Decode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}
