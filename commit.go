package main

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func commitContainer(imgageName string, cwd string) {
	cwd, _ = filepath.Abs(cwd)
	mntURL := filepath.Join(cwd, "rootfs")
	imageTar := filepath.Join(cwd, imgageName+".tar")
	fmt.Println("镜像路径是", imageTar)
	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntURL, ".").CombinedOutput(); err != nil {
		logrus.Errorf("Tar folder %s error %v", mntURL, err)
	}
}
