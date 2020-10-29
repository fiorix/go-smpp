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
	"net"
	"os"
	"strings"

	"github.com/tsocial/go-smpp/smpp"
	"github.com/tsocial/go-smpp/smpp/pdu/pdufield"
	"github.com/tsocial/go-smpp/smpp/pdu/pdutext"
	"github.com/urfave/cli"
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
	_ = app.Run(os.Args)
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
		cli.StringFlag{
			Name:  "service-type",
			Usage: "set service_type PDU (optional)",
			Value: "",
		},
		cli.IntFlag{
			Name:  "source-addr-ton",
			Usage: "set source_addr_ton PDU (optional)",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "source-addr-npi",
			Usage: "set source_addr_npi PDU (optional)",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "dest-addr-ton",
			Usage: "set dest_addr_ton PDU (optional)",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "dest-addr-npi",
			Usage: "set dest_addr_npi PDU (optional)",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "esm-class",
			Usage: "set esm_class PDU (optional)",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "protocol-id",
			Usage: "set protocol_id PDU (optional)",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "priority-flag",
			Usage: "set priority_flag PDU (optional)",
			Value: 0,
		},
		cli.StringFlag{
			Name:  "schedule-delivery-time",
			Usage: "set schedule_delivery_time PDU (optional)",
			Value: "",
		},
		cli.IntFlag{
			Name:  "replace-if-present-flag",
			Usage: "set replace_if_present_flag PDU (optional)",
			Value: 0,
		},
		cli.IntFlag{
			Name:  "sm-default-msg-id",
			Usage: "set sm_default_msg_id PDU (optional)",
			Value: 0,
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
		defer func() {
			_ = tx.Close()
		}()
		log.Println("Connected to", tx.Addr)
		sender := c.Args()[0]
		recipient := c.Args()[1]
		text := strings.Join(c.Args()[2:], " ")
		log.Printf("Command: send %q %q %q", sender, recipient, text)
		var register pdufield.DeliverySetting
		if c.Bool("register") {
			register = pdufield.FinalDeliveryReceipt
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
			Src:                  sender,
			Dst:                  recipient,
			Text:                 codec,
			Register:             register,
			ServiceType:          c.String("service-type"),
			SourceAddrTON:        uint8(c.Int("source-addr-ton")),
			SourceAddrNPI:        uint8(c.Int("source-addr-npi")),
			DestAddrTON:          uint8(c.Int("dest-addr-ton")),
			DestAddrNPI:          uint8(c.Int("dest-addr-npi")),
			ESMClass:             uint8(c.Int("esm-class")),
			ProtocolID:           uint8(c.Int("protocol-id")),
			PriorityFlag:         uint8(c.Int("priority-flag")),
			ScheduleDeliveryTime: c.String("schedule-delivery-time"),
			ReplaceIfPresentFlag: uint8(c.Int("replace-if-present-flag")),
			SMDefaultMsgID:       uint8(c.Int("sm-default-msg-id")),
		})
		if err != nil {
			log.Fatalln("Failed:", err)
		}
		log.Printf("Message ID: %q", sm.RespID())
	},
}

var cmdQueryMessage = cli.Command{
	Name:  "query",
	Usage: "status of short message",
	Action: func(c *cli.Context) {
		if len(c.Args()) != 2 {
			fmt.Println("usage: query [sender] [message ID]")
			return
		}
		log.Println("Connecting...")
		tx := newTransmitter(c)
		defer func() {
			_ = tx.Close()
		}()
		log.Println("Connected to", tx.Addr)
		sender, msgid := c.Args()[0], c.Args()[1]
		log.Printf("Command: query %q %q", sender, msgid)
		qr, err := tx.QuerySM(
			sender,
			msgid,
			uint8(c.Int("source-addr-ton")),
			uint8(c.Int("source-addr-npi")),
		)
		if err != nil {
			log.Fatalln("Failed:", err)
		}
		log.Printf("Status: %#v", *qr)
	},
}

func newTransmitter(c *cli.Context) *smpp.Transmitter {
	tx := &smpp.Transmitter{
		Addr:   c.GlobalString("addr"),
		User:   os.Getenv("SMPP_USER"),
		Passwd: os.Getenv("SMPP_PASSWD"),
	}
	if s := c.GlobalString("user"); s != "" {
		tx.User = s
	}
	if s := c.GlobalString("passwd"); s != "" {
		tx.Passwd = s
	}
	if c.GlobalBool("tls") {
		host, _, _ := net.SplitHostPort(tx.Addr)
		tx.TLS = &tls.Config{ //nolint:gosec
			ServerName: host,
		}
		if c.GlobalBool("precaire") {
			tx.TLS.InsecureSkipVerify = true
		}
	}
	conn := <-tx.Bind()
	if conn.Status() != smpp.Connected {
		log.Fatalln("Connection failed:", conn.Error())
	}
	return tx
}
