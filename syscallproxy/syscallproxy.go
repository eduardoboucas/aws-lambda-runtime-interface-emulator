// +build !windows

package syscallproxy

import (
	"syscall"
)

func CloseOnExec(sock int) {
	syscall.CloseOnExec(sock)
}

func CreateNewProcessGroup() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}

func KillProcess(pid int) {
	syscall.Kill(pid, syscall.SIGTERM)
}

func KillProcessOrProcessGroup(pid int) {
	pgid, err := syscall.Getpgid(pid)
	if err == nil {
		syscall.Kill(-pgid, 9) // Negative pid sends signal to all in process group
	} else {
		syscall.Kill(pid, 9)
	}
}
