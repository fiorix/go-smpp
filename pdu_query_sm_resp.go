package smpp34

import (
	"bytes"
)

var (
	// Required QuerySmResp Fields
	reqQSMRespFields = []string{
		MESSAGE_ID,
		FINAL_DATE,
		MESSAGE_STATE,
		ERROR_CODE,
	}
)

type QuerySmResp struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       map[uint16]*TLVField
}

func NewQuerySmResp(hdr *Header, b []byte) (*QuerySmResp, error) {
	r := bytes.NewBuffer(b)

	fields, _, err := create_pdu_fields(reqQSMRespFields, r)

	if err != nil {
		return nil, err
	}

	s := &QuerySmResp{Header: hdr, mandatoryFields: fields}

	return s, nil
}

func (s *QuerySmResp) GetField(f string) Field {
	return s.mandatoryFields[f]
}

func (s *QuerySmResp) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *QuerySmResp) MandatoryFieldsList() []string {
	return reqQSMRespFields
}

func (s *QuerySmResp) Ok() bool {
	return true
}

func (s *QuerySmResp) GetHeader() *Header {
	return s.Header
}

func (s *QuerySmResp) SetField(f string, v interface{}) error {
	if s.validate_field(f, v) {
		field := NewField(f, v)

		if field != nil {
			s.mandatoryFields[f] = field

			return nil
		}
	}

	return FieldValueErr
}

func (s *QuerySmResp) SetSeqNum(i uint32) {
	s.Header.Sequence = i
}

func (s *QuerySmResp) SetTLVField(t, l int, v []byte) error {
	return TLVFieldPduErr
}

func (s *QuerySmResp) validate_field(f string, v interface{}) bool {
	if included_check(s.MandatoryFieldsList(), f) && validate_pdu_field(f, v) {
		return true
	}
	return false
}

func (s *QuerySmResp) TLVFields() map[uint16]*TLVField {
	return s.tlvFields
}

func (s *QuerySmResp) writeFields() []byte {
	b := []byte{}

	for _, i := range s.MandatoryFieldsList() {
		v := s.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (s *QuerySmResp) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(uint32(s.Header.Id))...)
	h = append(h, packUi32(uint32(s.Header.Status))...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
