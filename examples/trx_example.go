package main

import (
	"fmt"
	smpp "github.com/CodeMonkeyKevin/smpp34"
)

func main() {
	// connect and bind
	trx, err := smpp.NewTransceiver(
		"localhost",
		9000,
		5,
		smpp.Params{
			"system_type": "CMT",
			"system_id":   "hugo",
			"password":    "ggoohu",
		},
	)
	if err != nil {
		fmt.Println("Connection Err:", err)
		return
	}

	// Send SubmitSm
	seq, err := trx.SubmitSm("test", "test2", "msg", &smpp.Params{})

	// Pdu gen errors
	if err != nil {
		fmt.Println("SubmitSm err:", err)
	}

	// Should save this to match with message_id
	fmt.Println("seq:", seq)

	// start reading PDUs
	for {
		pdu, err := trx.Read() // This is blocking
		if err != nil {
			break
		}

		// Transceiver auto handles EnquireLinks
		switch pdu.GetHeader().Id {
		case smpp.SUBMIT_SM_RESP:
			// message_id should match this with seq message
			fmt.Println("MSG ID:", pdu.GetField("message_id").Value())
		case smpp.DELIVER_SM:
			// received Deliver Sm

			// Print all fields
			for _, v := range pdu.MandatoryFieldsList() {
				f := pdu.GetField(v)
				fmt.Println(v, ":", f)
			}

			// Respond back to Deliver SM with Deliver SM Resp
			err := trx.DeliverSmResp(pdu.GetHeader().Sequence, smpp.ESME_ROK)

			if err != nil {
				fmt.Println("DeliverSmResp err:", err)
			}
		default:
			fmt.Println("PDU ID:", pdu.GetHeader().Id)
		}
	}

	fmt.Println("ending...")
}
