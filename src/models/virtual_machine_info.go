package models

import (
	"github.com/clarechu/infra-pulse/src/utils/netstat"
	"net"
)

type VirtualMachineInfo struct {
	*OSInfo
	ID            string                    `json:"id,omitempty"`
	Agent         string                    `json:"agent"`
	State         int                       `json:"state"`
	Hostname      string                    `json:"hostname"`
	Nets          map[string]*NetworkConfig `json:"nets"`
	InterfaceName string                    `json:"interfaceName"`
	// StartTime 操作系统的启动时间
	// time of last system boot
	SystemStartTime string `json:"sysStartTime,omitempty"`
	// AgentStartTime 节点的创建时间
	AgentStartTime string                 `json:"agentStartTime,omitempty"`
	Processes      any                    `json:"processes,omitempty"`
	Ports          []netstat.SockTabEntry `json:"ports,omitempty"`
	Version        string                 `json:"version,omitempty" yaml:"version,omitempty"`
}

// Processes 进程详情
type Processes struct {
	Pid               int32    `json:"pid,omitempty"`
	Environ           []string `json:"environ,omitempty"`
	WorkDir           string   `json:"workDir"`
	Name              string   `json:"name,omitempty"`
	State             int      `json:"state,omitempty"`
	Status            []string `json:"status,omitempty"`
	CPUUsedPercent    float64  `json:"cpuUsedPercent,omitempty"`
	MemoryUsedPercent float32  `json:"memoryUsedPercent,omitempty"`
	Ppid              int32    `json:"ppid,omitempty"`
	StartTime         string   `json:"startTime,omitempty"`
	UserName          string   `json:"userName,omitempty"`
	Command           string   `json:"command,omitempty"`
	Result            string   `json:"result,omitempty"`
}

type SockTabEntries []SockTabEntry

// ProcessesState 进程详情
type ProcessesState struct {
	*Processes
	SockTabEntries `json:"ports"`
}

// SockTabEntry type represents each line of the /proc/net/[tcp|udp]
type SockTabEntry struct {
	LocalAddr  *SockAddr `json:"localAddr,omitempty"`
	RemoteAddr *SockAddr `json:"remoteAddr,omitempty"`
	State      SkState   `json:"state,omitempty"`
}

// SkState type represents socket connection state
type SkState uint8

// SockAddr represents an ip:port pair
type SockAddr struct {
	IP   net.IP `json:"ip,omitempty"`
	Port uint16 `json:"port,omitempty"`
}

func ToSockTabEntry(entries []netstat.SockTabEntry) SockTabEntries {
	sockTabEntries := make([]SockTabEntry, 0)
	for _, en := range entries {
		sockTabEntry := SockTabEntry{
			LocalAddr: &SockAddr{
				IP:   en.LocalAddr.IP,
				Port: en.LocalAddr.Port,
			},
			RemoteAddr: &SockAddr{
				IP:   en.RemoteAddr.IP,
				Port: en.RemoteAddr.Port,
			},
			State: SkState(en.State),
		}
		sockTabEntries = append(sockTabEntries, sockTabEntry)
	}
	return sockTabEntries
}
