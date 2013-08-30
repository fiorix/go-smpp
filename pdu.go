package smpp34

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

type PduReadErr string

type Pdu interface {
	Fields() map[string]Field
	MandatoryFieldsList() []string
	GetField(string) Field
	GetHeader() *Header
	TLVFields() map[uint16]*TLVField
	Writer() []byte
	SetField(f string, v interface{}) error
	SetTLVField(t, l int, v []byte) error
	SetSeqNum(uint32)
	Ok() bool
}

func (p PduReadErr) Error() string {
	return string(p)
}

func ParsePdu(data []byte) (Pdu, error) {
	if len(data) < 16 {
		return nil, PduReadErr("Invalid PDU. Length under 16 bytes")
	}

	header := ParsePduHeader(data[:16])

	switch header.Id {
	case SUBMIT_SM:
		n, err := NewSubmitSm(header, data[16:])
		return Pdu(n), err
	case SUBMIT_SM_RESP:
		n, err := NewSubmitSmResp(header, data[16:])
		return Pdu(n), err
	case DELIVER_SM:
		n, err := NewDeliverSm(header, data[16:])
		return Pdu(n), err
	case DELIVER_SM_RESP:
		n, err := NewDeliverSmResp(header, data[16:])
		return Pdu(n), err
	case BIND_TRANSCEIVER, BIND_RECEIVER, BIND_TRANSMITTER:
		n, err := NewBind(header, data[16:])
		return Pdu(n), err
	case BIND_TRANSCEIVER_RESP, BIND_RECEIVER_RESP, BIND_TRANSMITTER_RESP:
		n, err := NewBindResp(header, data[16:])
		return Pdu(n), err
	case ENQUIRE_LINK:
		n, err := NewEnquireLink(header)
		return Pdu(n), err
	case ENQUIRE_LINK_RESP:
		n, err := NewEnquireLinkResp(header)
		return Pdu(n), err
	case UNBIND:
		n, err := NewUnbind(header)
		return Pdu(n), err
	case UNBIND_RESP:
		n, err := NewUnbindResp(header)
		return Pdu(n), err
	default:
		return nil, PduReadErr(header.Id.Error())
	}
}

func ParsePduHeader(data []byte) *Header {
	return NewPduHeader(
		unpackUi32(data[:4]),
		CMDId(unpackUi32(data[4:8])),
		CMDStatus(unpackUi32(data[8:12])),
		unpackUi32(data[12:16]),
	)
}

func create_pdu_fields(fieldNames []string, r *bytes.Buffer) (map[string]Field, map[uint16]*TLVField, error) {

	fields := make(map[string]Field)
	eof := false
	for _, k := range fieldNames {
		switch k {
		case SERVICE_TYPE, SOURCE_ADDR, DESTINATION_ADDR, SCHEDULE_DELIVERY_TIME, VALIDITY_PERIOD, SYSTEM_ID, PASSWORD, SYSTEM_TYPE, ADDRESS_RANGE, MESSAGE_ID:
			t, err := r.ReadBytes(0x00)

			if err == io.EOF {
				eof = true
			} else if err != nil {
				return nil, nil, err
			}

			fields[k] = NewVariableField(t)
		case SOURCE_ADDR_TON, SOURCE_ADDR_NPI, DEST_ADDR_TON, DEST_ADDR_NPI, ESM_CLASS, PROTOCOL_ID, PRIORITY_FLAG, REGISTERED_DELIVERY, REPLACE_IF_PRESENT_FLAG, DATA_CODING, SM_DEFAULT_MSG_ID, INTERFACE_VERSION, ADDR_TON, ADDR_NPI:
			t, err := r.ReadByte()

			if err == io.EOF {
				eof = true
			} else if err != nil {
				return nil, nil, err
			}

			fields[k] = NewFixedField(t)
		case SM_LENGTH:
			// Short Message Length
			t, err := r.ReadByte()

			if err == io.EOF {
				eof = true
			} else if err != nil {
				return nil, nil, err
			}

			fields[k] = NewFixedField(t)

			// Short Message
			p := make([]byte, t)

			_, err = r.Read(p)
			if err == io.EOF {
				eof = true
			} else if err != nil {
				return nil, nil, err
			}

			fields[SHORT_MESSAGE] = NewVariableField(p)
		case SHORT_MESSAGE:
			continue
		}
	}

	// Optional Fields
	tlvs := map[uint16]*TLVField{}
	var err error

	if !eof {
		tlvs, err = parse_tlv_fields(r)

		if err != nil {
			return nil, nil, err
		}
	}

	return fields, tlvs, nil
}

func parse_tlv_fields(r *bytes.Buffer) (map[uint16]*TLVField, error) {
	tlvs := map[uint16]*TLVField{}

	for {
		p := make([]byte, 4)
		_, err := r.Read(p)

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		// length
		l := unpackUi16(p[2:4])

		// Get Value
		v := make([]byte, l)

		_, err = r.Read(v)
		if err != nil {
			return nil, err
		}

		tlvs[unpackUi16(p[0:2])] = &TLVField{
			unpackUi16(p[0:2]),
			unpackUi16(p[2:4]),
			v,
		}
	}

	return tlvs, nil
}

func validate_pdu_field(f string, v interface{}) bool {
	switch f {
	case SOURCE_ADDR_TON, SOURCE_ADDR_NPI, DEST_ADDR_TON, DEST_ADDR_NPI, ESM_CLASS, PROTOCOL_ID, PRIORITY_FLAG, REGISTERED_DELIVERY, REPLACE_IF_PRESENT_FLAG, DATA_CODING, SM_DEFAULT_MSG_ID, INTERFACE_VERSION, ADDR_TON, ADDR_NPI, SM_LENGTH:
		if validate_pdu_field_type(0x00, v) {
			return true
		}
	case SERVICE_TYPE, SOURCE_ADDR, DESTINATION_ADDR, SCHEDULE_DELIVERY_TIME, VALIDITY_PERIOD, SYSTEM_ID, PASSWORD, SYSTEM_TYPE, ADDRESS_RANGE, MESSAGE_ID, SHORT_MESSAGE:
		if validate_pdu_field_type("string", v) {
			return true
		}
	}
	return false
}

func validate_pdu_field_type(t interface{}, v interface{}) bool {
	if reflect.TypeOf(t) == reflect.TypeOf(v) {
		return true
	}

	return false
}

func included_check(a []string, v string) bool {
	for _, k := range a {
		if k == v {
			return true
		}
	}
	return false
}

func unpackUi32(b []byte) (n uint32) {
	n = binary.BigEndian.Uint32(b)
	return
}

func packUi32(n uint32) (b []byte) {
	b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return
}

func unpackUi16(b []byte) (n uint16) {
	n = binary.BigEndian.Uint16(b)
	return
}

func packUi16(n uint16) (b []byte) {
	b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return
}

func packUi8(n uint8) (b []byte) {
	b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, uint16(n))
	return b[1:]
}
