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
// cwd 是当前工作目录, 应该有 busybox.tar 等文件
func Run(tty bool, comArray []string, res *subsystems.ResourceConfig, cwd string, volume string) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	// 指定当前运行的工作目录
	cwd, _ = filepath.Abs(cwd)
	mntURL := filepath.Join(cwd, "rootfs")
	// 构建隔离空间
	container.NewWorkSpace(cwd, mntURL, volume)
	parent.Dir = mntURL

	if err := parent.Start(); err != nil {
		logrus.Error("运行 Run 时发生错误", err)
	}

	cgroupManager := cgroups.NewCgroupManager("dodocker")
	defer cgroupManager.Destroy()
	cgroupManager.Set(res)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(comArray, writePipe)

	parent.Wait()

	// 到了退出的时候了, 清理目录
	container.DeleteWorkSpace(cwd, mntURL, volume)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
