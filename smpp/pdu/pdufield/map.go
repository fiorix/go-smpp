// Copyright 2015 go-smpp authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package pdufield

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/veoo/go-smpp/smpp/pdu/pdutext"
)

// Map is a collection of PDU field data indexed by name.
type Map map[Name]Body

// Set updates the PDU map with the given key and value, and
// returns error if the value cannot be converted to type Data.
//
// This is a shortcut for m[k] = New(k, v) converting v properly.
//
// If k is ShortMessage and v is of type pdutext.Codec, text is
// encoded and data_coding PDU and sm_length PDUs are set.
func (m Map) Set(k Name, v interface{}) error {
	switch v.(type) {
	case nil:
		m[k] = New(k, nil) // use default value
	case uint8:
		m[k] = New(k, []byte{v.(uint8)})
	case int:
		m[k] = New(k, []byte{uint8(v.(int))})
	case string:
		m[k] = New(k, []byte(v.(string)))
	case []byte:
		m[k] = New(k, []byte(v.([]byte)))
	case Body:
		m[k] = v.(Body)
	case pdutext.Codec:
		c := v.(pdutext.Codec)
		m[k] = New(k, c.Encode())
		if k == ShortMessage {
			m[DataCoding] = &Fixed{Data: uint8(c.Type())}
		}
	default:
		return fmt.Errorf("unsupported field data: %#v", v)
	}
	if k == ShortMessage {
		m[SMLength] = &Fixed{Data: uint8(m[k].Len())}
	}
	return nil
}

func (m Map) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	length := len(m)
	count := 0
	for k, v := range m {
		data := v.Raw()
		switch data.(type) {
		case []uint8:
			// Marshall the bytes as-is and also the string in two different
			// fields for readability
			jsonValue, _ := json.Marshal(v.String())
			buffer.WriteString(fmt.Sprintf("\"%v\":%s", k+"_text", jsonValue))
			buffer.WriteString(",")
			jsonValue, _ = json.Marshal(hex.EncodeToString(data.([]byte)))
			buffer.WriteString(fmt.Sprintf("\"%v\":%s", k, jsonValue))
		case uint8:
			jsonValue, _ := json.Marshal(data.(uint8))
			buffer.WriteString(fmt.Sprintf("\"%v\":%s", k, string(jsonValue)))
		default:
			jsonValue, _ := json.Marshal(v)
			buffer.WriteString(fmt.Sprintf("\"%v\":%s", k, string(jsonValue)))
		}

		count++
		if count < length {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

func (m *Map) UnmarshalJSON(b []byte) error {
	if *m == nil {
		*m = Map{}
	}
	tmp := m
	var d map[string]interface{}
	err := json.Unmarshal(b, &d)
	if err != nil {
		return err
	}
	for k, v := range d {
		// These were only put for readability
		if strings.Contains(k, "_text") {
			continue
		}
		var err error
		switch v.(type) {
		case string:
			s := v.(string)
			// Decode the string from hex
			bytes, err := hex.DecodeString(s)
			if err != nil {
				return err
			}
			err = tmp.Set(Name(k), bytes)
		case float64:
			err = tmp.Set(Name(k), uint8(v.(float64)))
		default:
			return fmt.Errorf("unsupported field type: %#v", v)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// TLVMap is a collection of PDU TLV field data indexed by type.
type TLVMap map[TLVType]*TLVBody

func (m TLVMap) Set(k TLVType, v *TLVBody) {
	m[k] = v
}

func (m TLVMap) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	length := len(m)
	count := 0
	for k, v := range m {
		jsonValue, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(fmt.Sprintf("\"%d\":%s", k, string(jsonValue)))
		count++
		if count < length {
			buffer.WriteString(",")
		}
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

func (m TLVMap) UnmarshalJSON(b []byte) error {
	var tmp map[string]*TLVBody
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	for k, v := range tmp {
		numericKey, err := strconv.Atoi(k)
		if err != nil {
			return err
		}
		m[TLVType(numericKey)] = v
	}
	return nil
}
