package smpp34

const (
	TLVFieldLenErr TLVFieldErr = "Invalid TLV value lenght"
	TLVFieldPduErr TLVFieldErr = "PDU Type does not support TLV"
)

type TLVField struct {
	Tag    uint16
	Length uint16
	value  []byte
}

type TLVFieldErr string

func (t TLVFieldErr) Error() string {
	return string(t)
}

func (t *TLVField) Value() []byte {
	return t.value
}

func (t *TLVField) String() string {
	return string(t.Value())
}

func (t *TLVField) Writer() []byte {
	b := []byte{}
	b = append(b, packUi16(t.Tag)...)
	b = append(b, packUi16(t.Length)...)
	b = append(b, t.Value()...)
	return b
}
