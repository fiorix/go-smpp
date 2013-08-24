package smpp34

import (
	"bytes"
)

var (
	reqBindRespFields = []string{SYSTEM_ID}
)

type BindResp struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       []*TLVField
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

func (s *BindResp) GetHeader() *Header {
	return s.Header
}

func (s *BindResp) TLVFields() []*TLVField {
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
	h = append(h, packUi32(s.Header.Id)...)
	h = append(h, packUi32(s.Header.Status)...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
