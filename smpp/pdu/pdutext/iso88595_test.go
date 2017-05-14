// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

import (
	"bytes"
	"io/ioutil"
	"testing"
)

var (
	ISO88595_Bytes []byte = readBytesFromFile("test_text/iso88595_test.txt")
	Utf8_Bytes     []byte = readBytesFromFile("test_text/iso88595_test_utf8.txt")
)

func TestISO88595Encoder(t *testing.T) {
	want := []byte(ISO88595_Bytes)
	text := []byte(Utf8_Bytes)
	s := ISO88595(text)
	if s.Type() != 0x06 {
		t.Fatalf("Unexpected data type; want 0x03, have %d", s.Type())
	}
	have := s.Encode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}

func TestISO88595Decoder(t *testing.T) {
	want := []byte(Utf8_Bytes)
	text := []byte(ISO88595_Bytes)
	s := ISO88595(text)
	if s.Type() != 0x06 {
		t.Fatalf("Unexpected data type; want 0x03, have %d", s.Type())
	}
	have := s.Decode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}

func readBytesFromFile(aFileName string) []byte {
	dat, err := ioutil.ReadFile(aFileName)
	if err != nil {
		return nil
	} else {
		return dat
	}
}
