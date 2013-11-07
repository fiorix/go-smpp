package smpp34

import "strconv"

const FieldValueErr FieldErr = "Invalid field value"

type Field interface {
	Length() interface{}
	Value() interface{}
	String() string
	ByteArray() []byte
}

type FieldErr string

type SMField struct {
	value []byte
}

type VariableField struct {
	value []byte
}

type FixedField struct {
	size  uint8
	value uint8
}

func NewField(f string, v interface{}) Field {
	switch f {
	case SOURCE_ADDR_TON, SOURCE_ADDR_NPI, DEST_ADDR_TON, DEST_ADDR_NPI, ESM_CLASS, PROTOCOL_ID, PRIORITY_FLAG, REGISTERED_DELIVERY, REPLACE_IF_PRESENT_FLAG, DATA_CODING, SM_DEFAULT_MSG_ID, INTERFACE_VERSION, ADDR_TON, ADDR_NPI, SM_LENGTH, MESSAGE_STATE, ERROR_CODE:
		return NewFixedField(uint8(v.(int)))
	case SERVICE_TYPE, SOURCE_ADDR, DESTINATION_ADDR, SCHEDULE_DELIVERY_TIME, VALIDITY_PERIOD, SYSTEM_ID, PASSWORD, SYSTEM_TYPE, ADDRESS_RANGE, MESSAGE_ID, FINAL_DATE:
		return NewVariableField([]byte(v.(string)))
	case SHORT_MESSAGE:
		return NewSMField([]byte(v.(string)))
	}
	return nil
}

func NewSMField(v []byte) Field {
	i := &SMField{v}
	f := Field(i)
	return f
}

func NewVariableField(v []byte) Field {
	i := &VariableField{v}
	f := Field(i)
	return f
}

func NewFixedField(v uint8) Field {
	i := &FixedField{1, v}
	f := Field(i)
	return f
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

func (v *VariableField) ByteArray() []byte {
	return append(v.value, 0x00)
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

func (f *FixedField) ByteArray() []byte {
	return packUi8(f.value)
}

func (f FieldErr) Error() string {
	return string(f)
}

func (v *SMField) Length() interface{} {
	l := len(v.value)
	return l
}

func (v *SMField) Value() interface{} {
	return v.value
}

func (v *SMField) String() string {
	return string(v.value)
}

func (v *SMField) ByteArray() []byte {
	return v.value
}
