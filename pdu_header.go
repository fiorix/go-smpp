package smpp34

type Header struct {
	Length   uint32
	Id       uint32
	Status   uint32
	Sequence uint32
}

func NewPduHeader(l, id, status, seq uint32) *Header {
	return &Header{l, id, status, seq}
}
