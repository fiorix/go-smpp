// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutext

// Silent text codec, no encoding, but uses 192 as data_coding
type Silent []byte

// Type implements the Codec interface.
func (s Silent) Type() DataCoding {
	return 0xC0
}

// Encode raw text.
func (s Silent) Encode() []byte {
	return []byte{}
}

// Decode raw text.
func (s Silent) Decode() []byte {
	return s
}
