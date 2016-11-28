// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"strconv"
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
	_, err := w.Write(b)
	return err
}

type tlvBodyJSON struct {
	Tag  TLVType `json:"tag"`
	Len  uint16  `json:"len"`
	Data []byte  `json:"data"`
	Text string  `json:"text"`
}

func (tlv TLVBody) MarshalJSON() ([]byte, error) {
	s := tlvBodyJSON{
		Tag:  tlv.Tag,
		Len:  tlv.Len,
		Data: tlv.Bytes(),
		Text: string(tlv.Bytes()),
	}
	return json.Marshal(s)
}

func (tlv *TLVBody) UnmarshalJSON(b []byte) error {
	s := tlvBodyJSON{}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	tlv.Tag = s.Tag
	tlv.Len = s.Len
	tlv.data = s.Data
	return nil
}

var tlvTypeMap = map[TLVType]string{
	DestAddrSubunit:          "dest_addr_subunit",
	DestNetworkType:          "dest_network_type",
	DestBearerType:           "dest_bearer_type",
	DestTelematicsID:         "dest_telematics_id",
	SourceAddrSubunit:        "source_addr_subunit",
	SourceNetworkType:        "source_network_type",
	SourceBearerType:         "source_bearer_type",
	SourceTelematicsID:       "source_telematics_id",
	QosTimeToLive:            "qos_time_to_live",
	PayloadType:              "payload_type",
	AdditionalStatusInfoText: "additional_status_info_text",
	ReceiptedMessageID:       "receipted_message_id",
	MsMsgWaitFacilities:      "ms_msg_wait_facilities",
	PrivacyIndicator:         "privacy_indicator",
	SourceSubaddress:         "source_subaddress",
	DestSubaddress:           "dest_subaddress",
	UserMessageReference:     "user_message_reference",
	UserResponseCode:         "user_response_code",
	SourcePort:               "source_port",
	DestinationPort:          "destination_port",
	SarMsgRefNum:             "sar_msg_ref_num",
	LanguageIndicator:        "language_indicator",
	SarTotalSegments:         "sar_total_segments",
	SarSegmentSeqnum:         "sar_segment_seqnum",
	CallbackNumPresInd:       "callback_num_pres_ind",
	CallbackNumAtag:          "callback_num_atag",
	NumberOfMessages:         "number_of_messages",
	CallbackNum:              "callback_num",
	DpfResult:                "dpf_result",
	SetDpf:                   "set_dpf",
	MsAvailabilityStatus:     "ms_availability_status",
	NetworkErrorCode:         "network_error_code",
	MessagePayload:           "message_payload",
	DeliveryFailureReason:    "delivery_failure_reason",
	MoreMessagesToSend:       "more_messages_to_send",
	MessageStateOption:       "message_state_option",
	UssdServiceOp:            "ussd_service_op",
	DisplayTime:              "display_time",
	SmsSignal:                "sms_signal",
	MsValidity:               "ms_validity",
	AlertOnMessageDelivery:   "alert_on_message_delivery",
	ItsReplyType:             "its_reply_type",
	ItsSessionInfo:           "its_session_info",
}

func (t TLVType) String() string {
	s := tlvTypeMap[t]
	if s == "" {
		s = strconv.Itoa(int(t))
	}
	return s
}
