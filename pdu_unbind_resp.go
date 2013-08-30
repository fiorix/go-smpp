package smpp34

import (
	"errors"
)

var (
	reqUnbindRespFields = []string{}
)

type UnbindResp struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       map[uint16]*TLVField
}

func NewUnbindResp(hdr *Header) (*UnbindResp, error) {
	s := &UnbindResp{Header: hdr}

	return s, nil
}

func (s *UnbindResp) GetField(f string) Field {
	return nil
}

func (s *UnbindResp) SetField(f string, v interface{}) error {
	return errors.New("Invalid field value")
}

func (s *UnbindResp) SetSeqNum(i uint32) {
	s.Header.Sequence = i
}

func (s *UnbindResp) SetTLVField(t, l int, v []byte) error {
	return errors.New("Invalid TLV value lenght")
}

func (s *UnbindResp) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *UnbindResp) MandatoryFieldsList() []string {
	return reqUnbindRespFields
}

func (s *UnbindResp) Ok() bool {
	if s.Header.Status == ESME_ROK {
		return true
	}

	return false
}

func (s *UnbindResp) GetHeader() *Header {
	return s.Header
}

func (s *UnbindResp) TLVFields() map[uint16]*TLVField {
	return s.tlvFields
}

func (s *UnbindResp) writeFields() []byte {
	return []byte{}
}

func (s *UnbindResp) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(s.Header.Id)...)
	h = append(h, packUi32(s.Header.Status)...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
