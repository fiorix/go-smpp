// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"sync/atomic"

	"github.com/veoo/go-smpp/smpp/pdu/pdufield"
)

var nextSeq uint32

// Codec is the base type of all PDUs.
// It implements the PDU interface and provides a generic encoder.
type Codec struct {
	h *Header
	l pdufield.List
	f pdufield.Map
	t pdufield.TLVMap
}

// init initializes the Codec's list and maps and sets the header
// sequence number.
func (pdu *Codec) init() {
	if pdu.l == nil {
		pdu.l = pdufield.List{}
	}
	pdu.f = make(pdufield.Map)
	pdu.t = make(pdufield.TLVMap)
	pdu.h.Seq = atomic.AddUint32(&nextSeq, 1)
}

// setup replaces the Codec's current maps with the given ones.
func (pdu *Codec) setup(f pdufield.Map, t pdufield.TLVMap) {
	pdu.f, pdu.t = f, t
}

// Header implements the PDU interface.
func (pdu *Codec) Header() *Header {
	return pdu.h
}

// Len implements the PDU interface.
func (pdu *Codec) Len() int {
	l := HeaderLen
	for _, f := range pdu.f {
		l += f.Len()
	}
	for _, t := range pdu.t {
		l += int(t.Len)
	}
	return l
}

// FieldList implements the PDU interface.
func (pdu *Codec) FieldList() pdufield.List {
	return pdu.l
}

// Fields implement the PDU interface.
func (pdu *Codec) Fields() pdufield.Map {
	return pdu.f
}

// TLVFields implement the PDU interface.
func (pdu *Codec) TLVFields() pdufield.TLVMap {
	return pdu.t
}

// SerializeTo implements the PDU interface.
func (pdu *Codec) SerializeTo(w io.Writer) error {
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
	for k, v := range pdu.TLVFields() {
		if err := v.SerializeTo(&b); err != nil {
			return err
		}
		pdu.t.Set(k, nil)
	}

	pdu.h.Len = uint32(pdu.Len())
	err := pdu.h.SerializeTo(w)
	if err != nil {
		return err
	}
	_, err = io.Copy(w, &b)
	return err
}

type CodecJSON struct {
	Header    *Header         `json:"header"`
	FieldList pdufield.List   `json:"fieldList"`
	Fields    pdufield.Map    `json:"fields"`
	TLVFields pdufield.TLVMap `json:"tlvFields"`
}

func (c *Codec) MarshalJSON() ([]byte, error) {
	j := CodecJSON{
		c.Header(),
		c.FieldList(),
		c.Fields(),
		c.TLVFields(),
	}
	return json.Marshal(j)
}

// Since Codec is a private struct, we expose the Unmarshal function
// for other packages to use it
func (c *Codec) UnmarshalJSON(b []byte) error {
	j := CodecJSON{}

	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	c.h = j.Header
	c.l = j.FieldList
	c.f = j.Fields
	c.t = j.TLVFields
	return nil
}

// decoder wraps a PDU (e.g. Bind) and the Codec together and is
// used for initializing new PDUs with map data decoded off the wire.
type decoder interface {
	Body
	setup(f pdufield.Map, t pdufield.TLVMap)
}

func decodeFields(pdu decoder, b []byte) (Body, error) {
	l := pdu.FieldList()
	r := bytes.NewBuffer(b)
	f, err := l.Decode(r)
	if err != nil {
		return nil, err
	}
	t, err := l.DecodeTLV(r)
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
	// TODO(fiorix): Implement AlertNotification.
	case BindReceiverID, BindTransceiverID, BindTransmitterID:
		return decodeFields(newBind(hdr), b)
	case BindReceiverRespID, BindTransceiverRespID, BindTransmitterRespID:
		return decodeFields(newBindResp(hdr), b)
	case CancelSMID:
	// TODO(fiorix): Implement CancelSM.
	case CancelSMRespID:
	// TODO(fiorix): Implement CancelSMResp.
	case DataSMID:
	// TODO(fiorix): Implement DataSM.
	case DataSMRespID:
	// TODO(fiorix): Implement DataSMResp.
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
	// TODO(fiorix): Implement Outbind.
	case QuerySMID:
		return decodeFields(newQuerySM(hdr), b)
	case QuerySMRespID:
		return decodeFields(newQuerySMResp(hdr), b)
	case ReplaceSMID:
	// TODO(fiorix): Implement ReplaceSM.
	case ReplaceSMRespID:
	// TODO(fiorix): Implement ReplaceSMResp.
	case SubmitMultiID:
	// TODO(fiorix): Implement SubmitMulti.
	case SubmitMultiRespID:
	// TODO(fiorix): Implement SubmitMultiResp.
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
	return nil, fmt.Errorf("PDU not implemented: %#x", hdr.ID)
}
