// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

import (
    "testing"
    "bytes"
)

func TestGSM7PackedEncoder(t *testing.T) {
    want := []byte("\xC8\x32\x9B\xFD\x06\xDD\xDF\x72\x36\x19")
    text := []byte("Hello world")
    s := GSM7Packed(text)
    if s.Type() != 0x00 {
        t.Fatalf("Unexpected data type; want 0x00, have %d", s.Type())
    }
    have := s.Encode()
    if !bytes.Equal(want, have) {
        t.Fatalf("Unexpected text; want %q, have %q", want, have)
    }
}

func TestGSM7PackedDecoder(t *testing.T) {
    want := []byte("Hello world")
    text := []byte("\xC8\x32\x9B\xFD\x06\xDD\xDF\x72\x36\x19")
    s := GSM7Packed(text)
    if s.Type() != 0x00 {
        t.Fatalf("Unexpected data type; want 0x00, have %d", s.Type())
    }
    have := s.Decode()
    if !bytes.Equal(want, have) {
        t.Fatalf("Unexpected text; want %q, have %q", want, have)
    }
}
