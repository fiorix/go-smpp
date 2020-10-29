// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"fmt"
	"testing"
	"time"

	"github.com/tsocial/go-smpp/smpp/pdu"
	"github.com/tsocial/go-smpp/smpp/pdu/pdufield"
	"github.com/tsocial/go-smpp/smpp/pdu/pdutext"
	"github.com/tsocial/go-smpp/smpp/smpptest"
	"golang.org/x/time/rate"
)

func TestTransceiver(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		if p.Header().ID == pdu.SubmitSMID {
			r := pdu.NewSubmitSMResp()
			r.Header().Seq = p.Header().Seq
			_ = r.Fields().Set(pdufield.MessageID, "foobar")
			_ = c.Write(r)
			pf := p.Fields()
			rd := pf[pdufield.RegisteredDelivery]
			if rd.Bytes()[0] == 0 {
				return
			}
			r = pdu.NewDeliverSM()
			f := r.Fields()
			_ = f.Set(pdufield.SourceAddr, pf[pdufield.SourceAddr])
			_ = f.Set(pdufield.DestinationAddr, pf[pdufield.DestinationAddr])
			_ = f.Set(pdufield.ShortMessage, pf[pdufield.ShortMessage])
			_ = c.Write(r)
		} else {
			smpptest.EchoHandler(c, p)
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
		Addr:        s.Addr(),
		User:        smpptest.DefaultUser,
		Passwd:      smpptest.DefaultPasswd,
		Handler:     receiver,
		RateLimiter: rate.NewLimiter(rate.Limit(10), 1),
	}
	defer func() {
		_ = tc.Close()
	}()
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
		Register: pdufield.FinalDeliveryReceipt,
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
