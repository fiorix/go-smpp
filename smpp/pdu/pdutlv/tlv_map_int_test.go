// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutlv

import (
	"bytes"
	"testing"
)

// uint16 and uint32 values must encode as big-endian multi-octet TLVs, not be
// dropped (the previous Map.Set had no case for them and returned an error).
func TestMapSet_Uint16And32(t *testing.T) {
	m := Map{}
	if err := m.Set(TagSarMsgRefNum, uint16(0x1234)); err != nil {
		t.Fatalf("Set uint16: %v", err)
	}
	if got := m[TagSarMsgRefNum].Bytes(); !bytes.Equal(got, []byte{0x12, 0x34}) {
		t.Fatalf("uint16 TLV = % X; want 12 34", got)
	}

	if err := m.Set(TagQosTimeToLive, uint32(0x01020304)); err != nil {
		t.Fatalf("Set uint32: %v", err)
	}
	if got := m[TagQosTimeToLive].Bytes(); !bytes.Equal(got, []byte{0x01, 0x02, 0x03, 0x04}) {
		t.Fatalf("uint32 TLV = % X; want 01 02 03 04", got)
	}
}
