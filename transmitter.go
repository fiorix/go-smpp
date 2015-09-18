package smpp34

import (
	"crypto/tls"
	"time"
)

type Transmitter struct {
	Smpp
	eLTicker     *time.Ticker // Enquire Link ticker
	eLCheckTimer *time.Timer  // Enquire Link Check timer
	eLDuration   int          // Enquire Link Duration
	Err          error        // Errors generated in go routines that lead to conn close
}

// NewTransmitter creates and initializes a new Transmitter.
// The eli parameter is for EnquireLink interval, in seconds.
func NewTransmitter(host string, port int, eli int, bindParams Params) (*Transmitter, error) {
	return newTransmitter(host, port, eli, bindParams, nil)
}

// NewTransmitterTLs creates and initializes a new Transmitter using TLS.
// The eli parameter is for EnquireLink interval, in seconds.
func NewTransmitterTLS(host string, port int, eli int, bindParams Params, config *tls.Config) (*Transmitter, error) {
	if config == nil {
		config = &tls.Config{}
	}
	return newTransmitter(host, port, eli, bindParams, config)
}

// eli = EnquireLink Interval in Seconds
func newTransmitter(host string, port int, eli int, bindParams Params, config *tls.Config) (*Transmitter, error) {
	tx := &Transmitter{}
	var err error
	if config == nil {
		err = tx.Connect(host, port)
	} else {
		err = tx.ConnectTLS(host, port, config)
	}
	if err != nil {
		return nil, err
	}

	sysId := bindParams[SYSTEM_ID].(string)
	pass := bindParams[PASSWORD].(string)

	if err := tx.Bind(sysId, pass, &bindParams); err != nil {
		return nil, err
	}

	// EnquireLinks should not be less 10seconds
	if eli < 10 {
		eli = 10
	}

	tx.eLDuration = eli

	go tx.startEnquireLink(eli)

	return tx, nil
}

func (t *Transmitter) Bind(system_id string, password string, params *Params) error {
	pdu, err := t.Smpp.Bind(BIND_TRANSMITTER, system_id, password, params)
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

	if pdu.GetHeader().Id != BIND_TRANSMITTER_RESP {
		return SmppBindRespErr
	}

	if !pdu.Ok() {
		return SmppBindAuthErr("Bind auth failed. " + pdu.GetHeader().Status.Error())
	}

	t.Bound = true

	return nil
}

func (t *Transmitter) SubmitSm(source_addr, destination_addr, short_message string, params *Params) (seq uint32, err error) {
	p, err := t.Smpp.SubmitSm(source_addr, destination_addr, short_message, params)

	if err != nil {
		return 0, err
	}

	if err := t.Write(p); err != nil {
		return 0, err
	}

	return p.GetHeader().Sequence, nil
}

func (t *Transmitter) QuerySm(message_id, source_addr string, params *Params) (seq uint32, err error) {
	p, err := t.Smpp.QuerySm(message_id, source_addr, params)

	if err != nil {
		return 0, err
	}

	if err := t.Write(p); err != nil {
		return 0, err
	}

	return p.GetHeader().Sequence, nil
}

func (t *Transmitter) DeliverSmResp(seq, status uint32) error {
	return SmppPduErr
}

func (t *Transmitter) Unbind() error {
	p, _ := t.Smpp.Unbind()

	if err := t.Write(p); err != nil {
		return err
	}

	return nil
}

func (t *Transmitter) UnbindResp(seq uint32) error {
	p, _ := t.Smpp.UnbindResp(seq)

	if err := t.Write(p); err != nil {
		return err
	}

	t.Bound = false

	return nil
}

func (t *Transmitter) bindCheck() {
	// Block
	<-time.After(time.Duration(5 * time.Second))
	if !t.Bound {
		// send error to t.err? So it can be read before closing.
		t.Err = SmppBindRespErr
		t.Close()
	}
}

func (t *Transmitter) startEnquireLink(eli int) {
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

func (t *Transmitter) Read() (Pdu, error) {
	pdu, err := t.Smpp.Read()
	if err != nil {
		if _, ok := err.(PduCmdIdErr); ok {
			// Invalid PDU Command ID, should send back GenericNack
			t.GenericNack(uint32(0), ESME_RINVCMDID)
		} else if SmppPduLenErr == err {
			// Invalid PDU, PDU read or Len error
			t.GenericNack(uint32(0), ESME_RINVCMDLEN)
		}

		return nil, err
	}

	switch pdu.GetHeader().Id {
	case SUBMIT_SM_RESP:
		return pdu, nil
	case QUERY_SM_RESP:
		return pdu, nil
	case ENQUIRE_LINK:
		p, _ := t.Smpp.EnquireLinkResp(pdu.GetHeader().Sequence)

		if err := t.Write(p); err != nil {
			return nil, err
		}
	case ENQUIRE_LINK_RESP:
		// Stop EnquireLink Check.
		t.eLCheckTimer.Stop()
	case UNBIND:
		t.UnbindResp(pdu.GetHeader().Sequence)
		t.Close()
	default:
		// Should not have received these PDUs on a TX bind
		return nil, SmppPduErr
	}

	return pdu, nil
}

func (t *Transmitter) Close() {
	// Check timers exists incase we Close() before timers are created
	if t.eLCheckTimer != nil {
		t.eLCheckTimer.Stop()
	}

	if t.eLTicker != nil {
		t.eLTicker.Stop()
	}

	t.Smpp.Close()
}

func (t *Transmitter) Write(p Pdu) error {
	err := t.Smpp.Write(p)

	return err
}
