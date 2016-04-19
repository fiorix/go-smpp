// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"testing"
	"time"

	"github.com/veoo/go-smpp/smpp/pdu"
)

func TestReceiver(t *testing.T) {
	port := 0 // any port
	s := NewServer(DefaultUser, DefaultPasswd, NewLocalListener(port))
	defer s.Close()
	rc := make(chan pdu.Body)
	r := &Receiver{
		Addr:    s.Addr(),
		User:    DefaultUser,
		Passwd:  DefaultPasswd,
		Handler: func(p pdu.Body) { rc <- p },
	}
	defer r.Close()
	conn := <-r.Bind()
	switch conn.Status() {
	case Connected:
	default:
		t.Fatal(conn.Error())
	}
	// cheat: inject GenericNACK PDU for the server to echo back.
	p := pdu.NewGenericNACK()
	r.conn.Lock()
	r.conn.Write(p)
	r.conn.Unlock()
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
