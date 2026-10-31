// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdu

import (
	"sync/atomic"
	"testing"
)

// Sequence numbers must stay within SMPP 3.4 §5.1.4's 1..0x7FFFFFFF range,
// including across the wrap at 0x7FFFFFFF.
func TestSequenceNumber_StaysInSpecRange(t *testing.T) {
	// Drive the counter across the 0x7FFFFFFF boundary.
	atomic.StoreUint32(&nextSeq, 0x7FFFFFFE)
	for i := 0; i < 5; i++ {
		c := &codec{h: &Header{}}
		c.init()
		if c.h.Seq == 0 || c.h.Seq > 0x7FFFFFFF {
			t.Fatalf("seq %#x out of range 1..0x7FFFFFFF", c.h.Seq)
		}
	}
}
