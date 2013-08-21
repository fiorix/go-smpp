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

func (s *EnquireLinkResp) writeFields() []byte {
	b := []byte{}

	for i, _ := range s.MandatoryFieldsList() {
		v := s.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (s *EnquireLinkResp) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(ENQUIRE_LINK_RESP)...)
	h = append(h, packUi32(s.Header.Status)...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
