// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdu

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/v2/smpp/pdu/pdutlv"
)

var nextSeq uint32

// codec is the base type of all PDUs.
// It implements the PDU interface and provides a generic encoder.
type codec struct {
	h *Header
	l pdufield.List
	f pdufield.Map
	t pdutlv.Map
}

// init initializes the codec's list and maps and sets the header
// sequence number.
func (pdu *codec) init() {
	if pdu.l == nil {
		pdu.l = pdufield.List{}
	}
	pdu.f = make(pdufield.Map)
	pdu.t = make(pdutlv.Map)
	if pdu.h.Seq == 0 { // If Seq not set
		// SMPP 3.4 §5.1.4 restricts sequence_number to 0x00000001..0x7FFFFFFF.
		// Wrap within that range (skipping 0) so a long-lived connection never
		// emits a value strict SMSCs would reject as an invalid header.
		next := atomic.AddUint32(&nextSeq, 1) & 0x7FFFFFFF
		if next == 0 {
			next = atomic.AddUint32(&nextSeq, 1) & 0x7FFFFFFF
		}
		pdu.h.Seq = next
	}
}

// setup replaces the codec's current maps with the given ones.
func (pdu *codec) setup(f pdufield.Map, t pdutlv.Map) {
	pdu.f, pdu.t = f, t
}

// Header implements the PDU interface.
func (pdu *codec) Header() *Header {
	return pdu.h
}

// Len implements the PDU interface.
func (pdu *codec) Len() int {
	l := HeaderLen
	for _, f := range pdu.f {
		l += f.Len()
	}
	for _, t := range pdu.t {
		l += t.Len()
	}
	return l
}

// FieldList implements the PDU interface.
func (pdu *codec) FieldList() pdufield.List {
	return pdu.l
}

// Fields implement the PDU interface.
func (pdu *codec) Fields() pdufield.Map {
	return pdu.f
}

// Fields implement the PDU interface.
func (pdu *codec) TLVFields() pdutlv.Map {
	return pdu.t
}

// SerializeTo implements the PDU interface.
func (pdu *codec) SerializeTo(w io.Writer) error {
	var b bytes.Buffer
	for _, k := range pdu.FieldList() {
		f, ok := pdu.f[k]
		if !ok {
			pdu.f.Set(k, nil)
			f = pdu.f[k]
		}
		if err := f.SerializeTo(&b); err != nil {
			return err
		}
	}
	for _, f := range pdu.TLVFields() {
		if err := f.SerializeTo(&b); err != nil {
			return err
		}
	}
	// Header length must reflect what we are about to write, not the
	// raw size of the field map: Map.Set on ShortMessage adds entries
	// like DataCoding that are not necessarily in the field list, so
	// counting the map would overstate the wire length on PDUs that
	// omit DataCoding (e.g. ReplaceSM).
	pdu.h.Len = uint32(HeaderLen + b.Len())
	err := pdu.h.SerializeTo(w)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, &b)
	return err
}

// decoder wraps a PDU (e.g. Bind) and the codec together and is
// used for initializing new PDUs with map data decoded off the wire.
type decoder interface {
	Body
	setup(f pdufield.Map, t pdutlv.Map)
}

func decodeFields(pdu decoder, b []byte) (Body, error) {
	l := pdu.FieldList()
	r := bytes.NewBuffer(b)
	f, err := l.Decode(r)
	if err != nil {
		// A response may legitimately come without a body: SMPP 3.4
		// omits the body of e.g. bind_*_resp and submit_sm_resp when
		// command_status is non-zero, and some SMSCs do the same on
		// success. A request has no such allowance, and neither has a
		// body that stops part-way through its fields.
		var te *pdufield.TruncatedError
		if !errors.As(err, &te) || len(b) != 0 || !isResponse(pdu.Header().ID) {
			return nil, err
		}
		f = make(pdufield.Map)
	}
	t, err := pdutlv.DecodeTLV(r)
	if err != nil {
		return nil, err
	}
	pdu.setup(f, t)
	return pdu, nil
}

// Decode decodes binary PDU data. It returns a new PDU object, e.g. Bind,
// with header and all fields decoded. The returned PDU can be modified
// and re-serialized to its binary form.
func Decode(r io.Reader) (Body, error) {
	hdr, err := DecodeHeader(r)
	if err != nil {
		return nil, err
	}
	b := make([]byte, hdr.Len-HeaderLen)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return nil, err
	}
	switch hdr.ID {
	case AlertNotificationID:
		return decodeFields(newAlertNotification(hdr), b)
	case BindReceiverID, BindTransceiverID, BindTransmitterID:
		return decodeFields(newBind(hdr), b)
	case BindReceiverRespID, BindTransceiverRespID, BindTransmitterRespID:
		return decodeFields(newBindResp(hdr), b)
	case CancelSMID:
		return decodeFields(newCancelSM(hdr), b)
	case CancelSMRespID:
		return decodeFields(newCancelSMResp(hdr), b)
	case DataSMID:
		return decodeFields(newDataSM(hdr), b)
	case DataSMRespID:
		return decodeFields(newDataSMResp(hdr), b)
	case DeliverSMID:
		return decodeFields(newDeliverSM(hdr), b)
	case DeliverSMRespID:
		return decodeFields(newDeliverSMResp(hdr), b)
	case EnquireLinkID:
		return decodeFields(newEnquireLink(hdr), b)
	case EnquireLinkRespID:
		return decodeFields(newEnquireLinkResp(hdr), b)
	case GenericNACKID:
		return decodeFields(newGenericNACK(hdr), b)
	case OutbindID:
		return decodeFields(newOutbind(hdr), b)
	case QuerySMID:
		return decodeFields(newQuerySM(hdr), b)
	case QuerySMRespID:
		return decodeFields(newQuerySMResp(hdr), b)
	case ReplaceSMID:
		return decodeFields(newReplaceSM(hdr), b)
	case ReplaceSMRespID:
		return decodeFields(newReplaceSMResp(hdr), b)
	case SubmitMultiID:
		return decodeFields(newSubmitMulti(hdr), b)
	case SubmitMultiRespID:
		return decodeFields(newSubmitMultiResp(hdr), b)
	case SubmitSMID:
		return decodeFields(newSubmitSM(hdr), b)
	case SubmitSMRespID:
		return decodeFields(newSubmitSMResp(hdr), b)
	case UnbindID:
		return decodeFields(newUnbind(hdr), b)
	case UnbindRespID:
		return decodeFields(newUnbindResp(hdr), b)
	default:
		return nil, fmt.Errorf("unknown PDU type: %#x", hdr.ID)
	}
}

// isResponse reports whether id is a response command: bit 31 set, as
// SMPP 3.4 section 5.1.2.1 defines.
func isResponse(id ID) bool {
	return uint32(id)&0x80000000 != 0
}
