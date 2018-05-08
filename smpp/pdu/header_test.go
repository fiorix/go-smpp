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
		{SubmitSMID, 0x04},
		{SubmitSMRespID, 0x04},
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
	testCases := []struct {
		id  ID
		seq uint32
		key string
	}{
		{GenericNACKID, 1, "0-1"},
		{BindReceiverID, 2, "1-2"},
		{BindReceiverRespID, 3, "1-3"},
		{BindTransmitterID, 4, "2-4"},
		{BindTransmitterRespID, 5, "2-5"},
		{QuerySMID, 6, "3-6"},
		{QuerySMRespID, 7, "3-7"},
		{SubmitSMID, 8, "4-8"},
		{SubmitSMRespID, 9, "4-9"},
		{DeliverSMID, 10, "5-10"},
		{DeliverSMRespID, 1, "5-1"},
		{UnbindID, 1, "6-1"},
		{UnbindRespID, 1, "6-1"},
		{ReplaceSMID, 1, "7-1"},
		{ReplaceSMRespID, 1, "7-1"},
		{CancelSMID, 1, "8-1"},
		{CancelSMRespID, 1, "8-1"},
		{BindTransceiverID, 1, "9-1"},
		{BindTransceiverRespID, 1, "9-1"},
		{OutbindID, 1, "11-1"},
		{EnquireLinkID, 1, "21-1"},
		{EnquireLinkRespID, 1, "21-1"},
		{SubmitMultiID, 1, "33-1"},
		{SubmitMultiRespID, 1, "33-1"},
		{AlertNotificationID, 1, "258-1"},
		{DataSMID, 1, "259-1"},
		{DataSMRespID, 1, "259-1"},
	}

	for _, tc := range testCases {
		header := Header{HeaderLen, tc.id, Status(0), tc.seq}
		key := header.Key()

		if key != tc.key {
			t.Fatalf("expected: %s, actual: %s", tc.key, key)
		}
	}
}
