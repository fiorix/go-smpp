// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"fmt"
	"testing"
	"time"

	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
	"github.com/fiorix/go-smpp/smpp/smpptest"
)

func TestTransceiver(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, m pdu.Body) {
		switch m.Header().ID {
		case pdu.SubmitSMID:
			p := pdu.NewSubmitSMResp()
			p.Header().Seq = m.Header().Seq
			p.Fields().Set(pdufield.MessageID, "foobar")
			c.Write(p)
			mf := m.Fields()
			rd := mf[pdufield.RegisteredDelivery]
			if rd.Bytes()[0] == 0 {
				return
			}
			p = pdu.NewDeliverSM()
			f := p.Fields()
			f.Set(pdufield.SourceAddr, mf[pdufield.SourceAddr])
			f.Set(pdufield.DestinationAddr, mf[pdufield.DestinationAddr])
			f.Set(pdufield.ShortMessage, mf[pdufield.ShortMessage])
			c.Write(p)
		default:
			smpptest.EchoHandler(c, m)
		}
	}
	s.Start()
	defer s.Close()
	ack := make(chan error)
	receiver := func(p pdu.Body) {
		defer close(ack)
		if p.Header().ID != pdu.DeliverSMID {
			ack <- fmt.Errorf("unexpected PDU: %s", p.Header().ID)
		}
	}
	tc := &Transceiver{
		Addr:    s.Addr(),
		User:    smpptest.DefaultUser,
		Passwd:  smpptest.DefaultPasswd,
		Handler: receiver,
	}
	defer tc.Close()
	conn := <-tc.Bind()
	switch conn.Status() {
	case Connected:
	default:
		t.Fatal(conn.Error())
	}
	sm, err := tc.Submit(&ShortMessage{
		Src:      "root",
		Dst:      "foobar",
		Text:     pdutext.Raw("Lorem ipsum"),
		Register: FinalDeliveryReceipt,
	})
	if err != nil {
		t.Fatal(err)
	}
	msgid := sm.RespID()
	if msgid == "" {
		t.Fatalf("pdu does not contain msgid: %#v", sm.Resp())
	}
	if msgid != "foobar" {
		t.Fatalf("unexpected msgid: want foobar, have %q", msgid)
	}
	select {
	case err := <-ack:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for ack")
	}
}
