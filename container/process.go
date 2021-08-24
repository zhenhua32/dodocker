package container

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/sirupsen/logrus"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
		return nil, nil
	}
	// 这里就是调用自己
	cmd := exec.Command("/proc/self/exe", "init")
	// 在新的 namespace 中执行
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
	// TZH:
	// go 文档中说 If non-nil, entry i becomes file descriptor 3+i.
	// 所以, readUserCommand 函数里就使用了 uintptr(3)
	// https://stackoverflow.com/questions/29528756/how-can-i-read-from-exec-cmd-extrafiles-fd-in-child-process
	cmd.ExtraFiles = []*os.File{readPipe}
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

// NewWorkSpace 创建一个 overlay2 filesystem 作为 container root workspace
func NewWorkSpace(rootURL string, mntURL string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
}

func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := filepath.Join(rootURL, "busybox")
	busyboxTarURL := filepath.Join(rootURL, "busybox.tar")
	exist, err := PathExists(busyboxURL)
	if err != nil {
		logrus.Infof("Fail to judge whether dir %s exists. %v", busyboxURL, err)
	}
	// 路径不存在, 就先解压内容 "busybox.tar" 到 "busybox/"
	if !exist {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			logrus.Errorf("Mkdir dir %s error. %v", busyboxURL, err)
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			logrus.Errorf("Untar dir %s error %v", busyboxURL, err)
		}
	}
}

func CreateWriteLayer(rootURL string) {
	writeURL := filepath.Join(rootURL, "writeLayer")
	if err := os.MkdirAll(writeURL, 0777); err != nil {
		logrus.Errorf("Mkdir dir %s error. %v", writeURL, err)
	}
	// 顺便也建立个 workdir
	workURL := filepath.Join(rootURL, "workLayer")
	if err := os.MkdirAll(workURL, 0777); err != nil {
		logrus.Errorf("Mkdir dir %s error. %v", workURL, err)
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.MkdirAll(mntURL, 0777); err != nil {
		logrus.Errorf("Mkdir dir %s error. %v", mntURL, err)
	}
	lowerdir := filepath.Join(rootURL, "busybox")
	upperdir := filepath.Join(rootURL, "writeLayer/")
	workdir := filepath.Join(rootURL, "workLayer/")
	dirs := "lowerdir=" + lowerdir + ",upperdir=" + upperdir + ",workdir=" + workdir
	logrus.Infof("dirs is %s", dirs)
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", dirs, mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("%v", err)
	}
}

// DeleteWorkSpace 删除 overlay2 filesystem 当 container 退出时
func DeleteWorkSpace(rootURL string, mntURL string) {
	DeleteMountPoint(rootURL, mntURL)
	DeleteWriteLayer(rootURL)
}

func DeleteMountPoint(rootURL string, mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		logrus.Errorf("%v", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		logrus.Errorf("Remove dir %s error %v", mntURL, err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := filepath.Join(rootURL, "writeLayer")
	if err := os.RemoveAll(writeURL); err != nil {
		logrus.Errorf("Remove dir %s error %v", writeURL, err)
	}
	workURL := filepath.Join(rootURL, "workLayer")
	if err := os.RemoveAll(workURL); err != nil {
		logrus.Errorf("Remove dir %s error %v", writeURL, err)
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}
