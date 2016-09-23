// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp_test

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/veoo/go-smpp/smpp"
	"github.com/veoo/go-smpp/smpp/pdu"
	"github.com/veoo/go-smpp/smpp/pdu/pdufield"
	"github.com/veoo/go-smpp/smpp/pdu/pdutext"
)

func ExampleReceiver() {
	f := func(p pdu.Body) {
		switch p.Header().ID {
		case pdu.DeliverSMID:
			f := p.Fields()
			src := f[pdufield.SourceAddr]
			dst := f[pdufield.DestinationAddr]
			txt := f[pdufield.ShortMessage]
			log.Printf("Short message from=%q to=%q: %q",
				src, dst, txt)
		}
	}
	r := &smpp.Receiver{
		Addr:    "localhost:2775",
		User:    "foobar",
		Passwd:  "secret",
		Handler: f,
	}
	conn := r.Bind() // make persistent connection.
	time.AfterFunc(10*time.Second, func() { r.Close() })
	for c := range conn {
		log.Println("SMPP connection status:", c.Status())
	}
}

func ExampleTransmitter() {
	tx := &smpp.Transmitter{
		Addr:   "localhost:2775",
		User:   "foobar",
		Passwd: "secret",
	}
	conn := <-tx.Bind() // make persistent connection.
	switch conn.Status() {
	case smpp.Connected:
		sm, err := tx.Submit(&smpp.ShortMessage{
			Src:      "sender",
			Dst:      "recipient",
			Text:     pdutext.Latin1("Olá mundo"),
			Register: smpp.NoDeliveryReceipt,
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Message ID:", sm.RespID())
	default:
		log.Fatal(conn.Error())
	}
}

func ExampleTransceiver() {
	f := func(p pdu.Body) {
		switch p.Header().ID {
		case pdu.DeliverSMID:
			f := p.Fields()
			src := f[pdufield.SourceAddr]
			dst := f[pdufield.DestinationAddr]
			txt := f[pdufield.ShortMessage]
			log.Printf("Short message from=%q to=%q: %q",
				src, dst, txt)
		}
	}
	tx := &smpp.Transceiver{
		Addr:    "localhost:2775",
		User:    "foobar",
		Passwd:  "secret",
		Handler: f,
	}
	conn := tx.Bind() // make persistent connection.
	go func() {
		for c := range conn {
			log.Println("SMPP connection status:", c.Status())
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sm, err := tx.Submit(&smpp.ShortMessage{
			Src:      r.FormValue("src"),
			Dst:      r.FormValue("dst"),
			Text:     pdutext.Raw(r.FormValue("text")),
			Register: smpp.FinalDeliveryReceipt,
		})
		if err == smpp.ErrNotConnected {
			http.Error(w, "Oops.", http.StatusServiceUnavailable)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		io.WriteString(w, sm.RespID())
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
