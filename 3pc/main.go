package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	flag.String("role", "cohort", "Add a new node in the system.")
	flag.Parse()

	ctx, cancelFun := context.WithCancel(context.Background())
	defer cancelFun()

	n := NewNode(ctx, CoordinatorRole)
	go n.Start()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		cmd := scanner.Text()

		switch ParseCommand(cmd) {
		case COMMIT:
			fmt.Println("Commit received.")
		case ABORT:
			fmt.Println("Abort received.")
		case QUIT:
			fmt.Println("Quit Received")
			os.Exit(0)
		case UNKNOWN:
			fmt.Println("Unknown command.")
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}
