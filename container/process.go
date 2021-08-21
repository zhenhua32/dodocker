package container

import (
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	// init 是子命令, 就是 main_command.go 定义的 initCommand 函数
	args := []string{"init", command}
	// 这里就是调用自己
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}
