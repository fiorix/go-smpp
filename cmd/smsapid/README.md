# HTTP API for sending SMS

The smsapid tool is a web server that connects to an SMSC and provide
endpoints for sending short messages, querying their status when
supported by the SMSC. It also supports sending delivery receipts
via [Server-Sent Events](http://www.w3schools.com/html/html5_serversentevents.asp).

The HTTP server supports HTTP/1 and HTTP/2, and TLS.

## Usage

With the server running, send a message:

	curl localhost:8080/v1/send -X POST -F src=root -F dst=foobar -F msg=hi

In case of success, the server returns a JSON document containing a
message ID that can be used for querying its delivery status later.
This functionality is not always available on some SMSCs.

	curl localhost:8080/v1/query -X GET -F src=root -F msgid=1234

For collecting incoming SMS, or delivery receipts:

	curl localhost:8080/v1/delivery

This is the Server-Sent Events (SSE) endpoint that deliver messages
as events, as they arrive on the server.

## Send parameters

The `/v1/send` endpoint supports the following parameters:

- src: number of sender (optional)
- dst: number of recipient
- msg: text message, encoded as UTF-8
- enc: text encoding for SMS delivery: latin1 or ucs2 (optional)
- register: register for delivery: final, failure (optional)

If an encoding is not provided, data is sent as a binary blob and may
not display well on devices.

For special characters, try:

	curl localhost:8080/v1/send -X POST -F dst=foobar -F enc=ucs2 -F msg="é nóis"

Have fun!
