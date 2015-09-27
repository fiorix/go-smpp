// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"net/rpc"
	"strings"
	"testing"

	"github.com/fiorix/go-smpp/smpp"
)

func TestSM_Submit_BadRequest(t *testing.T) {
	tx := newTransceiver()
	defer tx.Close()
	<-tx.Bind()
	sm := NewSM(tx, rpc.NewServer())
	var resp ShortMessageResp
	err := sm.Submit(&ShortMessage{}, &resp)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "400 Bad Request") {
			t.Fatal(err)
		}
		return
	}
	t.Fatal("submit with no params is not supposed to work")
}

func TestSM_Submit_BadGateway(t *testing.T) {
	tx := newTransceiver()
	defer tx.Close()
	<-tx.Bind()
	sm := NewSM(tx, rpc.NewServer())
	req := &ShortMessage{
		Src:  "root", // causes failure
		Dst:  "root",
		Text: "gotcha",
	}
	var resp ShortMessageResp
	err := sm.Submit(req, &resp)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "502 Bad Gateway") {
			t.Fatal(err)
		}
		return
	}
	t.Fatal("submit with bad params is not supposed to work")
}

func TestSM_Submit_ServiceUnavailable(t *testing.T) {
	tx := smpp.Transceiver{Addr: ":0"}
	defer tx.Close()
	<-tx.Bind()
	sm := NewSM(&tx, rpc.NewServer())
	req := &ShortMessage{
		Dst:  "root",
		Text: "gotcha",
	}
	var resp ShortMessageResp
	err := sm.Submit(req, &resp)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "503 Service Unavailable") {
			t.Fatal(err)
		}
		return
	}
	t.Fatal("submit with no server is not supposed to work")
}

func TestSM_Submit_EncParams(t *testing.T) {
	tx := newTransceiver()
	defer tx.Close()
	<-tx.Bind()
	sm := NewSM(tx, rpc.NewServer())
	for _, enc := range []string{"latin1", "ucs2", "fail-me"} {
		req := &ShortMessage{
			Dst:  "root",
			Text: "gotcha",
			Enc:  enc,
		}
		var resp ShortMessageResp
		err := sm.Submit(req, &resp)
		if err != nil && enc != "fail-me" {
			t.Fatal(err)
		}
	}
}

func TestSM_Submit_RegisterParam(t *testing.T) {
	tx := newTransceiver()
	defer tx.Close()
	<-tx.Bind()
	sm := NewSM(tx, rpc.NewServer())
	for _, reg := range []string{"final", "failure", "fail-me"} {
		req := &ShortMessage{
			Dst:      "root",
			Text:     "gotcha",
			Register: reg,
		}
		var resp ShortMessageResp
		err := sm.Submit(req, &resp)
		if err != nil && reg != "fail-me" {
			t.Fatal(err)
		}
	}
}

func TestSM_Query_BadRequest(t *testing.T) {
	tx := newTransceiver()
	defer tx.Close()
	<-tx.Bind()
	sm := NewSM(tx, rpc.NewServer())
	var resp QueryMessageResp
	err := sm.Query(&QueryMessage{}, &resp)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "400 Bad Request") {
			t.Fatal(err)
		}
		return
	}
	t.Fatal("submit with no params is not supposed to work")
}

func TestSM_Query_BadGateway(t *testing.T) {
	tx := newTransceiver()
	defer tx.Close()
	<-tx.Bind()
	sm := NewSM(tx, rpc.NewServer())
	req := &QueryMessage{
		Src:       "root", // causes failure
		MessageID: "13",
	}
	var resp QueryMessageResp
	err := sm.Query(req, &resp)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "502 Bad Gateway") {
			t.Fatal(err)
		}
		return
	}
	t.Fatal("query with bad params is not supposed to work")
}

func TestSM_Query_ServiceUnavailable(t *testing.T) {
	tx := smpp.Transceiver{Addr: ":0"}
	defer tx.Close()
	<-tx.Bind()
	sm := NewSM(&tx, rpc.NewServer())
	req := &QueryMessage{
		MessageID: "13",
	}
	var resp QueryMessageResp
	err := sm.Query(req, &resp)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "503 Service Unavailable") {
			t.Fatal(err)
		}
		return
	}
	t.Fatal("query with no server is not supposed to work")
}
