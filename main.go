package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("bad command")
	}
}

func run() {
	fmt.Printf("Running %v\n", os.Args[2:])

	// ... Unpacking a slices
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set OS specific attribute for command execution
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS, /* Create the Process in a new Namespace */
	}

	cmd.Run()
}

func child() {
	fmt.Printf("Running %v\n", os.Args[2:])

	syscall.Sethostname([]byte("container"))

	// ... Unpacking a slices
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set OS specific attribute for command execution
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS, /* Create the Process in a new Namespace */
	}

	cmd.Run()
}
