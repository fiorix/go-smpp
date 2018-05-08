// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpptest

import (
	"bufio"
	"bytes"
	"io"
	"net"

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
	return pdu.Decode(c.r)
}

// Write implements the Conn interface.
func (c *conn) Write(p pdu.Body) error {
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
	return c.rwc.Close()
}
