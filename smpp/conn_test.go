// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"testing"

	"github.com/veoo/go-smpp/smpp/pdu"
	"github.com/veoo/go-smpp/smpp/pdu/pdufield"
)

func TestConn(t *testing.T) {
	pass	:= "secret"
	user    := "client"
	port 	:= 0 // any port
	s := NewServer(user, pass, NewLocalListener(port))
	defer s.Close()
	c, err := Dial(s.Addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	p := pdu.NewBindTransmitter()
	f := p.Fields()
	f.Set(pdufield.SystemID, user)
	f.Set(pdufield.Password, pass)
	f.Set(pdufield.InterfaceVersion, 0x34)
	if err = c.Write(p); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Read(); err != nil {
		t.Fatal(err)
	}
}
