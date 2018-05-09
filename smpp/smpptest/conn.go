// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpptest

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"sync"

	"github.com/tsocial/go-smpp/smpp/pdu"
)

// Conn implements a server side connection.
type Conn interface {
	// Write serializes the given PDU and writes to the connection.
	Write(p pdu.Body) error

	// Close terminates the connection.
	Close() error

	// RemoteAddr returns the peer address.
	RemoteAddr() net.Addr
}

// conn provides the basics of an SMPP connection.
type conn struct {
	rwc net.Conn
	r   *bufio.Reader
	w   *bufio.Writer
	l   sync.RWMutex
}

func newConn(c net.Conn) *conn {
	return &conn{
		rwc: c,
		r:   bufio.NewReader(c),
		w:   bufio.NewWriter(c),
	}
}

// RemoteAddr implements the Conn interface.
func (c *conn) RemoteAddr() net.Addr {
	return c.rwc.RemoteAddr()
}

// Read reads PDU off the wire.
func (c *conn) Read() (pdu.Body, error) {
	c.l.RLock()
	defer c.l.RUnlock()

	return pdu.Decode(c.r)
}

// Write implements the Conn interface.
func (c *conn) Write(p pdu.Body) error {
	c.l.Lock()
	defer c.l.Unlock()

	var b bytes.Buffer
	err := p.SerializeTo(&b)
	if err != nil {
		return err
	}
	_, err = io.Copy(c.w, &b)
	if err != nil {
		return err
	}
	return c.w.Flush()
}

// Close implements the Conn interface.
func (c *conn) Close() error {
	c.l.Lock()
	defer c.l.Unlock()

	return c.rwc.Close()
}
