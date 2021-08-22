package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// 发现特定 cgroups 的挂载路径
func FindCgroupMountpoint(subsystem string) string {
	// 这基本就相当于 mount 命令
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	// 解析的是类似下面这样的字符串
	// 46 32 0:41 / /sys/fs/cgroup/cpu,cpuacct rw,nosuid,nodev,noexec,relatime shared:24 - cgroup cgroup rw,cpu,cpuacct
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		// 使用 空格 分隔
		fields := strings.Split(txt, " ")
		// 只要 fields 的最后一个, 然后用 逗号 分隔
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			// 找到名字相同的
			if opt == subsystem {
				// 返回的第五个值是 mount path
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

// 组合下路径
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}
