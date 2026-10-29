// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdu

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"

	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdufield"
)

func rawPDU(id ID, status Status, body []byte) []byte {
	b := make([]byte, HeaderLen+len(body))
	binary.BigEndian.PutUint32(b[0:], uint32(len(b)))
	binary.BigEndian.PutUint32(b[4:], uint32(id))
	binary.BigEndian.PutUint32(b[8:], uint32(status))
	binary.BigEndian.PutUint32(b[12:], 1)
	copy(b[HeaderLen:], body)
	return b
}

var fullBind = []byte("sys\x00pw\x00\x00\x34\x00\x00\x00")

func TestDecodeTruncatedRequest(t *testing.T) {
	for _, id := range []ID{BindReceiverID, BindTransmitterID, BindTransceiverID} {
		for cut := 0; cut < len(fullBind); cut++ {
			_, err := Decode(bytes.NewReader(rawPDU(id, 0, fullBind[:cut])))
			var te *pdufield.TruncatedError
			if !errors.As(err, &te) {
				t.Fatalf("%s cut at %d: want *pdufield.TruncatedError, have %v", id, cut, err)
			}
		}
		p, err := Decode(bytes.NewReader(rawPDU(id, 0, fullBind)))
		if err != nil {
			t.Fatalf("%s complete: %v", id, err)
		}
		if p.Fields()[pdufield.SystemID].String() != "sys" {
			t.Fatalf("%s complete: unexpected system_id %q", id, p.Fields()[pdufield.SystemID])
		}
	}

	// submit_sm cut after source_addr: the body ends between fields.
	sm := []byte("\x00\x01\x01src\x00")
	_, err := Decode(bytes.NewReader(rawPDU(SubmitSMID, 0, sm)))
	var te *pdufield.TruncatedError
	if !errors.As(err, &te) || te.Field != pdufield.DestAddrTON || te.Offset != len(sm) {
		t.Fatalf("submit_sm: want *pdufield.TruncatedError on %q at %d, have %v", pdufield.DestAddrTON, len(sm), err)
	}
}

func TestDecodeResponseWithoutBody(t *testing.T) {
	for _, tc := range []struct {
		id     ID
		status Status
	}{
		{BindTransceiverRespID, 0x0E},
		{BindTransmitterRespID, 0x0D},
		{SubmitSMRespID, 0x58},
		{SubmitSMRespID, 0},
		{DeliverSMRespID, 0},
	} {
		p, err := Decode(bytes.NewReader(rawPDU(tc.id, tc.status, nil)))
		if err != nil {
			t.Fatalf("%s status %#x without body: %v", tc.id, uint32(tc.status), err)
		}
		if p.Header().Status != tc.status {
			t.Fatalf("%s: status %#x, want %#x", tc.id, uint32(p.Header().Status), uint32(tc.status))
		}
	}
}

func TestDecodeTruncatedResponseBody(t *testing.T) {
	// A body that starts and then stops is truncated, response or not.
	_, err := Decode(bytes.NewReader(rawPDU(BindTransceiverRespID, 0, []byte("smsc"))))
	var te *pdufield.TruncatedError
	if !errors.As(err, &te) || te.Field != pdufield.SystemID {
		t.Fatalf("want *pdufield.TruncatedError on system_id, have %v", err)
	}
}
