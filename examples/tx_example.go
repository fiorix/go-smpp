package main

import (
	"fmt"
	smpp "github.com/CodeMonkeyKevin/smpp34"
)

func main() {
	// connect and bind
	tx, err := smpp.NewTransmitter(
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
	seq, err := tx.SubmitSm("test", "test2", "msg", &smpp.Params{})

	// Pdu gen errors
	if err != nil {
		fmt.Println("SubmitSm err:", err)
	}

	// Should save this to match with message_id
	fmt.Println("seq:", seq)

	for {
		pdu, err := tx.Read() // This is blocking
		if err != nil {
			fmt.Println("Read Err:", err)
			break
		}

		// EnquireLinks are auto handles
		switch pdu.GetHeader().Id {
		case smpp.SUBMIT_SM_RESP:
			// message_id should match this with seq message
			fmt.Println("MSG ID:", pdu.GetField("message_id").Value())
		default:
			// ignore all other PDUs or do what you link with them
			fmt.Println("PDU ID:", pdu.GetHeader().Id)
		}
	}

	fmt.Println("ending...")
}
