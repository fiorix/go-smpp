// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

import (
	"bytes"
	"io/ioutil"
	"testing"
)

const (
	iso88595TypeCode = 0x06
	testDataDir      = "testdata"
)

var (
	iso88595Bytes     = readBytesFromFile(testDataDir + "/iso88595_test.txt")
	iso88595UTF8Bytes = readBytesFromFile(testDataDir + "/iso88595_test_utf8.txt")
)

func TestISO88595Encoder(t *testing.T) {
	want := iso88595Bytes
	text := iso88595UTF8Bytes
	s := ISO88595(text)
	if s.Type() != iso88595TypeCode {
		t.Fatalf("Unexpected data type; want %d, have %d", iso88595TypeCode, s.Type())
	}
	have := s.Encode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}

func TestISO88595Decoder(t *testing.T) {
	want := iso88595UTF8Bytes
	text := iso88595Bytes
	s := ISO88595(text)
	if s.Type() != iso88595TypeCode {
		t.Fatalf("Unexpected data type; want %d, have %d", iso88595TypeCode, s.Type())
	}
	have := s.Decode()
	if !bytes.Equal(want, have) {
		t.Fatalf("Unexpected text; want %q, have %q", want, have)
	}
}

func readBytesFromFile(filename string) []byte {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Error reading testdata file; " + filename + ", err " + err.Error())
	} else {
		return dat
	}
}
