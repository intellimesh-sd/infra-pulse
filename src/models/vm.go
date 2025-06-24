package models

type OSInfo struct {
	CPU    *CPU        `json:"cpu"`
	Memory *Memory     `json:"memory"`
	Disks  []*DiskInfo `json:"disks"`
}

type CPU struct {
	// 总占用率
	Percent float64 `json:"percent"`
	// 总核数
	Cores int32     `json:"cores"`
	Info  []CPUInfo `json:"info"`
}

type Memory struct {
	// Total 一共使用量
	Total uint64 `json:"total"`
	// Available 可使用率
	Available uint64 `json:"available"`
	// Used 已经使用率
	Used uint64 `json:"used"`
	// UsedPercent 使用率
	UsedPercent float64 `json:"usedPercent"`
	Free        uint64  `json:"free"`
}

type CPUInfo struct {
	// Cores cpu 核数
	Cores int32 `json:"cores"`
	// CPU 型号
	ModelName string `json:"modelName"`
	// HZ
	Mhz float64 `json:"mhz"`
}

type DiskInfo struct {
	Device      string  `json:"device"`
	Path        string  `json:"path"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

type Release struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Platform string `json:"platform"`
	Family   string `json:"family"`
	Version  string `json:"version"`
}

type CommandOptions struct {
	Command string   `json:"command,omitempty" validate:"required"`
	Envs    []string `json:"envs,omitempty"`
	WorkDir string   `json:"workDir,omitempty"`
}

type SoftwareInfo struct {
	Id        string            `json:"id"`
	Processes []*ProcessesState `json:"processes"`
}

type SoftwareOptions struct {
	Id string `json:"id" validate:"required"`
	// Name 特征名称
	Name    string   `json:"name" validate:"required"`
	Command string   `json:"command"`
	Envs    []string `json:"envs,omitempty"`
	WorkDir string   `json:"workDir"`
}
