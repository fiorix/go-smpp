// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
)

// Transceiver implements an SMPP transceiver.
//
// The API is a combination of the Transmitter and Receiver.
type Transceiver struct {
	Addr        string
	User        string
	Passwd      string
	SystemType  string
	EnquireLink time.Duration
	TLS         *tls.Config
	Handler     HandlerFunc

	Transmitter
}

// Bind implements the ClientConn interface.
func (t *Transceiver) Bind() <-chan ConnStatus {
	t.conn.Lock()
	defer t.conn.Unlock()
	if t.conn.client != nil {
		return t.conn.Status
	}
	t.tx.Lock()
	t.tx.inflight = make(map[uint32]chan *tx)
	t.tx.Unlock()
	c := &client{
		Addr:        t.Addr,
		TLS:         t.TLS,
		EnquireLink: t.EnquireLink,
		Status:      make(chan ConnStatus, 1),
		BindFunc:    t.bindFunc,
	}
	t.conn.client = c
	go c.Bind()
	return c.Status
}

func (t *Transceiver) bindFunc(c Conn) error {
	p := pdu.NewBindTransceiver()
	f := p.Fields()
	f.Set(pdufield.SystemID, t.User)
	f.Set(pdufield.Password, t.Passwd)
	f.Set(pdufield.SystemType, t.SystemType)
	resp, err := bind(c, p)
	if err != nil {
		return err
	}
	if resp.Header().ID != pdu.BindTransceiverRespID {
		return fmt.Errorf("unexpected response for BindTransceiver: %s",
			resp.Header().ID)
	}
	go t.handlePDU(t.Handler)
	return nil
}
