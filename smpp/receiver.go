// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
)

// Receiver implements an SMPP client receiver.
type Receiver struct {
	Addr               string
	User               string
	Passwd             string
	SystemType         string
	EnquireLink        time.Duration
	EnquireLinkTimeout time.Duration // Time after last EnquireLink response when connection considered down
	TLS                *tls.Config
	Handler            HandlerFunc

	conn struct {
		sync.Mutex
		*client
	}
}

// HandlerFunc is the handler function that a Receiver calls
// when a new PDU arrives.
type HandlerFunc func(p pdu.Body)

// Bind starts the Receiver. It creates a persistent connection
// to the server, update its status via the returned channel,
// and calls the registered Handler when new PDU arrives.
//
// Bind implements the ClientConn interface.
func (r *Receiver) Bind() <-chan ConnStatus {
	r.conn.Lock()
	defer r.conn.Unlock()
	if r.conn.client != nil {
		return r.conn.Status
	}
	c := &client{
		Addr:               r.Addr,
		TLS:                r.TLS,
		EnquireLink:        r.EnquireLink,
		EnquireLinkTimeout: r.EnquireLinkTimeout,
		Status:             make(chan ConnStatus, 1),
		BindFunc:           r.bindFunc,
	}
	r.conn.client = c
	c.init()
	go c.Bind()
	return c.Status
}

func (r *Receiver) bindFunc(c Conn) error {
	p := pdu.NewBindReceiver()
	f := p.Fields()
	f.Set(pdufield.SystemID, r.User)
	f.Set(pdufield.Password, r.Passwd)
	f.Set(pdufield.SystemType, r.SystemType)
	resp, err := bind(c, p)
	if err != nil {
		return err
	}
	if resp.Header().ID != pdu.BindReceiverRespID {
		return fmt.Errorf("unexpected response for BindReceiver: %s",
			resp.Header().ID)
	}
	if r.Handler != nil {
		go r.handlePDU()
	}
	return nil
}

func (r *Receiver) handlePDU() {
	for {
		pdu, err := r.conn.Read()
		if err != nil {
			break
		}
		r.Handler(pdu)
	}
}

// Close implements the ClientConn interface.
func (r *Receiver) Close() error {
	r.conn.Lock()
	defer r.conn.Unlock()
	if r.conn.client == nil {
		return ErrNotConnected
	}
	return r.conn.Close()
}
