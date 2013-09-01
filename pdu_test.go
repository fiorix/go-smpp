package smpp34

import (
	"encoding/hex"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func (s *MySuite) Test_BindPdu(c *C) {
	h := "000000240000000900000000000000016875676f0067676f6f687500434d540034000000"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x24, BIND_TRANSCEIVER, ESME_ROK, uint32(1)))
	c.Check(p.GetField(SYSTEM_ID).String(), Equals, "hugo")
	c.Check(p.GetField(PASSWORD).String(), Equals, "ggoohu")
	c.Check(p.GetField(SYSTEM_TYPE).String(), Equals, "CMT")
	c.Check(p.GetField(INTERFACE_VERSION).Value(), Equals, uint8(0x34))
	c.Check(p.Writer(), DeepEquals, data)

	// Change values
	p.SetField(SYSTEM_ID, "test1")
	c.Check(p.GetField(SYSTEM_ID).String(), Equals, "test1")
	c.Check(hex.EncodeToString(p.Writer()), DeepEquals, "0000002500000009000000000000000174657374310067676f6f687500434d540034000000")
}

func (s *MySuite) Test_BindRespPdu(c *C) {
	h := "0000001d80000009000000000000000474657374696e67000210000134"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x1d, BIND_TRANSCEIVER_RESP, ESME_ROK, uint32(4)))
	c.Check(p.GetField(SYSTEM_ID).String(), Equals, "testing")
	c.Check(p.TLVFields()[0x210].Value(), DeepEquals, []uint8{0x34})
	c.Check(p.Writer(), DeepEquals, data)
}

func (s *MySuite) Test_DeliverSmPdu(c *C) {
	h := "0000004d000000050000000052227280000001746573743200010174657374000000010000010000002338393261386563303634633064373639666134353366373762343a2074657374206d6f"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x4d, DELIVER_SM, ESME_ROK, uint32(0x52227280)))
	c.Check(p.Writer(), DeepEquals, data)

	// Change Short Message
	p.SetField(SHORT_MESSAGE, "test1")
	c.Check(p.GetField(SHORT_MESSAGE).String(), Equals, "test1")
	c.Check(hex.EncodeToString(p.Writer()), DeepEquals, "0000002f00000005000000005222728000000174657374320001017465737400000001000001000000057465737431")
}

func (s *MySuite) Test_DeliverSmRespPdu(c *C) {
	h := "0000001180000005000000005222728000"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x11, DELIVER_SM_RESP, ESME_ROK, uint32(0x52227280)))
	c.Check(p.Writer(), DeepEquals, data)
}

func (s *MySuite) Test_EnquireLinkPdu(c *C) {
	h := "00000010000000150000000000000005"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x10, ENQUIRE_LINK, ESME_ROK, uint32(5)))
	c.Check(p.Writer(), DeepEquals, data)
}

func (s *MySuite) Test_EnquireLinkRespPdu(c *C) {
	h := "00000010800000150000000000000005"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x10, ENQUIRE_LINK_RESP, ESME_ROK, uint32(5)))
	c.Check(p.Writer(), DeepEquals, data)
}

func (s *MySuite) Test_GenericNackPdu(c *C) {
	h := "00000010800000000000000200000000"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x10, GENERIC_NACK, ESME_RINVCMDLEN, uint32(0)))
	c.Check(p.Writer(), DeepEquals, data)
}

func (s *MySuite) Test_SubmitSmPdu(c *C) {
	h := "0000002d00000004000000000000000200000074657374000000746573743200000000000000000000036d7367"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x2d, SUBMIT_SM, ESME_ROK, uint32(2)))
	c.Check(p.Writer(), DeepEquals, data)

	// Change Short Message
	p.SetField(SHORT_MESSAGE, "test1")
	c.Check(p.GetField(SHORT_MESSAGE).String(), Equals, "test1")
	c.Check(hex.EncodeToString(p.Writer()), DeepEquals, "0000002f00000004000000000000000200000074657374000000746573743200000000000000000000057465737431")
}

func (s *MySuite) Test_SubmitSmRespPdu(c *C) {
	h := "0000003580000004000000005221ac3831303039343665342d356138662d343835642d386536342d65646639616133373761323200"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x35, SUBMIT_SM_RESP, ESME_ROK, uint32(0x5221ac38)))
	c.Check(p.Writer(), DeepEquals, data)
	c.Check(p.GetField(MESSAGE_ID).String(), Equals, "100946e4-5a8f-485d-8e64-edf9aa377a22")
}

func (s *MySuite) Test_UnbindPdu(c *C) {
	h := "00000010000000060000000000000003"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x10, UNBIND, ESME_ROK, uint32(3)))
}

func (s *MySuite) Test_UnbindRespPdu(c *C) {
	h := "00000010800000060000000000000003"
	data, _ := hex.DecodeString(h)
	p, err := ParsePdu(data)

	c.Check(err, IsNil)
	c.Check(p.GetHeader(), DeepEquals, NewPduHeader(0x10, UNBIND_RESP, ESME_ROK, uint32(3)))
}
