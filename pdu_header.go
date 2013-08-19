package smpp34

type Header struct {
	Length   uint32
	Id       uint32
	Status   uint32
	Sequence uint32
}
