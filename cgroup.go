package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"syscall"
)

const cgroupMemoryHierarchy = "/sys/fs/cgroup/memory"

// /proc/self/exe 是指当前程序

func main() {
	if os.Args[0] == "/proc/self/exe" {
		// 容器进程, 这个会运行失败, 因为内存不足
		fmt.Printf("current pid %d \n", syscall.Getegid())
		cmd := exec.Command("bash", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 得到 fork 出来进程映射在外部命名空间的 pid
	cmdPid := cmd.Process.Pid
	cgPath := path.Join(cgroupMemoryHierarchy, "test_meomery")

	fmt.Printf("cmd.Process.Pid: %v \n", cmdPid)

	// 在系统默认创建挂载了 meomery subsystem 的 hierarchy 上创建 cgroup
	os.Mkdir(cgPath, 0755)
	// 将容器进程加入到这个 cgroups 中
	os.WriteFile(path.Join(cgPath, "tasks"), []byte(fmt.Sprintf("%v", cmdPid)), 0644)
	// 限制内存
	os.WriteFile(path.Join(cgPath, "memory.limit_in_bytes"), []byte("100m"), 0644)
	cmd.Wait()
}
