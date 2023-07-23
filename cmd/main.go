package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		log.Printf("usage: %s <command>\n", os.Args[0])
		os.Exit(0)
	}

	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		log.Fatal("invalid or missing command")
	}
}

func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func child() {
	//must("mount", syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
	//must("mkdir", os.MkdirAll("rootfs/oldrootfs", 0700))
	//must("pivot", syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
	syscall.Sethostname([]byte("container"))
	must("chdir", os.Chdir("/"))
	must("mount", syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal("ERROR: %v\n", err)
		os.Exit(1)
	}
}

func must(label string, err error) {
	if err != nil {
		log.Fatal(label, ": ", err)
		os.Exit(1)
	}
}
