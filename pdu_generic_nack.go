package smpp34

var (
	reqGNFields = []string{}
)

type GenericNack struct {
	*Header
	mandatoryFields map[string]Field
	tlvFields       map[uint16]*TLVField
}

func NewGenericNack(hdr *Header) (*GenericNack, error) {
	s := &GenericNack{Header: hdr}

	return s, nil
}

func (s *GenericNack) GetField(f string) Field {
	return nil
}

func (s *GenericNack) Fields() map[string]Field {
	return s.mandatoryFields
}

func (s *GenericNack) MandatoryFieldsList() []string {
	return reqGNFields
}

func (s *GenericNack) Ok() bool {
	return true
}

func (s *GenericNack) GetHeader() *Header {
	return s.Header
}

func (s *GenericNack) SetField(f string, v interface{}) error {
	return FieldValueErr
}

func (s *GenericNack) SetSeqNum(i uint32) {
	s.Header.Sequence = i
}

func (s *GenericNack) SetTLVField(t, l int, v []byte) error {
	return TLVFieldPduErr
}

func (s *GenericNack) TLVFields() map[uint16]*TLVField {
	return s.tlvFields
}

func (s *GenericNack) writeFields() []byte {
	return []byte{}
}

func (s *GenericNack) Writer() []byte {
	b := s.writeFields()
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(uint32(GENERIC_NACK))...)
	h = append(h, packUi32(uint32(s.Header.Status))...)
	h = append(h, packUi32(s.Header.Sequence)...)
	return append(h, b...)
}
