package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
	"time"

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

// Taken from netns package
// SYS_SETNS syscall allows changing the namespace of the current process.
var SYS_SETNS = map[string]uintptr{
	"386":     346,
	"amd64":   308,
	"arm64":   268,
	"arm":     375,
	"mips":    4344,
	"mipsle":  4344,
	"ppc64":   350,
	"ppc64le": 350,
	"s390x":   339,
}[runtime.GOARCH]

// This method creates a new process with a new namespace
// and we start another gorouting that tries to join that
// namespace.
func joinExistingNS() {
	// STEP 1 START
	// create a new sleep process with new namespaces.
	cmd := exec.Command("/bin/sleep", "500")

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
	if err := cmd.Start(); err != nil {
		fmt.Println("ah man.. failed to run the command.")
	}
	pid := cmd.Process.Pid
	// STEP 1 END

	// STEP 2 CREATE A NEW PROCESS AND SETNS OF THE PREVIOUS PROCESS
	// Make sure we don't change the namespace of the current process.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	fmt.Println("running the process as ", pid)
	time.Sleep(1 * time.Second)

	// Enter pid and mount namespace of the sleep process.
	for _, ns := range []string{"pid", "uts", "net"} {
		nsPath := fmt.Sprintf("/proc/%d/ns/%s", pid, ns)
		// Opening the namespace changes your namespac.e
		fd, err := syscall.Open(nsPath, syscall.O_RDONLY, 0)
		defer syscall.Close(fd)

		if err != nil {
			fmt.Println("failed to open the ns", err)
			os.Exit(1)
		}

		var f int
		switch ns {
		case "pid":
			f = syscall.CLONE_NEWPID
		default:
			f = 0
		}
		_ = f
		fmt.Println("joining the namespace: ", ns)
		_, _, e1 := syscall.RawSyscall(SYS_SETNS, uintptr(fd), uintptr(f), 0)
		if e1 != 0 {
			fmt.Println("setns failed. ", e1)
			os.Exit(1)
		}
	}
	// STEP 2 CREATE A NEW PROCESS AND SETNS OF THE PREVIOUS PROCESS

	// Now that we have called setns. Call a command to execute
	// bash again. without specifying namespace again.
	// You should be dropped into the namespace of the sleep process.

	// Get shell in the new process.
	cmd = exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=[in-namespace2]- # "}
	if err := cmd.Run(); err != nil {
		fmt.Println("run failed in the new namespace.")
		os.Exit(1)
	}

	/*

		[angoyal@angoyal-ld3 go-stuff]$ sudo ns/bin/ns 2
		running the process as  19825
		joining the namespace:  pid
		joining the namespace:  uts
		joining the namespace:  net
		[in-namespace2]- # echo $$
		2

		[angoyal@angoyal-ld3 go-stuff]$ sudo ls -l /proc/19829/ns/uts
		lrwxrwxrwx 1 root root 0 May  5 18:17 /proc/19829/ns/uts -> uts:[4026532639]
		[angoyal@angoyal-ld3 go-stuff]$ sudo ls -l /proc/19825/ns/uts
		lrwxrwxrwx 1 root root 0 May  5 18:16 /proc/19825/ns/uts -> uts:[4026532639]

	*/
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
