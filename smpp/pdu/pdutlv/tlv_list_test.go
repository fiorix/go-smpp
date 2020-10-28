package pdutlv

import (
	"bytes"
	"testing"
)

func TestDecodeTLV(t *testing.T) {
	f := NewTLV(TagDestAddrSubunit, []byte("hello"))
	var b bytes.Buffer
	if err := f.SerializeTo(&b); err != nil {
		t.Fatalf("serialization failed: %s", err)
	}
	m, err := DecodeTLV(&b)
	if err != nil {
		t.Fatal(err)
	}
	f, ok := m[TagDestAddrSubunit]
	if !ok {
		t.Fatalf("missing %q key: %#v", TagDestAddrSubunit.Hex(), m)
	}
	v, ok := f.(*Field)
	if !ok {
		t.Fatalf("field is not type Field: %#v", f)
	}
	want := []byte("hello")
	if !bytes.Equal(want, v.Bytes()) {
		t.Fatalf("unexpected data: want %q, have %q", want, v)
	}
}

func TestDecodeTLV_Error(t *testing.T) {
	want := []byte("hello")
	b := bytes.NewBuffer([]byte{0x00, 0x05, 0x00, 0x08})
	b.Write(want)

	m, err := DecodeTLV(b)
	if err == nil {
		t.Fatalf("expected not enough data error to be raised")
	} else if m != nil {
		t.Fatalf("expected returned Map to be nil: %#v", m)
	}
}
