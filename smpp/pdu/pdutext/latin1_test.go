// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

import (
	"bytes"
	"testing"
)

func TestLatin1Encoder(t *testing.T) {
	want := []byte("Ol\xe1 mund\xe3o")
	text := []byte("Olá mundão")
	s := Latin1(text)
	if s.Type() != 0x03 {
		t.Fatalf("Unexpected data type; want 0x03, have %d", s.Type())
	}
	have := s.Encode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}

func TestLatin1Decoder(t *testing.T) {
	want := []byte("Olá mundão")
	text := []byte("Ol\xe1 mund\xe3o")
	s := Latin1(text)
	if s.Type() != 0x03 {
		t.Fatalf("Unexpected data type; want 0x03, have %d", s.Type())
	}
	have := s.Decode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}
