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
		Cloneflags: syscall.CLONE_NEWUTS | /* Create the Process in a new UTS (UNIX Timsharing System) Namespace. Necessary to set domain and hostname */
			syscall.CLONE_NEWPID | /* Create New PID Namespace */
			syscall.CLONE_NEWNS, /* Start child in a new mount namespace */
		Unshareflags: syscall.CLONE_NEWNS, /* Dont share mount namespace with host */
	}

	cmd.Run()
}

func child() {
	fmt.Printf("Running %v with PID %d\n", os.Args[2:], os.Getpid())

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/home/fabian/ContainersFromScratch/UbuntuRootFS/ubuntuIMG") // chroot in base image root
	syscall.Chdir("/")                                                          // Set where you end up after chroot
	syscall.Mount("proc", "proc", "proc", 0, "")                                // Arguments: src, target, festype, flags, data. Mount proc to the "container"

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

	syscall.Unmount("/proc", 0) // Unmount after we are finished
}
