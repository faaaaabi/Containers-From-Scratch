package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "spawnChild":
		spawnChild()
	default:
		panic("bad command")
	}
}

func run() {
	fmt.Printf("Running %v with PID %d\n", os.Args[2:], os.Getpid())

	// ... Unpacking a slices
	cmd := exec.Command("/proc/self/exe", append([]string{"spawnChild"}, os.Args[2:]...)...) // Execute oneself again but now with child argument
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

func spawnChild() {
	fmt.Printf("Running %v with PID %d\n", os.Args[2:], os.Getpid())

	// execute cgroup configuration
	configureControlGroups()

	syscall.Sethostname([]byte("container"))
	syscall.Chroot("/home/fabian/ContainersFromScratch/UbuntuRootFS/ubuntuIMG") // chroot in base image root
	syscall.Chdir("/")                                                          // Set where you end up after chroot
	syscall.Mount("proc", "proc", "proc", 0, "")                                // Arguments: src, target, festype, flags, data. Mount proc to the "container"

	// ... Unpacking a slice
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	syscall.Unmount("/proc", 0) // Unmount proc after we are finished
}

func configureControlGroups() {
	cgroupDir := "/sys/fs/cgroup"                                 // Path to cgroup folder
	pidDir := filepath.Join(cgroupDir, "pids")                    // pid child fodler
	err := os.Mkdir(filepath.Join(pidDir, "fabicontainer"), 0775) // specific pid conf folder for the "container"
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	must(ioutil.WriteFile(filepath.Join(pidDir, "fabicontainer/pids.max"), []byte("20"), 0700))                          // max number of processes inside the container
	must(ioutil.WriteFile(filepath.Join(pidDir, "fabicontainer/notify_on_release"), []byte("1"), 0700))                  // Removes cgroup after "container" quits
	must(ioutil.WriteFile(filepath.Join(pidDir, "fabicontainer/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)) // write PID of parent and child process to cgroup.procs of the "container"
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
