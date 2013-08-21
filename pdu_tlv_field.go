package smpp34

type TLVField struct {
	Tag    uint16
	Length uint16
	value  []byte
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
