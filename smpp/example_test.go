// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp_test

import (
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
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
	// Create persistent connection.
	conn := r.Bind()
	time.AfterFunc(10*time.Second, func() { r.Close() })
	// Print connection status (Connected, Disconnected, etc).
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
	// Create persistent connection, wait for the first status.
	conn := <-tx.Bind()
	if conn.Status() != smpp.Connected {
		log.Fatal(conn.Error())
	}
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
	lm := rate.NewLimiter(rate.Limit(10), 1) // Max rate of 10/s.
	tx := &smpp.Transceiver{
		Addr:        "localhost:2775",
		User:        "foobar",
		Passwd:      "secret",
		Handler:     f,  // Handle incoming SM or delivery receipts.
		RateLimiter: lm, // Optional rate limiter.
	}
	// Create persistent connection.
	conn := tx.Bind()
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
			Register: pdufield.FinalDeliveryReceipt,
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
