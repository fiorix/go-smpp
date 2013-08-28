package smpp34

import (
	"bytes"
	"errors"
)

var (
	// Required SubmitSm Fields
	reqSSMFields = []string{
		SERVICE_TYPE,
		SOURCE_ADDR_TON,
		SOURCE_ADDR_NPI,
		SOURCE_ADDR,
		DEST_ADDR_TON,
		DEST_ADDR_NPI,
		DESTINATION_ADDR,
		ESM_CLASS,
		PROTOCOL_ID,
		PRIORITY_FLAG,
		SCHEDULE_DELIVERY_TIME,
		VALIDITY_PERIOD,
		REGISTERED_DELIVERY,
		REPLACE_IF_PRESENT_FLAG,
		DATA_CODING,
		SM_DEFAULT_MSG_ID,
		SM_LENGTH,
		SHORT_MESSAGE,
	}
)

type SubmitSm struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       map[uint16]*TLVField
}

func NewSubmitSm(hdr *Header, b []byte) (*SubmitSm, error) {
	r := bytes.NewBuffer(b)

	fields, tlvs, err := create_pdu_fields(reqSSMFields, r)

	if err != nil {
		return nil, err
	}

	s := &SubmitSm{hdr, fields, tlvs}

	return s, nil
}

func (s *SubmitSm) GetField(f string) Field {
	return s.mandatoryFields[f]
}

func (s *SubmitSm) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *SubmitSm) MandatoryFieldsList() []string {
	return reqSSMFields
}

func (s *SubmitSm) GetHeader() *Header {
	return s.Header
}

func (s *SubmitSm) SetField(f string, v interface{}) error {
	if s.validate_field(f, v) {
		field := NewField(f, v)

		if field != nil {
			s.mandatoryFields[f] = field

			return nil
		}
	}

	return errors.New("Invalid field value")
}

func (s *SubmitSm) SetSeqNum(i uint32) {
	s.Header.Sequence = i
}

func (s *SubmitSm) SetTLVField(t, l int, v []byte) error {
	if l != len(v) {
		return errors.New("Invalid TLV value lenght")
	}

	s.tlvFields[uint16(t)] = &TLVField{uint16(t), uint16(l), v}

	return nil
}

func (s *SubmitSm) validate_field(f string, v interface{}) bool {
	if included_check(s.MandatoryFieldsList(), f) && validate_pdu_field(f, v) {
		return true
	}
	return false
}

func (s *SubmitSm) TLVFields() map[uint16]*TLVField {
	return s.tlvFields
}

func (s *SubmitSm) writeFields() []byte {
	b := []byte{}

	for _, i := range s.MandatoryFieldsList() {
		v := s.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (s *SubmitSm) writeTLVFields() []byte {
	b := []byte{}

	for _, v := range s.tlvFields {
		b = append(b, v.Writer()...)
	}

	return b
}

func (s *SubmitSm) Writer() []byte {
	// Set SM_LENGTH
	sm := len(s.GetField(SHORT_MESSAGE).ByteArray())
	s.SetField(SM_LENGTH, sm)

	b := append(s.writeFields(), s.writeTLVFields()...)
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(SUBMIT_SM)...)
	h = append(h, packUi32(s.Header.Status)...)
	h = append(h, packUi32(s.Header.Sequence)...)

	return append(h, b...)
}
