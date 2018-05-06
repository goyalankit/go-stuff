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

	inputC := make(chan Command)

	n := NewNode(ctx)
	go n.Start(inputC)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		cmd := scanner.Text()

		switch parseCommand(cmd) {
		case COMMIT:
			fmt.Println("Commit received.")
			inputC <- COMMIT
		case ABORT:
			fmt.Println("Abort received.")
			inputC <- ABORT
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
