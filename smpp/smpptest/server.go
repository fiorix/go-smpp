package main

import (
	"bufio"
	"bytes"
	//	"crypto/tls"
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"github.com/fiorix/go-smpp/smpp/pdu"
)

//TODO:TLS, Hijack.

//ResponseWriter is interface for Write answer to SMPP client.
type ResponseWriter interface {
	io.Writer

	WritePDU(pdu pdu.Body) error
}

//Handler created for API to handle the requests.
type Handler interface {
	ServerSMPP(pdu.Body, ResponseWriter)
}

//Server for smpp proto.
type Server struct {
	l net.Listener

	Network, Address string
	Handler          Handler
	ErrLog           *log.Logger
}

//ListenAndServe open new socket and start serving.
func (srv *Server) ListenAndServe() (err error) {
	if srv.Address == "" {
		return errors.New("Empty Address")
	}
	if srv.Handler == nil {
		return errors.New("Empty Handler")
	}
	if srv.Network == "" {
		srv.Network = "tcp"
	}
	srv.l, err = net.Listen(srv.Network, srv.Address)
	if err != nil {
		return err
	}
	return srv.Serve()
}

//Serve is loop for new accept new connections.
func (srv *Server) Serve() error {
	for {
		cli, err := srv.l.Accept()
		if err != nil {
			return err
		}
		c := srv.newConn(cli)
		go c.serve()
	}
}

func (srv *Server) newConn(rwc net.Conn) *conn {
	return &conn{
		srv:  srv,
		rwc:  rwc,
		mu:   sync.Mutex{},
		bufw: bufio.NewWriter(rwc),
		bufr: bufio.NewReader(rwc),
	}
}

type conn struct {
	srv  *Server
	rwc  net.Conn
	mu   sync.Mutex
	bufw *bufio.Writer
	bufr *bufio.Reader

	//	tlsState   *tls.ConnectionState
}

type connRespWriter struct {
	c *conn
}

func (crw *connRespWriter) WritePDU(p pdu.Body) error {
	var b bytes.Buffer
	err := p.SerializeTo(&b)
	if err != nil {
		return err
	}
	_, err = io.Copy(crw.c.bufw, &b)
	if err != nil {
		return err
	}
	return crw.c.bufw.Flush()
}

func (crw *connRespWriter) Write(p []byte) (int, error) {
	n, err := crw.c.bufw.Write(p)
	if err != nil {
		return n, err
	}
	return n, crw.c.bufw.Flush()
}

func (c *conn) serve() {
	for {
		pdus, err := pdu.Decode(c.bufr)
		if err != nil && c.srv.ErrLog != nil {
			c.srv.ErrLog.Println(err)
		}
		wrtr := &connRespWriter{c}
		c.srv.Handler.ServerSMPP(pdus, wrtr)
	}
}
