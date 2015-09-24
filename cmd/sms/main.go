// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// SMPP client for the command line.
//
// We bind to the SMSC as a transmitter, therefore can do SubmitSM
// (send Short Message) or QuerySM (query for message status). The
// latter may not be available depending on the SMSC.
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

// Version of smppcli.
var Version = "tip"

// Author of smppcli.
var Author = "go-smpp authors"

func main() {
	app := cli.NewApp()
	app.Name = "smppcli"
	app.Usage = "SMPP client for the command line"
	app.Version = Version
	app.Author = Author
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr",
			Value: "localhost:2775",
			Usage: "Set SMPP server host:port",
		},
		cli.StringFlag{
			Name:  "user",
			Value: "",
			Usage: "Set SMPP username",
		},
		cli.StringFlag{
			Name:  "passwd",
			Value: "",
			Usage: "Set SMPP password",
		},
		cli.BoolFlag{
			Name:  "tls",
			Usage: "Use client TLS connection",
		},
		cli.BoolFlag{
			Name:  "precaire",
			Usage: "Accept invalid TLS certificate",
		},
	}
	app.Commands = []cli.Command{
		cmdShortMessage,
		cmdQueryMessage,
	}
	app.Run(os.Args)
}

var cmdShortMessage = cli.Command{
	Name:  "send",
	Usage: "send short message",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "register",
			Usage: "register for delivery receipt",
		},
		cli.StringFlag{
			Name:  "encoding",
			Usage: "set text encoding: raw, ucs2 or latin1",
			Value: "raw",
		},
	},
	Action: func(c *cli.Context) {
		if len(c.Args()) < 3 {
			fmt.Println("usage: send [options] <sender> <recipient> <message...>")
			fmt.Println("example: send --register foobar 011-236-0873 é nóis")
			return
		}
		log.Println("Connecting...")
		tx := newTransmitter(c)
		defer tx.Close()
		log.Println("Connected to", tx.Addr)
		sender := c.Args()[0]
		recipient := c.Args()[1]
		text := strings.Join(c.Args()[2:], " ")
		log.Printf("Command: send %q %q %q", sender, recipient, text)
		var register smpp.DeliverySetting
		if c.Bool("register") {
			register = smpp.FinalDeliveryReceipt
		}
		var codec pdutext.Codec
		switch c.String("encoding") {
		case "ucs2", "ucs-2":
			codec = pdutext.UCS2(text)
		case "latin1", "latin-1":
			codec = pdutext.Latin1(text)
		default:
			codec = pdutext.Raw(text)
		}
		sm, err := tx.Submit(&smpp.ShortMessage{
			Src:      sender,
			Dst:      recipient,
			Text:     codec,
			Register: register,
		})
		if err != nil {
			log.Println("Failed:", err)
		}
		log.Printf("Message ID: %q", sm.RespID())
	},
}

var cmdQueryMessage = cli.Command{
	Name:  "status",
	Usage: "status of short message",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 2 {
			fmt.Println("usage: status [sender] [message ID]")
			return
		}
		log.Println("Connecting...")
		tx := newTransmitter(c)
		defer tx.Close()
		log.Println("Connected to", tx.Addr)
		sender, msgid := c.Args()[0], c.Args()[1]
		log.Printf("Command: status %q %q", sender, msgid)
		qr, err := tx.QuerySM(sender, msgid)
		if err != nil {
			log.Fatalln("Failed:", err)
		}
		log.Printf("Status: %#v", *qr)
	},
}

func newTransmitter(c *cli.Context) *smpp.Transmitter {
	tx := &smpp.Transmitter{
		Addr:   c.GlobalString("addr"),
		User:   c.GlobalString("user"),
		Passwd: c.GlobalString("passwd"),
	}
	if c.GlobalBool("tls") {
		tx.TLS = &tls.Config{}
		if c.GlobalBool("precaire") {
			tx.TLS.InsecureSkipVerify = true
		}
	}
	conn := <-tx.Bind()
	switch conn.Status() {
	case smpp.Connected:
	default:
		log.Fatalln("Connection failed:", conn.Error())
	}
	return tx
}
