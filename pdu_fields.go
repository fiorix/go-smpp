package smpp34

type Field interface {
	Length() interface{}
	Value() interface{}
	V2() []byte
}

type VariableField struct {
	value []byte
}

type FixedField struct {
	size  uint8
	value interface{}
}

func (v *VariableField) Length() interface{} {
	l := len(v.value)
	return l
}

func (v *VariableField) Value() interface{} {
	return v.value
}

func (v *VariableField) V2() []byte {
	return v.value
}

func (f *FixedField) Length() interface{} {
	return uint8(1)
}

func (f *FixedField) Value() interface{} {
	return f.value
}

func (f *FixedField) V2() []byte {
	return make([]byte, 0)
}
