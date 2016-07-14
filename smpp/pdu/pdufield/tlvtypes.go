package pdufield

type MessageStateType uint8

const (
	Enroute       MessageStateType = 1 // The message is in enroute
	Delivered     MessageStateType = 2 // Message is delivered
	Expired       MessageStateType = 3 // Message validity period has expired
	Deleted       MessageStateType = 4 // Message has been deleted
	Undeliverable MessageStateType = 5 // Message is undeliverable
	Accepted      MessageStateType = 6 // Message is in accepted state
	Unknown       MessageStateType = 7 // Message is in invalid state
	Rejected      MessageStateType = 8 // Message is in a rejected state
)
