package container

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

// 这是容器内部执行的第一个进程
func RunContainerInitProcess(command string, args []string) error {
	logrus.Infof("command is %s", command)

	// 挂载 proc 文件系统
	// TZH: 加上 syscall.MS_PRIVATE 就使得挂载是私有的, 就不会在退出之后, 变更本地环境的 proc
	// https://github.com/xianlubird/mydocker/issues/33
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV | syscall.MS_PRIVATE
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	argv := []string{command}
	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		logrus.Errorf("执行失败 %s\n", err)
	}
	return nil
}
