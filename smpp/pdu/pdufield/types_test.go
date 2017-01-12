// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"strconv"
	"testing"
)

func TestFixed(t *testing.T) {
	f := &Fixed{Data: 0x34}
	if f.Len() != 1 {
		t.Fatalf("unexpected len: want 1, have %d", f.Len())
	}
	if v, ok := f.Raw().(uint8); !ok {
		t.Fatalf("unexpected type: want uint8, have %#v", v)
	}
	ws := strconv.Itoa(0x34)
	if v := f.String(); v != string(ws) {
		t.Fatalf("unexpected string: want %q, have %q", ws, v)
	}
	wb := []byte{0x34}
	if v := f.Bytes(); !bytes.Equal(wb, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", wb, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(wb, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", wb, v)
	}
}

func TestVariable(t *testing.T) {
	want := []byte("foobar")
	f := &Variable{Data: want}
	lw := len(want) + 1
	if f.Len() != lw {
		t.Fatalf("unexpected len: want %d, have %d", lw, f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string(want) {
		t.Fatalf("unexpected string: want %q have %q", want, v)
	}
	want = []byte("foobar\x00")
	if v := f.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", want, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", want, v)
	}
}

func TestSM(t *testing.T) {
	want := []byte("foobar")
	f := &SM{Data: want}
	if f.Len() != len(want) {
		t.Fatalf("unexpected len: want %d, have %d", len(want), f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string(want) {
		t.Fatalf("unexpected string: want %q have %q", want, v)
	}
	if v := f.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", want, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", want, v)
	}
}

func TestDestSme(t *testing.T) {
	var want []byte
	want = append(want, byte(0x01))        // flag
	want = append(want, byte(0x01))        // ton
	want = append(want, byte(0x01))        // npi
	want = append(want, []byte("1234")...) // Address
	want = append(want, byte(0x00))        // null terminator

	flag := Fixed{Data: byte(0x01)}
	ton := Fixed{Data: byte(0x01)}
	npi := Fixed{Data: byte(0x01)}
	destAddr := Variable{Data: []byte("1234")}
	fieldLen := flag.Len() + ton.Len() + npi.Len() + destAddr.Len()
	strRep := flag.String() + "," + ton.String() + "," + npi.String() + "," + destAddr.String()

	f := &DestSme{Flag: flag, Ton: ton, Npi: npi, DestAddr: destAddr}
	if f.Len() != fieldLen {
		t.Fatalf("unexpected len: want %d, have %d", fieldLen, f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string(strRep) {
		t.Fatalf("unexpected string: want %q have %q", strRep, v)
	}
	if v := f.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", want, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", want, v)
	}
}

func TestDestSmeList(t *testing.T) {
	flag := Fixed{Data: byte(0x01)}
	ton := Fixed{Data: byte(0x01)}
	npi := Fixed{Data: byte(0x01)}
	destAddr := Variable{Data: []byte("1234")}
	destAddr2 := Variable{Data: []byte("5678")}

	sme1 := DestSme{Flag: flag, Ton: ton, Npi: npi, DestAddr: destAddr}
	sme2 := DestSme{Flag: flag, Ton: ton, Npi: npi, DestAddr: destAddr2}
	fieldLen := sme1.Len() + sme2.Len()
	strRep := sme1.String() + ";" + sme2.String() + ";"
	var bytesRep []byte
	bytesRep = append(bytesRep, sme1.Bytes()...)
	bytesRep = append(bytesRep, sme2.Bytes()...)

	f := &DestSmeList{Data: []DestSme{sme1, sme2}}
	if f.Len() != fieldLen {
		t.Fatalf("unexpected len: want %d, have %d", fieldLen, f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string(strRep) {
		t.Fatalf("unexpected string: want %q have %q", strRep, v)
	}
	if v := f.Bytes(); !bytes.Equal(bytesRep, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", bytesRep, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(bytesRep, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", bytesRep, v)
	}
}

func TestUnSme(t *testing.T) {
	err := []byte{0x00, 0x00, 0x00, 0x11}
	var want []byte
	want = append(want, byte(0x01))       // TON
	want = append(want, byte(0x01))       // NPI
	want = append(want, []byte("123")...) // Address
	want = append(want, byte(0x00))       // null terminator
	want = append(want, err...)           // Error
	want = append(want, byte(0x00))       // null terminator

	ton := Fixed{Data: byte(0x01)}
	npi := Fixed{Data: byte(0x01)}
	destAddr := Variable{Data: []byte("123")}
	errCode := Variable{Data: err}
	fieldLen := ton.Len() + npi.Len() + destAddr.Len() + errCode.Len()
	strRep := ton.String() + "," + npi.String() + "," + destAddr.String() + "," + strconv.Itoa(17) // convertion to uint

	f := UnSme{Ton: ton, Npi: npi, DestAddr: destAddr, ErrCode: errCode}
	if f.Len() != fieldLen {
		t.Fatalf("unexpected len: want %d, have %d", fieldLen, f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string(strRep) {
		t.Fatalf("unexpected string: want %q have %q", strRep, v)
	}
	if v := f.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", want, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(want, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", want, v)
	}
}

func TestUnSmeList(t *testing.T) {
	err := []byte{0x00, 0x00, 0x00, 0x11}
	ton := Fixed{Data: byte(0x01)}
	npi := Fixed{Data: byte(0x01)}
	destAddr := Variable{Data: []byte("123")}
	destAddr2 := Variable{Data: []byte("456")}
	errCode := Variable{Data: err}

	unSme1 := UnSme{Ton: ton, Npi: npi, DestAddr: destAddr, ErrCode: errCode}
	unSme2 := UnSme{Ton: ton, Npi: npi, DestAddr: destAddr2, ErrCode: errCode}
	fieldLen := unSme1.Len() + unSme2.Len()
	strRep := unSme1.String() + ";" + unSme2.String() + ";"
	var bytesRep []byte
	bytesRep = append(bytesRep, unSme1.Bytes()...)
	bytesRep = append(bytesRep, unSme2.Bytes()...)

	f := &UnSmeList{Data: []UnSme{unSme1, unSme2}}
	if f.Len() != fieldLen {
		t.Fatalf("unexpected len: want %d, have %d", fieldLen, f.Len())
	}
	if v, ok := f.Raw().([]byte); !ok {
		t.Fatalf("unexpected type: want []byte, have %#v", v)
	}
	if v := f.String(); v != string(strRep) {
		t.Fatalf("unexpected string: want %q have %q", strRep, v)
	}
	if v := f.Bytes(); !bytes.Equal(bytesRep, v) {
		t.Fatalf("unexpected bytes: want %q, have %q", bytesRep, v)
	}
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	if v := b.Bytes(); !bytes.Equal(bytesRep, v) {
		t.Fatalf("unexpected serialized bytes: want %q, have %q", bytesRep, v)
	}
}
