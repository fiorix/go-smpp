package smpp34

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Smpp struct {
	mu       sync.Mutex
	conn     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer
	Sequence uint32
	Bound    bool
}

type SmppErr string
type SmppBindAuthErr string

type Params map[string]interface{}

const (
	SmppBindRespErr SmppErr = "BIND Resp not received"
	SmppPduErr      SmppErr = "PDU out of spec for this connection type"
	SmppPduSizeErr  SmppErr = "PDU Len larger than MAX_PDU_SIZE"
	SmppPduLenErr   SmppErr = "PDU Len different than read bytes"
	SmppELWriteErr  SmppErr = "Error writing ELR PDU"
	SmppELRespErr   SmppErr = "No enquire link response"
)

func (p SmppErr) Error() string {
	return string(p)
}

func (p SmppBindAuthErr) Error() string {
	return string(p)
}

func NewSmppConnect(host string, port int) (*Smpp, error) {
	s := &Smpp{}

	err := s.Connect(host, port)

	return s, err
}

func (s *Smpp) Connect(host string, port int) (err error) {
	s.conn, err = net.Dial("tcp", host+":"+strconv.Itoa(port))

	return err
}

func (s *Smpp) NewSeqNum() uint32 {
	defer s.mu.Unlock()

	s.mu.Lock()
	s.Sequence++
	return s.Sequence
}

func (s *Smpp) Bind(cmdId CMDId, system_id string, password string, params *Params) (Pdu, error) {
	b, _ := NewBind(
		&Header{Id: cmdId},
		[]byte{},
	)

	b.SetField(INTERFACE_VERSION, 0x34)
	b.SetField(SYSTEM_ID, system_id)
	b.SetField(PASSWORD, password)
	b.SetSeqNum(s.NewSeqNum())

	for f, v := range *params {
		err := b.SetField(f, v)

		if err != nil {
			return nil, err
		}
	}

	return Pdu(b), nil
}

func (s *Smpp) BindResp(cmdId CMDId, seq uint32, status CMDStatus, sysId string) (Pdu, error) {
	p, _ := NewBindResp(
		&Header{
			Id:       cmdId,
			Status:   status,
			Sequence: seq,
		},
		[]byte{},
	)

	p.SetField(SYSTEM_ID, sysId)
	p.SetTLVField(0x0210, 1, []byte{0x34}) // sc_interface_version TLV

	return Pdu(p), nil
}

func (s *Smpp) EnquireLink() (Pdu, error) {
	p, _ := NewEnquireLink(
		&Header{
			Id:       ENQUIRE_LINK,
			Sequence: s.NewSeqNum(),
		},
	)

	return Pdu(p), nil
}

func (s *Smpp) EnquireLinkResp(seq uint32) (Pdu, error) {
	p, _ := NewEnquireLinkResp(
		&Header{
			Id:       ENQUIRE_LINK_RESP,
			Status:   ESME_ROK,
			Sequence: seq,
		},
	)

	return Pdu(p), nil
}

func (s *Smpp) SubmitSm(source_addr, destination_addr, short_message string, params *Params) (Pdu, error) {

	p, _ := NewSubmitSm(
		&Header{
			Id:       SUBMIT_SM,
			Sequence: s.NewSeqNum(),
		},
		[]byte{},
	)

	p.SetField(SOURCE_ADDR, source_addr)
	p.SetField(DESTINATION_ADDR, destination_addr)
	p.SetField(SHORT_MESSAGE, short_message)

	for f, v := range *params {
		err := p.SetField(f, v)

		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func (s *Smpp) Unbind() (Pdu, error) {
	p, _ := NewUnbind(
		&Header{
			Id:       UNBIND,
			Sequence: s.NewSeqNum(),
		},
	)

	return Pdu(p), nil
}

func (s *Smpp) UnbindResp(seq uint32) (Pdu, error) {
	p, _ := NewUnbindResp(
		&Header{
			Id:       UNBIND_RESP,
			Sequence: seq,
		},
	)

	return Pdu(p), nil
}

func (s *Smpp) DeliverSmResp(seq uint32, status CMDStatus) (Pdu, error) {
	p, _ := NewDeliverSmResp(
		&Header{
			Id:       DELIVER_SM_RESP,
			Status:   status,
			Sequence: seq,
		},
		[]byte{},
	)

	return Pdu(p), nil
}

func (s *Smpp) Read() (Pdu, error) {
	l := make([]byte, 4)
	_, err := s.conn.Read(l)
	if err != nil {
		return nil, err
	}

	pduLength := unpackUi32(l) - 4
	if pduLength > MAX_PDU_SIZE {
		return nil, SmppPduSizeErr
	}

	data := make([]byte, pduLength)

	i, err := s.conn.Read(data)
	if err != nil {
		return nil, err
	}

	if i != int(pduLength) {
		return nil, SmppPduLenErr
	}

	pkt := append(l, data...)
	fmt.Println(hex.Dump(pkt))

	pdu, err := ParsePdu(pkt)
	if err != nil {
		return nil, err
	}

	return pdu, nil
}

func (s *Smpp) Write(p Pdu) error {
	_, err := s.conn.Write(p.Writer())

	fmt.Println(hex.Dump(p.Writer()))

	return err
}

func (s *Smpp) Close() {
	s.conn.Close()
}
