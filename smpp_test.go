package smpp34

import (
	"encoding/hex"
	"fmt"
	. "launchpad.net/gocheck"
	"net"
	"time"
)

func (x *MySuite) Test_Smpp(c *C) {
	go startServer(c)
	time.Sleep(1 * time.Second)

	s, err := NewSmppConnect("localhost", 8080)

	if err != nil {
		fmt.Println("Connect Err:", err)
		return
	}

	p, err := s.Bind(BIND_TRANSCEIVER, "test", "pass", &Params{})
	if err != nil {
		fmt.Println("Bind:", err)
		return
	}
	s.Write(p)

	p, err = s.Read()
	if err != nil {
		fmt.Println("BindResp:", err)
		return
	}

	c.Check(p.GetHeader().Id, Equals, BIND_TRANSCEIVER_RESP)
	c.Check(p.GetField(SYSTEM_ID).String(), Equals, "hugo")

	_, err = s.Read()
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "Unknown PDU Type. ID:2415919125")

	_, err = s.Read()
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "PDU Len different than read bytes")

	_, err = s.Read()
	c.Check(err, NotNil)
	c.Check(err.Error(), Equals, "PDU Len larger than MAX_PDU_SIZE")

}

func startServer(c *C) {
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		// handle error
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn, c)
	}
}

func handleConnection(n net.Conn, c *C) {
	l := make([]byte, 1024)
	_, err := n.Read(l)

	if err != nil {
		fmt.Println("server read error")
	}

	// Bind Resp
	d, _ := hex.DecodeString("000000158000000900000000000000016875676f00")
	n.Write(d)

	time.Sleep(100 * time.Millisecond)

	// Invalid CMD ID
	d, _ = hex.DecodeString("0000001090000015000000005222b523")
	n.Write(d)

	time.Sleep(100 * time.Millisecond)

	// PDU Len different than read bytes
	d, _ = hex.DecodeString("000000178000000900000000000000016875676f00")
	n.Write(d)

	time.Sleep(100 * time.Millisecond)

	// Max PDU Len err
	d, _ = hex.DecodeString("0000F01080000015000000005222b526")
	n.Write(d)
}
