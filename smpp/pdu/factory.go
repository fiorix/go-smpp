// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdu

import (
	"fmt"
	"sync/atomic"
)

// Factory is used to instantiate PDUs in a more controllable fashion. Its main purpose
// is to handle sequence number generation in a contained way and without it being
// a global variable in the package
type Factory interface {
	CreatePDU(id ID) (Body, error)
	CreatePDUResp(id ID, seq uint32) (Body, error)
}

type factory struct {
	nextSeq uint32
}

func NewFactory() Factory {
	return &factory{}
}

func (f *factory) CreatePDU(id ID) (Body, error) {
	var c *Codec
	switch id {
	case AlertNotificationID:
		// TODO(cesar0094): Implement AlertNotification.
	case BindReceiverID, BindTransceiverID, BindTransmitterID:
		c = newBind(&Header{ID: id})
	case CancelSMID:
		// TODO(cesar0094): Implement CancelSM.
	case DataSMID:
		// TODO(cesar0094): Implement DataSM.
	case DeliverSMID:
		c = newDeliverSM(&Header{ID: id})
	case EnquireLinkID:
		c = newEnquireLink(&Header{ID: id})
	case OutbindID:
		// TODO(cesar0094): Implement Outbind.
	case QuerySMID:
		c = newQuerySM(&Header{ID: id})
	case ReplaceSMID:
		// TODO(cesar0094): Implement ReplaceSM.
	case SubmitMultiID:
		c = newSubmitMulti(&Header{ID: id})
	case SubmitSMID:
		c = newSubmitSM(&Header{ID: id})
	case UnbindID:
		c = newUnbind(&Header{ID: id})
	default:
		return nil, fmt.Errorf("unknown PDU type: %#x", id)
	}
	if c == nil {
		return nil, fmt.Errorf("PDU not implemented: %#x", id)
	}
	if nextSeq >= 0x7FFFFFFF {
		nextSeq = 0
	}
	c.h.Seq = atomic.AddUint32(&nextSeq, 1)
	c.init()
	return c, nil
}

func (f *factory) CreatePDUResp(id ID, seq uint32) (Body, error) {
	var c *Codec
	switch id {
	case BindReceiverRespID, BindTransceiverRespID, BindTransmitterRespID:
		c = newBindResp(&Header{ID: id})
	case CancelSMRespID:
		// TODO(cesar0094): Implement CancelSMResp.
	case DataSMRespID:
		// TODO(cesar0094): Implement DataSMResp.
	case DeliverSMRespID:
		c = newDeliverSMResp(&Header{ID: id})
	case EnquireLinkRespID:
		c = newEnquireLinkResp(&Header{ID: id})
	case GenericNACKID:
		c = newGenericNACK(&Header{ID: id})
	case QuerySMRespID:
		c = newQuerySMResp(&Header{ID: id})
	case ReplaceSMRespID:
		// TODO(cesar0094): Implement ReplaceSMResp.
	case SubmitMultiRespID:
		c = newSubmitMultiResp(&Header{ID: id})
	case SubmitSMRespID:
		c = newSubmitSMResp(&Header{ID: id})
	case UnbindRespID:
		c = newUnbindResp(&Header{ID: id})
	default:
		return nil, fmt.Errorf("unknown PDU type: %#x", id)
	}
	if c == nil {
		return nil, fmt.Errorf("PDU not implemented: %#x", id)
	}
	c.h.Seq = seq
	c.init()
	return c, nil
}
