// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdutlv

import (
	"testing"
)

func TestMapSet(t *testing.T) {
	m := make(Map)
	test := []struct {
		k  Tag
		v  interface{}
		ok bool
	}{
		{TagDestAddrSubunit, nil, true},
		{TagDestAddrSubunit, "hello", true},
		{TagDestAddrSubunit, []byte("hello"), true},
		{TagDestBearerType, nil, true},
		{TagDestBearerType, uint8(1), true},
		{TagDestBearerType, int(1), true},
		{TagDestBearerType, t, false},
		{TagDestBearerType, String("hello"), true},
		{TagDestBearerType, CString("hello\x00"), true},
		{TagDestBearerType, CString("hello"), true},
		{TagDestBearerType, NewTLV(TagDestBearerType, []byte{0x03}), true},
	}
	for _, el := range test {
		if err := m.Set(el.k, el.v); el.ok && err != nil {
			t.Fatal(err)
		} else if !el.ok && err == nil {
			t.Fatalf("unexpected set of %q=%#v", el.k, el.v)
		}
	}
}
