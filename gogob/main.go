/*
This program demonstrates the simple use of gob package
and the simplicity of Reader/Writer interfaces.

Server programs simply listens on a tcp socket and the
clients send a struct encoded using gob. Server decodes
them and sends an Ack encoded using gob as well.

Usage:
$ go run main.go -role server
Received message:  hello
Sent ack back.

$ go run main.go -role client
hello
Received ack from server:  Got your message. Thanks.
*/
package main

import (
	"bufio"
	"encoding/gob"
	"flag"
	"fmt"
	"net"
	"os"
)

// Msg is the message sent by client
// to the server.
type Msg struct {
	Data string
}

// AckMsg is the ack of the Msg from
// server to the client.
type AckMsg struct {
	Data string
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	dec := gob.NewDecoder(conn)
	var msg Msg
	if err := dec.Decode(&msg); err != nil {
		fmt.Println("Decode error: ", err)
	}

	fmt.Println("Received message: ", msg.Data)

	enc := gob.NewEncoder(conn)
	if err := enc.Encode(AckMsg{Data: "Got your message. Thanks."}); err != nil {
		fmt.Println("Encode error: ", err)
		os.Exit(1)
	}
	fmt.Println("Sent ack back.")

}

func startServer() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("error in listening.", err)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("error during accept.", err)
			continue
		}
		go handleConnection(conn)
	}
}

func sendMsgToServer(text string) {
	conn, err := net.Dial("tcp", ":8080")
	defer conn.Close()
	if err != nil {
		fmt.Println("Error connecting to server.")
		os.Exit(1)
	}

	enc := gob.NewEncoder(conn)
	if err := enc.Encode(Msg{Data: text}); err != nil {
		fmt.Println("Decode error: ", err)
		os.Exit(1)
	}

	dec := gob.NewDecoder(conn)
	var ackMsg AckMsg
	if err := dec.Decode(&ackMsg); err != nil {
		fmt.Println("Decode error: ", err)
		os.Exit(1)
	}
	fmt.Println("Received ack from server: ", ackMsg.Data)
}

func startClient() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		sendMsgToServer(text)
	}
}

func main() {
	var role string
	flag.StringVar(&role, "role", "client", "server or client.")
	flag.Parse()

	if role == "server" {
		startServer()
	}
	startClient()
}
