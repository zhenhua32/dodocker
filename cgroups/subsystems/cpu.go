package subsystems

import (
	"fmt"
	"os"
	"path"
	"strconv"
)

type CpuSubsystem struct{}

func (s *CpuSubsystem) Name() string {
	return "cpu"
}

func (s *CpuSubsystem) Set(cgroupPath string, res *ResourceConfig) (err error) {
	var subsysCgroupPath string
	if subsysCgroupPath, err = GetCgroupPath(s.Name(), cgroupPath, true); err != nil {
		return err
	}
	if res.CpuShare != "" {
		if err := os.WriteFile(path.Join(subsysCgroupPath, "cpu.shares"), []byte(res.CpuShare), 0644); err != nil {
			return fmt.Errorf("set cgroup share fail %v", err)
		}
	}
	return nil
}

func (s *CpuSubsystem) Remove(cgroupPath string) (err error) {
	var subsysCgroupPath string
	if subsysCgroupPath, err = GetCgroupPath(s.Name(), cgroupPath, false); err != nil {
		return
	}
	return os.RemoveAll(subsysCgroupPath)
}

func (s *CpuSubsystem) Apply(cgroupPath string, pid int) (err error) {
	var subsysCgroupPath string
	if subsysCgroupPath, err = GetCgroupPath(s.Name(), cgroupPath, false); err != nil {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
	if err := os.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return fmt.Errorf("set cgroup proc fail %v", err)
	}
	return nil
}
