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

func TestListDecoder_DestinationList(t *testing.T) {
	l := List{NumberDests, DestinationList}
	want := []byte{0x02, 0x01, 0x01, 0x01, '1', '2', '3', 0x00, 0x01, 0x01, 0x01, '5', '6', '7', 0x00}

	b := bytes.NewBuffer(want)
	m, err := l.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	f, ok := m[DestinationList]
	if !ok {
		t.Fatalf("missing %q key: %#v", DestinationList, m)
	}
	v, ok := f.(*DestSmeList)
	if !ok {
		t.Fatalf("field is not type DestSmeList: %#v", f)
	}

	flag := Fixed{Data: byte(0x01)}
	ton := Fixed{Data: byte(0x01)}
	npi := Fixed{Data: byte(0x01)}
	destAddr := Variable{Data: []byte("123")}
	destAddr2 := Variable{Data: []byte("567")}
	sme1 := DestSme{Flag: flag, Ton: ton, Npi: npi, DestAddr: destAddr}
	sme2 := DestSme{Flag: flag, Ton: ton, Npi: npi, DestAddr: destAddr2}
	resSmeList := &DestSmeList{Data: []DestSme{sme1, sme2}}

	if !bytes.Equal(resSmeList.Bytes(), v.Bytes()) {
		t.Fatalf("unexpected data: want %q, have %q, len %d", resSmeList, v, len(v.Data))
	}
}

func TestListDecoder_UnSmeList(t *testing.T) {
	l := List{NoUnsuccess, UnsuccessSme}
	want := []byte{0x01, 0x01, 0x01, '1', '2', '3', 0x00, 0x00, 0x00, 0x00, 0x11, 0x00}

	b := bytes.NewBuffer(want)
	m, err := l.Decode(b)
	if err != nil {
		t.Fatal(err)
	}
	f, ok := m[UnsuccessSme]
	if !ok {
		t.Fatalf("missing %q key: %#v", UnsuccessSme, m)
	}
	v, ok := f.(*UnSmeList)
	if !ok {
		t.Fatalf("field is not type UnSmeList: %#v", f)
	}

	errC := []byte{0x00, 0x00, 0x00, 0x11}
	ton := Fixed{Data: byte(0x01)}
	npi := Fixed{Data: byte(0x01)}
	destAddr := Variable{Data: []byte("123")}
	errCode := Variable{Data: errC}
	unSme1 := UnSme{Ton: ton, Npi: npi, DestAddr: destAddr, ErrCode: errCode}
	resUnSmeList := &UnSmeList{Data: []UnSme{unSme1}}
	if !bytes.Equal(resUnSmeList.Bytes(), v.Bytes()) {
		t.Fatalf("unexpected data: want %q, have %q, len %d", resUnSmeList, v, len(v.Data))
	}
}
