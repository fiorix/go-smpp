// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// List is a list of PDU fields.
type List []Name

// Decode decodes binary data in the given buffer to build a Map.
//
// If the ShortMessage field is present, and DataCoding as well,
// we attempt to decode text automatically. See pdutext package
// for more information.
func (l List) Decode(r *bytes.Buffer) (Map, error) {
	var unsuccessCount int
	var numDest int
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
			NumberDests,
			NoUnsuccess,
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
			if k == NoUnsuccess {
				unsuccessCount = int(b)
			} else if k == NumberDests {
				numDest = int(b)
			}
		case DestinationList:
			var destList []DestSme
			for i := 0; i < numDest; i++ {
				var dest DestSme
				// Read DestFlag
				b, err := r.ReadByte()
				if err == io.EOF {
					break loop
				}
				if err != nil {
					return nil, err
				}
				dest.Flag = Fixed{Data: b}
				// Read Ton
				b, err = r.ReadByte()
				if err == io.EOF {
					break loop
				}
				if err != nil {
					return nil, err
				}
				dest.Ton = Fixed{Data: b}
				// Read npi
				b, err = r.ReadByte()
				if err == io.EOF {
					break loop
				}
				if err != nil {
					return nil, err
				}
				dest.Npi = Fixed{Data: b}
				// Read address
				bt, err := r.ReadBytes(0x00)
				if err == io.EOF {
					break loop
				}
				if err != nil {
					return nil, err
				}
				dest.DestAddr = Variable{Data: bt}
				destList = append(destList, dest)
			}
			f[k] = &DestSmeList{Data: destList}
		case UnsuccessSme:
			var unsList []UnSme
			for i := 0; i < unsuccessCount; i++ {
				var uns UnSme
				// Read Ton
				b, err := r.ReadByte()
				if err == io.EOF {
					break loop
				}
				if err != nil {
					return nil, err
				}
				uns.Ton = Fixed{Data: b}
				// Read npi
				b, err = r.ReadByte()
				if err == io.EOF {
					break loop
				}
				if err != nil {
					return nil, err
				}
				uns.Npi = Fixed{Data: b}
				// Read address
				bt, err := r.ReadBytes(0x00)
				if err == io.EOF {
					break loop
				}
				if err != nil {
					return nil, err
				}
				uns.DestAddr = Variable{Data: bt}
				// Read error code
				uns.ErrCode = Variable{Data: r.Next(4)}
				// Add unSme to the list
				unsList = append(unsList, uns)
			}
			f[k] = &UnSmeList{Data: unsList}
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
