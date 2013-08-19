package smpp34

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

var (
	defaultMFields = []string{
		"service_type",
		"source_addr_ton",
		"source_addr_npi",
		"source_addr",
		"dest_addr_ton",
		"dest_addr_npi",
		"destination_addr",
		"ems_class",
		"protocol_id",
		"priority_flag",
		"schedule_delivery_time",
		"validity_period",
		"registered_delivery",
		"replace_if_present_flag",
		"data_coding",
		"sm_default_msg_id",
		"sm_length",
		"short_message",
	}
)

type SubmitSm struct {
	Header
	FieldsOrder []string
	Fields      map[int]Field
	SequenceNum uint32
}

func NewSubmitSm(b []byte) {
	r := bytes.NewBuffer(b)
	// var a Field
	// f := &FixedField{size: 1, value: 3}
	// v := &VariableField{[]byte("string")}
	// a = f

	// fmt.Println(a.Length(), a.Value())

	// a = v
	// fmt.Println(a.Length(), a.Value())

	p := make([]byte, 16)
	_, err := r.Read(p)
	if err != nil {
		return
	}

	// fmt.Println(p)

	fields := make(map[int]Field)
	var f Field
	for i, k := range defaultMFields {
		switch k {
		case "service_type", "source_addr", "destination_addr":
			t, _ := r.ReadBytes(0x00)
			// fmt.Println(t)
			v := &VariableField{t}
			f = v
			fields[i] = f
		case "source_addr_ton", "source_addr_npi", "dest_addr_ton", "dest_addr_npi":
			t, _ := r.ReadByte()
			v := &FixedField{1, t}
			f = v
			fields[i] = f
		}
	}

	s := &SubmitSm{Fields: fields}

	z := s.Fields[6].V2()
	i := hex.Dump(z)
	fmt.Println(i)
}
