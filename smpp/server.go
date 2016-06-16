// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/veoo/go-smpp/smpp/pdu"
	"github.com/veoo/go-smpp/smpp/pdu/pdufield"
)

// Default settings.
var (
	DefaultSystemID = "sys_id"
	IDLen           = 16
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"

// RequestHandlerFunc is the signature of a function passed to Server instances,
// that is called when client PDU messages arrive.
type RequestHandlerFunc func(Session, pdu.Body)

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

type Session interface {
	Reader
	Writer
	Closer
	ID() string
}

// NOTE: should handler funcs be session methods?
type session struct {
	conn *connSwitch
	id   string
}

// TODO(cesar0094): Make sure Read(), Write() and Close() are working as expected

// Read reads PDU binary data off the wire and returns it.
func (s *session) Read() (pdu.Body, error) {
	return s.conn.Read()
}

// Write serializes the given PDU and writes to the connection.
func (s *session) Write(w pdu.Body) error {
	return s.conn.Write(w)
}

// Close terminates the current connection and stop any further attempts.
func (s *session) Close() error {
	return s.conn.Close()
}

func (s *session) ID() string {
	return s.id
}

func NewLocalListener(port int) net.Listener {
	// Try the default port first
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err == nil {
		return l
	}
	if l, err = net.Listen("tcp", "127.0.0.1:0"); err == nil {
		return l
	}
	if l, err = net.Listen("tcp6", "[::1]:0"); err == nil {
		return l
	}
	panic(fmt.Sprintf("%s: failed to listen on a port: %v", DefaultSystemID, err))
}

type Server interface {
	Addr() string
	Close()
	Handle(id pdu.ID, h RequestHandlerFunc)
	Start()
	Serve()
	Session(id string) Session
}

// Server is an SMPP server for testing purposes. By default it authenticate
// clients with the configured credentials, and echoes any other PDUs
// back to the client.
type server struct {
	User     string
	Passwd   string
	systemId string
	TLS      *tls.Config

	m  map[pdu.ID]RequestHandlerFunc
	s  map[string]Session
	mu sync.Mutex
	l  net.Listener
}

// NewServer creates and initializes a new Server. Callers are supposed
// to call Close on that server later.
func NewServer(user, password string, listener net.Listener) *Server {
	s := NewUnstartedServer(user, password, listener)
	s.Start()
	return &s
}

// NewUnstartedServer creates a new Server with default settings, and
// does not start it. Callers are supposed to call Start and Close later.
func NewUnstartedServer(user, password string, listener net.Listener) Server {
	s := &server{
		User:   user,
		Passwd: password,
		m:      map[pdu.ID]RequestHandlerFunc{},
		s:      map[string]Session{},
		l:      listener,
	}
	return s
}

// Start starts the server.
func (srv *server) Start() {
	go srv.Serve()
}

// Addr returns the local address of the server, or an empty string
// if the server hasn't been started yet.
func (srv *server) Addr() string {
	if srv.l == nil {
		return ""
	}
	return srv.l.Addr().String()
}

// Close stops the server, causing the accept loop to break out.
func (srv *server) Close() {
	if srv.l == nil {
		panic("smpptest: server is not started")
	}
	srv.l.Close()
}

// Session returns the session provided the id from the map of sessions
func (srv *server) Session(id string) Session {
	return srv.s[id]
}

// Serve accepts new clients and handle them by authenticating the
// first PDU, expected to be a Bind PDU, then echoing all other PDUs.
func (srv *server) Serve() {
	for {
		cli, err := srv.l.Accept()
		if err != nil {
			log.Println("Closing server:", err)
			break // on srv.l.Close
		}
		log.Println("New client", cli.RemoteAddr())
		go srv.handle(newConn(cli))
	}
}

// handle new clients.
func (srv *server) handle(c *conn) {
	defer c.Close()
	if err := srv.auth(c); err != nil {
		if err != io.EOF {
			log.Println("Server auth failed:", err)
		}
		return
	}
	// Use connSwitch to have synced read/write
	s := &session{conn: &connSwitch{}}
	s.conn.Set(c)
	s.id = randomString(IDLen)
	srv.mu.Lock()
	srv.s[s.id] = s
	srv.mu.Unlock()
	for {
		p, err := s.Read()
		if err != nil {
			if err != io.EOF {
				log.Println("Read failed:", err)
			}
			break
		}
		h, ok := srv.m[p.Header().ID]
		if ok {
			h(s, p)
		} else {
			log.Println("Handler not found for:", p.Header().ID)
		}
	}
	srv.mu.Lock()
	delete(srv.s, s.id)
	srv.mu.Unlock()
}

func (srv *server) Handle(id pdu.ID, h RequestHandlerFunc) {
	srv.m[id] = h
}

// auth authenticate new clients.
func (srv *server) auth(c *conn) error {
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
	resp.Fields().Set(pdufield.SystemID, DefaultSystemID)
	if err = c.Write(resp); err != nil {
		return err
	}
	return nil
}

// EchoHandler is the default Server RequestHandlerFunc, and echoes back
// any PDUs received.
func EchoHandler(s Session, m pdu.Body) {
	// logger.Server.Printf("smpptest: echo PDU from %s: %#v", s.RemoteAddr(), m)
	//
	// Real servers will reply with at least the same sequence number
	// from the request:
	//     resp := pdu.NewSubmitSMResp()
	//     resp.Header().Seq = m.Header().Seq
	//     resp.Fields().Set(pdufield.MessageID, "1234")
	//     s.Write(resp)
	//
	// We just echo m back:
	s.Write(m)
}
