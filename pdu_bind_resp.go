package smpp34

import (
	"bytes"
	"errors"
)

var (
	reqBindRespFields = []string{"system_id"}
)

type BindResp struct {
	*Header
	mandatoryFields map[int]Field
}

func NewBindResp(hdr *Header, b []byte) (*BindResp, error) {
	r := bytes.NewBuffer(b)

	fields, err := create_pdu_fields(reqBindRespFields, r)

	if err != nil {
		return nil, err
	}

	s := &BindResp{Header: hdr, mandatoryFields: fields}

	return s, nil
}

func (s *BindResp) GetField(f string) (Field, error) {
	for i, v := range s.MandatoryFieldsList() {
		if v == f {
			return s.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (s *BindResp) Fields() map[int]Field {
	return s.mandatoryFields
}

func (s *BindResp) MandatoryFieldsList() []string {
	return reqBindRespFields
}
