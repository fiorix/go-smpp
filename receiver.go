package smpp34

import (
	"time"
)

type Receiver struct {
	Smpp
	eLTicker     *time.Ticker // Enquire Link ticker
	eLCheckTimer *time.Timer  // Enquire Link Check timer
	eLDuration   int          // Enquire Link Duration
	Err          error        // Errors generated in go routines that lead to conn close
}

// eli = EnquireLink Interval in Seconds
func NewReceiver(host string, port int, eli int, bindParams Params) (*Receiver, error) {
	rx := &Receiver{}
	if err := rx.Connect(host, port); err != nil {
		return nil, err
	}

	sysId := bindParams[SYSTEM_ID].(string)
	pass := bindParams[PASSWORD].(string)

	if err := rx.Bind(sysId, pass, &bindParams); err != nil {
		return nil, err
	}

	// EnquireLinks should not be less 10seconds
	if eli < 10 {
		eli = 10
	}

	rx.eLDuration = eli

	go rx.startEnquireLink(eli)

	return rx, nil
}

func (t *Receiver) Bind(system_id string, password string, params *Params) error {
	pdu, err := t.Smpp.Bind(BIND_RECEIVER, system_id, password, params)
	if err := t.Write(pdu); err != nil {
		return err
	}

	// If BindResp NOT received in 5secs close connection
	go t.bindCheck()

	// Read (blocking)
	pdu, err = t.Smpp.Read()

	if err != nil {
		return err
	}

	if pdu.GetHeader().Id != BIND_RECEIVER_RESP {
		return SmppBindRespErr
	}

	if !pdu.Ok() {
		return SmppBindAuthErr("Bind auth failed. " + pdu.GetHeader().Status.Error())
	}

	t.Bound = true

	return nil
}

func (t *Receiver) SubmitSm(source_addr, destination_addr, short_message string, params *Params) (seq uint32, err error) {
	return 0, SmppPduErr
}

func (t *Receiver) DeliverSmResp(seq uint32, status CMDStatus) error {
	p, err := t.Smpp.DeliverSmResp(seq, status)

	if err != nil {
		return err
	}

	if err := t.Write(p); err != nil {
		return err
	}

	return nil
}

func (t *Receiver) Unbind() error {
	p, _ := t.Smpp.Unbind()

	if err := t.Write(p); err != nil {
		return err
	}

	return nil
}

func (t *Receiver) UnbindResp(seq uint32) error {
	p, _ := t.Smpp.UnbindResp(seq)

	if err := t.Write(p); err != nil {
		return err
	}

	t.Bound = false

	return nil
}

func (t *Receiver) bindCheck() {
	// Block
	<-time.After(time.Duration(5 * time.Second))
	if !t.Bound {
		// send error to t.err? So it can be read before closing.
		t.Err = SmppBindRespErr
		t.Close()
	}
}

func (t *Receiver) startEnquireLink(eli int) {
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
				t.Err = SmppELWriteErr
				t.Close()
				return
			}

			t.eLCheckTimer.Reset(d)
		case <-t.eLCheckTimer.C:
			t.Err = SmppELRespErr
			t.Close()
			return
		}
	}
}

func (t *Receiver) Read() (Pdu, error) {
	pdu, err := t.Smpp.Read()
	if err != nil {
		return nil, err
	}

	switch pdu.GetHeader().Id {
	case DELIVER_SM:
		return pdu, nil
	case ENQUIRE_LINK:
		p, _ := t.Smpp.EnquireLinkResp(pdu.GetHeader().Sequence)

		if err := t.Write(p); err != nil {
			return nil, err
		}
	case ENQUIRE_LINK_RESP:
		// Reset EnquireLink Check
		t.eLCheckTimer.Reset(time.Duration(t.eLDuration) * time.Second)
	case UNBIND:
		t.UnbindResp(pdu.GetHeader().Sequence)
		t.Close()
	default:
		// Should not have received these PDUs on a RX bind
		return nil, SmppPduErr
	}

	return pdu, nil
}

func (t *Receiver) Close() {
	// Check timers exists incase we Close() before timers are created
	if t.eLCheckTimer != nil {
		t.eLCheckTimer.Stop()
	}

	if t.eLTicker != nil {
		t.eLTicker.Stop()
	}

	t.Smpp.Close()
}

func (t *Receiver) Write(p Pdu) error {
	err := t.Smpp.Write(p)

	return err
}
