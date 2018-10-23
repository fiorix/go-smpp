package pdu

import (
	"github.com/tsocial/go-smpp/smpp/pdu/pdufield"
)

// AddressRange is an ESME address served via the SMPP session.
type AddressRange struct {
	TON     uint8
	NPI     uint8
	Address string
}

// SetFields set address range to field map
func (a *AddressRange) SetFields(f pdufield.Map) {
	if a == nil {
		return
	}

	f.Set(pdufield.AddrNPI, a.NPI)
	f.Set(pdufield.AddrTON, a.TON)
	f.Set(pdufield.AddressRange, a.Address)
}
