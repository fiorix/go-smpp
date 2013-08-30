package smpp34

import (
	"errors"
)

var (
	reqELFields = []string{}
)

type EnquireLink struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       map[uint16]*TLVField
}

func NewEnquireLink(hdr *Header) (*EnquireLink, error) {
	s := &EnquireLink{Header: hdr}

	return s, nil
}

func (s *EnquireLink) GetField(f string) Field {
	return nil
}

func (s *EnquireLink) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *EnquireLink) MandatoryFieldsList() []string {
	return reqELFields
}

func (s *EnquireLink) Ok() bool {
	return true
}

func (s *EnquireLink) GetHeader() *Header {
	return s.Header
}

func (s *EnquireLink) SetField(f string, v interface{}) error {
	return errors.New("Invalid field value")
}

func (s *EnquireLink) SetSeqNum(i uint32) {
	s.Header.Sequence = i
}

func (s *EnquireLink) SetTLVField(t, l int, v []byte) error {
	return errors.New("Invalid TLV value lenght")
}

func (s *EnquireLink) TLVFields() map[uint16]*TLVField {
	return s.tlvFields
}

func (s *EnquireLink) writeFields() []byte {
	return []byte{}
}

func (s *EnquireLink) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(uint32(ENQUIRE_LINK))...)
	h = append(h, packUi32(uint32(s.Header.Status))...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
