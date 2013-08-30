package smpp34

import (
	"errors"
)

var (
	reqUnbindFields = []string{}
)

type Unbind struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       map[uint16]*TLVField
}

func NewUnbind(hdr *Header) (*Unbind, error) {
	s := &Unbind{Header: hdr}

	return s, nil
}

func (s *Unbind) GetField(f string) Field {
	return s.mandatoryFields[f]
}

func (s *Unbind) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *Unbind) MandatoryFieldsList() []string {
	return reqUnbindFields
}

func (s *Unbind) Ok() bool {
	return true
}

func (s *Unbind) GetHeader() *Header {
	return s.Header
}

func (s *Unbind) SetField(f string, v interface{}) error {
	if s.validate_field(f, v) {
		field := NewField(f, v)

		if field != nil {
			s.mandatoryFields[f] = field

			return nil
		}
	}

	return errors.New("Invalid field value")
}

func (s *Unbind) SetSeqNum(i uint32) {
	s.Header.Sequence = i
}

func (s *Unbind) SetTLVField(t, l int, v []byte) error {
	return errors.New("Invalid TLV value lenght")
}

func (s *Unbind) validate_field(f string, v interface{}) bool {
	if included_check(s.MandatoryFieldsList(), f) && validate_pdu_field(f, v) {
		return true
	}
	return false
}

func (s *Unbind) TLVFields() map[uint16]*TLVField {
	return s.tlvFields
}

func (s *Unbind) writeFields() []byte {
	b := []byte{}

	for _, i := range s.MandatoryFieldsList() {
		v := s.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (s *Unbind) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(uint32(UNBIND))...)
	h = append(h, packUi32(uint32(s.Header.Status))...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
