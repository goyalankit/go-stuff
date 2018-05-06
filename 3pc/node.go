package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

type Role int

type Node interface {
	Start(<-chan Command) error
}

type node struct {
}

// NewNode creates a new node in the system.
func NewNode(ctx context.Context) Node {
	return &node{}
}

func listener() error {
	return nil
}

func startServer(connC chan<- net.Conn) {
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
		connC <- conn
	}
}

func (n *node) Start(inputC <-chan Command) error {
	// Start a local server to listen for
	// updates.
	servConnC := make(chan net.Conn)
	go startServer(servConnC)

	for {
		select {

		case <-time.After(time.Duration(rand.Int31n(5)+5) * time.Second):
			// initiate coordinator election.
			// fmt.Println("Doing a heartbeat.")
		case conn := <-servConnC:
			_ = conn
		case cmd := <-inputC:
			fmt.Println("User sent a cmd: ", cmd)
		}
	}
	return nil
}
