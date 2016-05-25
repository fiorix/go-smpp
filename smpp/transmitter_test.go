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
	port := 0 // any port
	s := NewUnstartedServer(DefaultUser, DefaultPasswd, NewLocalListener(port))
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
		User:   DefaultUser,
		Passwd: DefaultPasswd,
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

func TestShortMessageWindowSize(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		time.Sleep(200 * time.Millisecond)
		r := pdu.NewSubmitSMResp()
		r.Header().Seq = p.Header().Seq
		r.Fields().Set(pdufield.MessageID, "foobar")
		c.Write(r)
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:        s.Addr(),
		User:        smpptest.DefaultUser,
		Passwd:      smpptest.DefaultPasswd,
		WindowSize:  2,
		RespTimeout: time.Second,
	}
	defer tx.Close()
	conn := <-tx.Bind()
	switch conn.Status() {
	case Connected:
	default:
		t.Fatal(conn.Error())
	}
	msgc := make(chan *ShortMessage, 3)
	defer close(msgc)
	errc := make(chan error, 3)
	for i := 0; i < 3; i++ {
		go func(msgc chan *ShortMessage, errc chan error) {
			m := <-msgc
			if m == nil {
				return
			}
			_, err := tx.Submit(m)
			errc <- err
		}(msgc, errc)
		msgc <- &ShortMessage{
			Src:      "root",
			Dst:      "foobar",
			Text:     pdutext.Raw("Lorem ipsum"),
			Validity: 10 * time.Minute,
			Register: NoDeliveryReceipt,
		}
	}
	nerr := 0
	for i := 0; i < 3; i++ {
		if <-errc == ErrMaxWindowSize {
			nerr++
		}
	}
	if nerr != 1 {
		t.Fatalf("unexpected # of errors. want 1, have %d", nerr)
	}
}

func TestLongMessage(t *testing.T) {
	port := 0 // any port
	s := NewUnstartedServer(DefaultUser, DefaultPasswd, NewLocalListener(port))
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
		User:   DefaultUser,
		Passwd: DefaultPasswd,
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
	port := 0 // any port
	s := NewUnstartedServer(DefaultUser, DefaultPasswd, NewLocalListener(port))
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
		User:   DefaultUser,
		Passwd: DefaultPasswd,
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
