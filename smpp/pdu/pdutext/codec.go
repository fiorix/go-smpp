// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

// DataCoding to define text codecs.
type DataCoding uint8

// Supported text codecs.
const (
	Latin1Type DataCoding = 0x03
	UCS2Type   DataCoding = 0x08
	SilentType DataCoding = 0xC0
)

// Codec defines a text codec.
type Codec interface {
	// Type returns the value for the data_coding PDU.
	Type() DataCoding

	// Encode text.
	Encode() []byte

	// Decode text.
	Decode() []byte
}

// Encode text.
func Encode(typ DataCoding, text []byte) []byte {
	switch typ {
	case Latin1Type:
		return Latin1(text).Encode()
	case UCS2Type:
		return UCS2(text).Encode()
	case SilentType:
		return Silent(text).Encode()
	default:
		return text
	}
}

// Decode text.
func Decode(typ DataCoding, text []byte) []byte {
	switch typ {
	case Latin1Type:
		return Latin1(text).Decode()
	case UCS2Type:
		return UCS2(text).Decode()
	default:
		return text
	}
}
