// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpptest

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"

	"github.com/tsocial/go-smpp/smpp/pdu"
	"github.com/tsocial/go-smpp/smpp/pdu/pdufield"
)

// Default settings.
var (
	DefaultUser     = "client"
	DefaultPasswd   = "secret"
	DefaultSystemID = "smpptest"
)

// HandlerFunc is the signature of a function passed to Server instances,
// that is called when client PDU messages arrive.
type HandlerFunc func(c Conn, m pdu.Body)

// Server is an SMPP server for testing purposes. By default it authenticate
// clients with the configured credentials, and echoes any other PDUs
// back to the client.
type Server struct {
	User    string
	Passwd  string
	TLS     *tls.Config
	Handler HandlerFunc

	conns []Conn
	mu    sync.Mutex
	l     net.Listener
}

// NewServer creates and initializes a new Server. Callers are supposed
// to call Close on that server later.
func NewServer() *Server {
	s := NewUnstartedServer()
	s.Start()
	return s
}

// NewUnstartedServer creates a new Server with default settings, and
// does not start it. Callers are supposed to call Start and Close later.
func NewUnstartedServer() *Server {
	return &Server{
		User:    DefaultUser,
		Passwd:  DefaultPasswd,
		Handler: EchoHandler,
		l:       newLocalListener(),
	}
}

func newLocalListener() net.Listener {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err == nil {
		return l
	}
	if l, err = net.Listen("tcp6", "[::1]:0"); err != nil {
		panic(fmt.Sprintf("smpptest: failed to listen on a port: %v", err))
	}
	return l
}

// Start starts the server.
func (srv *Server) Start() {
	go srv.Serve()
}

// Addr returns the local address of the server, or an empty string
// if the server hasn't been started yet.
func (srv *Server) Addr() string {
	if srv.l == nil {
		return ""
	}
	return srv.l.Addr().String()
}

// Close stops the server, causing the accept loop to break out.
func (srv *Server) Close() {
	if srv.l == nil {
		panic("smpptest: server is not started")
	}
	_ = srv.l.Close()
}

// Serve accepts new clients and handle them by authenticating the
// first PDU, expected to be a Bind PDU, then echoing all other PDUs.
func (srv *Server) Serve() {
	for {
		cli, err := srv.l.Accept()
		if err != nil {
			break // on srv.l.Close
		}

		c := newConn(cli)
		srv.conns = append(srv.conns, c)
		go srv.handle(c)
	}
}

// BroadcastMessage broadcasts a test PDU to the all bound clients
func (srv *Server) BroadcastMessage(p pdu.Body) {
	for i := range srv.conns {
		_ = srv.conns[i].Write(p)
	}
}

// handle new clients.
func (srv *Server) handle(c *conn) {
	defer func() {
		_ = c.Close()
	}()
	if err := srv.auth(c); err != nil {
		if !errors.Is(err, io.EOF) {
			log.Println("smpptest: server auth failed:", err)
		}
		return
	}
	for {
		p, err := c.Read()
		if err != nil {
			if !errors.Is(err, io.EOF) {
				log.Println("smpptest: read failed:", err)
			}
			break
		}
		srv.Handler(c, p)
	}
}

// auth authenticate new clients.
func (srv *Server) auth(c *conn) error {
	p, err := c.Read()
	if err != nil {
		return err
	}
	var resp pdu.Body
	switch p.Header().ID {
	case pdu.BindTransmitterID:
		resp = pdu.NewBindTransmitterResp()
	case pdu.BindReceiverID:
		resp = pdu.NewBindReceiverResp()
	case pdu.BindTransceiverID:
		resp = pdu.NewBindTransceiverResp()
	default:
		return errors.New("unexpected pdu, want bind")
	}
	f := p.Fields()
	user := f[pdufield.SystemID]
	passwd := f[pdufield.Password]
	if user == nil || passwd == nil {
		return errors.New("malformed pdu, missing system_id/password")
	}
	if user.String() != srv.User {
		return errors.New("invalid user")
	}
	if passwd.String() != srv.Passwd {
		return errors.New("invalid passwd")
	}
	_ = resp.Fields().Set(pdufield.SystemID, DefaultSystemID)

	return c.Write(resp)
}

// EchoHandler is the default Server HandlerFunc, and echoes back
// any PDUs received.
func EchoHandler(cli Conn, m pdu.Body) {
	// log.Printf("smpptest: echo PDU from %s: %#v", cli.RemoteAddr(), m)
	//
	// Real servers will reply with at least the same sequence number
	// from the request:
	//     resp := pdu.NewSubmitSMResp()
	//     resp.Header().Seq = m.Header().Seq
	//     resp.Fields().Set(pdufield.MessageID, "1234")
	//     cli.Write(resp)
	//
	// We just echo m back:
	_ = cli.Write(m)
}
