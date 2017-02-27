// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdu

import (
	"math/rand"
	"testing"
	"time"
)

func TestCreatePDU(t *testing.T) {
	tests := []struct {
		id      ID
		mustErr bool
	}{
		{
			id:      AlertNotificationID,
			mustErr: true,
		},
		{
			id:      BindReceiverID,
			mustErr: false,
		},
		{
			id:      BindTransceiverID,
			mustErr: false,
		},
		{
			id:      BindTransmitterID,
			mustErr: false,
		},
		{
			id:      CancelSMID,
			mustErr: true,
		},
		{
			id:      DataSMID,
			mustErr: true,
		},
		{
			id:      DeliverSMID,
			mustErr: false,
		},
		{
			id:      EnquireLinkID,
			mustErr: false,
		},
		{
			id:      OutbindID,
			mustErr: true,
		},
		{
			id:      QuerySMID,
			mustErr: false,
		},
		{
			id:      ReplaceSMID,
			mustErr: true,
		},
		{
			id:      SubmitMultiID,
			mustErr: false,
		},
		{
			id:      SubmitSMID,
			mustErr: false,
		},
		{
			id:      UnbindID,
			mustErr: false,
		},
		{
			id:      ID(999999999),
			mustErr: true,
		},
	}
	f := NewFactory()
	for _, test := range tests {
		p, err := f.CreatePDU(test.id)
		if err != nil != test.mustErr {
			t.Errorf("Incorrect creation of PDU %#v", test.id.String())
		}
		if p != nil && p.Header().ID != test.id {
			t.Errorf("Created incorrect PDU type %#v", p.Header().ID.String())
		}
	}
}

func TestCreatePDUResp(t *testing.T) {
	tests := []struct {
		id      ID
		mustErr bool
	}{
		{
			id:      BindReceiverRespID,
			mustErr: false,
		},
		{
			id:      BindTransceiverRespID,
			mustErr: false,
		},
		{
			id:      BindTransmitterRespID,
			mustErr: false,
		},
		{
			id:      CancelSMRespID,
			mustErr: true,
		},
		{
			id:      DataSMRespID,
			mustErr: true,
		},
		{
			id:      DeliverSMRespID,
			mustErr: false,
		},
		{
			id:      EnquireLinkRespID,
			mustErr: false,
		},
		{
			id:      GenericNACKID,
			mustErr: false,
		},
		{
			id:      QuerySMRespID,
			mustErr: false,
		},
		{
			id:      ReplaceSMRespID,
			mustErr: true,
		},
		{
			id:      SubmitMultiRespID,
			mustErr: false,
		},
		{
			id:      SubmitSMRespID,
			mustErr: false,
		},
		{
			id:      UnbindRespID,
			mustErr: false,
		},
	}
	f := NewFactory()
	rand.Seed(time.Now().UTC().UnixNano())

	for _, test := range tests {
		seq := uint32(rand.Intn(99999))
		p, err := f.CreatePDUResp(test.id, seq)
		if err != nil != test.mustErr {
			t.Errorf("Incorrect creation of PDU %#v", test.id.String())
		}
		if p == nil {
			continue
		}

		if p.Header().ID != test.id {
			t.Errorf("Created incorrect PDU type %#v", p.Header().ID.String())
		} else if p.Header().Seq != seq {
			t.Errorf("Mismatching sequence number expected: %#v, got: %#v", seq, p.Header().Seq)
		}
	}
}

func TestSequenceNumbers(t *testing.T) {
	f := &factory{}

	for i := uint32(1); i < 10; i++ {
		p, err := f.CreatePDU(SubmitSMID)
		if err != nil {
			t.Errorf("Unexpected error creating PDU: %#v", err)
		}
		want, got := i, p.Header().Seq
		if want != got {
			t.Errorf("Unexpected sequence number want: %#v, got: %#v", want, got)
		}
	}

	// test limits
	f.nextSeq = 0x7FFFFFFF
	p, err := f.CreatePDU(SubmitSMID)
	if err != nil {
		t.Errorf("Unexpected error creating PDU: %#v", err)
	}
	want, got := uint32(1), p.Header().Seq
	if want != got {
		t.Errorf("Unexpected sequence number want: %#v, got: %#v", want, got)
	}
}
