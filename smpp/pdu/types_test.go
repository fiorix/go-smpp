// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdu

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/tsocial/go-smpp/smpp/pdu/pdufield"
)

func TestBind(t *testing.T) {
	tx := []byte{
		0x00, 0x00, 0x00, 0x2A, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x73, 0x6D, 0x70, 0x70, 0x63, 0x6C, 0x69, 0x65,
		0x6E, 0x74, 0x31, 0x00, 0x70, 0x61, 0x73, 0x73,
		0x77, 0x6F, 0x72, 0x64, 0x00, 0x00, 0x34, 0x00,
		0x00, 0x00,
	}
	pdu := NewBindTransmitter()
	f := pdu.Fields()
	_ = f.Set(pdufield.SystemID, "smppclient1")
	_ = f.Set(pdufield.Password, "password")
	_ = f.Set(pdufield.InterfaceVersion, 0x34)
	pdu.Header().Seq = 1
	var b bytes.Buffer
	if err := pdu.SerializeTo(&b); err != nil {
		t.Fatal(err)
	}
	l := uint32(b.Len())
	if l != pdu.Header().Len {
		t.Fatalf("unexpected len: want %d, have %d", l, pdu.Header().Len)
	}
	if !bytes.Equal(tx, b.Bytes()) {
		t.Fatalf("unexpected bytes:\nwant:\n%s\nhave:\n%s",
			hex.Dump(tx), hex.Dump(b.Bytes()))
	}
	pdu, err := Decode(&b)
	if err != nil {
		t.Fatal(err)
	}
	h := pdu.Header()
	if h.ID != BindTransmitterID {
		t.Fatalf("unexpected ID: want %d, have %d",
			BindTransmitterID, h.ID)
	}
	if h.Seq != 1 {
		t.Fatalf("unexpected Seq: want 1, have %d", h.Seq)
	}
	test := []struct {
		n pdufield.Name
		v string
	}{
		{pdufield.SystemID, "smppclient1"},
		{pdufield.Password, "password"},
		{pdufield.InterfaceVersion, strconv.Itoa(0x34)},
	}
	for _, el := range test {
		f := pdu.Fields()[el.n]
		if f == nil {
			t.Fatalf("missing field: %s", el.n)
		}
		if f.String() != el.v {
			t.Fatalf("unexpected value for %q: want %q, have %q",
				el.n, el.v, f.String())
		}
	}
}

/*
func TestBindResp(t *testing.T) {
	tx := []byte{
		0x00, 0x00, 0x00, 0x18, 0x80, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x53, 0x4D, 0x50, 0x50, 0x53, 0x69, 0x6D, 0x00,
	}
	t.Log(tx)
}
*/
