// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"bufio"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/textproto"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/smpptest"
)

var ts *smpptest.Server

func TestMain(m *testing.M) {
	ts = smpptest.NewUnstartedServer()
	ts.Handler = pduHandler
	ts.Start()
	defer ts.Close()
	os.Exit(m.Run())
}

func pduHandler(c smpptest.Conn, p pdu.Body) {
	src := p.Fields()[pdufield.SourceAddr]
	fail := src != nil && src.String() == "root"
	switch p.Header().ID {
	case pdu.SubmitSMID:
		r := pdu.NewSubmitSMResp()
		r.Header().Seq = p.Header().Seq
		if fail {
			r.Header().Status = 0x00000045 // submitsm failed
		} else {
			r.Fields().Set(pdufield.MessageID, "foobar")
		}
		c.Write(r)
		rd := p.Fields()[pdufield.RegisteredDelivery]
		if rd == nil || rd.Bytes()[0] == 0 {
			return
		}
		r = pdu.NewDeliverSM()
		rf := r.Fields()
		pf := p.Fields()
		rf.Set(pdufield.SourceAddr, pf[pdufield.SourceAddr])
		rf.Set(pdufield.DestinationAddr, pf[pdufield.DestinationAddr])
		rf.Set(pdufield.ShortMessage, "delivery receipt here")
		c.Write(r)
	case pdu.QuerySMID:
		r := pdu.NewQuerySMResp()
		r.Header().Seq = p.Header().Seq
		if fail {
			r.Header().Status = 0x00000067 // querysm failed
		} else {
			pf := p.Fields()
			rf := r.Fields()
			rf.Set(pdufield.MessageID, pf[pdufield.MessageID])
			rf.Set(pdufield.MessageState, 2) // DELIVERED
		}
		c.Write(r)
	default:
		smpptest.EchoHandler(c, p)
	}
}

func newTransceiver() *smpp.Transceiver {
	return &smpp.Transceiver{
		Addr:   ts.Addr(),
		User:   smpptest.DefaultUser,
		Passwd: smpptest.DefaultPasswd,
	}
}

func TestHandler_Version(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{
		VersionTag: "v2",
		Tx:         newTransceiver(),
	}
	h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	resp, err := http.Get(s.URL + "/v2/send") // causes 405 not 404
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatal("unexpected status:", resp.Status)
	}
}

func TestSend_BadRequest(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	resp, err := http.PostForm(s.URL+"/v1/send", url.Values{})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("unexpected status:", resp.Status)
	}
}

func TestSend_EncParam(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: &smpp.Transceiver{Addr: ":0"}}
	h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	for _, enc := range []string{"latin1", "ucs2"} {
		resp, err := http.PostForm(s.URL+"/v1/send", url.Values{
			"dst": {"root"},
			"msg": {"gotcha"},
			"enc": {enc},
		})
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Fatal("unexpected status:", resp.Status)
		}
	}
}

func TestSend_RegisterParam(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: &smpp.Transceiver{Addr: ":0"}}
	h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	for _, reg := range []string{"final", "failure"} {
		resp, err := http.PostForm(s.URL+"/v1/send", url.Values{
			"dst":      {"root"},
			"msg":      {"gotcha"},
			"register": {reg},
		})
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusServiceUnavailable {
			t.Fatal("unexpected status:", resp.Status)
		}
	}
}

func TestSend_BadGateway(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	resp, err := http.PostForm(s.URL+"/v1/send", url.Values{
		"src": {"root"},
		"dst": {"root"},
		"msg": {"gotcha"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatal("unexpected status:", resp.Status)
	}
}

func TestSend_OK(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	<-h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	resp, err := http.PostForm(s.URL+"/v1/send", url.Values{
		"dst": {"root"},
		"msg": {"gotcha"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("unexpected status:", resp.Status)
	}
}

func TestQuery_BadRequest(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	resp, err := http.Get(s.URL + "/v1/query")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatal("unexpected status:", resp.Status)
	}
}

func TestQuery_ServiceUnavailable(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: &smpp.Transceiver{Addr: ":0"}}
	h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	p := url.Values{
		"src":        {"root"},
		"message_id": {"foobar"},
	}
	resp, err := http.Get(s.URL + "/v1/query?" + p.Encode())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatal("unexpected status:", resp.Status)
	}
}

func TestQuery_BadGateway(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	p := url.Values{
		"src":        {"root"},
		"message_id": {"foobar"},
	}
	resp, err := http.Get(s.URL + "/v1/query?" + p.Encode())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadGateway {
		t.Fatal("unexpected status:", resp.Status)
	}
}

func TestQuery_OK(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	<-h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	p := url.Values{
		"src":        {"nobody"},
		"message_id": {"foobar"},
	}
	resp, err := http.Get(s.URL + "/v1/query?" + p.Encode())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("unexpected status:", resp.Status)
	}
}

func TestDeliveryReceipt(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	<-h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	// cheat: register ourselves for delivery
	dr := h.ds.Register("foobar")
	defer h.ds.Unregister("foobar")
	// make request
	resp, err := http.PostForm(s.URL+"/v1/send", url.Values{
		"dst":      {"root"},
		"msg":      {"gotcha"},
		"register": {"final"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("unexpected status:", resp.Status)
	}
	select {
	case r := <-dr:
		if r.Msg != "delivery receipt here" {
			t.Fatalf("unexpected message: %#v", r)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for delivery receipt")
	}
}

type serverSentEvent struct {
	Event string
	Data  string
	Error error
}

// sseClient is a specialized SSE client that connects to a server and
// issues a request for the events handler, then waits for events to be
// returned from the server and puts them in the returned channel. It
// only handles the initial connect event and one subsequent event.
// This client supports HTTP/1.1 on non-TLS sockets.
func sseClient(serverURL string) (chan *serverSentEvent, error) {
	u, err := url.Parse(serverURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "http" {
		return nil, errors.New("Unsupported URL scheme")
	}
	ev := make(chan *serverSentEvent, 2)
	tp, err := textproto.Dial("tcp", u.Host)
	if err != nil {
		return nil, err
	}
	tp.Cmd("GET %s HTTP/1.1\r\n", u.Path)
	line, err := tp.ReadLine()
	if err != nil {
		tp.Close()
		return nil, err
	}
	if line != "HTTP/1.1 200 OK" {
		tp.Close()
		return nil, errors.New("Unexpected response:" + line)
	}
	m, err := tp.ReadMIMEHeader()
	if err != nil {
		tp.Close()
		return nil, err
	}
	if v := m.Get("Content-Type"); v != "text/event-stream" {
		tp.Close()
		return nil, errors.New("Unexpected Content-Type: " + v)
	}
	if m.Get("Transfer-Encoding") == "chunked" {
		tp.R = bufio.NewReader(httputil.NewChunkedReader(tp.R))
	}
	go func() {
		defer close(ev)
		defer tp.Close()
		m, err = tp.ReadMIMEHeader()
		if err != nil {
			ev <- &serverSentEvent{Error: err}
			return
		}
		ev <- &serverSentEvent{
			Event: m.Get("Event"),
			Data:  m.Get("Data"),
		}
		if m.Get("Event") != "connect" {
			return
		}
		// If the first event is connect, we proceed and ship
		// the next one in line.
		m, err = tp.ReadMIMEHeader()
		if err != nil {
			ev <- &serverSentEvent{Error: err}
			return
		}
		ev <- &serverSentEvent{
			Event: m.Get("Event"),
			Data:  m.Get("Data"),
		}
	}()
	return ev, nil
}

func TestSSE(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	<-h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	sse, err := sseClient(s.URL + "/v1/delivery")
	if err != nil {
		t.Fatal(err)
	}
	// make request
	resp, err := http.PostForm(s.URL+"/v1/send", url.Values{
		"dst":      {"root"},
		"msg":      {"gotcha"},
		"register": {"final"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("unexpected status:", resp.Status)
	}
	// handle delivery via sse
	select {
	case m := <-sse:
		if m == nil {
			t.Fatal("unexpected receipt: empty")
		}
		var dr deliveryReceipt
		err := json.Unmarshal([]byte(m.Data), &dr)
		if err != nil {
			t.Fatal(err)
		}
		test := []struct {
			Field, Want, Have string
		}{
			{"src", "", dr.Src},
			{"dst", "root", dr.Dst},
			{"msg", "delivery receipt here", dr.Msg},
		}
		for _, el := range test {
			if el.Want != el.Have {
				t.Fatalf("unexpected value for %q: want %q, have %q",
					el.Field, el.Want, el.Have)
			}
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for delivery receipt")
	}
}
