// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"encoding/binary"
	"io"
)

//TLV Tags
const (
	DestAddrSubunit          TLVType = 0x0005
	DestNetworkType          TLVType = 0x0006
	DestBearerType           TLVType = 0x0007
	DestTelematicsID         TLVType = 0x0008
	SourceAddrSubunit        TLVType = 0x000D
	SourceNetworkType        TLVType = 0x000E
	SourceBearerType         TLVType = 0x000F
	SourceTelematicsID       TLVType = 0x0010
	QosTimeToLive            TLVType = 0x0017
	PayloadType              TLVType = 0x0019
	AdditionalStatusInfoText TLVType = 0x001D
	ReceiptedMessageID       TLVType = 0x001E
	MsMsgWaitFacilities      TLVType = 0x0030
	PrivacyIndicator         TLVType = 0x0201
	SourceSubaddress         TLVType = 0x0202
	DestSubaddress           TLVType = 0x0203
	UserMessageReference     TLVType = 0x0204
	UserResponseCode         TLVType = 0x0205
	SourcePort               TLVType = 0x020A
	DestinationPort          TLVType = 0x020B
	SarMsgRefNum             TLVType = 0x020C
	LanguageIndicator        TLVType = 0x020D
	SarTotalSegments         TLVType = 0x020E
	SarSegmentSeqnum         TLVType = 0x020F
	CallbackNumPresInd       TLVType = 0x0302
	CallbackNumAtag          TLVType = 0x0303
	NumberOfMessages         TLVType = 0x0304
	CallbackNum              TLVType = 0x0381
	DpfResult                TLVType = 0x0420
	SetDpf                   TLVType = 0x0421
	MsAvailabilityStatus     TLVType = 0x0422
	NetworkErrorCode         TLVType = 0x0423
	MessagePayload           TLVType = 0x0424
	DeliveryFailureReason    TLVType = 0x0425
	MoreMessagesToSend       TLVType = 0x0426
	MessageStateOption       TLVType = 0x0427
	UssdServiceOp            TLVType = 0x0501
	DisplayTime              TLVType = 0x1201
	SmsSignal                TLVType = 0x1203
	MsValidity               TLVType = 0x1204
	AlertOnMessageDelivery   TLVType = 0x130C
	ItsReplyType             TLVType = 0x1380
	ItsSessionInfo           TLVType = 0x1383
)

// TLVType is the Tag Length Value.
type TLVType uint16

// TLVBody represents data of a TLV field.
type TLVBody struct {
	Tag  TLVType
	Len  uint16
	data []byte
}

// Bytes return raw TLV binary data.
func (tlv *TLVBody) Bytes() []byte {
	return tlv.data
}

// SerializeTo serializes TLV data to its binary form.
func (tlv *TLVBody) SerializeTo(w io.Writer) error {
	b := make([]byte, 4+len(tlv.data))
	binary.BigEndian.PutUint16(b[0:2], uint16(tlv.Tag))
	binary.BigEndian.PutUint16(b[2:4], tlv.Len)
	copy(b[4:], tlv.data)
	return nil
}
