// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/veoo/go-smpp/smpp/pdu/pdutext"
)

// List is a list of PDU fields.
type List []Name

// Decode decodes binary data in the given buffer to build a Map.
//
// If the ShortMessage field is present, and DataCoding as well,
// we attempt to decode text automatically. See pdutext package
// for more information.
func (l List) Decode(r *bytes.Buffer) (Map, error) {
	f := make(Map)
loop:
	for _, k := range l {
		switch k {
		case
			AddressRange,
			DestinationAddr,
			ErrorCode,
			FinalDate,
			MessageID,
			MessageState,
			Password,
			ScheduleDeliveryTime,
			ServiceType,
			SourceAddr,
			SystemID,
			SystemType,
			ValidityPeriod:
			b, err := r.ReadBytes(0x00)
			if err == io.EOF {
				break loop
			}
			if err != nil {
				return nil, err
			}
			f[k] = &Variable{Data: b}
		case
			AddrNPI,
			AddrTON,
			DataCoding,
			DestAddrNPI,
			DestAddrTON,
			ESMClass,
			InterfaceVersion,
			PriorityFlag,
			ProtocolID,
			RegisteredDelivery,
			ReplaceIfPresentFlag,
			SMDefaultMsgID,
			SourceAddrNPI,
			SourceAddrTON:
			b, err := r.ReadByte()
			if err == io.EOF {
				break loop
			}
			if err != nil {
				return nil, err
			}
			f[k] = &Fixed{Data: b}
		case SMLength:
			b, err := r.ReadByte()
			if err == io.EOF {
				break loop
			}
			if err != nil {
				return nil, err
			}
			l := int(b)
			f[k] = &Fixed{Data: b}
			if r.Len() < l {
				return nil, fmt.Errorf("short read for smlength: want %d, have %d",
					l, r.Len())
			}
			f[ShortMessage] = &SM{Data: r.Next(l)}
		case ShortMessage:
			sm, exists := f[ShortMessage].(*SM)
			if !exists {
				continue
			}
			c, exists := f[DataCoding].(*Fixed)
			if !exists {
				continue
			}
			sm.Data = pdutext.Decode(pdutext.DataCoding(c.Data), sm.Data)
		}
	}
	return f, nil
}

// DecodeTLV scans the given byte slice to build a TLVMap from binary data.
func (l List) DecodeTLV(r *bytes.Buffer) (TLVMap, error) {
	t := make(TLVMap)
	for r.Len() >= 4 {
		b := r.Next(4)
		ft := TLVType(binary.BigEndian.Uint16(b[0:2]))
		fl := binary.BigEndian.Uint16(b[2:4])
		if r.Len() < int(fl) {
			return nil, fmt.Errorf("not enough data for tag %#x: want %d, have %d",
				ft, fl, r.Len())
		}
		b = r.Next(int(fl))
		t[ft] = &TLVBody{
			Tag:  ft,
			Len:  fl,
			data: b,
		}
	}
	return t, nil
}
