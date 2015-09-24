// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

// Handler is an HTTP handler that provides the endpoints of this service.
// It registers itself onto an existing ServeMux via Register.
type Handler struct {
	http.Handler

	// Prefix of all endpoints served by the handler.
	// Defaults to "/" if not set.
	Prefix string

	// VersionTag that follows the prefix.
	// Defaults to "v1" of not set.
	VersionTag string

	// SMPP Transceiver for sending and receiving SMS.
	// Register will update its Handler and Bind it.
	Tx *smpp.Transceiver

	// clients registered for receipt
	ds *deliveryStore
}

func (h *Handler) init() <-chan smpp.ConnStatus {
	// TODO: handle nil h.Tx
	h.ds = &deliveryStore{m: make(map[string]chan *deliveryReceipt)}
	h.Tx.Handler = h.delivery
	return h.Tx.Bind()
}

// Register add the endpoints of this service to the given ServeMux,
// and binds Handler.Tx. Returns the ConnStatus channel from Bind.
//
// Must be called once, before the server is started.
func (h *Handler) Register(mux *http.ServeMux) <-chan smpp.ConnStatus {
	conn := h.init()
	p := urlprefix(h)
	mux.Handle(p+"/send", h.send())
	mux.Handle(p+"/query", h.query())
	mux.Handle(p+"/delivery", h.sse())
	h.Handler = mux
	return conn
}

func urlprefix(h *Handler) string {
	path := "/" + h.Prefix + "/"
	if h.VersionTag == "" {
		path += "v1"
	} else {
		path += h.VersionTag
	}
	return filepath.Clean(strings.TrimRight(path, "/"))
}

func (h *Handler) send() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		sm := &smpp.ShortMessage{}
		var msg, enc, reg string
		form := Form{
			{"src", "number of sender", false, nil, &sm.Src},
			{"dst", "number of recipient", true, nil, &sm.Dst},
			{"msg", "text message", true, nil, &msg},
			{"enc", "text encoding", false, []string{"latin1", "ucs2"}, &enc},
			{"register", "registered delivery", false, []string{"final", "failure"}, &reg},
		}
		if err := form.Validate(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		switch enc {
		case "":
			sm.Text = pdutext.Raw(msg)
		case "latin1", "latin-1":
			sm.Text = pdutext.Latin1(msg)
		case "ucs2", "ucs-2":
			sm.Text = pdutext.UCS2(msg)
		}
		switch reg {
		case "final":
			sm.Register = smpp.FinalDeliveryReceipt
		case "failure":
			sm.Register = smpp.FailureDeliveryReceipt
		}
		sm, err := h.Tx.Submit(sm)
		if err == smpp.ErrNotConnected {
			http.Error(w,
				http.StatusText(http.StatusServiceUnavailable),
				http.StatusServiceUnavailable)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			MessageID string `json:"message_id"`
		}{sm.RespID()})
	}
	return auth(cors(f, "PUT", "POST"))
}

func (h *Handler) query() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		var src, msgid string
		form := Form{
			{"src", "number of sender", false, nil, &src},
			{"message_id", "message id from send", true, nil, &msgid},
		}
		if err := form.Validate(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		qr, err := h.Tx.QuerySM(src, msgid)
		if err != nil {
		}
		if err == smpp.ErrNotConnected {
			http.Error(w,
				http.StatusText(http.StatusServiceUnavailable),
				http.StatusServiceUnavailable)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			MsgState  string `json:"message_state"`
			FinalDate string `json:"final_date"`
			ErrCode   uint8  `json:"error_code"`
		}{qr.MsgState, qr.FinalDate, qr.ErrCode})
	}
	return auth(cors(f, "HEAD", "GET"))
}

func (h *Handler) sse() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		stop, ok := w.(http.CloseNotifier)
		if !ok {
			http.Error(w, "Notifier not supported",
				http.StatusInternalServerError)
			return
		}
		conn, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Flusher not supported",
				http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		conn.Flush()
		dr := h.ds.Register(r.RemoteAddr)
		defer h.ds.Unregister(r.RemoteAddr)
		j := json.NewEncoder(w)
		for {
			select {
			case r := <-dr:
				fmt.Fprintf(w, "data: ")
				j.Encode(&r)
				fmt.Fprintf(w, "\n")
				conn.Flush()
			case <-stop.CloseNotify():
				return
			}
		}
	}
	return auth(cors(f, "GET"))
}

func (h *Handler) delivery(p pdu.Body) {
	switch p.Header().ID {
	case pdu.DeliverSMID:
		f := p.Fields()
		dr := &deliveryReceipt{
			Src: f[pdufield.SourceAddr].String(),
			Dst: f[pdufield.DestinationAddr].String(),
			Msg: f[pdufield.ShortMessage].String(),
		}
		h.ds.Broadcast(dr)
	}
}

type deliveryReceipt struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
	Msg string `json:"message"`
}

type deliveryStore struct {
	mu sync.Mutex
	m  map[string]chan *deliveryReceipt
}

func (ds *deliveryStore) Register(k string) chan *deliveryReceipt {
	c := make(chan *deliveryReceipt, 10)
	ds.mu.Lock()
	ds.m[k] = c
	ds.mu.Unlock()
	return c
}

func (ds *deliveryStore) Unregister(k string) {
	ds.mu.Lock()
	c := ds.m[k]
	if c != nil {
		delete(ds.m, k)
		close(c)
	}
	ds.mu.Unlock()
}

func (ds *deliveryStore) Broadcast(r *deliveryReceipt) {
	ds.mu.Lock()
	for _, c := range ds.m {
		select {
		case c <- r:
		default:
		}
	}
	ds.mu.Unlock()
}
