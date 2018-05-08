// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"testing"
	"time"

	"github.com/tsocial/go-smpp/smpp/pdu"
	"github.com/tsocial/go-smpp/smpp/smpptest"
)

func TestReceiver(t *testing.T) {
	s := smpptest.NewServer()
	defer s.Close()
	rc := make(chan pdu.Body)
	r := &Receiver{
		Addr:    s.Addr(),
		User:    smpptest.DefaultUser,
		Passwd:  smpptest.DefaultPasswd,
		Handler: func(p pdu.Body) { rc <- p },
	}
	defer r.Close()
	conn := <-r.Bind()
	switch conn.Status() {
	case Connected:
	default:
		t.Fatal(conn.Error())
	}
	// trigger inbound message from server
	p := pdu.NewGenericNACK()
	s.BroadcastMessage(p)
	// check response.
	select {
	case m := <-rc:
		want, have := *p.Header(), *m.Header()
		if want != have {
			t.Fatalf("unexpected PDU: want %#v, have %#v",
				want, have)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for server to echo")
	}
}
