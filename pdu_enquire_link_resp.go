package smpp34

import (
	"errors"
)

var (
	reqELRespFields = []string{}
)

type EnquireLinkResp struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       []*TLVField
}

func NewEnquireLinkResp(hdr *Header) (*EnquireLinkResp, error) {
	s := &EnquireLinkResp{Header: hdr}

	return s, nil
}

func (s *EnquireLinkResp) GetField(f string) Field {
	return nil
}

func (s *EnquireLinkResp) SetField(f string, v interface{}) error {
	return errors.New("Invalid field value")
}

func (s *EnquireLinkResp) Fields() map[string]Field {
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
	return []byte{}
}

func (s *EnquireLinkResp) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(ENQUIRE_LINK_RESP)...)
	h = append(h, packUi32(s.Header.Status)...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
