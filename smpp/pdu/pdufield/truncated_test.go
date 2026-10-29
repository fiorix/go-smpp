// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestListDecoder_Truncated(t *testing.T) {
	bind := List{SystemID, Password, InterfaceVersion}
	for _, tc := range []struct {
		l      List
		in     string
		field  Name
		offset int
	}{
		{bind, "", SystemID, 0},
		{bind, "abc", SystemID, 0},
		{bind, "abc\x00", Password, 4},
		{bind, "abc\x00pw", Password, 4},
		{bind, "abc\x00pw\x00", InterfaceVersion, 7},
		// esm_class with UDHI, sm_length 10, UDH length 5, IE 00 03
		// but only one byte of its three data bytes.
		{List{ESMClass, SMLength, UDHLength, GSMUserData, ShortMessage},
			"\x40\x0a\x05\x00\x03\x2a", GSMUserData, 3},
		// sm_length 5, two bytes of short_message.
		{List{SMLength, ShortMessage}, "\x05hi", ShortMessage, 1},
		// One unsuccess_sme whose 4-byte error_status_code is cut.
		{List{NoUnsuccess, UnsuccessSme}, "\x01\x01\x01a\x00\x00\x00", UnsuccessSme, 1},
	} {
		m, err := tc.l.Decode(bytes.NewBufferString(tc.in))
		var te *TruncatedError
		if !errors.As(err, &te) {
			t.Fatalf("%q: want *TruncatedError, have map %#v, err %v", tc.in, m, err)
		}
		if te.Field != tc.field || te.Offset != tc.offset {
			t.Fatalf("%q: want field %q offset %d, have %q %d", tc.in, tc.field, tc.offset, te.Field, te.Offset)
		}
		if !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("%q: error does not wrap io.ErrUnexpectedEOF: %v", tc.in, err)
		}
		if m != nil {
			t.Fatalf("%q: want no map with the error, have %#v", tc.in, m)
		}
	}
	m, err := bind.Decode(bytes.NewBufferString("abc\x00pw\x00\x34"))
	if err != nil {
		t.Fatalf("complete body: %v", err)
	}
	if len(m) != 3 {
		t.Fatalf("complete body: want 3 fields, have %#v", m)
	}
}
