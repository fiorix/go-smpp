package smpp34

import (
	"errors"
	"fmt"
	"time"
)

type Transceiver struct {
	Smpp
	eLTicker     *time.Ticker // Enquire Link ticker
	eLCheckTimer *time.Timer  // Enquire Link Check timer
	eLDuration   int          // Enquire Link Duration
}

// eli = EnquireLink Interval in Seconds
func NewTransceiver(host string, port int, eli int, bindParams Params) (*Transceiver, error) {
	trx := &Transceiver{}
	if err := trx.Connect(host, port); err != nil {
		return nil, err
	}

	sysId := bindParams[SYSTEM_ID].(string)
	pass := bindParams[PASSWORD].(string)

	if err := trx.Bind(sysId, pass, &bindParams); err != nil {
		return nil, err
	}

	// EnquireLinks should not be less 10seconds
	if eli < 10 {
		eli = 10
	}

	trx.eLDuration = eli

	go trx.startEnquireLink(eli)

	return trx, nil
}

func (t *Transceiver) Bind(system_id string, password string, params *Params) error {
	pdu, err := t.Smpp.Bind(system_id, password, params)
	if err := t.Write(pdu); err != nil {
		return err
	}

	pdu, err = t.Smpp.Read()

	if err != nil {
		fmt.Println("pdu read err in bind:", err)
		return err
	}

	if pdu.GetHeader().Id != BIND_TRANSCEIVER_RESP {
		fmt.Println("TRX BIND Resp not received")
		return errors.New("TRX BIND Resp not received")
	}

	if !pdu.Ok() {
		return errors.New("Bind failed with status code" + string(pdu.GetHeader().Id))
	}

	t.Bound = true

	return nil
}

func (t *Transceiver) startEnquireLink(eli int) {
	t.eLTicker = time.NewTicker(time.Duration(eli) * time.Second)
	// check delay is half the time of enquire link intervel
	d := time.Duration(eli/2) * time.Second
	t.eLCheckTimer = time.NewTimer(d)
	t.eLCheckTimer.Stop()

	for {
		select {
		case <-t.eLTicker.C:

			p, _ := t.EnquireLink()
			if err := t.Write(p); err != nil {
				fmt.Println("Err writing ELR PDU:", err)
				t.Close()
				return
			}

			t.eLCheckTimer.Reset(d)
		case <-t.eLCheckTimer.C:
			fmt.Println("No enquire link response")
			t.Close()
			return
		}
	}
}

func (t *Transceiver) SubmitSm(source_addr, destination_addr, short_message string, params *Params) (seq uint32, err error) {
	p, err := t.Smpp.SubmitSm(source_addr, destination_addr, short_message, params)

	if err != nil {
		return 0, err
	}

	if err := t.Write(p); err != nil {
		return 0, err
	}

	return p.GetHeader().Sequence, nil
}

func (t *Transceiver) DeliverSmResp(seq, status uint32) error {
	p, err := t.Smpp.DeliverSmResp(seq, status)

	if err != nil {
		return err
	}

	if err := t.Write(p); err != nil {
		return err
	}

	return nil
}

func (t *Transceiver) Read() (Pdu, error) {
	pdu, err := t.Smpp.Read()
	if err != nil {
		return nil, err
	}

	switch pdu.GetHeader().Id {
	case SUBMIT_SM, SUBMIT_SM_RESP, DELIVER_SM_RESP, DELIVER_SM:
		return pdu, nil
	case ENQUIRE_LINK:
		p, _ := t.Smpp.EnquireLinkResp(pdu.GetHeader().Sequence)

		if err := t.Write(p); err != nil {
			return nil, err
		}
	case ENQUIRE_LINK_RESP:
		// Reset EnquireLink Check
		t.eLCheckTimer.Reset(time.Duration(t.eLDuration) * time.Second)
	}

	return pdu, nil
}

func (t *Transceiver) Close() {
	t.eLCheckTimer.Stop()
	t.eLTicker.Stop()

	t.Smpp.Close()
}

func (t *Transceiver) Write(p Pdu) error {
	err := t.Smpp.Write(p)

	return err
}
