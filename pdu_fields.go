package smpp34

import "strconv"

type Field interface {
	Length() interface{}
	Value() interface{}
	String() string
}

type VariableField struct {
	value []byte
}

type FixedField struct {
	size  uint8
	value uint8
}

func (v *VariableField) Length() interface{} {
	l := len(v.value)
	return l
}

func (v *VariableField) Value() interface{} {
	return v.value
}

func (v *VariableField) String() string {
	return string(v.value)
}

func (f *FixedField) Length() interface{} {
	return uint8(1)
}

func (f *FixedField) Value() interface{} {
	return f.value
}

func (f *FixedField) String() string {
	return strconv.Itoa(int(f.value))
}
