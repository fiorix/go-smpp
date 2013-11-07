package smpp34

import (
	"bytes"
)

var (
	// Required QuerySm Fields
	reqQSMFields = []string{
		MESSAGE_ID,
		SOURCE_ADDR_TON,
		SOURCE_ADDR_NPI,
		SOURCE_ADDR,
	}
)

type QuerySm struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       map[uint16]*TLVField
}

func NewQuerySm(hdr *Header, b []byte) (*QuerySm, error) {
	r := bytes.NewBuffer(b)

	fields, _, err := create_pdu_fields(reqQSMFields, r)

	if err != nil {
		return nil, err
	}

	s := &QuerySm{Header: hdr, mandatoryFields: fields}

	return s, nil
}

func (s *QuerySm) GetField(f string) Field {
	return s.mandatoryFields[f]
}

func (s *QuerySm) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *QuerySm) MandatoryFieldsList() []string {
	return reqQSMFields
}

func (s *QuerySm) Ok() bool {
	return true
}

func (s *QuerySm) GetHeader() *Header {
	return s.Header
}

func (s *QuerySm) SetField(f string, v interface{}) error {
	if s.validate_field(f, v) {
		field := NewField(f, v)

		if field != nil {
			s.mandatoryFields[f] = field

			return nil
		}
	}

	return FieldValueErr
}

func (s *QuerySm) SetSeqNum(i uint32) {
	s.Header.Sequence = i
}

func (s *QuerySm) SetTLVField(t, l int, v []byte) error {
	return TLVFieldPduErr
}

func (s *QuerySm) validate_field(f string, v interface{}) bool {
	if included_check(s.MandatoryFieldsList(), f) && validate_pdu_field(f, v) {
		return true
	}
	return false
}

func (s *QuerySm) TLVFields() map[uint16]*TLVField {
	return s.tlvFields
}

func (s *QuerySm) writeFields() []byte {
	b := []byte{}

	for _, i := range s.MandatoryFieldsList() {
		v := s.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (s *QuerySm) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(uint32(s.Header.Id))...)
	h = append(h, packUi32(uint32(s.Header.Status))...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
