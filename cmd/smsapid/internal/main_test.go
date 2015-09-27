// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package internal

import (
	"os"
	"testing"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/smpptest"
)

var ts *smpptest.Server

func TestMain(m *testing.M) {
	ts = smpptest.NewUnstartedServer()
	ts.Handler = pduHandler
	ts.Start()
	defer ts.Close()
	os.Exit(m.Run())
}

func pduHandler(c smpptest.Conn, p pdu.Body) {
	src := p.Fields()[pdufield.SourceAddr]
	fail := src != nil && src.String() == "root"
	switch p.Header().ID {
	case pdu.SubmitSMID:
		r := pdu.NewSubmitSMResp()
		r.Header().Seq = p.Header().Seq
		if fail {
			r.Header().Status = 0x00000045 // submitsm failed
		} else {
			r.Fields().Set(pdufield.MessageID, "foobar")
		}
		c.Write(r)
		rd := p.Fields()[pdufield.RegisteredDelivery]
		if rd == nil || rd.Bytes()[0] == 0 {
			return
		}
		r = pdu.NewDeliverSM()
		rf := r.Fields()
		pf := p.Fields()
		rf.Set(pdufield.SourceAddr, pf[pdufield.SourceAddr])
		rf.Set(pdufield.DestinationAddr, pf[pdufield.DestinationAddr])
		rf.Set(pdufield.ShortMessage, "delivery receipt here")
		c.Write(r)
	case pdu.QuerySMID:
		r := pdu.NewQuerySMResp()
		r.Header().Seq = p.Header().Seq
		if fail {
			r.Header().Status = 0x00000067 // querysm failed
		} else {
			pf := p.Fields()
			rf := r.Fields()
			rf.Set(pdufield.MessageID, pf[pdufield.MessageID])
			rf.Set(pdufield.MessageState, 2) // DELIVERED
		}
		c.Write(r)
	default:
		smpptest.EchoHandler(c, p)
	}
}

func newTransceiver() *smpp.Transceiver {
	return &smpp.Transceiver{
		Addr:   ts.Addr(),
		User:   smpptest.DefaultUser,
		Passwd: smpptest.DefaultPasswd,
	}
}
