// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"errors"
	"testing"
	"time"

	"github.com/tsocial/go-smpp/smpp/pdu"
	"github.com/tsocial/go-smpp/smpp/pdu/pdufield"
	"github.com/tsocial/go-smpp/smpp/pdu/pdutext"
	"github.com/tsocial/go-smpp/smpp/smpptest"
	"golang.org/x/time/rate"
)

func TestShortMessage(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		if p.Header().ID == pdu.SubmitSMID {
			r := pdu.NewSubmitSMResp()
			r.Header().Seq = p.Header().Seq
			_ = r.Fields().Set(pdufield.MessageID, "foobar")
			_ = c.Write(r)
		} else {
			smpptest.EchoHandler(c, p)
		}
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:        s.Addr(),
		User:        smpptest.DefaultUser,
		Passwd:      smpptest.DefaultPasswd,
		RateLimiter: rate.NewLimiter(rate.Limit(10), 1),
	}
	defer func() {
		_ = tx.Close()
	}()
	conn := <-tx.Bind()
	if conn.Status() != Connected {
		t.Fatal(conn.Error())
	}
	sm, err := tx.Submit(&ShortMessage{
		Src:      "root",
		Dst:      "foobar",
		Text:     pdutext.Raw("Lorem ipsum"),
		Validity: 10 * time.Minute,
		Register: pdufield.NoDeliveryReceipt,
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
		_ = r.Fields().Set(pdufield.MessageID, "foobar")
		_ = c.Write(r)
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
	defer func() {
		_ = tx.Close()
	}()
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
			Register: pdufield.NoDeliveryReceipt,
		}
	}
	nerr := 0
	for i := 0; i < 3; i++ {
		if errors.Is(<-errc, ErrMaxWindowSize) {
			nerr++
		}
	}
	if nerr != 1 {
		t.Fatalf("unexpected # of errors. want 1, have %d", nerr)
	}
}

func TestLongMessage(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		if p.Header().ID == pdu.SubmitSMID {
			r := pdu.NewSubmitSMResp()
			r.Header().Seq = p.Header().Seq
			_ = r.Fields().Set(pdufield.MessageID, "foobar")
			_ = c.Write(r)
		} else {
			smpptest.EchoHandler(c, p)
		}
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:   s.Addr(),
		User:   smpptest.DefaultUser,
		Passwd: smpptest.DefaultPasswd,
	}
	defer func() {
		_ = tx.Close()
	}()
	conn := <-tx.Bind()
	if conn.Status() != Connected {
		t.Fatal(conn.Error())
	}
	sm, err := tx.SubmitLongMsg(&ShortMessage{
		Src:      "root",
		Dst:      "foobar",
		Text:     pdutext.Raw("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam consequat nisl enim, vel finibus neque aliquet sit amet. Interdum et malesuada fames ac ante ipsum primis in faucibus."),
		Validity: 10 * time.Minute,
		Register: pdufield.NoDeliveryReceipt,
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

func TestLongMessageAsUCS2(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	var receivedMsg string
	shortMsg := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam consequat nisl enim, vel finibus neque aliquet sit amet. Interdum et malesuada fames ac ante ipsum primis in faucibus. ✓"
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		switch p.Header().ID {
		case pdu.SubmitSMID:
			r := pdu.NewSubmitSMResp()
			r.Header().Seq = p.Header().Seq
			_ = r.Fields().Set(pdufield.MessageID, "foobar")
			smByts := p.Fields()[pdufield.ShortMessage].Bytes()
			switch pdutext.DataCoding(p.Fields()[pdufield.DataCoding].Raw().(uint8)) {
			case pdutext.Latin1Type:
				receivedMsg += string(pdutext.Latin1(smByts)[7:].Decode())
			case pdutext.UCS2Type:
				receivedMsg += string(pdutext.UCS2(smByts)[7:].Decode())
			default:
				receivedMsg += string(smByts[7:])
			}
			_ = c.Write(r)
		default:
			smpptest.EchoHandler(c, p)
		}
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:   s.Addr(),
		User:   smpptest.DefaultUser,
		Passwd: smpptest.DefaultPasswd,
	}
	defer func() {
		_ = tx.Close()
	}()
	conn := <-tx.Bind()
	if conn.Status() != Connected {
		t.Fatal(conn.Error())
	}
	sm, err := tx.SubmitLongMsg(&ShortMessage{
		Src:      "root",
		Dst:      "foobar",
		Text:     pdutext.UCS2(shortMsg),
		Validity: 10 * time.Minute,
		Register: pdufield.NoDeliveryReceipt,
	})
	if err != nil {
		t.Fatal(err)
	}
	msgid := sm.RespID()
	if msgid == "" {
		t.Fatalf("pdu does not contain msgid: %#v", sm.Resp())
	}

	if receivedMsg != shortMsg {
		t.Fatalf("receivedMsg: %v, does not match shortMsg: %v", receivedMsg, shortMsg)
	}

	if msgid != "foobar" {
		t.Fatalf("unexpected msgid: want foobar, have %q", msgid)
	}
}

func TestQuerySM(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		r := pdu.NewQuerySMResp()
		r.Header().Seq = p.Header().Seq
		_ = r.Fields().Set(pdufield.MessageID, p.Fields()[pdufield.MessageID])
		_ = r.Fields().Set(pdufield.MessageState, 2)
		_ = c.Write(r)
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:   s.Addr(),
		User:   smpptest.DefaultUser,
		Passwd: smpptest.DefaultPasswd,
	}
	defer func() {
		_ = tx.Close()
	}()
	conn := <-tx.Bind()
	switch conn.Status() {
	case Connected:
	default:
		t.Fatal(conn.Error())
	}
	qr, err := tx.QuerySM("root", "13", uint8(5), uint8(0))
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

func TestSubmitMulti(t *testing.T) {
	// construct a byte array with the UnsuccessSme
	var bArray []byte
	bArray = append(bArray, byte(0x00))       // TON
	bArray = append(bArray, byte(0x00))       // NPI
	bArray = append(bArray, []byte("123")...) // Address
	bArray = append(bArray, byte(0x00))       // Error
	bArray = append(bArray, byte(0x00))       // Error
	bArray = append(bArray, byte(0x00))       // Error
	bArray = append(bArray, byte(0x11))       // Error
	bArray = append(bArray, byte(0x00))       // null terminator

	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		if p.Header().ID == pdu.SubmitMultiID {
			r := pdu.NewSubmitMultiResp()
			r.Header().Seq = p.Header().Seq
			_ = r.Fields().Set(pdufield.MessageID, "foobar")
			_ = r.Fields().Set(pdufield.NoUnsuccess, uint8(1))
			_ = r.Fields().Set(pdufield.UnsuccessSme, bArray)
			_ = c.Write(r)
		} else {
			smpptest.EchoHandler(c, p)
		}
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:   s.Addr(),
		User:   smpptest.DefaultUser,
		Passwd: smpptest.DefaultPasswd,
	}
	defer func() {
		_ = tx.Close()
	}()
	conn := <-tx.Bind()
	if conn.Status() != Connected {
		t.Fatal(conn.Error())
	}
	sm, err := tx.Submit(&ShortMessage{
		Src:      "root",
		DstList:  []string{"123", "2233", "32322", "4234234"},
		DLs:      []string{"DistributionList1"},
		Text:     pdutext.Raw("Lorem ipsum"),
		Validity: 10 * time.Minute,
		Register: pdufield.NoDeliveryReceipt,
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
	noUncess, _ := sm.NumbUnsuccess()
	if noUncess != 1 {
		t.Fatalf("unexpected number of unsuccess %d", noUncess)
	}
	uncessSmes, _ := sm.UnsuccessSmes()
	if len(uncessSmes) != 1 {
		t.Fatalf("unsucess sme list should have a size of 1, has %d", len(uncessSmes))
	}
}

func TestNotConnected(t *testing.T) {
	s := smpptest.NewUnstartedServer()
	s.Handler = func(c smpptest.Conn, p pdu.Body) {
		if p.Header().ID == pdu.SubmitSMID {
			r := pdu.NewSubmitSMResp()
			r.Header().Seq = p.Header().Seq
			_ = r.Fields().Set(pdufield.MessageID, "foobar")
			_ = c.Write(r)
		} else {
			smpptest.EchoHandler(c, p)
		}
	}
	s.Start()
	defer s.Close()
	tx := &Transmitter{
		Addr:   s.Addr(),
		User:   smpptest.DefaultUser,
		Passwd: smpptest.DefaultPasswd,
	}
	// Open connection and then close it
	conn := <-tx.Bind()
	if conn.Status() != Connected {
		t.Fatal(conn.Error())
	}
	_ = tx.Close()
	_, err := tx.Submit(&ShortMessage{
		Src:      "root",
		Dst:      "foobar",
		Text:     pdutext.Raw("Lorem ipsum"),
		Validity: 10 * time.Minute,
		Register: pdufield.NoDeliveryReceipt,
	})
	if !errors.Is(err, ErrNotConnected) {
		t.Fatalf("Error should be not connect, got %s", err.Error())
	}
}
