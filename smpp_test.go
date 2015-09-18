package smpp34

import (
	"crypto/tls"
	"encoding/hex"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	testClient(t, nil)
}

func TestClientTLS(t *testing.T) {
	cert, err := tls.X509KeyPair(localhostCert, localhostKey)
	if err != nil {
		t.Fatal(err)
	}
	tc := tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	testClient(t, &tc)
}

func testClient(t *testing.T, config *tls.Config) {
	s, err := testServer(config)
	if err != nil {
		t.Fatal("Failed to parse host and port:", err)
	}
	defer s.Close()
	host, port, err := s.Addr()
	if err != nil {
		t.Fatal("Failed parsing server addr:", err)
	}

	var c *Smpp
	if config == nil {
		c, err = NewSmppConnect(host, port)
	} else {
		config.InsecureSkipVerify = true
		c, err = NewSmppConnectTLS(host, port, config)
	}
	if err != nil {
		t.Fatal("Failed to connect:", err)
	}

	p, err := c.Bind(BIND_TRANSCEIVER, "test", "pass", &Params{})
	if err != nil {
		t.Fatal("Failed to bind:", err)
	}
	c.Write(p)

	p, err = c.Read()
	if err != nil {
		t.Fatal("Failed to read bind response:", err)
	}

	if v := p.GetHeader().Id; v != BIND_TRANSCEIVER_RESP {
		t.Fatalf("Unexpected Id: want %q, have %q",
			v, BIND_TRANSCEIVER_RESP)
	}

	if v := p.GetField(SYSTEM_ID).String(); v != "hugo" {
		t.Fatalf("Unexpected SYSTEM_ID: want \"hugo\", have %q", v)
	}
	_, err = c.Read()
	if err != nil && err.Error() != "Unknown PDU Type. ID:2415919125" {
		t.Fatal("Unexpected error:", err)
	}
	_, err = c.Read()
	if err != nil && err.Error() != "PDU Len different than read bytes" {
		t.Fatal("Unexpected error:", err)
	}
	_, err = c.Read()
	if err != nil && err.Error() != "PDU Len larger than MAX_PDU_SIZE" {
		t.Fatal("Unexpected error:", err)
	}
}

type server struct {
	l net.Listener
}

func testServer(config *tls.Config) (*server, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	if config != nil {
		l = tls.NewListener(l, config)
	}
	go func(l net.Listener) {
		for {
			conn, err := l.Accept()
			if err != nil {
				return // expected by Close
			}
			go handleConnection(conn)
		}
	}(l)
	return &server{l}, nil
}

func (s *server) Addr() (host string, port int, err error) {
	h, p, err := net.SplitHostPort(s.l.Addr().String())
	if err != nil {
		return "", 0, err
	}
	pn, err := strconv.Atoi(p)
	if err != nil {
		return "", 0, err
	}
	return h, pn, nil
}

func (s *server) Close() error {
	return s.l.Close()
}

func handleConnection(n net.Conn) {
	l := make([]byte, 1024)
	_, err := n.Read(l)
	if err != nil {
		panic("server read error")
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

// localhostCert is a PEM-encoded TLS cert with SAN IPs
// "127.0.0.1" and "[::1]", expiring at the last second of 2049 (the end
// of ASN.1 time).
// generated from src/crypto/tls:
// go run generate_cert.go  --rsa-bits 1024 --host 127.0.0.1,::1,example.com --ca --start-date "Jan 1 00:00:00 1970" --duration=1000000h
var localhostCert = []byte(`-----BEGIN CERTIFICATE-----
MIICEzCCAXygAwIBAgIQMIMChMLGrR+QvmQvpwAU6zANBgkqhkiG9w0BAQsFADAS
MRAwDgYDVQQKEwdBY21lIENvMCAXDTcwMDEwMTAwMDAwMFoYDzIwODQwMTI5MTYw
MDAwWjASMRAwDgYDVQQKEwdBY21lIENvMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCB
iQKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9SjY1bIw4
iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZBl2+XsDul
rKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQABo2gwZjAO
BgNVHQ8BAf8EBAMCAqQwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0TAQH/BAUw
AwEB/zAuBgNVHREEJzAlggtleGFtcGxlLmNvbYcEfwAAAYcQAAAAAAAAAAAAAAAA
AAAAATANBgkqhkiG9w0BAQsFAAOBgQCEcetwO59EWk7WiJsG4x8SY+UIAA+flUI9
tyC4lNhbcF2Idq9greZwbYCqTTTr2XiRNSMLCOjKyI7ukPoPjo16ocHj+P3vZGfs
h1fIw3cSS2OolhloGw/XM6RWPWtPAlGykKLciQrBru5NAPvCMsb/I1DAceTiotQM
fblo6RBxUQ==
-----END CERTIFICATE-----`)

// localhostKey is the private key for localhostCert.
var localhostKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXgIBAAKBgQDuLnQAI3mDgey3VBzWnB2L39JUU4txjeVE6myuDqkM/uGlfjb9
SjY1bIw4iA5sBBZzHi3z0h1YV8QPuxEbi4nW91IJm2gsvvZhIrCHS3l6afab4pZB
l2+XsDulrKBxKKtD1rGxlG4LjncdabFn9gvLZad2bSysqz/qTAUStTvqJQIDAQAB
AoGAGRzwwir7XvBOAy5tM/uV6e+Zf6anZzus1s1Y1ClbjbE6HXbnWWF/wbZGOpet
3Zm4vD6MXc7jpTLryzTQIvVdfQbRc6+MUVeLKwZatTXtdZrhu+Jk7hx0nTPy8Jcb
uJqFk541aEw+mMogY/xEcfbWd6IOkp+4xqjlFLBEDytgbIECQQDvH/E6nk+hgN4H
qzzVtxxr397vWrjrIgPbJpQvBsafG7b0dA4AFjwVbFLmQcj2PprIMmPcQrooz8vp
jy4SHEg1AkEA/v13/5M47K9vCxmb8QeD/asydfsgS5TeuNi8DoUBEmiSJwma7FXY
fFUtxuvL7XvjwjN5B30pNEbc6Iuyt7y4MQJBAIt21su4b3sjXNueLKH85Q+phy2U
fQtuUE9txblTu14q3N7gHRZB4ZMhFYyDy8CKrN2cPg/Fvyt0Xlp/DoCzjA0CQQDU
y2ptGsuSmgUtWj3NM9xuwYPm+Z/F84K6+ARYiZ6PYj013sovGKUFfYAqVXVlxtIX
qyUBnu3X9ps8ZfjLZO7BAkEAlT4R5Yl6cGhaJQYZHOde3JEMhNRcVFMO8dJDaFeo
f9Oeos0UUothgiDktdQHxdNEwLjQf7lJJBzV+5OtwswCWA==
-----END RSA PRIVATE KEY-----`)
