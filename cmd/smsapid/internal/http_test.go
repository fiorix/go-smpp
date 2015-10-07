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
	"net/rpc"
	"net/rpc/jsonrpc"
	"net/textproto"
	"net/url"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/websocket"

	"github.com/fiorix/go-smpp/smpp"
)

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

func TestSend_Error(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: &smpp.Transceiver{Addr: ":0"}}
	<-h.Register(mux)
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

func TestSend_OK(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	<-h.Register(mux)
	defer h.Tx.Close()
	s := httptest.NewServer(mux)
	defer s.Close()
	resp, err := http.PostForm(s.URL+"/v1/send", url.Values{
		"dst":  {"root"},
		"text": {"gotcha"},
	})
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("unexpected status:", resp.Status)
	}
}

func TestQuery_Error(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: &smpp.Transceiver{Addr: ":0"}}
	<-h.Register(mux)
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
	id, dr := h.pool.Register()
	defer h.pool.Unregister(id)
	// make request
	resp, err := http.PostForm(s.URL+"/v1/send", url.Values{
		"dst":      {"root"},
		"text":     {"gotcha"},
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
		if r.Text != "delivery receipt here" {
			t.Fatalf("unexpected message: %#v", r)
		}
	case <-time.After(2 * time.Second):
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
	sse, err := sseClient(s.URL + "/v1/sse")
	if err != nil {
		t.Fatal(err)
	}
	// make request
	resp, err := http.PostForm(s.URL+"/v1/send", url.Values{
		"dst":      {"root"},
		"text":     {"gotcha"},
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
		var dr DeliveryReceipt
		err := json.Unmarshal([]byte(m.Data), &dr)
		if err != nil {
			t.Fatal(err)
		}
		test := []struct {
			Field, Want, Have string
		}{
			{"src", "", dr.Src},
			{"dst", "root", dr.Dst},
			{"msg", "delivery receipt here", dr.Text},
		}
		for _, el := range test {
			if el.Want != el.Have {
				t.Fatalf("unexpected value for %q: want %q, have %q",
					el.Field, el.Want, el.Have)
			}
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for delivery receipt")
	}
}

func TestWebSocket_Send(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	<-h.Register(mux)
	s := httptest.NewServer(mux)
	defer s.Close()
	url := strings.Replace(s.URL, "http:", "ws:", -1)
	ws, err := websocket.Dial(url+"/v1/ws/jsonrpc", "", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()
	cli := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(ws))
	args := &ShortMessage{
		Dst:  "root",
		Text: "hello world",
	}
	var resp ShortMessageResp
	err = cli.Call("SM.Submit", args, &resp)
	if err != nil {
		t.Fatal(err)
	}
	want := "foobar"
	if resp.MessageID != want {
		t.Fatalf("unexpected message id: want %q, have %q",
			want, resp.MessageID)
	}
}

func TestWebSocket_Query(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: newTransceiver()}
	<-h.Register(mux)
	s := httptest.NewServer(mux)
	defer s.Close()
	url := strings.Replace(s.URL, "http:", "ws:", 1)
	ws, err := websocket.Dial(url+"/v1/ws/jsonrpc", "", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()
	cli := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(ws))
	args := &QueryMessage{
		Src:       "nobody",
		MessageID: "foobar",
	}
	var resp QueryMessageResp
	err = cli.Call("SM.Query", args, &resp)
	if err != nil {
		t.Fatal(err)
	}
	want := "DELIVERED"
	if resp.MsgState != want {
		t.Fatalf("unexpected message state: want %q, have %q",
			want, resp.MsgState)
	}
}
