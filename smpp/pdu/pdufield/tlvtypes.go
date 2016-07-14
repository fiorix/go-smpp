package pdufield

// Message States.
const (
	Enroute       uint8 = 1 // The message is in enroute
	Delivered     uint8 = 2 // Message is delivered
	Expired       uint8 = 3 // Message validity period has expired
	Deleted       uint8 = 4 // Message has been deleted
	Undeliverable uint8 = 5 // Message is undeliverable
	Accepted      uint8 = 6 // Message is in accepted state
	Unknown       uint8 = 7 // Message is in invalid state
	Rejected      uint8 = 8 // Message is in a rejected state
)
