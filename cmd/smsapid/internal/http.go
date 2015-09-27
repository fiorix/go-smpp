// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"path/filepath"
	"strings"

	"golang.org/x/net/websocket"

	"github.com/fiorix/go-smpp/smpp"
)

// Handler is an HTTP handler that provides the endpoints of this service.
// It registers itself onto an existing ServeMux via Register.
type Handler struct {
	http.Handler

	// Prefix of all endpoints served by the handler.
	// Defaults to "/" if not set.
	Prefix string

	// VersionTag that follows the prefix.
	// Defaults to "v1" if not set.
	VersionTag string

	// SMPP Transceiver for sending and receiving SMS.
	// Register will update its Handler and Bind it.
	Tx *smpp.Transceiver

	// clients registered for receipt
	pool *deliveryPool

	// sm for smpp functionality
	sm *SM
}

func (h *Handler) init() {
	// TODO: handle nil h.Tx
	h.pool = &deliveryPool{m: make(map[string]chan *DeliveryReceipt)}
	h.sm = NewSM(h.Tx, rpc.NewServer())
	h.Tx.Handler = h.pool.Handler
}

// Register add the endpoints of this service to the given ServeMux,
// and binds Tx. Returns the ConnStatus channel from Tx.Bind.
//
// Must be called once, before the server is started.
func (h *Handler) Register(mux *http.ServeMux) <-chan smpp.ConnStatus {
	h.init()
	p := urlprefix(h)
	mux.Handle(p+"/send", h.send())
	mux.Handle(p+"/query", h.query())
	mux.Handle(p+"/sse", h.sse())
	mux.Handle(p+"/ws/jsonrpc", h.wsrpc())
	mux.Handle(p+"/ws/jsonrpc/events", h.wsrpcEvents())
	h.Handler = mux
	return h.Tx.Bind()
}

func urlprefix(h *Handler) string {
	path := "/" + h.Prefix + "/"
	if h.VersionTag == "" {
		path += "v1"
	} else {
		path += h.VersionTag
	}
	return strings.TrimRight(filepath.Clean(path), "/")
}

func (h *Handler) send() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(1 << 20)
		resp, status, err := h.sm.submit(r.Form)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
	return auth(cors(f, "PUT", "POST"))
}

func (h *Handler) query() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		resp, status, err := h.sm.query(r.Form)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
	return auth(cors(f, "HEAD", "GET"))
}

func (h *Handler) sse() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		n, ok := w.(http.CloseNotifier)
		if !ok {
			http.Error(w, "Notifier not supported",
				http.StatusInternalServerError)
			return
		}
		stop := n.CloseNotify()
		conn, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Flusher not supported",
				http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		conn.Flush()
		dr := h.pool.Register(r.RemoteAddr)
		defer h.pool.Unregister(r.RemoteAddr)
		j := json.NewEncoder(w)
		for {
			select {
			case r := <-dr:
				fmt.Fprintf(w, "Data: ")
				j.Encode(&r)
				fmt.Fprintf(w, "\n")
				conn.Flush()
			case <-stop:
				return
			}
		}
	}
	return auth(cors(f, "GET"))
}

// WebSocket handler for JSON RPC exposing functions from the SM type.
func (h *Handler) wsrpc() http.Handler {
	f := func(ws *websocket.Conn) {
		h.sm.rpc.ServeCodec(jsonrpc.NewServerCodec(ws))
	}
	return auth(cors(websocket.Handler(f).ServeHTTP, "GET"))
}

// WebSocket handler for JSON RPC events, we call the client.
func (h *Handler) wsrpcEvents() http.Handler {
	type conn struct {
		io.Reader
		io.WriteCloser
	}
	f := func(ws *websocket.Conn) {
		k := ws.Request().RemoteAddr
		dr := h.pool.Register(k)
		defer h.pool.Unregister(k)
		stop := make(chan struct{})
		r, w := io.Pipe()
		defer w.Close()
		go func() {
			io.Copy(w, ws)
			close(stop)
		}()
		rwc := &conn{Reader: r, WriteCloser: ws}
		cli := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(rwc))
		for {
			select {
			case r := <-dr:
				err := cli.Call("SM.Deliver", r, nil)
				if err != nil {
					return
				}
			case <-stop:
				return
			}
		}
	}
	return auth(cors(websocket.Handler(f).ServeHTTP, "GET"))
}
