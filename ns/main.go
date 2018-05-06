package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"

	specs "github.com/opencontainers/runtime-spec/specs-go"
)

/*

Syscalls that are relevant for namespaces are:

- clone(2): creates a child process.
- setns(2): reassociate thread with a namespace.
- unshare(2): moves the calling process to a new namespace.

GO has syscall/SysProcAttr which allows one to specify Cloneflags:
	`Cloneflags   uintptr        // Flags for clone calls (Linux only)`

SysProcAttr are part of ProcAttr that holds attributes that will be applied
to a new process started by StartProcess.

*/

// This is an example of a simple clone call that can be
// used to run a command in a brand new set of namespaces.
func getShellInNewNS() {
	cmd := exec.Command("/bin/sh")

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=[in-namespace]- # "}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error in running cmd: %v\n", cmd)
		os.Exit(1)
	}
}

// This method creates a new process with a new namespace
// and we start another gorouting that tries to join that
// namespace.
func joinExistingNS() {

}

func runP(c int) func() {
	switch c {
	// Get shell in brand new namespaces. This is similar to
	// unshare.
	case 1:
		return getShellInNewNS
	case 2:
		return joinExistingNS
	}
	// no-op method.
	return func() { fmt.Println("choose some other number.") }
}

func main() {
	_ = specs.CgroupNamespace

	if len(os.Args) < 2 {
		fmt.Println("Usage: ns <number>")
		os.Exit(1)
	}

	var n int
	var err error
	if n, err = strconv.Atoi(os.Args[1]); err != nil {
		fmt.Println("2nd argument should be a number.")
	}

	runP(n)()
}
