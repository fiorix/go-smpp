package smpp34

import (
	"bytes"
	"errors"
)

var (
	reqBindFields = []string{
		"system_id",
		"password",
		"system_type",
		"interface_version",
		"addr_ton",
		"addr_npi",
		"address_range",
	}
)

type Bind struct {
	*Header
	mandatoryFields map[int]Field
	tlvFields       []*TLVField
}

func NewBind(hdr *Header, b []byte) (*Bind, error) {
	r := bytes.NewBuffer(b)

	fields, _, err := create_pdu_fields(reqBindFields, r)

	if err != nil {
		return nil, err
	}

	s := &Bind{Header: hdr, mandatoryFields: fields}

	return s, nil
}

func (s *Bind) GetField(f string) (Field, error) {
	for i, v := range s.MandatoryFieldsList() {
		if v == f {
			return s.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (s *Bind) Fields() map[int]Field {
	return s.mandatoryFields
}

func (s *Bind) MandatoryFieldsList() []string {
	return reqBindFields
}

func (s *Bind) GetHeader() *Header {
	return s.Header
}

func (s *Bind) TLVFields() []*TLVField {
	return s.tlvFields
}
