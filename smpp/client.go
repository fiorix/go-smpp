// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"crypto/tls"
	"io"
	"math"
	"sync"
	"time"

	"github.com/veoo/go-smpp/smpp/pdu"
	"github.com/veoo/go-smpp/smpp/pdu/pdufield"
)

// ConnStatus is an abstract interface for a connection status change.
type ConnStatus interface {
	Status() ConnStatusID
	Error() error
}

type connStatus struct {
	s   ConnStatusID
	err error
}

func (c *connStatus) Status() ConnStatusID { return c.s }
func (c *connStatus) Error() error         { return c.err }

// ConnStatusID represents a connection status change.
type ConnStatusID uint8

// Supported connection statuses.
const (
	Connected ConnStatusID = iota + 1
	Disconnected
	ConnectionFailed
	BindFailed
)

var connStatusText = map[ConnStatusID]string{
	Connected:        "Connected",
	Disconnected:     "Disconnected",
	ConnectionFailed: "Connection failed",
	BindFailed:       "Bind failed",
}

// String implements the Stringer interface.
func (cs ConnStatusID) String() string {
	return connStatusText[cs]
}

// ClientConn provides a persistent client connection that handles
// reconnection with a back-off algorithm.
type ClientConn interface {
	// Bind starts the client connection and returns a
	// channel that is triggered every time the connection
	// status changes.
	Bind() <-chan ConnStatus

	// Closer embeds the Closer interface. When Close is
	// called, client sends the Unbind command first and
	// terminates the connection upon response, or 1s timeout.
	Closer
}

// ConnMiddleware is useful for intercepting the traffic in/out happening
// for the SMPP server. It should be used in a read-only way since might
// effect the internal work of the SMPP server itself.
type ConnMiddleware func(conn Conn) Conn

// client provides a persistent client connection.
type client struct {
	Addr            string
	TLS             *tls.Config
	Status          chan ConnStatus
	BindFunc        func(c Conn) error
	EnquireLink     time.Duration
	RespTimeout     time.Duration
	WindowSize      uint
	ConnInterceptor ConnMiddleware

	// internal stuff.
	inbox chan pdu.Body
	conn  *connSwitch
	stop  chan struct{}
	once  sync.Once
}

func (c *client) init() {
	c.inbox = make(chan pdu.Body)
	c.conn = &connSwitch{}
	c.stop = make(chan struct{})
	if c.EnquireLink < 10*time.Second {
		c.EnquireLink = 10 * time.Second
	}
}

// Bind starts the connection manager and blocks until Close is called.
// It must be called in a goroutine.
func (c *client) Bind() {
	delay := 1.0
	const maxdelay = 120.0
	for !c.closed() {
		eli := make(chan struct{})
		conn, err := Dial(c.Addr, c.TLS)
		if err != nil {
			c.notify(&connStatus{
				s:   ConnectionFailed,
				err: err,
			})
			goto retry
		}
		if c.ConnInterceptor != nil {
			conn = c.ConnInterceptor(conn)
		}

		c.conn.Set(conn)
		if err = c.BindFunc(c.conn); err != nil {
			c.notify(&connStatus{s: BindFailed, err: err})
			goto retry
		}
		go c.enquireLink(eli)
		c.notify(&connStatus{s: Connected})
		delay = 1
		for {
			p, err := c.conn.Read()
			if err != nil {
				c.notify(&connStatus{
					s:   Disconnected,
					err: err,
				})
				break
			}
			switch p.Header().ID {
			case pdu.EnquireLinkID:
				pResp := pdu.NewEnquireLinkRespSeq(p.Header().Seq)
				err := c.conn.Write(pResp)
				if err != nil {
					break
				}
			case pdu.EnquireLinkRespID:
				// TODO: don't just ignore
			default:
				c.inbox <- p
			}
		}
	retry:
		close(eli)
		c.conn.Close()
		delay = math.Min(delay*math.E, maxdelay)
		c.trysleep(time.Duration(delay) * time.Second)
	}
	close(c.Status)
}

func (c *client) enquireLink(stop chan struct{}) {
	for {
		select {
		case <-time.After(c.EnquireLink):
			err := c.conn.Write(pdu.NewEnquireLink())
			if err != nil {
				return
			}
		case <-stop:
			return
		case <-c.stop:
			return
		}
	}
}

func (c *client) notify(ev ConnStatus) {
	select {
	case c.Status <- ev:
	default:
	}
}

// Read reads PDU binary data off the wire and returns it.
func (c *client) Read() (pdu.Body, error) {
	select {
	case pdu := <-c.inbox:
		return pdu, nil
	case <-c.stop:
		return nil, io.EOF
	}
}

// Write serializes the given PDU and writes to the connection.
func (c *client) Write(w pdu.Body) error {
	return c.conn.Write(w)
}

// Close terminates the current connection and stop any further attempts.
func (c *client) Close() error {
	c.once.Do(func() {
		close(c.stop)
		if err := c.conn.Write(pdu.NewUnbind()); err == nil {
			select {
			case <-c.inbox: // TODO: validate UnbindResp
			case <-time.After(time.Second):
			}
		}
		c.conn.Close()
	})
	return nil
}

// trysleep for the given duration, or return if Close is called.
func (c *client) trysleep(d time.Duration) {
	select {
	case <-time.After(d):
	case <-c.stop:
	}
}

// closed returns true after Close is called once.
func (c *client) closed() bool {
	select {
	case <-c.stop:
		return true
	default:
		return false
	}
}

// respTimeout returns a channel that fires based on the configured
// response timeout, or the default 1s.
func (c *client) respTimeout() <-chan time.Time {
	if c.RespTimeout == 0 {
		return time.After(time.Second)
	}
	return time.After(c.RespTimeout)
}

// bind attempts to bind the connection.
func bind(c Conn, p pdu.Body) (pdu.Body, error) {
	f := p.Fields()
	f.Set(pdufield.InterfaceVersion, 0x34)
	err := c.Write(p)
	if err != nil {
		return nil, err
	}
	resp, err := c.Read()
	if err != nil {
		return nil, err
	}
	h := resp.Header()
	if h.Status != 0 {
		return nil, h.Status
	}
	return resp, nil
}
