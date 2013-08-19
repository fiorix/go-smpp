package smpp34

import (
	"bytes"
	"encoding/binary"
	"errors"
	// "fmt"
	"io"
	"strconv"
)

type Pdu interface {
	Fields() map[int]Field
	MandatoryFieldsList() []string
	GetField(string) (Field, error)
}

func ParsePdu(data []byte) (Pdu, error) {
	if len(data) < 16 {
		return nil, errors.New("Invalid PDU. Length under 16 bytes")
	}

	header := &Header{
		unpackUi32(data[:4]),
		unpackUi32(data[4:8]),
		unpackUi32(data[8:12]),
		unpackUi32(data[12:16]),
	}

	switch header.Id {
	case SUBMIT_SM:
		n, err := NewSubmitSm(header, data[16:])

		if err != nil {
			return nil, err
		}

		return Pdu(n), nil
	case DELIVER_SM:
		n, err := NewDeliverSm(header, data[16:])

		if err != nil {
			return nil, err
		}

		return Pdu(n), nil
	default:
		return nil, errors.New("Unknown PDU Command ID: " + strconv.Itoa(int(header.Id)))
	}
}

func create_pdu_fields(fieldNames []string, r *bytes.Buffer) (map[int]Field, error) {

	fields := make(map[int]Field)
	var f Field
	for i, k := range fieldNames {
		switch k {
		case "service_type", "source_addr", "destination_addr", "schedule_delivery_time", "validity_period", "short_message":
			t, err := r.ReadBytes(0x00)

			if err == io.EOF {
				// continue
			} else if err != nil {
				return nil, err
			}

			v := &VariableField{t}
			f = v
			fields[i] = f
		case "source_addr_ton", "source_addr_npi", "dest_addr_ton", "dest_addr_npi", "esm_class", "protocol_id", "priority_flag", "registered_delivery", "replace_if_present_flag", "data_coding", "sm_default_msg_id", "sm_length":
			t, err := r.ReadByte()

			if err == io.EOF {
				// continue
			} else if err != nil {
				return nil, err
			}

			v := &FixedField{1, t}
			f = v
			fields[i] = f
		}
	}

	return fields, nil
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
