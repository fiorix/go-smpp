// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"bytes"
	"net"
	"testing"

	"github.com/veoo/go-smpp/smpp/pdu"
	"github.com/veoo/go-smpp/smpp/pdu/pdufield"
	"github.com/veoo/go-smpp/smpp/pdu/pdutext"
)

func TestServer(t *testing.T) {
	port := 0 // any port
	s := NewServer(DefaultUser, DefaultPasswd, NewLocalListener(port))
	defer s.Close()
	c, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	rw := newConn(c)
	// bind
	p := pdu.NewBindTransmitter()
	f := p.Fields()
	f.Set(pdufield.SystemID, DefaultUser)
	f.Set(pdufield.Password, DefaultPasswd)
	f.Set(pdufield.InterfaceVersion, 0x34)
	if err = rw.Write(p); err != nil {
		t.Fatal(err)
	}
	// bind resp
	resp, err := rw.Read()
	if err != nil {
		t.Fatal(err)
	}
	id, ok := resp.Fields()[pdufield.SystemID]
	if !ok {
		t.Fatalf("missing system_id field: %#v", resp)
	}
	if id.String() != "smpptest" {
		t.Fatalf("unexpected system_id: want smpptest, have %q", id)
	}
	// submit_sm
	p = pdu.NewSubmitSM()
	f = p.Fields()
	f.Set(pdufield.SourceAddr, "foobar")
	f.Set(pdufield.DestinationAddr, "bozo")
	f.Set(pdufield.ShortMessage, pdutext.Latin1("Lorem ipsum"))
	if err = rw.Write(p); err != nil {
		t.Fatal(err)
	}
	// same submit_sm
	r, err := rw.Read()
	if err != nil {
		t.Fatal(err)
	}
	want, have := *p.Header(), *r.Header()
	if want != have {
		t.Fatalf("unexpected header: want %#v, have %#v", want, have)
	}
	for k, v := range p.Fields() {
		vv, exists := r.Fields()[k]
		if !exists {
			t.Fatalf("unexpected fields: want %#v, have %#v",
				p.Fields(), r.Fields())
		}
		if !bytes.Equal(v.Bytes(), vv.Bytes()) {
			t.Fatalf("unexpected field data: want %#v, have %#v",
				v, vv)
		}
	}
}
