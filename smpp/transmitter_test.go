// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"testing"
	"time"

	"github.com/veoo/go-smpp/smpp/pdu"
	"github.com/veoo/go-smpp/smpp/pdu/pdufield"
	"github.com/veoo/go-smpp/smpp/pdu/pdutext"
)

func TestShortMessage(t *testing.T) {
	pass	:= "secret"
	user    := "client"
	port := 0 // any port
	s := NewUnstartedServer(user, pass, NewLocalListener(port))
	s.Handler = func(s Session, p pdu.Body) {
		switch p.Header().ID {
		case pdu.SubmitSMID:
			r := pdu.NewSubmitSMResp()
			r.Header().Seq = p.Header().Seq
			r.Fields().Set(pdufield.MessageID, "foobar")
			s.Write(r)
		default:
			EchoHandler(s, p)
		}
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:   s.Addr(),
		User:   user,
		Passwd: pass,
	}
	defer tx.Close()
	conn := <-tx.Bind()
	switch conn.Status() {
	case Connected:
	default:
		t.Fatal(conn.Error())
	}
	sm, err := tx.Submit(&ShortMessage{
		Src:      "root",
		Dst:      "foobar",
		Text:     pdutext.Raw("Lorem ipsum"),
		Validity: 10 * time.Minute,
		Register: NoDeliveryReceipt,
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
}

func TestLongMessage(t *testing.T) {
	pass	:= "secret"
	user    := "client"
	port 	:= 0 // any port
	s := NewUnstartedServer(user, pass, NewLocalListener(port))
	s.Handler = func(s Session, p pdu.Body) {
		switch p.Header().ID {
		case pdu.SubmitSMID:
			r := pdu.NewSubmitSMResp()
			r.Header().Seq = p.Header().Seq
			r.Fields().Set(pdufield.MessageID, "foobar")
			s.Write(r)
		default:
			EchoHandler(s, p)
		}
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:   s.Addr(),
		User:  	user,
		Passwd: pass,
	}
	defer tx.Close()
	conn := <-tx.Bind()
	switch conn.Status() {
	case Connected:
	default:
		t.Fatal(conn.Error())
	}
	sm, err := tx.SubmitLongMsg(&ShortMessage{
		Src:      "root",
		Dst:      "foobar",
		Text:     pdutext.Raw("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam consequat nisl enim, vel finibus neque aliquet sit amet. Interdum et malesuada fames ac ante ipsum primis in faucibus."),
		Validity: 10 * time.Minute,
		Register: NoDeliveryReceipt,
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
}

func TestQuerySM(t *testing.T) {
	pass	:= "secret"
	user    := "client"
	port 	:= 0 // any port
	s := NewUnstartedServer(user, pass, NewLocalListener(port))
	s.Handler = func(s Session, p pdu.Body) {
		r := pdu.NewQuerySMResp()
		r.Header().Seq = p.Header().Seq
		r.Fields().Set(pdufield.MessageID, p.Fields()[pdufield.MessageID])
		r.Fields().Set(pdufield.MessageState, 2)
		s.Write(r)
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:   s.Addr(),
		User:   user,
		Passwd: pass,
	}
	defer tx.Close()
	conn := <-tx.Bind()
	switch conn.Status() {
	case Connected:
	default:
		t.Fatal(conn.Error())
	}
	qr, err := tx.QuerySM("root", "13")
	if err != nil {
		t.Fatal(err)
	}
	if qr.MsgID != "13" {
		t.Fatalf("unexpected msgid: want 13, have %s", qr.MsgID)
	}
	if qr.MsgState != "DELIVERED" {
		t.Fatalf("unexpected state: want DELIVERED, have %q", qr.MsgState)
	}
}
