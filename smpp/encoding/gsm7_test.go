package encoding

import (
	"encoding/hex"
	"fmt"
	"reflect"
	"testing"

	"golang.org/x/text/transform"
)

var validationStringTests = []struct {
	Text     string
	Expected []rune
}{
	{Text: "12345678", Expected: []rune{}},
	{Text: "12345[6]", Expected: []rune{}},
	{Text: "@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞ\f^{}\\[~]|€ÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà", Expected: []rune{}},
	{Text: "你", Expected: []rune{'你'}},
}

var validationBufferTests = []struct {
	Buffer   []byte
	Expected []byte
}{
	{Buffer: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38}, Expected: []byte{}},
	{Buffer: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x1B, 0x3C, 0x36, 0x1B, 0x3E}, Expected: []byte{}},
	{Buffer: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x1B}, Expected: []byte{0x1B}},
	{Buffer: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x1B, 0x00}, Expected: []byte{0x1B, 0x00}},
	{Buffer: []byte{0x80, 0x81, 0x82, 0x83}, Expected: []byte{0x80, 0x81, 0x82, 0x83}},
}

var packedTests = []struct {
	Text string
	Buff []byte
}{
	{Text: "", Buff: []byte{}},
	{Text: "1", Buff: []byte{0x31}},
	{Text: "12", Buff: []byte{0x31, 0x19}},
	{Text: "123", Buff: []byte{0x31, 0xD9, 0x0C}},
	{Text: "1234", Buff: []byte{0x31, 0xD9, 0x8C, 0x06}},
	{Text: "12345", Buff: []byte{0x31, 0xD9, 0x8C, 0x56, 0x03}},
	{Text: "123456", Buff: []byte{0x31, 0xD9, 0x8C, 0x56, 0xB3, 0x01}},
	{Text: "1234567", Buff: []byte{0x31, 0xD9, 0x8C, 0x56, 0xB3, 0xDD, 0x00}},
	{Text: "12345678", Buff: []byte{0x31, 0xD9, 0x8C, 0x56, 0xB3, 0xDD, 0x70}},
	{Text: "123456789", Buff: []byte{0x31, 0xD9, 0x8C, 0x56, 0xB3, 0xDD, 0x70, 0x39}},
	{Text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur nec nunc venenatis, ultricies ipsum id, volutpat ante. Sed pretium ac metus a interdum metus.", Buff: []byte("\xCC\xB7\xBC\xDC\x06\xA5\xE1\xF3\x7A\x1B\x44\x7E\xB3\xDF\x72\xD0\x3C\x4D\x07\x85\xDB\x65\x3A\x0B\x34\x7E\xBB\xE7\xE5\x31\xBD\x4C\xAF\xCB\x41\x61\x72\x1A\x9E\x9E\x8F\xD3\xEE\x33\xA8\xCC\x4E\xD3\x5D\xA0\x61\x5D\x1E\x16\xA7\xE9\x75\x39\xC8\x5D\x1E\x83\xDC\x75\xF7\x18\x64\x2F\xBB\xCB\xEE\x30\x3D\x3D\x67\x81\xEA\x6C\xBA\x3C\x3D\x4E\x97\xE7\xA0\x34\x7C\x5E\x6F\x83\xD2\x64\x16\xC8\xFE\x66\xD7\xE9\xF0\x30\x1D\x14\x76\xD3\xCB\x2E\xD0\xB4\x4C\x06\xC1\xE5\x65\x7A\xBA\xDE\x06\x85\xC7\xA0\x76\x99\x5E\x9F\x83\xC2\xA0\xB4\x9B\x5E\x96\x93\xEB\x6D\x50\xBB\x4C\xAF\xCF\x5D")},
	{Text: "\n", Buff: []byte{0x0A}},
	{Text: "\r", Buff: []byte{0x0D}},
	{Text: "\f", Buff: []byte{0x1B, 0x05}},
	{Text: "^{}\\[~]|€", Buff: []byte{0x1B, 0xCA, 0x06, 0xB5, 0x49, 0x6D, 0x5E, 0x1B, 0xDE, 0xA6, 0xB7, 0xF1, 0x6D, 0x80, 0x9B, 0x32}},
	{Text: "@£$¥èéùìòÇØøÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà", Buff: []byte("\x80\x80\x60\x40\x28\x18\x0E\x88\xC4\x82\xE1\x78\x40\x22\x92\x09\xA5\x62\xB9\x60\x32\x1A\x4E\xC7\xF3\x01\x85\x44\x23\x52\xC9\x74\x42\xA5\x54\x2B\x56\xCB\xF5\x82\xC5\x64\x33\x5A\xCD\x76\xC3\xE5\x74\x3B\x5E\xCF\xF7\x03\x06\x85\x43\x62\xD1\x78\x44\x26\x95\x4B\x66\xD3\xF9\x84\x46\xA5\x53\x6A\xD5\x7A\xC5\x66\xB5\x5B\x6E\xD7\xFB\x05\x87\xC5\x63\x72\xD9\x7C\x46\xA7\xD5\x6B\x76\xDB\xFD\x86\xC7\xE5\x73\x7A\xDD\x7E\xC7\xE7\xF5\x7B\x7E\xDF\xFF\x07")},
}

var unpackedTests = []struct {
	Text string
	Buff []byte
}{
	{Text: "", Buff: []byte{}},
	{Text: "1", Buff: []byte{0x31}},
	{Text: "12", Buff: []byte{0x31, 0x32}},
	{Text: "123", Buff: []byte{0x31, 0x32, 0x33}},
	{Text: "1234", Buff: []byte{0x31, 0x32, 0x33, 0x34}},
	{Text: "12345", Buff: []byte{0x31, 0x32, 0x33, 0x34, 0x35}},
	{Text: "123456", Buff: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36}},
	{Text: "1234567", Buff: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37}},
	{Text: "12345678", Buff: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38}},
	{Text: "123456789", Buff: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39}},
	{Text: "12345[6", Buff: []byte{0x31, 0x32, 0x33, 0x34, 0x35, 0x1B, 0x3C, 0x36}},
	{Text: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur nec nunc venenatis, ultricies ipsum id, volutpat ante. Sed pretium ac metus a interdum metus.", Buff: []byte("\x4C\x6F\x72\x65\x6D\x20\x69\x70\x73\x75\x6D\x20\x64\x6F\x6C\x6F\x72\x20\x73\x69\x74\x20\x61\x6D\x65\x74\x2C\x20\x63\x6F\x6E\x73\x65\x63\x74\x65\x74\x75\x72\x20\x61\x64\x69\x70\x69\x73\x63\x69\x6E\x67\x20\x65\x6C\x69\x74\x2E\x20\x43\x75\x72\x61\x62\x69\x74\x75\x72\x20\x6E\x65\x63\x20\x6E\x75\x6E\x63\x20\x76\x65\x6E\x65\x6E\x61\x74\x69\x73\x2C\x20\x75\x6C\x74\x72\x69\x63\x69\x65\x73\x20\x69\x70\x73\x75\x6D\x20\x69\x64\x2C\x20\x76\x6F\x6C\x75\x74\x70\x61\x74\x20\x61\x6E\x74\x65\x2E\x20\x53\x65\x64\x20\x70\x72\x65\x74\x69\x75\x6D\x20\x61\x63\x20\x6D\x65\x74\x75\x73\x20\x61\x20\x69\x6E\x74\x65\x72\x64\x75\x6D\x20\x6D\x65\x74\x75\x73\x2E")},
	{Text: "\n", Buff: []byte{0x0A}},
	{Text: "\r", Buff: []byte{0x0D}},
	{Text: "\f", Buff: []byte{0x1B, 0x0A}},
	{Text: "^{}\\[~]|€", Buff: []byte{0x1B, 0x14, 0x1B, 0x28, 0x1B, 0x29, 0x1B, 0x2F, 0x1B, 0x3C, 0x1B, 0x3D, 0x1B, 0x3E, 0x1B, 0x40, 0x1B, 0x65}},
	{Text: "@£$¥èéùìòÇØøÅåΔ_ΦΓΛΩΠΨΣΘΞÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà", Buff: []byte("\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0B\x0C\x0E\x0F\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1A\x1C\x1D\x1E\x1F\x20\x21\x22\x23\x24\x25\x26\x27\x28\x29\x2A\x2B\x2C\x2D\x2E\x2F\x30\x31\x32\x33\x34\x35\x36\x37\x38\x39\x3A\x3B\x3C\x3D\x3E\x3F\x40\x41\x42\x43\x44\x45\x46\x47\x48\x49\x4A\x4B\x4C\x4D\x4E\x4F\x50\x51\x52\x53\x54\x55\x56\x57\x58\x59\x5A\x5B\x5C\x5D\x5E\x5F\x60\x61\x62\x63\x64\x65\x66\x67\x68\x69\x6A\x6B\x6C\x6D\x6E\x6F\x70\x71\x72\x73\x74\x75\x76\x77\x78\x79\x7A\x7B\x7C\x7D\x7E\x7F")},
}

var invalidCharacterTests = []struct {
	Packed bool
	Text   string
}{
	{Packed: true, Text: "你"},
	{Packed: false, Text: "你"},
}

var invalidByteTests = []struct {
	Packed bool
	Buff   []byte
}{
	{Packed: false, Buff: []byte{0x80}},
	{Packed: false, Buff: []byte{0x1B}},
	{Packed: false, Buff: []byte{0x1B, 0x80}},
}

func TestGSM7EncodingString(t *testing.T) {
	tests := []struct {
		Packed   bool
		Expected string
	}{
		{Packed: true, Expected: "GSM 7-bit (Packed)"},
		{Packed: false, Expected: "GSM 7-bit (Unpacked)"},
	}

	for index, row := range tests {
		actual := fmt.Sprint(GSM7(row.Packed))
		if actual != row.Expected {
			t.Fatalf("%d: expected '%s' but got '%s'", index, row.Expected, actual)
		}
	}
}

func TestValidateGSM7String(t *testing.T) {
	for index, row := range validationStringTests {
		actual := ValidateGSM7String(row.Text)
		if !reflect.DeepEqual(actual, row.Expected) {
			t.Fatalf("%2d: actual did not equal expected.\nactual: %s\nexpect: %s", index, string(actual), string(row.Expected))
		}
	}
}

func TestValidateGSM7Buffer(t *testing.T) {
	for index, row := range validationBufferTests {
		actual := ValidateGSM7Buffer(row.Buffer)
		if !reflect.DeepEqual(actual, row.Expected) {
			t.Fatalf("%2d: actual did not equal expected.\nactual: %s\nexpect: %s", index, hex.EncodeToString(actual), hex.EncodeToString(row.Expected))
		}
	}
}

func TestPackedEncoder(t *testing.T) {
	encoder := GSM7(true).NewEncoder()
	for index, row := range packedTests {
		es, _, err := transform.Bytes(encoder, []byte(row.Text))
		if err != nil {
			t.Fatalf("%2d: unexpected error: '%s'", index, err.Error())
		}
		if !reflect.DeepEqual(es, row.Buff) {
			t.Fatalf("%2d: actual did not equal expected.\nactual: %s\nexpect: %s", index, hex.EncodeToString(es), hex.EncodeToString(row.Buff))
		}
	}
}

func TestUnpackedEncoder(t *testing.T) {
	encoder := GSM7(false).NewEncoder()
	for index, row := range unpackedTests {
		es, _, err := transform.Bytes(encoder, []byte(row.Text))
		if err != nil {
			t.Fatalf("%2d: unexpected error: '%s'", index, err.Error())
		}
		if !reflect.DeepEqual(es, row.Buff) {
			t.Fatalf("%2d: actual did not equal expected.\nactual: %s\nexpect: %s", index, hex.EncodeToString(es), hex.EncodeToString(row.Buff))
		}
	}
}

func TestPackedDecoder(t *testing.T) {
	decoder := GSM7(true).NewDecoder()
	for index, row := range packedTests {
		es, _, err := transform.Bytes(decoder, row.Buff)
		if err != nil {
			t.Fatalf("%2d: unexpected error: '%s'", index, err.Error())
		}
		if string(es) != row.Text {
			t.Fatalf("%2d: actual did not equal expected.\nactual: %s\nexpect: %s", index, string(es), row.Text)
		}
	}
}

func TestUnpackedDecoder(t *testing.T) {
	encoder := GSM7(false).NewDecoder()
	for index, row := range unpackedTests {
		es, _, err := transform.Bytes(encoder, row.Buff)
		if err != nil {
			t.Fatalf("%2d: unexpected error: '%s'", index, err.Error())
		}
		if string(es) != row.Text {
			t.Fatalf("%2d: actual did not equal expected.\nactual: %s\nexpect: %s", index, string(es), row.Text)
		}
	}
}

func TestInvalidCharacter(t *testing.T) {
	for index, row := range invalidCharacterTests {
		encoder := GSM7(row.Packed).NewEncoder()
		_, _, err := transform.Bytes(encoder, []byte(row.Text))
		if err == nil {
			t.Fatalf("%2d: expected error but got no error", index)
		}
		if err != ErrInvalidCharacter {
			t.Fatalf("%2d: expected '%s' but got '%s'", index, ErrInvalidCharacter, err.Error())
		}
	}
}

func TestInvalidByte(t *testing.T) {
	for index, row := range invalidByteTests {
		decoder := GSM7(row.Packed).NewDecoder()
		_, _, err := transform.Bytes(decoder, row.Buff)
		if err == nil {
			t.Fatalf("%2d: expected error but got no error", index)
		}
		if err != ErrInvalidByte {
			t.Fatalf("%2d: expected '%s' but got '%s'", index, ErrInvalidByte, err.Error())
		}
	}
}
