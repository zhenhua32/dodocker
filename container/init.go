package container

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// 这是容器内部执行的第一个进程
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if len(cmdArray) == 0 {
		return fmt.Errorf("run container get user command error, cmdArray is nil")
	}
	logrus.Infof("command is %s", strings.Join(cmdArray, " "))

	setupMount()

	// 扩展路径
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("Exec loop path error %v", err)
		return err
	}
	logrus.Infof("find path %s", path)

	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logrus.Errorf("Exec error %v", err)
	}
	return nil
}

func readUserCommand() []string {
	// TZH: 这个就是读取了 pipe 文件, 从 NewPipe 分配的
	// index 为 3 的文件描述符, 也就是传递进来的管道的一端
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

// 其实下面的步骤就是 https://manpages.ubuntu.com/manpages/focal/en/man2/pivot_root.2.html 的示例的 Go 版本

// setupMount 设置 mount, 包含更改根目录, 以及加载 proc 等
func setupMount() {
	// 这里获取的当前路径, 就是 cmd.Dir 配置的
	pwd, err := os.Getwd()
	if err != nil {
		logrus.Errorf("get current workdir error %v", err)
		return
	}
	logrus.Infof("current workdir is %s", pwd)

	// https://github.com/xianlubird/mydocker/issues/62
	// systemd 加入linux之后, mount namespace 就变成 shared by default, 所以你必须显示
	//声明你要这个新的mount namespace独立。
	err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		logrus.Errorf("mount / error %v", err)
		return
	}

	if err := pivotRoot(pwd); err != nil {
		logrus.Errorf("pivot root error: %v", err)
		return
	}

	// 挂载 proc 文件系统
	// TZH: 加上 syscall.MS_PRIVATE 就使得挂载是私有的, 就不会在退出之后, 变更本地环境的 proc
	// https://github.com/xianlubird/mydocker/issues/33
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	err = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	if err != nil {
		logrus.Errorf("mount proc error %v", err)
		return
	}

	// tmpfs是一种基于内存的文件系统，可以使用RAM或swap分区来存储
	err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	if err != nil {
		logrus.Infof("mount tmpfs error %v", err)
		return
	}
}

// pivotRoot 替换根目录为指定目录
func pivotRoot(root string) error {
	/**
	  为了使当前root的老 root 和新 root 不在同一个文件系统下，我们把root重新mount了一次
	  bind mount是把相同的内容换了一个挂载点的挂载方法
	*/
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount rootfs to itself error: %v", err)
	}
	// 创建 rootfs/.pivot_root 存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	os.RemoveAll(pivotDir)
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("mkdir pivotDir error :%v", err)
	}
	// pivot_root 到新的rootfs, 现在老的 old_root 是挂载在rootfs/.pivot_root
	// 挂载点现在依然可以在mount命令中看到
	logrus.Infof("root: %s, pivotDir: %s", root, pivotDir)
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("call PivotRoot error %v", err)
	}
	// 修改当前的工作目录到根目录
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	// umount rootfs/.pivot_root
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmount pivot_root dir %v", err)
	}
	// 删除临时文件夹
	return os.Remove(pivotDir)
}
