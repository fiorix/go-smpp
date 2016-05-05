// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/veoo/go-smpp/smpp/logger"
	"github.com/veoo/go-smpp/smpp/pdu"
	"github.com/veoo/go-smpp/smpp/pdu/pdufield"
)

// Default settings.
var (
	DefaultUser           = "client"
	DefaultPasswd         = "secret"
	DefaultSystemID       = "sys_id"
	DeliverDelay          = 1 * time.Second
	msgIdCounter    int64 = 0
)

var ServerStubHandlers = map[pdu.ID]func(pdu.Body) pdu.Body{
	pdu.EnquireLinkID:     handleEnquireLink,
	pdu.EnquireLinkRespID: handleEnquireLinkResp,
	pdu.SubmitSMID:        handleSubmitSM,
	pdu.SubmitSMRespID:    handleInvalidCommand,
	pdu.DeliverSMRespID:   handleDeliverSMResp,
}

// RequestHandlerFunc is the signature of a function passed to Server instances,
// that is called when client PDU messages arrive.
type RequestHandlerFunc func(s *session, m pdu.Body)

// Server is an SMPP server for testing purposes. By default it authenticate
// clients with the configured credentials, and echoes any other PDUs
// back to the client.
type Server struct {
	User    string
	Passwd  string
	TLS     *tls.Config
	Handler RequestHandlerFunc

	mu sync.Mutex
	l  net.Listener
}

func NextMessageId() string {
	return strconv.FormatInt(atomic.AddInt64(&msgIdCounter, 1), 10)
}

// NOTE: should handler funcs be session methods?
type session struct {
	conn         *connSwitch
	msgIDCounter int64
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

// NewServer creates and initializes a new Server. Callers are supposed
// to call Close on that server later.
func NewServer(user, password string, listener net.Listener) *Server {
	s := NewUnstartedServer(user, password, listener)
	s.Start()
	return s
}

// NewUnstartedServer creates a new Server with default settings, and
// does not start it. Callers are supposed to call Start and Close later.
func NewUnstartedServer(user, password string, listener net.Listener) *Server {
	if user == "" {
		user = DefaultUser
	}
	if password == "" {
		password = DefaultPasswd
	}
	return &Server{
		User:    user,
		Passwd:  password,
		Handler: EchoHandler,
		l:       listener,
	}
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
	srv.l.Close()
}

// Serve accepts new clients and handle them by authenticating the
// first PDU, expected to be a Bind PDU, then echoing all other PDUs.
func (srv *Server) Serve() {
	for {
		cli, err := srv.l.Accept()
		if err != nil {
			logger.Server.Error("Closing server:", err)
			break // on srv.l.Close
		}
		logger.Server.WithFields(log.Fields{
			"address": cli.RemoteAddr(),
		}).Info("New client")
		go srv.handle(newConn(cli))
	}
}

// handle new clients.
func (srv *Server) handle(c *conn) {
	defer c.Close()
	if err := srv.auth(c); err != nil {
		if err != io.EOF {
			logger.Server.Error("Server auth failed:", err)
		}
		return
	}
	// Use connSwitch to have synced read/write
	s := &session{conn: &connSwitch{}}
	s.conn.Set(c)

	for {
		pdu, err := s.Read()
		if err != nil {
			if err != io.EOF {
				logger.Server.Error("Read failed:", err)
			}
			break
		}
		srv.Handler(s, pdu)
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
	resp.Fields().Set(pdufield.SystemID, DefaultSystemID)
	if err = c.Write(resp); err != nil {
		return err
	}
	return nil
}

// EchoHandler is the default Server RequestHandlerFunc, and echoes back
// any PDUs received.
func EchoHandler(s *session, m pdu.Body) {
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

// StubHandler is a RequestHandlerFunc that returns compliant but dummy PDUs that are useful
// for testing clients
func StubHandler(s *session, m pdu.Body) {
	logger.Server.WithFields(log.Fields{
		"pudId": m.Header().ID.String(),
		"seq":   m.Header().Seq,
	}).Info("Processing incoming PDU")

	var resp pdu.Body
	if handler, ok := ServerStubHandlers[m.Header().ID]; ok {
		resp = handler(m)
	} else {
		logger.Server.Info(
			"Could not find handler matching PDU ID. Falling back to EchoHandler.")
		EchoHandler(s, m)
		return
	}

	if resp == nil {
		return
	}
	err := s.Write(resp)
	if err != nil {
		logger.Server.Error("Failed sending response:", err)
	}
}

func handleSubmitSM(m pdu.Body) pdu.Body {
	resp := pdu.NewSubmitSMResp()
	resp.Header().Seq = m.Header().Seq

	messageId := NextMessageId()
	resp.Fields().Set(pdufield.MessageID, messageId)
	m.Fields().Set(pdufield.MessageID, messageId)

	// TODO(cesar0094): "send" message and return deliverySM
	return resp
}

func handleEnquireLink(m pdu.Body) pdu.Body {
	resp := pdu.NewEnquireLinkResp()
	resp.Header().Seq = m.Header().Seq
	return resp
}

func handleEnquireLinkResp(m pdu.Body) pdu.Body {
	// TODO(cesar0094): what should happen if this is not received after request
	return nil
}

func handleDeliverSMResp(m pdu.Body) pdu.Body {
	// TODO(cesar0094): what should happen if this is not received after request
	return nil
}

func handleInvalidCommand(m pdu.Body) pdu.Body {
	resp := pdu.NewGenericNACK()
	resp.Header().Status = pdu.InvalidCommandID
	return resp
}

func processShortMessage(s *session, submitSmPdu pdu.Body) {
	submitDate := time.Now()
	// Pretend to be sending the SM
	time.Sleep(DeliverDelay)
	doneDate := time.Now()

	reqFields := submitSmPdu.Fields()
	respPdu := pdu.NewDeliverSM()
	respFields := respPdu.Fields()

	// Source and Destination info are reversed
	respFields.Set(pdufield.SourceAddrTON, reqFields[pdufield.DestAddrTON])
	respFields.Set(pdufield.SourceAddrNPI, reqFields[pdufield.DestAddrNPI])
	respFields.Set(pdufield.SourceAddr, reqFields[pdufield.DestinationAddr])

	respFields.Set(pdufield.DestAddrTON, reqFields[pdufield.SourceAddrTON])
	respFields.Set(pdufield.DestAddrNPI, reqFields[pdufield.SourceAddrNPI])
	respFields.Set(pdufield.DestinationAddr, reqFields[pdufield.SourceAddr])

	respFields.Set(pdufield.ServiceType, DefaultSystemID)
	respFields.Set(pdufield.ESMClass, reqFields[pdufield.ESMClass])
	respFields.Set(pdufield.ProtocolID, reqFields[pdufield.ProtocolID])
	respFields.Set(pdufield.PriorityFlag, reqFields[pdufield.PriorityFlag])
	respFields.Set(pdufield.RegisteredDelivery, FinalDeliveryReceipt)
	respFields.Set(pdufield.DataCoding, reqFields[pdufield.DataCoding])

	id := reqFields[pdufield.MessageID].String()
	// TODO(cesar0094): handle submitted and delivered ID.
	sub := "001"
	dlvrd := "001"
	stat := "DELIVRD"
	errTxt := "000"
	shortMessage := fmt.Sprintf("id:%s sub:%s dlvrd:%s submit date:%d done date:%d stat:%s err:%s Text:%s", id, sub, dlvrd, submitDate.Unix(), doneDate.Unix(), stat, errTxt, reqFields[pdufield.ShortMessage])
	respFields.Set(pdufield.ShortMessage, shortMessage)

	err := s.Write(respPdu)
	if err != nil {
		logger.Server.Error("Failed sending delivery_sm: ", err)
	}
}
