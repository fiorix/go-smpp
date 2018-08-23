// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

import (
    "testing"
    "bytes"
)

func TestGSM7Encoder(t *testing.T) {
    want := []byte("\x48\x65\x6C\x6C\x6F \x77\x6F\x72\x6C\x64")
    text := []byte("Hello world")
    s := GSM7(text)
    if s.Type() != 0x00 {
        t.Fatalf("Unexpected data type; want 0x00, have %d", s.Type())
    }
    have := s.Encode()
    if !bytes.Equal(want, have) {
        t.Fatalf("Unexpected text; want %q, have %q", want, have)
    }
}

func TestGSM7Decoder(t *testing.T) {
    want := []byte("Hello world")
    text := []byte("\x48\x65\x6C\x6C\x6F \x77\x6F\x72\x6C\x64")
    s := GSM7(text)
    if s.Type() != 0x00 {
        t.Fatalf("Unexpected data type; want 0x00, have %d", s.Type())
    }
    have := s.Decode()
    if !bytes.Equal(want, have) {
        t.Fatalf("Unexpected text; want %q, have %q", want, have)
    }
}
