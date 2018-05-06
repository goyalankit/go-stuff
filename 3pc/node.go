package main

import (
	"context"
	"time"
)

type Role int

const (
	CohortRole Role = iota
	CoordinatorRole
)

type Node interface {
	Start() error
}

type node struct {
	role Role
}

// NewNode creates a new node in the system.
func NewNode(ctx context.Context, r Role) Node {
	return &node{role: r}
}

func listner() error {

	return nil
}

func (n *node) Start() error {

	for {
		select {
		case <-time.After(2 * time.Second):
			// initiate coordinator election.
			// fmt.Println("Doing a heartbeat.")
		}
	}

	return nil
}
