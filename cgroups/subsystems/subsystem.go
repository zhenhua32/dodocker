package subsystems

// 定义资源配置
type ResourceConfig struct {
	MemoryLimit string
	CpuShare    string
	CpuSet      string
}

// 定义一个子系统应该实现的方法
type Subsystem interface {
	// 名字很关键, 这是用于查找 mount 路径的
	Name() string
	// 添加 cgroup
	Set(path string, res *ResourceConfig) error
	// 将 pid 加入到 cgroups
	Apply(path string, pid int) error
	// 移除 cgroup
	Remove(path string) error
}

var SubsystemsIns = []Subsystem{
	&CpusetSubsystem{},
	&MemorySubsystem{},
	&CpuSubsystem{},
}
