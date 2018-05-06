package main

import (
	"strings"
)

type Command string

const (
	// ABORT command for a single round.
	ABORT Command = "abort"

	// COMMIT command for a single round.
	COMMIT Command = "commit"

	// QUIT exists the node.
	QUIT Command = "quit"

	// UNKNOWN command.
	UNKNOWN Command = "unknown"
)

// parseCommand parses the user passed string
// to one of the commands.
func parseCommand(s string) Command {

	switch strings.ToLower(s) {
	case string(ABORT):
		return ABORT
	case string(COMMIT):
		return ABORT
	case string(QUIT):
		return QUIT
	default:
		return UNKNOWN
	}
}
