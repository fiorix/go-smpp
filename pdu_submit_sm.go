package smpp34

import (
	"bytes"
	"errors"
)

var (
	// Required SubmitSm Fields
	reqSSMFields = []string{
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

type SubmitSm struct {
	*Header
	mandatoryFields map[int]Field
	tlvFields       []*TLVField
}

func NewSubmitSm(hdr *Header, b []byte) (*SubmitSm, error) {
	r := bytes.NewBuffer(b)

	fields, tlvs, err := create_pdu_fields(reqSSMFields, r)

	if err != nil {
		return nil, err
	}

	s := &SubmitSm{hdr, fields, tlvs}

	return s, nil
}

func (s *SubmitSm) GetField(f string) (Field, error) {
	for i, v := range s.MandatoryFieldsList() {
		if v == f {
			return s.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (s *SubmitSm) Fields() map[int]Field {
	return s.mandatoryFields
}

func (s *SubmitSm) MandatoryFieldsList() []string {
	return reqSSMFields
}

func (s *SubmitSm) GetHeader() *Header {
	return s.Header
}

func (s *SubmitSm) TLVFields() []*TLVField {
	return s.tlvFields
}
