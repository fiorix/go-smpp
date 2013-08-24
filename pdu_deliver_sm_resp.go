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
	mandatoryFields map[int]Field
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

func (s *DeliverSmResp) GetField(f string) (Field, error) {
	for i, v := range s.MandatoryFieldsList() {
		if v == f {
			return s.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (s *DeliverSmResp) Fields() map[int]Field {
	return s.mandatoryFields
}

func (s *DeliverSmResp) MandatoryFieldsList() []string {
	return reqDSMRespFields
}

func (s *DeliverSmResp) GetHeader() *Header {
	return s.Header
}

func (s *DeliverSmResp) TLVFields() []*TLVField {
	return s.tlvFields
}

func (s *DeliverSmResp) writeFields() []byte {
	b := []byte{}

	for i, _ := range s.MandatoryFieldsList() {
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
