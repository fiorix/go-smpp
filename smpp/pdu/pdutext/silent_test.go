// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

import (
	"bytes"
	"testing"
)

func TestSilentEncoder(t *testing.T) {
	text := []byte("Olá mundão")
	want := []byte("")
	s := Silent(text)
	if s.Type() != 0xC0 {
		t.Fatalf("Unexpected data type; want 0x00, have %d", s.Type())
	}
	have := s.Encode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}

func TestSilentDecoder(t *testing.T) {
	text := []byte("Olá mundão")
	want := text
	s := Silent(text)
	if s.Type() != 0xC0 {
		t.Fatalf("Unexpected data type; want 0x00, have %d", s.Type())
	}
	have := s.Decode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}
