// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"io"
	"strconv"
)

// Name is the name of a PDU field.
type Name string

// Supported PDU field names.
const (
	AddrNPI              Name = "addr_npi"
	AddrTON              Name = "addr_ton"
	AddressRange         Name = "address_range"
	DataCoding           Name = "data_coding"
	DestAddrNPI          Name = "dest_addr_npi"
	DestAddrTON          Name = "dest_addr_ton"
	DestinationAddr      Name = "destination_addr"
	ESMClass             Name = "esm_class"
	ErrorCode            Name = "error_code"
	FinalDate            Name = "final_date"
	InterfaceVersion     Name = "interface_version"
	MessageID            Name = "message_id"
	MessageState         Name = "message_state"
	Password             Name = "password"
	PriorityFlag         Name = "priority_flag"
	ProtocolID           Name = "protocol_id"
	RegisteredDelivery   Name = "registered_delivery"
	ReplaceIfPresentFlag Name = "replace_if_present_flag"
	SMDefaultMsgID       Name = "sm_default_msg_id"
	SMLength             Name = "sm_length"
	ScheduleDeliveryTime Name = "schedule_delivery_time"
	ServiceType          Name = "service_type"
	ShortMessage         Name = "short_message"
	SourceAddr           Name = "source_addr"
	SourceAddrNPI        Name = "source_addr_npi"
	SourceAddrTON        Name = "source_addr_ton"
	SystemID             Name = "system_id"
	SystemType           Name = "system_type"
	ValidityPeriod       Name = "validity_period"
)

// Fixed is a PDU of fixed length.
type Fixed struct {
	Data uint8
}

// Len implements the Data interface.
func (f *Fixed) Len() int {
	return 1
}

// Raw implements the Data interface.
func (f *Fixed) Raw() interface{} {
	return f.Data
}

// String implements the Data interface.
func (f *Fixed) String() string {
	return strconv.Itoa(int(f.Data))
}

// Bytes implements the Data interface.
func (f *Fixed) Bytes() []byte {
	return []byte{f.Data}
}

// SerializeTo implements the Data interface.
func (f *Fixed) SerializeTo(w io.Writer) error {
	_, err := w.Write(f.Bytes())
	return err
}

// Variable is a PDU field of variable length.
type Variable struct {
	Data []byte
}

// Len implements the Data interface.
func (v *Variable) Len() int {
	return len(v.Bytes())
}

// Raw implements the Data interface.
func (v *Variable) Raw() interface{} {
	return v.Data
}

// String implements the Data interface.
func (v *Variable) String() string {
	if l := len(v.Data); l > 0 && v.Data[l-1] == 0x00 {
		return string(v.Data[:l-1])
	}
	return string(v.Data)
}

// Bytes implements the Data interface.
func (v *Variable) Bytes() []byte {
	if len(v.Data) > 0 && v.Data[len(v.Data)-1] == 0x00 {
		return v.Data
	}
	return append(v.Data, 0x00)
}

// SerializeTo implements the Data interface.
func (v *Variable) SerializeTo(w io.Writer) error {
	_, err := w.Write(v.Bytes())
	return err
}

// SM is a PDU field used for Short Messages.
type SM struct {
	Data []byte
}

// Len implements the Data interface.
func (sm *SM) Len() int {
	return len(sm.Data)
}

// Raw implements the Data interface.
func (sm *SM) Raw() interface{} {
	return sm.Data
}

// String implements the Data interface.
func (sm *SM) String() string {
	return string(sm.Data)
}

// Bytes implements the Data interface.
func (sm *SM) Bytes() []byte {
	return sm.Data
}

// SerializeTo implements the Data interface.
func (sm *SM) SerializeTo(w io.Writer) error {
	_, err := w.Write(sm.Bytes())
	return err
}
