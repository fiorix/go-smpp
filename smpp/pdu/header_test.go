// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdu

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestHeader(t *testing.T) {
	want := []byte{
		0x00, 0x00, 0x00, 0x10, // 16 Len
		0x80, 0x00, 0x00, 0x00, // GenericNACK ID
		0x00, 0x00, 0x00, 0x01, // Invalid message length Status
		0x00, 0x00, 0x00, 0x0D, // 13 Seq
	}
	h, err := DecodeHeader(bytes.NewBuffer(want))
	if err != nil {
		t.Fatal(err)
	}
	if h.Len != 16 {
		t.Fatalf("unexpected Len: want 16, have %d", h.Len)
	}
	if h.ID != GenericNACKID {
		t.Fatalf("unexpected ID: want GenericNACK, have %d", h.ID)
	}
	if h.Status != 1 {
		t.Fatalf("unexpected Status: want 1, have %d", h.Status)
	}
	if h.Seq != 13 {
		t.Fatalf("unexpected Seq: want 13, have %d", h.Seq)
	}
	var b bytes.Buffer
	if err := h.SerializeTo(&b); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(want, b.Bytes()) {
		t.Fatalf("malformed header:\nwant:%s\nhave:\n%s",
			hex.Dump(want), hex.Dump(b.Bytes()))
	}
	we := "invalid message length"
	if have := h.Status.Error(); have != we {
		t.Fatalf("unexpected status: want %q, have %q", we, have)
	}
	h.Status = 0x2000
	we = "unknown status: 8192"
	if have := h.Status.Error(); have != we {
		t.Fatalf("unexpected status: want %q, have %q", we, have)
	}
}

func TestDecodeHeader(t *testing.T) {
	h, err := DecodeHeader(bytes.NewBuffer(nil))
	if err == nil {
		t.Fatalf("unexpected parsing of no data: %#v", h)
	}
	bin := []byte{
		0x00, 0x00, 0x00, 0x01, // 1 Len
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}
	h, err = DecodeHeader(bytes.NewBuffer(bin))
	if err == nil {
		t.Fatalf("unexpected parsing of short Len: %#v", h)
	}
	bin[2] = 0x20
	h, err = DecodeHeader(bytes.NewBuffer(bin))
	if err == nil {
		t.Fatalf("unexpected parsing of big Len: %#v", h)
	}
}

func TestGroup(t *testing.T) {
	testCases := []struct {
		id    ID
		group uint16
	}{
		{GenericNACKID, 0x00},
		{BindReceiverID, 0x01},
		{BindReceiverRespID, 0x01},
		{BindTransmitterID, 0x02},
		{BindTransmitterRespID, 0x02},
		{QuerySMID, 0x03},
		{QuerySMRespID, 0x03},
		{SubmitSMID, 0x0004},
		{SubmitSMRespID, 0x0004},
		{DeliverSMID, 0x05},
		{DeliverSMRespID, 0x05},
		{UnbindID, 0x06},
		{UnbindRespID, 0x06},
		{ReplaceSMID, 0x07},
		{ReplaceSMRespID, 0x07},
		{CancelSMID, 0x08},
		{CancelSMRespID, 0x08},
		{BindTransceiverID, 0x09},
		{BindTransceiverRespID, 0x09},
		{OutbindID, 0x0B},
		{EnquireLinkID, 0x15},
		{EnquireLinkRespID, 0x15},
		{SubmitMultiID, 0x21},
		{SubmitMultiRespID, 0x21},
		{AlertNotificationID, 0x102},
		{DataSMID, 0x103},
		{DataSMRespID, 0x103},
	}

	for _, tc := range testCases {
		group := tc.id.Group()
		if group != tc.group {
			t.Fatalf("expected: %o, actual: %o", tc.group, group)
		}
	}
}

func TestKey(t *testing.T) {
	sm := []byte{
		0x00, 0x00, 0x00, 0x10, // 16 Len
		0x00, 0x00, 0x00, 0x04, // SubmitSMID
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	smH, err := DecodeHeader(bytes.NewBuffer(sm))
	if err != nil {
		t.Fatal(err)
	}

	k := smH.Key()
	if k != "4-0" {
		t.Fatalf("unexpected key: %s", k)
	}
}
