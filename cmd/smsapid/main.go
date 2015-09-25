// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// HTTP API for sending SMS via an SMSC.
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	_ "net/http/pprof"

	"github.com/gorilla/handlers"

	"github.com/fiorix/go-smpp/cmd/smsapid/internal"
	"github.com/fiorix/go-smpp/smpp"
)

func main() {
	laddr := flag.String("http", ":8080", "host:port to listen on")
	logreq := flag.Bool("log", false, "log http requests to stderr")
	certf := flag.String("cert", "", "ssl certificate file for http server, optional")
	keyf := flag.String("key", "", "ssl key file for http server, optional")
	prefix := flag.String("prefix", "/", "prefix for http endpoints")
	cliaddr := flag.String("smpp", "localhost:2775", "host:port of the smsc to connect to via smpp 3.4")
	clitls := flag.Bool("tls", false, "connect to smsc using tls")
	cliprecaire := flag.Bool("precaire", false, "accept invalid ssl certificate from smsc")
	flag.Usage = func() {
		fmt.Printf("Usage: [env] %s [options]\n", os.Args[0])
		fmt.Printf("Environment variables:\n")
		fmt.Printf(" SMPP_USER: username for smpp client connection\n")
		fmt.Printf(" SMPP_PASSWD: password for smpp client connection\n")
		fmt.Printf("Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	tx := &smpp.Transceiver{
		Addr:   *cliaddr,
		User:   os.Getenv("SMPP_USER"),
		Passwd: os.Getenv("SMPP_PASSWD"),
	}
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)
	go func() {
		<-exit
		tx.Close()
		os.Exit(0)
	}()
	if *clitls {
		tx.TLS = &tls.Config{}
		if *cliprecaire {
			tx.TLS.InsecureSkipVerify = true
		}
	}
	api := &internal.Handler{
		Prefix: *prefix,
		Tx:     tx,
	}
	conn := api.Register(http.DefaultServeMux)
	go func() {
		for c := range conn {
			log.Println("SMPP connection status:", c.Status())
		}
	}()
	mux := http.Handler(http.DefaultServeMux)
	if *logreq {
		mux = handlers.LoggingHandler(os.Stderr, mux)
	}
	if *certf == "" || *keyf == "" {
		log.Fatal(http.ListenAndServe(*laddr, mux))
	}
	log.Fatal(http.ListenAndServeTLS(*laddr, *certf, *keyf, mux))
}
