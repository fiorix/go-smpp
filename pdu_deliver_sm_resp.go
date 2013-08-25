package smpp34

import (
	"bytes"
	"errors"
)

var (
	reqDSMRespFields = []string{MESSAGE_ID}
)

type DeliverSmResp struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       []*TLVField
}

func NewDeliverSmResp(hdr *Header, b []byte) (*DeliverSmResp, error) {
	r := bytes.NewBuffer(b)

	fields, _, err := create_pdu_fields(reqDSMRespFields, r)

	if err != nil {
		return nil, err
	}

	s := &DeliverSmResp{Header: hdr, mandatoryFields: fields}

	return s, nil
}

func (s *DeliverSmResp) GetField(f string) Field {
	return s.mandatoryFields[f]
}

func (s *DeliverSmResp) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *DeliverSmResp) MandatoryFieldsList() []string {
	return reqDSMRespFields
}

func (s *DeliverSmResp) GetHeader() *Header {
	return s.Header
}

func (s *DeliverSmResp) SetField(f string, v interface{}) error {
	if s.validate_field(f, v) {
		field := NewField(f, v)

		if field != nil {
			s.mandatoryFields[f] = field

			return nil
		}
	}

	return errors.New("Invalid field value")
}

func (s *DeliverSmResp) validate_field(f string, v interface{}) bool {
	if included_check(s.MandatoryFieldsList(), f) && validate_pdu_field(f, v) {
		return true
	}
	return false
}

func (s *DeliverSmResp) TLVFields() []*TLVField {
	return s.tlvFields
}

func (s *DeliverSmResp) writeFields() []byte {
	b := []byte{}

	for _, i := range s.MandatoryFieldsList() {
		v := s.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (s *DeliverSmResp) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(DELIVER_SM_RESP)...)
	h = append(h, packUi32(s.Header.Status)...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
