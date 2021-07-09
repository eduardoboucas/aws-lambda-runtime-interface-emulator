package syscallproxy

import (
	"os"
	"syscall"
)

func CloseOnExec(sock int) {
	syscall.CloseOnExec(syscall.Handle(sock))
}

func CreateNewProcessGroup() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
}

func KillProcess(pid int) {
	p, err := os.FindProcess(pid)

	if err != nil {
		return
	}

	p.Signal(syscall.SIGTERM)
}

func KillProcessOrProcessGroup(pid int) {
	KillProcess(pid)
}
