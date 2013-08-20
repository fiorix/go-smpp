package smpp34

import (
	"errors"
)

var (
	reqELFields = []string{}
)

type EnquireLink struct {
	*Header
	mandatoryFields map[int]Field
	tlvFields       []*TLVField
}

func NewEnquireLink(hdr *Header) (*EnquireLink, error) {
	s := &EnquireLink{Header: hdr}

	return s, nil
}

func (s *EnquireLink) GetField(f string) (Field, error) {
	for i, v := range s.MandatoryFieldsList() {
		if v == f {
			return s.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (s *EnquireLink) Fields() map[int]Field {
	return s.mandatoryFields
}

func (s *EnquireLink) MandatoryFieldsList() []string {
	return reqELFields
}

func (s *EnquireLink) GetHeader() *Header {
	return s.Header
}

func (s *EnquireLink) TLVFields() []*TLVField {
	return s.tlvFields
}
