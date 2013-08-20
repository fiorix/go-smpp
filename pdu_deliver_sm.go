package smpp34

import (
	"bytes"
	"errors"
)

var (
	reqDSMFields = []string{
		"service_type",
		"source_addr_ton",
		"source_addr_npi",
		"source_addr",
		"dest_addr_ton",
		"dest_addr_npi",
		"destination_addr",
		"esm_class",
		"protocol_id",
		"priority_flag",
		"schedule_delivery_time",
		"validity_period",
		"registered_delivery",
		"replace_if_present_flag",
		"data_coding",
		"sm_default_msg_id",
		"sm_length",
		"short_message",
	}
)

type DeliverSm struct {
	*Header
	mandatoryFields map[int]Field
	tlvFields       []*TLVField
}

func NewDeliverSm(hdr *Header, b []byte) (*DeliverSm, error) {
	r := bytes.NewBuffer(b)

	fields, tlvs, err := create_pdu_fields(reqDSMFields, r)

	if err != nil {
		return nil, err
	}

	s := &DeliverSm{hdr, fields, tlvs}

	return s, nil
}

func (s *DeliverSm) GetField(f string) (Field, error) {
	for i, v := range s.MandatoryFieldsList() {
		if v == f {
			return s.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (s *DeliverSm) Fields() map[int]Field {
	return s.mandatoryFields
}

func (s *DeliverSm) MandatoryFieldsList() []string {
	return reqDSMFields
}

func (s *DeliverSm) GetHeader() *Header {
	return s.Header
}

func (s *DeliverSm) TLVFields() []*TLVField {
	return s.tlvFields
}
