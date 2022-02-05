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
	"github.com/stretchr/testify/assert"
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

func TestGenericNACK(t *testing.T) {
	input := NewGenericNACK()
	assert.Equal(t, GenericNACKID, input.Header().ID)
	assert.Equal(t, pdufield.List{}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	assert.Equal(t, input.Fields(), output.Fields())
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestGenericNACKSeq(t *testing.T) {
	input := NewGenericNACKSeq(42)
	assert.Equal(t, GenericNACKID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	assert.Equal(t, input.Fields(), output.Fields())
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestBindReceiver(t *testing.T) {
	input := NewBindReceiver()
	assert.Equal(t, BindReceiverID, input.Header().ID)
	assert.Equal(t, pdufield.List{"system_id", "password", "system_type", "interface_version", "addr_ton", "addr_npi", "address_range"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestBindReceiverResp(t *testing.T) {
	input := NewBindReceiverRespSeq(42)
	assert.Equal(t, BindReceiverRespID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{"system_id"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestBindTransmitter(t *testing.T) {
	input := NewBindTransmitter()
	assert.Equal(t, BindTransmitterID, input.Header().ID)
	assert.Equal(t, pdufield.List{"system_id", "password", "system_type", "interface_version", "addr_ton", "addr_npi", "address_range"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestBindTransmitterResp(t *testing.T) {
	input := NewBindTransmitterRespSeq(42)
	assert.Equal(t, BindTransmitterRespID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{"system_id"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestQuerySM(t *testing.T) {
	input := NewQuerySM()
	assert.Equal(t, QuerySMID, input.Header().ID)
	assert.Equal(t, pdufield.List{"message_id", "source_addr_ton", "source_addr_npi", "source_addr"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestQuerySMResp(t *testing.T) {
	input := NewQuerySMRespSeq(42)
	assert.Equal(t, QuerySMRespID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{"message_id", "final_date", "message_state", "error_code"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestSubmitSM(t *testing.T) {
	input := NewSubmitSM(nil)
	assert.Equal(t, SubmitSMID, input.Header().ID)
	assert.Equal(t, pdufield.List{"service_type", "source_addr_ton", "source_addr_npi", "source_addr", "dest_addr_ton", "dest_addr_npi", "destination_addr", "esm_class", "protocol_id", "priority_flag", "schedule_delivery_time", "validity_period", "registered_delivery", "replace_if_present_flag", "data_coding", "sm_default_msg_id", "sm_length", "short_message"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestSubmitSMResp(t *testing.T) {
	input := NewSubmitSMRespSeq(42)
	assert.Equal(t, SubmitSMRespID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{"message_id"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestDeliverSM(t *testing.T) {
	input := NewDeliverSM()
	assert.Equal(t, DeliverSMID, input.Header().ID)
	assert.Equal(t, pdufield.List{"service_type", "source_addr_ton", "source_addr_npi", "source_addr", "dest_addr_ton", "dest_addr_npi", "destination_addr", "esm_class", "protocol_id", "priority_flag", "schedule_delivery_time", "validity_period", "registered_delivery", "replace_if_present_flag", "data_coding", "sm_default_msg_id", "sm_length", "short_message"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestDeliverSMResp(t *testing.T) {
	input := NewDeliverSMRespSeq(42)
	assert.Equal(t, DeliverSMRespID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{"message_id"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestUnbind(t *testing.T) {
	input := NewUnbind()
	assert.Equal(t, UnbindID, input.Header().ID)
	assert.Equal(t, pdufield.List{}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestUnbindResp(t *testing.T) {
	input := NewUnbindRespSeq(42)
	assert.Equal(t, UnbindRespID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestBindTransceiver(t *testing.T) {
	input := NewBindTransceiver()
	assert.Equal(t, BindTransceiverID, input.Header().ID)
	assert.Equal(t, pdufield.List{"system_id", "password", "system_type", "interface_version", "addr_ton", "addr_npi", "address_range"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestBindTransceiverResp(t *testing.T) {
	input := NewBindTransceiverRespSeq(42)
	assert.Equal(t, BindTransceiverRespID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{"system_id"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestEnquireLink(t *testing.T) {
	input := NewEnquireLink()
	assert.Equal(t, EnquireLinkID, input.Header().ID)
	assert.Equal(t, pdufield.List{}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestEnquireLinkResp(t *testing.T) {
	input := NewEnquireLinkRespSeq(42)
	assert.Equal(t, EnquireLinkRespID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestSubmitMulti(t *testing.T) {
	input := NewSubmitMulti(nil)
	assert.Equal(t, SubmitMultiID, input.Header().ID)
	assert.Equal(t, pdufield.List{"service_type", "source_addr_ton", "source_addr_npi", "source_addr", "number_of_dests", "dest_addresses", "esm_class", "protocol_id", "priority_flag", "schedule_delivery_time", "validity_period", "registered_delivery", "replace_if_present_flag", "data_coding", "sm_default_msg_id", "sm_length", "short_message"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}

func TestSubmitMultiResp(t *testing.T) {
	input := NewSubmitMultiRespSeq(42)
	assert.Equal(t, SubmitMultiRespID, input.Header().ID)
	assert.Equal(t, uint32(42), input.Header().Seq)
	assert.Equal(t, pdufield.List{"message_id", "no_unsuccess", "unsuccess_sme"}, input.FieldList())

	var buf bytes.Buffer
	err := input.SerializeTo(&buf)
	assert.NoError(t, err)

	output, err := Decode(&buf)
	assert.NoError(t, err)
	assert.Equal(t, input.Header(), output.Header())
	for _, name := range input.FieldList() {
		_, ok := output.Fields()[pdufield.Name(name)]
		assert.True(t, ok, name)
	}
	assert.Equal(t, input.TLVFields(), output.TLVFields())
}
