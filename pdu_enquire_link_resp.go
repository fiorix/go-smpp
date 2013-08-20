package smpp34

import (
	"errors"
)

var (
	reqELRespFields = []string{}
)

type EnquireLinkResp struct {
	*Header
	mandatoryFields map[int]Field
	tlvFields       []*TLVField
}

func NewEnquireLinkResp(hdr *Header) (*EnquireLinkResp, error) {
	s := &EnquireLinkResp{Header: hdr}

	return s, nil
}

func (s *EnquireLinkResp) GetField(f string) (Field, error) {
	for i, v := range s.MandatoryFieldsList() {
		if v == f {
			return s.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (s *EnquireLinkResp) Fields() map[int]Field {
	return s.mandatoryFields
}

func (s *EnquireLinkResp) MandatoryFieldsList() []string {
	return reqELRespFields
}

func (s *EnquireLinkResp) GetHeader() *Header {
	return s.Header
}

func (s *EnquireLinkResp) TLVFields() []*TLVField {
	return s.tlvFields
}
