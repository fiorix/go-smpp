package smpp34

import (
	"bytes"
	"errors"
)

var (
	reqDSMFields = []string{
		"service_type",
		"source_addr_ton",
		"source_addr_npi",
		"source_addr",
		"dest_addr_ton",
		"dest_addr_npi",
		"destination_addr",
		"esm_class",
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

type DeliverSm struct {
	*Header
	mandatoryFields map[int]Field
	tlvFields       []*TLVField
}

func NewDeliverSm(hdr *Header, b []byte) (*DeliverSm, error) {
	r := bytes.NewBuffer(b)

	fields, tlvs, err := create_pdu_fields(reqDSMFields, r)

	if err != nil {
		return nil, err
	}

	d := &DeliverSm{hdr, fields, tlvs}

	return d, nil
}

func (d *DeliverSm) GetField(f string) (Field, error) {
	for i, v := range d.MandatoryFieldsList() {
		if v == f {
			return d.mandatoryFields[i], nil
		}
	}

	return nil, errors.New("field not found")
}

func (d *DeliverSm) Fields() map[int]Field {
	return d.mandatoryFields
}

func (d *DeliverSm) MandatoryFieldsList() []string {
	return reqDSMFields
}

func (d *DeliverSm) GetHeader() *Header {
	return d.Header
}

func (d *DeliverSm) TLVFields() []*TLVField {
	return d.tlvFields
}

func (d *DeliverSm) writeFields() []byte {
	b := []byte{}

	for i, _ := range d.MandatoryFieldsList() {
		v := d.mandatoryFields[i].ByteArray()
		b = append(b, v...)
	}

	return b
}

func (d *DeliverSm) writeTLVFields() []byte {
	b := []byte{}

	for _, v := range d.tlvFields {
		b = append(b, v.Writer()...)
	}

	return b
}

func (d *DeliverSm) Writer() []byte {
	b := append(d.writeFields(), d.writeTLVFields()...)
	h := packUi32(uint32(len(b) + 16))
	h = append(h, packUi32(DELIVER_SM)...)
	h = append(h, packUi32(d.Header.Status)...)
	h = append(h, packUi32(d.Header.Sequence)...)

	return append(h, b...)
}
