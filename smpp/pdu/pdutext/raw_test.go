// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

import (
	"bytes"
	"testing"
)

func TestRawEncoder(t *testing.T) {
	text := []byte("Olá mundão")
	want := text
	s := Raw(text)
	if s.Type() != 0x00 {
		t.Fatalf("Unexpected data type; want 0x00, have %d", s.Type())
	}
	have := s.Encode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}

func TestRawDecoder(t *testing.T) {
	text := []byte("Olá mundão")
	want := text
	s := Raw(text)
	if s.Type() != 0x00 {
		t.Fatalf("Unexpected data type; want 0x00, have %d", s.Type())
	}
	have := s.Decode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}
