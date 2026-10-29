// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdu

import (
	"bytes"
	"encoding/hex"
	"strconv"
	"testing"

	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdufield"
)

func TestBind(t *testing.T) {
	tx := []byte{
		0x00, 0x00, 0x00, 0x2A, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x73, 0x6D, 0x70, 0x70, 0x63, 0x6C, 0x69, 0x65,
		0x6E, 0x74, 0x31, 0x00, 0x70, 0x61, 0x73, 0x73,
		0x77, 0x6F, 0x72, 0x64, 0x00, 0x00, 0x34, 0x00,
		0x00, 0x00,
	}
	pdu := NewBindTransmitter()
	f := pdu.Fields()
	f.Set(pdufield.SystemID, "smppclient1")
	f.Set(pdufield.Password, "password")
	f.Set(pdufield.InterfaceVersion, 0x34)
	pdu.Header().Seq = 1
	var b bytes.Buffer
	if err := pdu.SerializeTo(&b); err != nil {
		t.Fatal(err)
	}
	l := uint32(b.Len())
	if l != pdu.Header().Len {
		t.Fatalf("unexpected len: want %d, have %d", l, pdu.Header().Len)
	}
	if !bytes.Equal(tx, b.Bytes()) {
		t.Fatalf("unexpected bytes:\nwant:\n%s\nhave:\n%s",
			hex.Dump(tx), hex.Dump(b.Bytes()))
	}
	pdu, err := Decode(&b)
	if err != nil {
		t.Fatal(err)
	}
	h := pdu.Header()
	if h.ID != BindTransmitterID {
		t.Fatalf("unexpected ID: want %d, have %d",
			BindTransmitterID, h.ID)
	}
	if h.Seq != 1 {
		t.Fatalf("unexpected Seq: want 1, have %d", h.Seq)
	}
	test := []struct {
		n pdufield.Name
		v string
	}{
		{pdufield.SystemID, "smppclient1"},
		{pdufield.Password, "password"},
		{pdufield.InterfaceVersion, strconv.Itoa(0x34)},
	}
	for _, el := range test {
		f := pdu.Fields()[el.n]
		if f == nil {
			t.Fatalf("missing field: %s", el.n)
		}
		if f.String() != el.v {
			t.Fatalf("unexpected value for %q: want %q, have %q",
				el.n, el.v, f.String())
		}
	}
}

/*
func TestBindResp(t *testing.T) {
	tx := []byte{
		0x00, 0x00, 0x00, 0x18, 0x80, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x53, 0x4D, 0x50, 0x50, 0x53, 0x69, 0x6D, 0x00,
	}
	t.Log(tx)
}
*/

func roundTrip(t *testing.T, p Body, wantID ID) Body {
	t.Helper()
	var b bytes.Buffer
	if err := p.SerializeTo(&b); err != nil {
		t.Fatalf("serialize: %v", err)
	}
	if int(p.Header().Len) != b.Len() {
		t.Fatalf("header len %d, serialized %d", p.Header().Len, b.Len())
	}
	got, err := Decode(&b)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if got.Header().ID != wantID {
		t.Fatalf("decoded ID: want %s, have %s", wantID, got.Header().ID)
	}
	return got
}

func TestCancelSMRoundTrip(t *testing.T) {
	p := NewCancelSM()
	f := p.Fields()
	f.Set(pdufield.ServiceType, "")
	f.Set(pdufield.MessageID, "msg-7")
	f.Set(pdufield.SourceAddrTON, uint8(1))
	f.Set(pdufield.SourceAddrNPI, uint8(1))
	f.Set(pdufield.SourceAddr, "src")
	f.Set(pdufield.DestAddrTON, uint8(1))
	f.Set(pdufield.DestAddrNPI, uint8(1))
	f.Set(pdufield.DestinationAddr, "dst")

	got := roundTrip(t, p, CancelSMID)
	if v := got.Fields()[pdufield.MessageID].String(); v != "msg-7" {
		t.Fatalf("message_id: %q", v)
	}
	if v := got.Fields()[pdufield.SourceAddr].String(); v != "src" {
		t.Fatalf("source_addr: %q", v)
	}
	if v := got.Fields()[pdufield.DestinationAddr].String(); v != "dst" {
		t.Fatalf("destination_addr: %q", v)
	}
}

func TestCancelSMRespRoundTrip(t *testing.T) {
	roundTrip(t, NewCancelSMResp(), CancelSMRespID)
}

func TestReplaceSMRoundTrip(t *testing.T) {
	p := NewReplaceSM()
	f := p.Fields()
	f.Set(pdufield.MessageID, "msg-9")
	f.Set(pdufield.SourceAddrTON, uint8(1))
	f.Set(pdufield.SourceAddrNPI, uint8(1))
	f.Set(pdufield.SourceAddr, "sender")
	f.Set(pdufield.ScheduleDeliveryTime, "")
	f.Set(pdufield.ValidityPeriod, "")
	f.Set(pdufield.RegisteredDelivery, uint8(0))
	f.Set(pdufield.SMDefaultMsgID, uint8(0))
	f.Set(pdufield.ShortMessage, []byte("hi"))

	got := roundTrip(t, p, ReplaceSMID)
	if v := got.Fields()[pdufield.MessageID].String(); v != "msg-9" {
		t.Fatalf("message_id: %q", v)
	}
	if v := string(got.Fields()[pdufield.ShortMessage].Bytes()); v != "hi" {
		t.Fatalf("short_message: %q", v)
	}
}

func TestReplaceSMRespRoundTrip(t *testing.T) {
	roundTrip(t, NewReplaceSMResp(), ReplaceSMRespID)
}

func TestAlertNotificationRoundTrip(t *testing.T) {
	p := NewAlertNotification()
	f := p.Fields()
	f.Set(pdufield.SourceAddrTON, uint8(1))
	f.Set(pdufield.SourceAddrNPI, uint8(1))
	f.Set(pdufield.SourceAddr, "msisdn")
	f.Set(pdufield.ESMAddrTON, uint8(2))
	f.Set(pdufield.ESMAddrNPI, uint8(0))
	f.Set(pdufield.ESMAddr, "esme")

	got := roundTrip(t, p, AlertNotificationID)
	if v := got.Fields()[pdufield.SourceAddr].String(); v != "msisdn" {
		t.Fatalf("source_addr: %q", v)
	}
	if v := got.Fields()[pdufield.ESMAddr].String(); v != "esme" {
		t.Fatalf("esme_addr: %q", v)
	}
	if v := got.Fields()[pdufield.ESMAddrTON].String(); v != "2" {
		t.Fatalf("esme_addr_ton: %q", v)
	}
}

func TestOutbindRoundTrip(t *testing.T) {
	p := NewOutbind()
	f := p.Fields()
	f.Set(pdufield.SystemID, "smsc")
	f.Set(pdufield.Password, "secret")

	got := roundTrip(t, p, OutbindID)
	if v := got.Fields()[pdufield.SystemID].String(); v != "smsc" {
		t.Fatalf("system_id: %q", v)
	}
	if v := got.Fields()[pdufield.Password].String(); v != "secret" {
		t.Fatalf("password: %q", v)
	}
}

// message_state and error_code are 1-octet integers (SMPP 3.4 section
// 4.8.2); decoding them as C-octet strings loses both.
func TestQuerySMRespRoundTrip(t *testing.T) {
	p := NewQuerySMResp()
	f := p.Fields()
	f.Set(pdufield.MessageID, "13")
	f.Set(pdufield.FinalDate, "")
	f.Set(pdufield.MessageState, uint8(5))
	f.Set(pdufield.ErrorCode, uint8(7))

	got := roundTrip(t, p, QuerySMRespID)
	if v := got.Fields()[pdufield.MessageID].String(); v != "13" {
		t.Fatalf("message_id: %q", v)
	}
	if v := got.Fields()[pdufield.MessageState]; v == nil || !bytes.Equal(v.Bytes(), []byte{5}) {
		t.Fatalf("message_state: %v", v)
	}
	if v := got.Fields()[pdufield.ErrorCode]; v == nil || !bytes.Equal(v.Bytes(), []byte{7}) {
		t.Fatalf("error_code: %v", v)
	}
}
