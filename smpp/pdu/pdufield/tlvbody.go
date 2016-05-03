// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"encoding/binary"
	"encoding/json"
	"io"
)

// TODO(fiorix): Implement TLV parameters.

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

type TLVBodyJSON struct {
	Tag  TLVType `json:"tag"`
	Len  uint16  `json:"len"`
	Data []byte  `json:"data"`
}

func (tlv TLVBody) MarshalJSON() ([]byte, error) {
	s := TLVBodyJSON{
		Tag:  tlv.Tag,
		Len:  tlv.Len,
		Data: tlv.Bytes(),
	}
	return json.Marshal(s)
}

func (tlv TLVBody) UnmarshalJSON(b []byte) error {
	s := TLVBodyJSON{}
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	tlv.Tag = s.Tag
	tlv.Len = s.Len
	tlv.data = s.Data
	return nil
}
