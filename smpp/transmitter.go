// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package smpp

import (
	"crypto/tls"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
)

// ErrMaxWindowSize is returned when an operation (such as Submit) violates
// the maximum window size configured for the Transmitter or Transceiver.
var ErrMaxWindowSize = errors.New("reached max window size")

// Transmitter implements an SMPP client transmitter.
type Transmitter struct {
	Addr        string
	User        string
	Passwd      string
	SystemType  string
	EnquireLink time.Duration
	RespTimeout time.Duration
	TLS         *tls.Config
	WindowSize  uint
	r           *rand.Rand

	conn struct {
		sync.Mutex
		*client
	}
	tx struct {
		count int32
		sync.Mutex
		inflight map[uint32]chan *tx
	}
}

type tx struct {
	PDU pdu.Body
	Err error
}

// Bind implements the ClientConn interface.
//
// Any commands (e.g. Submit) attempted on a dead connection will
// return ErrNotConnected.
func (t *Transmitter) Bind() <-chan ConnStatus {
	t.r = rand.New(rand.NewSource(time.Now().UnixNano()))
	t.conn.Lock()
	defer t.conn.Unlock()
	if t.conn.client != nil {
		return t.conn.Status
	}
	t.tx.Lock()
	t.tx.inflight = make(map[uint32]chan *tx)
	t.tx.Unlock()
	c := &client{
		Addr:        t.Addr,
		TLS:         t.TLS,
		EnquireLink: t.EnquireLink,
		RespTimeout: t.RespTimeout,
		Status:      make(chan ConnStatus, 1),
		BindFunc:    t.bindFunc,
		WindowSize:  t.WindowSize,
	}
	t.conn.client = c
	c.init()
	go c.Bind()
	return c.Status
}

func (t *Transmitter) bindFunc(c Conn) error {
	p := pdu.NewBindTransmitter()
	f := p.Fields()
	f.Set(pdufield.SystemID, t.User)
	f.Set(pdufield.Password, t.Passwd)
	f.Set(pdufield.SystemType, t.SystemType)
	resp, err := bind(c, p)
	if err != nil {
		return err
	}
	if resp.Header().ID != pdu.BindTransmitterRespID {
		return fmt.Errorf("unexpected response for BindTransmitter: %s",
			resp.Header().ID)
	}
	go t.handlePDU(nil)
	return nil
}

// f is only set on transceiver.
func (t *Transmitter) handlePDU(f HandlerFunc) {
	for {
		p, err := t.conn.Read()
		if err != nil {
			break
		}
		seq := p.Header().Seq
		t.tx.Lock()
		rc := t.tx.inflight[seq]
		t.tx.Unlock()
		if rc != nil {
			rc <- &tx{PDU: p}
		} else if f != nil {
			f(p)
		}
		if p.Header().ID == pdu.DeliverSMID { // Send DeliverSMResp
			pResp := pdu.NewDeliverSMRespSeq(p.Header().Seq)
			t.conn.Write(pResp)
		}
	}
	t.tx.Lock()
	for _, rc := range t.tx.inflight {
		rc <- &tx{Err: ErrNotConnected}
	}
	t.tx.Unlock()
}

// Close implements the ClientConn interface.
func (t *Transmitter) Close() error {
	t.conn.Lock()
	defer t.conn.Unlock()
	if t.conn.client == nil {
		return ErrNotConnected
	}
	return t.conn.Close()
}

// DeliverySetting is used to configure registered delivery
// for short messages.
type DeliverySetting uint8

// Supported delivery settings.
const (
	NoDeliveryReceipt      DeliverySetting = 0x00
	FinalDeliveryReceipt   DeliverySetting = 0x01
	FailureDeliveryReceipt DeliverySetting = 0x02
)

// ShortMessage configures a short message that can be submitted via
// the Transmitter. When returned from Submit, the ShortMessage
// provides Resp and RespID.
type ShortMessage struct {
	Src       string
	Dst       string
	Text      pdutext.Codec
	Validity  time.Duration
	Register  DeliverySetting
	OptParams pdufield.TLVMap

	// Other fields, normally optional.
	ServiceType          string
	SourceAddrTON        uint8
	SourceAddrNPI        uint8
	DestAddrTON          uint8
	DestAddrNPI          uint8
	ESMClass             uint8
	ProtocolID           uint8
	PriorityFlag         uint8
	ScheduleDeliveryTime string
	ReplaceIfPresentFlag uint8
	SMDefaultMsgID       uint8

	resp struct {
		sync.Mutex
		p pdu.Body
	}
}

// Resp returns the response PDU, or nil if not set.
func (sm *ShortMessage) Resp() pdu.Body {
	sm.resp.Lock()
	defer sm.resp.Unlock()
	return sm.resp.p
}

// RespID is a shortcut to Resp().Fields()[pdufield.MessageID].
// Returns empty if the response PDU is not available, or does
// not contain the MessageID field.
func (sm *ShortMessage) RespID() string {
	sm.resp.Lock()
	defer sm.resp.Unlock()
	if sm.resp.p == nil {
		return ""
	}
	f := sm.resp.p.Fields()[pdufield.MessageID]
	if f == nil {
		return ""
	}
	return f.String()
}

func (t *Transmitter) do(p pdu.Body) (*tx, error) {
	t.conn.Lock()
	notbound := t.conn.client == nil
	t.conn.Unlock()
	if notbound {
		return nil, ErrNotBound
	}
	if t.conn.WindowSize > 0 {
		inflight := uint(atomic.AddInt32(&t.tx.count, 1))
		defer func(t *Transmitter) { atomic.AddInt32(&t.tx.count, -1) }(t)
		if inflight > t.conn.WindowSize {
			return nil, ErrMaxWindowSize
		}
	}
	rc := make(chan *tx, 1)
	seq := p.Header().Seq
	t.tx.Lock()
	t.tx.inflight[seq] = rc
	t.tx.Unlock()
	defer func() {
		close(rc)
		t.tx.Lock()
		delete(t.tx.inflight, seq)
		t.tx.Unlock()
	}()
	err := t.conn.Write(p)
	if err != nil {
		return nil, err
	}
	select {
	case resp := <-rc:
		return resp, nil
	case <-t.conn.respTimeout():
		return nil, errors.New("timeout waiting for response")
	}
}

// Submit sends a short message and returns and updates the given
// sm with the response status. It returns the same sm object.
func (t *Transmitter) Submit(sm *ShortMessage) (*ShortMessage, error) {
	p := pdu.NewSubmitSM()
	return t.submitMsg(sm, p, uint8(sm.Text.Type()))
}

// SubmitLongMsg sends a long message (more than 140 bytes)
// and returns and updates the given sm with the response status.
// It returns the same sm object.
func (t *Transmitter) SubmitLongMsg(sm *ShortMessage) (*ShortMessage, error) {
	maxLen := 134 // 140-6 (UDH)
	rawMsg := sm.Text.Encode()
	countParts := int((len(rawMsg)-1)/maxLen) + 1

	ri := uint8(t.r.Intn(128))
	UDHHeader := make([]byte, 6)
	UDHHeader[0] = 5
	UDHHeader[1] = 0
	UDHHeader[2] = 3
	UDHHeader[3] = ri
	UDHHeader[4] = uint8(countParts)
	for i := 0; i < countParts; i++ {
		UDHHeader[5] = uint8(i + 1)
		p := pdu.NewSubmitSM()
		f := p.Fields()
		f.Set(pdufield.SourceAddr, sm.Src)
		f.Set(pdufield.DestinationAddr, sm.Dst)
		if i != countParts-1 {
			f.Set(pdufield.ShortMessage, pdutext.Raw(append(UDHHeader, rawMsg[i*maxLen:(i+1)*maxLen]...)))
		} else {
			f.Set(pdufield.ShortMessage, pdutext.Raw(append(UDHHeader, rawMsg[i*maxLen:]...)))
		}
		f.Set(pdufield.RegisteredDelivery, uint8(sm.Register))
		if sm.Validity != time.Duration(0) {
			f.Set(pdufield.ValidityPeriod, convertValidity(sm.Validity))
		}
		f.Set(pdufield.ServiceType, sm.ServiceType)
		f.Set(pdufield.SourceAddrTON, sm.SourceAddrTON)
		f.Set(pdufield.SourceAddrNPI, sm.SourceAddrNPI)
		f.Set(pdufield.DestAddrTON, sm.DestAddrTON)
		f.Set(pdufield.DestAddrNPI, sm.DestAddrNPI)
		f.Set(pdufield.ESMClass, 0x40)
		f.Set(pdufield.ProtocolID, sm.ProtocolID)
		f.Set(pdufield.PriorityFlag, sm.PriorityFlag)
		f.Set(pdufield.ScheduleDeliveryTime, sm.ScheduleDeliveryTime)
		f.Set(pdufield.ReplaceIfPresentFlag, sm.ReplaceIfPresentFlag)
		f.Set(pdufield.SMDefaultMsgID, sm.SMDefaultMsgID)
		f.Set(pdufield.DataCoding, uint8(sm.Text.Type()))
		resp, err := t.do(p)
		if err != nil {
			return nil, err
		}
		sm.resp.Lock()
		sm.resp.p = resp.PDU
		sm.resp.Unlock()
		if id := resp.PDU.Header().ID; id != pdu.SubmitSMRespID {
			return sm, fmt.Errorf("unexpected PDU ID: %s", id)
		}
		if s := resp.PDU.Header().Status; s != 0 {
			return sm, s
		}
		if resp.Err != nil {
			return sm, resp.Err
		}
	}
	return sm, nil
}

func (t *Transmitter) submitMsg(sm *ShortMessage, p pdu.Body, dataCoding uint8) (*ShortMessage, error) {
	f := p.Fields()
	f.Set(pdufield.SourceAddr, sm.Src)
	f.Set(pdufield.DestinationAddr, sm.Dst)
	f.Set(pdufield.ShortMessage, sm.Text)
	f.Set(pdufield.RegisteredDelivery, uint8(sm.Register))
	// Check if the message has validity set.
	if sm.Validity != time.Duration(0) {
		f.Set(pdufield.ValidityPeriod, convertValidity(sm.Validity))
	}
	f.Set(pdufield.ServiceType, sm.ServiceType)
	f.Set(pdufield.SourceAddrTON, sm.SourceAddrTON)
	f.Set(pdufield.SourceAddrNPI, sm.SourceAddrNPI)
	f.Set(pdufield.DestAddrTON, sm.DestAddrTON)
	f.Set(pdufield.DestAddrNPI, sm.DestAddrNPI)
	f.Set(pdufield.ESMClass, sm.ESMClass)
	f.Set(pdufield.ProtocolID, sm.ProtocolID)
	f.Set(pdufield.PriorityFlag, sm.PriorityFlag)
	f.Set(pdufield.ScheduleDeliveryTime, sm.ScheduleDeliveryTime)
	f.Set(pdufield.ReplaceIfPresentFlag, sm.ReplaceIfPresentFlag)
	f.Set(pdufield.SMDefaultMsgID, sm.SMDefaultMsgID)
	f.Set(pdufield.DataCoding, dataCoding)
	//set the optional parameters in the submit pdu from sm
	optParams := p.TLVFields()
	for param, value := range sm.OptParams {
		optParams.Set(param, value)
	}

	resp, err := t.do(p)
	if err != nil {
		return nil, err
	}
	sm.resp.Lock()
	sm.resp.p = resp.PDU
	sm.resp.Unlock()
	if id := resp.PDU.Header().ID; id != pdu.SubmitSMRespID {
		return sm, fmt.Errorf("unexpected PDU ID: %s", id)
	}
	if s := resp.PDU.Header().Status; s != 0 {
		return sm, s
	}
	return sm, resp.Err
}

// QueryResp contains the parsed the response of a QuerySM request.
type QueryResp struct {
	MsgID     string
	MsgState  string
	FinalDate string
	ErrCode   uint8
}

// QuerySM queries the delivery status of a message. It requires the
// source address (sender) and message ID.
func (t *Transmitter) QuerySM(src, msgid string) (*QueryResp, error) {
	p := pdu.NewQuerySM()
	f := p.Fields()
	f.Set(pdufield.SourceAddr, src)
	f.Set(pdufield.MessageID, msgid)
	resp, err := t.do(p)
	if err != nil {
		return nil, err
	}
	if id := resp.PDU.Header().ID; id != pdu.QuerySMRespID {
		return nil, fmt.Errorf("unexpected PDU ID: %s", id)
	}
	if s := resp.PDU.Header().Status; s != 0 {
		return nil, s
	}
	f = resp.PDU.Fields()
	ms := f[pdufield.MessageState]
	if ms == nil {
		return nil, fmt.Errorf("no state available")
	}
	qr := &QueryResp{MsgID: msgid}
	switch ms.Bytes()[0] {
	case 0:
		qr.MsgState = "DELIVERED"
	case 1:
		qr.MsgState = "ENROUTE"
	case 2:
		qr.MsgState = "DELIVERED"
	case 3:
		qr.MsgState = "EXPIRED"
	case 4:
		qr.MsgState = "DELETED"
	case 5:
		qr.MsgState = "UNDELIVERABLE"
	case 6:
		qr.MsgState = "ACCEPTED"
	case 7:
		qr.MsgState = "UNKNOWN"
	case 8:
		qr.MsgState = "REJECTED"
	case 9:
		qr.MsgState = "SKIPPED"
	default:
		qr.MsgState = fmt.Sprintf("UNKNOWN (%d)", ms.Bytes()[0])
	}
	if fd := f[pdufield.FinalDate]; fd != nil {
		qr.FinalDate = fd.String()
	}
	if ec := f[pdufield.ErrorCode]; ec != nil {
		qr.ErrCode = ec.Bytes()[0]
	}
	return qr, nil
}

func convertValidity(d time.Duration) string {
	validity := time.Now().UTC().Add(d)
	// Absolute time format YYMMDDhhmmsstnnp, see SMPP3.4 spec 7.1.1.
	return validity.Format("060102150405") + "000+"
}
