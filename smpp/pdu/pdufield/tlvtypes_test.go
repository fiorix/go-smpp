package pdufield

import (
	"testing"
)

// TestTLVLen, checks the len of the TLV parameter data
func TestTLVLen(t *testing.T) {

	m := make(TLVMap)

	var testCases = []struct {
		//data len expected
		lenExpected uint16

		// identifier TLV parameter
		parameter TLVType

		// Data to send in the TLV parameter
		data interface{}
	}{
		{5, SubmitSMTLVParameter[SourcePort], "12345"},
		{1, SubmitSMTLVParameter[MSValidity], true},
		{1, SubmitSMTLVParameter[MoreMessagesToSend], 2},
		{2, DeliverSMTLVParameter[LanguageIndicator], "es"},
	}

	for _, test := range testCases {
		err := m.Set(test.parameter, test.data)
		if err != nil {
			t.Errorf("Expected nil but get: %v", err)
		}
		tlvb := m[test.parameter]
		if tlvb.Len != test.lenExpected {
			t.Errorf("The size calculation is wrong, expeted %d, but get: %d", 5, tlvb.Len)
		}
	}

}
