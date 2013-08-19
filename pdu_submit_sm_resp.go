package smpp34

import (
	"bytes"
	"errors"
)

var (
	reqSSMRespFields = []string{"message_id"}
)

type SubmitSmResp struct {
	*Header
	mandatoryFields map[int]Field
}

func NewSubmitSmResp(hdr *Header, b []byte) (*SubmitSmResp, error) {
	r := bytes.NewBuffer(b)

	fields, err := create_pdu_fields(reqSSMRespFields, r)

	if err != nil {
		return nil, err
	}

	s := &SubmitSmResp{Header: hdr, mandatoryFields: fields}

	return s, nil
}

func (s *SubmitSmResp) GetField(f string) (Field, error) {
	for i, v := range s.MandatoryFieldsList() {
		if v == f {
			return s.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (s *SubmitSmResp) Fields() map[int]Field {
	return s.mandatoryFields
}

func (s *SubmitSmResp) MandatoryFieldsList() []string {
	return reqSSMRespFields
}
