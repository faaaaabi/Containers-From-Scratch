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
	fmt.Printf("Running %v with PID %d\n", os.Args[2:], os.Getpid())

	// ... Unpacking a slices
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...) // Execute oneself again but now with child argument
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set OS specific attribute for command execution
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS /* Create the Process in a new Namespace */ | syscall.CLONE_NEWPID,
	}

	cmd.Run()
}

func child() {
	fmt.Printf("Running %v with PID %d\n", os.Args[2:], os.Getpid())

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/home/fabian/ContainersFromScratch/UbuntuRootFS/ubuntuIMG") // chroot in base image root
	syscall.Chdir("/")                                                          // Set where you end up after chroot

	// ... Unpacking a slice
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
