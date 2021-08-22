package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zhenhua32/dodocker/cgroups"
	"github.com/zhenhua32/dodocker/cgroups/subsystems"
	"github.com/zhenhua32/dodocker/container"
)

// 获取构建好的命令, 并运行
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, cwd string) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	parent.Dir, _ = filepath.Abs(cwd)
	if err := parent.Start(); err != nil {
		logrus.Error("运行 Run 时发生错误", err)
	}

	cgroupManager := cgroups.NewCgroupManager("dodocker")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)

	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
