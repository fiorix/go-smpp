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
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
)

// Default settings.
var (
	DefaultUser           = "client"
	DefaultPasswd         = "secret"
	DefaultSystemID       = "smpptest"
	DeliverDelay          = 1 * time.Second
	msgIDcounter    int64 = 0
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

	mu sync.Mutex
	l  net.Listener
}

func nextMessageId() string {
	return strconv.FormatInt(atomic.AddInt64(&msgIDcounter, 1), 10)
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
	srv.l.Close()
}

// Serve accepts new clients and handle them by authenticating the
// first PDU, expected to be a Bind PDU, then echoing all other PDUs.
func (srv *Server) Serve() {
	for {
		cli, err := srv.l.Accept()
		if err != nil {
			break // on srv.l.Close
		}
		go srv.handle(newConn(cli))
	}
}

// handle new clients.
func (srv *Server) handle(c *conn) {
	defer c.Close()
	if err := srv.auth(c); err != nil {
		if err != io.EOF {
			log.Println("smpptest: server auth failed:", err)
		}
		return
	}
	for {
		pdu, err := c.Read()
		if err != nil {
			if err != io.EOF {
				log.Println("smpptest: read failed:", err)
			}
			break
		}
		srv.Handler(c, pdu)
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
	cli.Write(m)
}

// StubHandler is a HandlerFunc that returns compliant but dummy PDUs that are useful
// for testing clients
func StubHandler(conn Conn, m pdu.Body) {
	log.Println("smpptest: processing:", m.Header().ID)
	var resp pdu.Body
	switch m.Header().ID {
	case pdu.EnquireLinkID:
		resp = handleEnquireLink(conn, m)
	case pdu.EnquireLinkRespID:
		// TODO(cesar0094): what should happen if this is not received after request
		return
	case pdu.SubmitSMID:
		resp = handleSubmitSM(conn, m)
	case pdu.SubmitSMRespID:
		resp = handleInvalidCommand()
	case pdu.DeliverSMID:
		resp = handleInvalidCommand()
	case pdu.DeliverSMRespID:
		// TODO(cesar0094): Good to go?
		return
	default:
		// falls back to echoing the response
		EchoHandler(conn, m)
	}

	err := conn.Write(resp)
	if err != nil {
		log.Println("smpptest: error sending response:", err)
	}
}

func handleSubmitSM(conn Conn, m pdu.Body) pdu.Body {
	resp := pdu.NewSubmitSMResp()
	resp.Header().Seq = m.Header().Seq

	messageId := nextMessageId()
	resp.Fields().Set(pdufield.MessageID, messageId)
	m.Fields().Set(pdufield.MessageID, messageId)

	go processShortMessage(conn, m)
	return resp
}

func handleEnquireLink(conn Conn, m pdu.Body) pdu.Body {
	resp := pdu.NewEnquireLinkResp()
	resp.Header().Seq = m.Header().Seq
	return resp
}

func handleInvalidCommand() pdu.Body {
	resp := pdu.NewGenericNACK()
	resp.Header().Status = pdu.InvalidCommandID
	return resp
}

func processShortMessage(conn Conn, submitSmPdu pdu.Body) {
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
	respFields.Set(pdufield.RegisteredDelivery, smpp.FinalDeliveryReceipt)
	respFields.Set(pdufield.DataCoding, reqFields[pdufield.DataCoding])

	id := reqFields[pdufield.MessageID].String()
	// TODO(cesar0094): handle submitted and delivered ID.
	sub := "001"
	dlvrd := "001"
	stat := "DELIVRD"
	errTxt := "000"
	shortMessage := fmt.Sprintf("id:%s sub:%s dlvrd:%s submit date:%d done date:%d stat:%s err:%s Text:%s", id, sub, dlvrd, submitDate.Unix(), doneDate.Unix(), stat, errTxt, reqFields[pdufield.ShortMessage])
	respFields.Set(pdufield.ShortMessage, shortMessage)

	err := conn.Write(respPdu)
	if err != nil {
		log.Println("smpptest: failed sending delivery_sm:", err)
	}
}
