// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"bytes"
	"errors"
	"net"
	"os"
	"testing"

	"github.com/veoo/go-smpp/smpp/pdu"
	"github.com/veoo/go-smpp/smpp/pdu/pdufield"
	"github.com/veoo/go-smpp/smpp/pdu/pdutext"
)

var (
	s          Server
	pass       = "secret"
	user       = "client"
	customUser = "customUser"
	customPass = "customPass"
	port       = 0 // any port
)

func BindTransceiverHandler(s Session, m pdu.Body) error {
	f := m.Fields()
	user := f[pdufield.SystemID]
	passwd := f[pdufield.Password]
	if user == nil || passwd == nil {
		return errors.New("malformed pdu, missing system_id/password")
	}
	if user.String() != customUser {
		return errors.New("invalid user")
	}
	if passwd.String() != customPass {
		return errors.New("invalid passwd")
	}
	resp := pdu.NewBindTransceiverResp()
	resp.Fields().Set(pdufield.SystemID, DefaultSystemID)
	return s.Write(resp)
}

func TestMain(m *testing.M) {
	s = NewServer(user, pass, NewLocalListener(port))
	s.Handle(pdu.BindTransmitterID, EchoHandler)
	s.Handle(pdu.SubmitSMID, EchoHandler)
	s.HandleAuth(pdu.BindTransceiverID, BindTransceiverHandler)

	defer s.Close()
	os.Exit(m.Run())
}

func TestServer(t *testing.T) {
	c, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	rw := newConn(c)
	// bind
	p := pdu.NewBindTransmitter()
	f := p.Fields()
	f.Set(pdufield.SystemID, user)
	f.Set(pdufield.Password, pass)
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
	if id.String() != "sys_id" {
		t.Fatalf("unexpected system_id: want sys_id, have %q", id)
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

func TestIncorrectAuth(t *testing.T) {
	c, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	rw := newConn(c)
	// bind
	p := pdu.NewBindTransceiver()
	f := p.Fields()

	// Test with incorrect credentials that would work with defaultAuth()
	f.Set(pdufield.SystemID, user)
	f.Set(pdufield.Password, pass)
	f.Set(pdufield.InterfaceVersion, 0x34)
	if err = rw.Write(p); err != nil {
		t.Fatal(err)
	}
	// bind resp
	resp, err := rw.Read()
	if err == nil {
		t.Fatalf("authenticated with incorrect credentials: %#v", resp)
	}
}

func TestCorrectAuth(t *testing.T) {
	c, err := net.Dial("tcp", s.Addr())
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	rw := newConn(c)
	// bind
	p := pdu.NewBindTransceiver()
	f := p.Fields()

	// Test with correct credentials that BindTransceiverHandler accepts
	f.Set(pdufield.SystemID, customUser)
	f.Set(pdufield.Password, customPass)
	if err := rw.Write(p); err != nil {
		t.Fatal(err)
	}
	// bind resp
	resp, err := rw.Read()
	if err != nil {
		t.Fatalf("failed to read resp: %v", err)
	}
	id, ok := resp.Fields()[pdufield.SystemID]
	if !ok {
		t.Fatalf("missing system_id field: %#v", resp)
	}
	if id.String() != "sys_id" {
		t.Fatalf("unexpected system_id: want sys_id, have %q", id)
	}
}
