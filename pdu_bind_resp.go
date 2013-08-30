package smpp34

import (
	"bytes"
	"errors"
)

var (
	reqBindRespFields = []string{SYSTEM_ID}
)

type BindResp struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       map[uint16]*TLVField
}

func NewBindResp(hdr *Header, b []byte) (*BindResp, error) {
	r := bytes.NewBuffer(b)

	fields, tlvs, err := create_pdu_fields(reqBindRespFields, r)

	if err != nil {
		return nil, err
	}

	s := &BindResp{hdr, fields, tlvs}

	return s, nil
}

func (s *BindResp) GetField(f string) Field {
	return s.mandatoryFields[f]
}

func (s *BindResp) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *BindResp) MandatoryFieldsList() []string {
	return reqBindRespFields
}

func (s *BindResp) Ok() bool {
	if s.Header.Status == ESME_ROK {
		return true
	}

	return false
}

func (s *BindResp) GetHeader() *Header {
	return s.Header
}

func (s *BindResp) SetField(f string, v interface{}) error {
	if s.validate_field(f, v) {
		field := NewField(f, v)

		if field != nil {
			s.mandatoryFields[f] = field

			return nil
		}
	}

	return errors.New("Invalid field value")
}

func (s *BindResp) SetSeqNum(i uint32) {
	s.Header.Sequence = i
}

func (s *BindResp) SetTLVField(t, l int, v []byte) error {
	if l != len(v) {
		return errors.New("Invalid TLV value lenght")
	}

	s.tlvFields[uint16(t)] = &TLVField{uint16(t), uint16(l), v}

	return nil
}

func (s *BindResp) validate_field(f string, v interface{}) bool {
	if included_check(s.MandatoryFieldsList(), f) && validate_pdu_field(f, v) {
		return true
	}
	return false
}

func (s *BindResp) TLVFields() map[uint16]*TLVField {
	return s.tlvFields
}

func (s *BindResp) writeFields() []byte {
	b := []byte{}

	for _, i := range s.MandatoryFieldsList() {
		v := s.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (s *BindResp) writeTLVFields() []byte {
	b := []byte{}

	for _, v := range s.tlvFields {
		b = append(b, v.Writer()...)
	}

	return b
}

func (s *BindResp) Writer() []byte {
	b := append(s.writeFields(), s.writeTLVFields()...)
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(uint32(s.Header.Id))...)
	h = append(h, packUi32(uint32(s.Header.Status))...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
