/*
This program demonstrates the simple use of gob package
and the simplicity of Reader/Writer interfaces.

Server programs simply listens on a tcp socket and the
clients send a struct encoded using gob. Server decodes
them and sends an Ack encoded using gob as well.

This also transfers Msg as the interface type to showcase
how this can be extended to multiple types.

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

// Message is the generic interface
// type to store the message in Gob.
type Message interface{}

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
	var message Message
	if err := dec.Decode(&message); err != nil {
		fmt.Println("Decode error: ", err)
	}
	msg := message.(*Msg)
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
	var msg2 Message = Msg{Data: text}
	if err := enc.Encode(&msg2); err != nil {
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

func init() {
	// Types that are shipped as implementation of an
	// interface need to be registered. The type has
	// to be exact (including the `&` part)
	gob.Register(&Msg{})
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
