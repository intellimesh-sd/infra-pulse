package models

// RegisterVirtualMachine 虚拟机基本详情
type RegisterVirtualMachine struct {
	ID    string `json:"id,omitempty"`
	Agent string `json:"agent,omitempty"`
	//当前主机的主机名称
	Hostname string `json:"hostname"`
	Port     int32  `json:"port"`
	// 主网络接口名称
	InterfaceName string `json:"interfaceName"`

	//网口基本信息
	Networks map[string]*NetworkConfig `json:"nets"`
	// Linux 发行版 ubuntu centos
	Distribution string `json:"distribution,omitempty"`
	// OperationSystem 操作系统 default: linux
	OperationSystem string `json:"operationSystem,omitempty"`
	// Arch CPU 架构 default: arch64 amd64 x86_64
	Arch string `json:"arch,omitempty"`
	// family debian
	Family string `json:"family,omitempty"`
	// Version 版本 18.04
	Version string `json:"version,omitempty"`
}

type NetworkConfig struct {
	Inet         string `json:"inet"`
	Inet6        string `json:"inet6"`
	HardwareAddr string `json:"hardwareAddr"`
}
