package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/zhenhua32/dodocker/container"
)

// 获取构建好的命令, 并运行
func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		logrus.Error("运行 Run 时发生错误", err)
	}
	parent.Wait()
	os.Exit(-1)
}
