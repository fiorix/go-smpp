# SMPP 3.4

This is an implementation of SMPP 3.4 for Go, based on the original
[smpp34](https://github.com/CodeMonkeyKevin/smpp34) from Kevin Patel.

The API has been refactored to idiomatic Go code with more tests
and documentation. There are also quite a few new features, such
as a test server (see smpptest package) and support for text
transformation for LATIN-1 and UCS-2.

It is not fully compliant, there are some TODOs in the code.

### Usage

[![GoDoc](https://godoc.org/github.com/fiorix/go-smpp?status.svg)](https://godoc.org/github.com/fiorix/go-smpp)

Following is an SMPP client transmitter wrapped by an HTTP server
that can send Short Messages (SMS):

```go
func main() {
	tx := &smpp.Transmitter{
		Addr:    "localhost:2775",
		User:    "foobar",
		Passwd:  "secret",
	}
	conn := tx.Bind() // make persistent connection.
	go func() {
		for c := range conn {
			log.Println("SMPP connection status:", c.Status())
		}
	}()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		sm, err := tx.Submit(&smpp.ShortMessage{
			Src:      r.FormValue("src"),
			Dst:      r.FormValue("dst"),
			Text:     pdutext.Raw(r.FormValue("text")),
			Register: smpp.NoDeliveryReceipt,
		})
		if err == smpp.ErrNotConnected {
			http.Error(w, "Oops.", http.StatusServiceUnavailable)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		io.WriteString(w, sm.RespID())
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

You can test from the command line:

	curl localhost:8080 -X GET -F src=bart -F dst=lisa -F text=hello

If you don't have an SMPP server to test, check out
[Selenium SMPPSim](http://www.seleniumsoftware.com/downloads.html).
It has been used for the development of this package.

## Supported PDUs

- [x] bind_transmitter
- [x] bind_transmitter_resp
- [x] bind_receiver
- [x] bind_receiver_resp
- [x] bind_transceiver
- [x] bind_transceiver_resp
- [ ] outbind
- [x] unbind
- [x] unbind_resp
- [x] submit_sm
- [x] submit_sm_resp
- [ ] submit_sm_multi
- [ ] submit_sm_multi_resp
- [ ] data_sm
- [ ] data_sm_resp
- [x] deliver_sm
- [x] deliver_sm_resp
- [x] query_sm
- [x] query_sm_resp
- [ ] cancel_sm
- [ ] cancel_sm_resp
- [ ] replace_sm
- [ ] replace_sm_resp
- [x] enquire_link
- [x] enquire_link_resp
- [ ] alert_notification
- [x] generic_nack

## Copyright

See LICENSE and AUTHORS files for details.
