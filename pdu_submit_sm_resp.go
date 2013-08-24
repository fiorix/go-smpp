package smpp34

import (
	"bytes"
)

var (
	reqSSMRespFields = []string{MESSAGE_ID}
)

type SubmitSmResp struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       []*TLVField
}

func NewSubmitSmResp(hdr *Header, b []byte) (*SubmitSmResp, error) {
	r := bytes.NewBuffer(b)

	fields, _, err := create_pdu_fields(reqSSMRespFields, r)

	if err != nil {
		return nil, err
	}

	s := &SubmitSmResp{Header: hdr, mandatoryFields: fields}

	return s, nil
}

func (s *SubmitSmResp) GetField(f string) Field {
	return s.mandatoryFields[f]
}

func (s *SubmitSmResp) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *SubmitSmResp) MandatoryFieldsList() []string {
	return reqSSMRespFields
}

func (s *SubmitSmResp) GetHeader() *Header {
	return s.Header
}

func (s *SubmitSmResp) TLVFields() []*TLVField {
	return s.tlvFields
}

func (s *SubmitSmResp) writeFields() []byte {
	b := []byte{}

	for _, i := range s.MandatoryFieldsList() {
		v := s.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (s *SubmitSmResp) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(SUBMIT_SM_RESP)...)
	h = append(h, packUi32(s.Header.Status)...)
	h = append(h, packUi32(s.Header.Sequence)...)

	return append(h, b...)
}
