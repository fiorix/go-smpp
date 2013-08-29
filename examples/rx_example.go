package main

import (
	"fmt"
	smpp "github.com/CodeMonkeyKevin/smpp34"
)

func main() {
	// connect and bind
	trx, err := smpp.NewTransmitter(
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

	for {
		pdu, err := trx.Read() // This is blocking
		if err != nil {
			fmt.Println("Read Err:", err)
			break
		}

		// EnquireLinks are auto handles
		switch pdu.GetHeader().Id {
		case smpp.DELIVER_SM:
			// received Deliver Sm
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
			// ignore all other PDUs or do what you link with them
			fmt.Println("PDU ID:", pdu.GetHeader().Id)
		}
	}

	fmt.Println("ending...")
}
