# gogob
This program demonstrates the simple use of gob package
and the simplicity of Reader/Writer interfaces.

Server programs simply listens on a tcp socket and the
clients send a struct encoded using gob. Server decodes
them and sends an Ack encoded using gob as well.

## Usage:
```
$ go run main.go -role server
Received message:  hello
Sent ack back.
```
```
$ go run main.go -role client
hello
Received ack from server:  Got your message. Thanks.
```
