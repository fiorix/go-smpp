// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/websocket"

	"github.com/fiorix/go-smpp/smpp"
)

var deliverErr = make(chan error, 1)

func (sm *SM) Deliver(args *DeliveryReceipt, resp *string) error {
	if args.Text != "delivery receipt here" {
		deliverErr <- errors.New("unexpected delivery receipt")
	}
	*resp = ""
	return nil
}

func TestWebSocket_Deliver(t *testing.T) {
	mux := http.NewServeMux()
	h := Handler{Tx: &smpp.Transceiver{Addr: ":0"}}
	<-h.Register(mux)
	s := httptest.NewServer(mux)
	defer s.Close()
	u := strings.Replace(s.URL, "http:", "ws:", 1)
	ws, err := websocket.Dial(u+"/v1/ws/jsonrpc/events", "", "http://localhost")
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()
	h.pool.Broadcast(&DeliveryReceipt{
		Src:  "bart",
		Dst:  "lisa",
		Text: "delivery receipt here",
	})
	srv := rpc.NewServer()
	NewSM(h.Tx, srv)
	go func() {
		deliverErr <- srv.ServeRequest(jsonrpc.NewServerCodec(ws))
	}()
	select {
	case err = <-deliverErr:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for delivery receipt")
	}
}
