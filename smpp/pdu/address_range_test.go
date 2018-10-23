package pdu

import (
	"testing"

	"github.com/tsocial/go-smpp/smpp/pdu/pdufield"
)

func TestSetFields(t *testing.T) {
	t.Run("not_set", func(t *testing.T) {
		f := make(pdufield.Map)
		var a *AddressRange
		a.SetFields(f)

		if len(f) != 0 {
			t.Fatalf("expected: 0 items. actual: %d", len(f))
		}
	})

	t.Run("success", func(t *testing.T) {
		f := make(pdufield.Map)
		a := &AddressRange{
			TON:     5,
			NPI:     1,
			Address: "09999999",
		}
		a.SetFields(f)

		if len(f) != 3 {
			t.Fatalf("expected: 3 items. actual: %d", len(f))
		}

		address, exist := f[pdufield.AddressRange]
		if !exist {
			t.Fatalf("missing key %s", pdufield.AddressRange)
		}

		if a.Address != address.String() {
			t.Fatalf("expected: %s. actual: %s", a.Address, address)
		}

		npi, exist := f[pdufield.AddrNPI]
		if !exist {
			t.Fatalf("missing key %s", pdufield.AddrNPI)
		}

		if a.NPI != npi.Bytes()[0] {
			t.Fatalf("expected: %d. actual: %s", a.NPI, npi)
		}

		ton, exist := f[pdufield.AddrTON]
		if !exist {
			t.Fatalf("missing key %s", pdufield.AddrTON)
		}

		if a.TON != ton.Bytes()[0] {
			t.Fatalf("expected: %d. actual: %s", a.TON, ton)
		}
	})
}
