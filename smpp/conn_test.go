// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"testing"

	"github.com/tsocial/go-smpp/smpp/pdu"
	"github.com/tsocial/go-smpp/smpp/pdu/pdufield"
	"github.com/tsocial/go-smpp/smpp/smpptest"
)

func TestConn(t *testing.T) {
	s := smpptest.NewServer()
	defer s.Close()
	c, err := Dial(s.Addr(), nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = c.Close()
	}()
	p := pdu.NewBindTransmitter()
	f := p.Fields()
	_ = f.Set(pdufield.SystemID, smpptest.DefaultUser)
	_ = f.Set(pdufield.Password, smpptest.DefaultPasswd)
	_ = f.Set(pdufield.InterfaceVersion, 0x34)
	if err = c.Write(p); err != nil {
		t.Fatal(err)
	}
	if _, err = c.Read(); err != nil {
		t.Fatal(err)
	}
}
