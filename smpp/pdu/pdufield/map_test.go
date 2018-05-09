// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"testing"

	"github.com/tsocial/go-smpp/smpp/pdu/pdutext"
)

func TestMapSet(t *testing.T) {
	m := make(Map)
	test := []struct {
		k  Name
		v  interface{}
		ok bool
	}{
		{SystemID, nil, true},
		{SystemID, "hello", true},
		{SystemID, []byte("hello"), true},
		{DataCoding, nil, true},
		{DataCoding, uint8(1), true},
		{DataCoding, int(1), true},
		{DataCoding, t, false},
		{DataCoding, New(DataCoding, []byte{0x03}), true},
	}
	for _, el := range test {
		if err := m.Set(el.k, el.v); el.ok && err != nil {
			t.Fatal(err)
		} else if !el.ok && err == nil {
			t.Fatalf("unexpected set of %q=%#v", el.k, el.v)
		}
	}
}

func TestMapSetTextCodec(t *testing.T) {
	m := make(Map)
	text := pdutext.Latin1("Olá mundo")
	err := m.Set(ShortMessage, text)
	if err != nil {
		t.Fatal(err)
	}
	dc, exists := m[DataCoding]
	if !exists {
		t.Fatal("missing data_coding pdu")
	}
	dv, ok := dc.(*Fixed)
	if !ok {
		t.Fatalf("unexpected type for data_coding: %#v", dc)
	}
	if dv.Data != uint8(text.Type()) {
		t.Fatalf("unexpected value for data_coding: want %d, have %d",
			text.Type(), dv.Data)
	}
	pt, exists := m[ShortMessage]
	if !exists {
		t.Fatal("missing short_message pdu")
	}
	nt := pdutext.Latin1(pt.Bytes()).Decode()
	if !bytes.Equal(text, nt) {
		t.Fatalf("unexpected text: want %q, have %q", text, nt)
	}
}
