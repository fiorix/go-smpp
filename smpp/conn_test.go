// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"testing"

	"github.com/fiorix/go-smpp/v2/smpp/pdu"
	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/v2/smpp/smpptest"
)

func TestConn(t *testing.T) {
	s := smpptest.NewServer()
	defer s.Close()
	c, err := Dial(s.Addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	p := pdu.NewBindTransmitter()
	f := p.Fields()
	f.Set(pdufield.SystemID, smpptest.DefaultUser)
	f.Set(pdufield.Password, smpptest.DefaultPasswd)
	f.Set(pdufield.InterfaceVersion, 0x34)
	if err = c.Write(p); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Read(); err != nil {
		t.Fatal(err)
	}
}
