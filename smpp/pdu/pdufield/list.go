// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"fmt"
	"io"
)

// TruncatedError reports a PDU body that ends inside its mandatory
// fields: Field is the first field that could not be read in full, and
// Offset is where that field starts in the body. It wraps
// io.ErrUnexpectedEOF.
type TruncatedError struct {
	Field  Name
	Offset int
}

func (e *TruncatedError) Error() string {
	return fmt.Sprintf("pdu body truncated: mandatory field %q at offset %d is incomplete", e.Field, e.Offset)
}

// Unwrap allows errors.Is(err, io.ErrUnexpectedEOF).
func (e *TruncatedError) Unwrap() error { return io.ErrUnexpectedEOF }

func truncated(k Name, off int) error {
	return &TruncatedError{Field: k, Offset: off}
}

// List is a list of PDU fields.
type List []Name

// Decode decodes binary data in the given buffer to build a Map.
//
// Every field in the list is mandatory. If the buffer ends before the
// last one has been read in full, Decode returns a *TruncatedError
// naming the field it could not read, instead of a Map that silently
// lacks it and every field after it.
//
// If the ShortMessage field is present, and DataCoding as well,
// we attempt to decode text automatically. See pdutext package
// for more information.
func (l List) Decode(r *bytes.Buffer) (Map, error) {
	var (
		unsuccessCount, numDest, udhLength, smLength int

		udhiFlag bool
	)
	f := make(Map)
	size := r.Len()
	for _, k := range l {
		off := size - r.Len()
		switch k {
		case
			AddressRange,
			DestinationAddr,
			ESMAddr,
			FinalDate,
			MessageID,
			Password,
			ScheduleDeliveryTime,
			ServiceType,
			SourceAddr,
			SystemID,
			SystemType,
			ValidityPeriod:
			b, err := r.ReadBytes(0x00)
			if err == io.EOF {
				return nil, truncated(k, off)
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
			ESMAddrNPI,
			ESMAddrTON,
			ESMClass,
			ErrorCode,
			InterfaceVersion,
			MessageState,
			NumberDests,
			NoUnsuccess,
			PriorityFlag,
			ProtocolID,
			RegisteredDelivery,
			ReplaceIfPresentFlag,
			SMDefaultMsgID,
			SourceAddrNPI,
			SourceAddrTON,
			SMLength:
			b, err := r.ReadByte()
			if err == io.EOF {
				return nil, truncated(k, off)
			}
			if err != nil {
				return nil, err
			}
			f[k] = &Fixed{Data: b}
			switch k {
			case NoUnsuccess:
				unsuccessCount = int(b)
			case NumberDests:
				numDest = int(b)
			case SMLength:
				smLength = int(b)
			case ESMClass:
				mask := byte(1 << 6)
				udhiFlag = mask == b&mask
			}
		case UDHLength:
			if !udhiFlag {
				continue
			}
			b, err := r.ReadByte()
			if err == io.EOF {
				return nil, truncated(k, off)
			}
			if err != nil {
				return nil, err
			}
			udhLength = int(b)
			f[k] = &Fixed{Data: b}
		case GSMUserData:
			if !udhiFlag {
				continue
			}
			var udhList []UDH
			var l int
			for i := udhLength; i > 0; i -= l + 2 {
				var udh UDH
				// Read IEI
				b, err := r.ReadByte()
				if err == io.EOF {
					return nil, truncated(k, off)
				}
				if err != nil {
					return nil, err
				}
				udh.IEI = Fixed{Data: b}
				// Read IELength
				b, err = r.ReadByte()
				if err == io.EOF {
					return nil, truncated(k, off)
				}
				if err != nil {
					return nil, err
				}
				l = int(b)
				udh.IELength = Fixed{Data: b}
				// Read IEData
				bt := r.Next(l)
				udh.IEData = SM{Data: bt}
				udhList = append(udhList, udh)
				if len(bt) != l {
					return nil, truncated(k, off)
				}
			}
			f[k] = &UDHList{Data: udhList}
		case DestinationList:
			var destList []DestSme
			for i := 0; i < numDest; i++ {
				var dest DestSme
				// Read DestFlag
				b, err := r.ReadByte()
				if err == io.EOF {
					return nil, truncated(k, off)
				}
				if err != nil {
					return nil, err
				}
				dest.Flag = Fixed{Data: b}
				// Read Ton
				b, err = r.ReadByte()
				if err == io.EOF {
					return nil, truncated(k, off)
				}
				if err != nil {
					return nil, err
				}
				dest.Ton = Fixed{Data: b}
				// Read npi
				b, err = r.ReadByte()
				if err == io.EOF {
					return nil, truncated(k, off)
				}
				if err != nil {
					return nil, err
				}
				dest.Npi = Fixed{Data: b}
				// Read address
				bt, err := r.ReadBytes(0x00)
				if err == io.EOF {
					return nil, truncated(k, off)
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
					return nil, truncated(k, off)
				}
				if err != nil {
					return nil, err
				}
				uns.Ton = Fixed{Data: b}
				// Read npi
				b, err = r.ReadByte()
				if err == io.EOF {
					return nil, truncated(k, off)
				}
				if err != nil {
					return nil, err
				}
				uns.Npi = Fixed{Data: b}
				// Read address
				bt, err := r.ReadBytes(0x00)
				if err == io.EOF {
					return nil, truncated(k, off)
				}
				if err != nil {
					return nil, err
				}
				uns.DestAddr = Variable{Data: bt}
				// Read error code (4 bytes, big-endian uint32)
				if r.Len() < 4 {
					return nil, truncated(k, off)
				}
				copy(uns.ErrCode[:], r.Next(4))
				// Add unSme to the list
				unsList = append(unsList, uns)
			}
			f[k] = &UnSmeList{Data: unsList}
		case ShortMessage:
			// Check UDHLength
			if udhLength > 0 {
				if smLength-udhLength-1 < 0 {
					return nil, fmt.Errorf("smLength is lesser than udhLength+1: have %d and %d",
						smLength, udhLength)
				}
				smLength -= udhLength + 1
				f[SMLength] = &Fixed{Data: byte(smLength)}
			}
			// Check SMLength
			if r.Len() < smLength {
				return nil, truncated(k, off)
			}
			f[ShortMessage] = &SM{Data: r.Next(smLength)}
		}
	}
	return f, nil
}
