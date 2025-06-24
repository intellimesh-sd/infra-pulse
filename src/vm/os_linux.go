//go:build linux
// +build linux

package vm

import (
	"context"
	"errors"
	"fmt"
	"github.com/clarechu/infra-pulse/src/models"
	"github.com/clarechu/infra-pulse/src/utils/netstat"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"k8s.io/klog/v2"
	"net"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	OperationSystem = "linux"
)

// GetHostname 获取当前虚拟机的hostname
func GetHostname() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return name, nil
}

// GetNetWorks 获取当前虚拟机的网卡信息
func GetNetWorks() (map[string]*models.NetworkConfig, error) {
	nets := make(map[string]*models.NetworkConfig, 0)
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}
		var ipv4Addr net.IP
		var ipv6Addr net.IP
		for _, addr := range addrs {
			if ipv4 := addr.(*net.IPNet).IP.To4(); ipv4 != nil {
				if ipv4Addr == nil {
					ipv4Addr = ipv4
				}
				continue
			}
			if ipv4 := addr.(*net.IPNet).IP.To4(); ipv4 == nil {
				if ipv6Addr == nil {
					ipv6Addr = addr.(*net.IPNet).IP
				}
				continue
			}
		}
		nets[i.Name] = &models.NetworkConfig{
			Inet:         IPToString(ipv4Addr),
			Inet6:        IPToString(ipv6Addr),
			HardwareAddr: i.HardwareAddr.String(),
		}
	}
	return nets, err
}

// IPToString 找到网卡IP nets 不存在的情况
func IPToString(ip net.IP) string {
	if ip.String() == "<nil>" {
		return ""
	}
	return ip.String()
}

func GetInfo() (os *models.OSInfo, err error) {
	os = &models.OSInfo{
		Disks: make([]*models.DiskInfo, 0),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c := &models.CPU{
		Info: []models.CPUInfo{},
	}

	percent, err := cpu.PercentWithContext(ctx, 0, false)
	if err != nil {
		klog.Errorf("get cpu percents:%s", err.Error())
		return nil, err
	}
	infos, err := cpu.InfoWithContext(ctx)
	if err != nil {
		klog.Errorf("get cpu info :%s", err.Error())
		return nil, err
	}
	for _, info := range infos {
		// klog.Info("cpu info: %s", info.String())
		in := models.CPUInfo{
			Cores:     info.Cores,
			ModelName: info.ModelName,
			Mhz:       info.Mhz,
		}
		c.Info = append(c.Info, in)
		c.Cores = c.Cores + info.Cores
		c.Percent = percent[0]
	}
	os.CPU = c
	v, err := mem.VirtualMemory()
	if err != nil {
		klog.Errorf("get memory percents:%s", err.Error())
		return nil, err
	}
	memory := models.Memory{
		Total:       v.Total,
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
		Free:        v.Free,
	}
	os.Memory = &memory
	fdisks, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return nil, err
	}
	for _, d := range fdisks {
		counters, err := disk.Usage(d.Mountpoint)
		if err != nil {
			return nil, err
		}
		di := &models.DiskInfo{
			Device:      d.Device,
			Path:        d.Mountpoint,
			Fstype:      d.Fstype,
			Total:       counters.Total,
			Free:        counters.Free,
			Used:        counters.Used,
			UsedPercent: counters.UsedPercent,
		}
		os.Disks = append(os.Disks, di)
	}
	return os, nil
}

// GetLastBootTime 操作系统最后启动时间
func GetLastBootTime() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	withContext, err := host.BootTimeWithContext(ctx)
	if err != nil {
		return "", err
	}
	return time.Unix(int64(withContext), 0).Format("2006-01-02 15:04:05"), nil
}

func GetMainIp() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	conn, err := net.Dial("udp", "8.8.8.8:8")
	if err != nil {
		return "", err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			if ipv4 := addr.(*net.IPNet).IP.To4(); ipv4 != nil {
				if localAddr.IP.String() == ipv4.String() {
					return i.Name, nil
				}
			}
		}
	}

	return "", errors.New("not found main ip")
}

// GetProcessesState 根据当前进程查询当前进程的详情
func GetProcessesState(hasListenPort bool) ([]*models.ProcessesState, error) {
	socks, err := Ports()
	if err != nil {
		return nil, err
	}
	processes := make([]*models.ProcessesState, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ps, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range ps {
		procName, _ := p.Cmdline()
		status, _ := p.Status()
		cpuPercentWithContext, _ := p.CPUPercentWithContext(ctx)
		memPercentWithContext, _ := p.MemoryPercentWithContext(ctx)
		ppid, _ := p.Ppid()
		createTime, _ := p.CreateTime()
		username, _ := p.Username()
		environ, _ := p.Environ()
		command, _ := p.Cmdline()
		workdir, _ := p.Cwd()
		if workdir == "" {
			workdir = os.Getenv("HOME")
		}
		proc := &models.Processes{
			Pid:               p.Pid,
			WorkDir:           workdir,
			Environ:           DeleteSliceString(environ, ""),
			UserName:          username,
			Command:           command,
			Name:              procName,
			State:             1,
			Status:            status,
			CPUUsedPercent:    cpuPercentWithContext,
			MemoryUsedPercent: memPercentWithContext,
			Ppid:              ppid,
			StartTime:         time.UnixMilli(createTime).Format("2006-01-02 15:04:05"),
		}
		sockTabEntries := make([]netstat.SockTabEntry, 0)
		for _, sock := range socks {
			if sock.Process != nil && sock.Process.Pid == int(p.Pid) {
				sockTabEntries = append(sockTabEntries, sock)
			}
		}
		if !hasListenPort || len(sockTabEntries) != 0 {
			processes = append(processes, &models.ProcessesState{
				Processes:      proc,
				SockTabEntries: models.ToSockTabEntry(sockTabEntries),
			})
		}
	}
	return processes, nil
}

type NoopFilter func(processes *models.Processes) error

// GetProcessesByName 根据当前进程查询当前进程的详情
func GetProcessesByName(name string, filter NoopFilter) ([]*models.ProcessesState, error) {
	processes := make([]*models.ProcessesState, 0)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ps, err := process.Processes()
	if err != nil {
		return nil, err
	}
	for _, p := range ps {
		procName, _ := p.Cmdline()
		if !strings.Contains(procName, name) {
			continue
		}
		status, _ := p.Status()
		cpuPercentWithContext, _ := p.CPUPercentWithContext(ctx)
		memPercentWithContext, _ := p.MemoryPercentWithContext(ctx)
		ppid, _ := p.Ppid()
		createTime, _ := p.CreateTime()
		username, _ := p.Username()
		environ, _ := p.Environ()
		command, _ := p.Cmdline()
		workdir, _ := p.Cwd()
		if workdir == "" {
			workdir = os.Getenv("HOME")
		}
		proc := &models.Processes{
			Pid:               p.Pid,
			WorkDir:           workdir,
			Environ:           DeleteSliceString(environ, ""),
			UserName:          username,
			Command:           command,
			Name:              name,
			State:             1,
			Status:            status,
			CPUUsedPercent:    cpuPercentWithContext,
			MemoryUsedPercent: memPercentWithContext,
			Ppid:              ppid,
			StartTime:         time.UnixMilli(createTime).Format("2006-01-02 15:04:05"),
		}
		if filter != nil {
			err = filter(proc)
			if err != nil {
				return nil, err
			}
		}
		sockTabEntries, err := PortsByPid(int(p.Pid))
		if err != nil {
			return nil, err
		}
		processes = append(processes, &models.ProcessesState{
			Processes:      proc,
			SockTabEntries: models.ToSockTabEntry(sockTabEntries),
		})
	}

	return processes, nil
}

// DeleteSliceString 删除指定元素。
func DeleteSliceString(s []string, elem string) []string {
	j := 0
	for _, v := range s {
		if v != elem {
			s[j] = v
			j++
		}
	}
	return s[:j]
}

func GetProcessesByPort(port int32) (int, error) {
	ports, err := Ports()
	if err != nil {
		return 0, err
	}
	for _, sockTabEntry := range ports {
		if sockTabEntry.LocalAddr == nil {
			continue
		}
		// 避免 Process 对象为空
		if int32(sockTabEntry.LocalAddr.Port) == port && sockTabEntry.Process != nil {
			return sockTabEntry.Process.Pid, nil
		}
	}
	return 0, fmt.Errorf("ports: %d not found ", port)
}

func GetProcessesById(pid int32) (*models.Processes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ps, err := process.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	name, _ := ps.Name()
	environ, _ := ps.Environ()
	status, _ := ps.Status()
	cpuPercentWithContext, _ := ps.CPUPercentWithContext(ctx)
	memPercentWithContext, _ := ps.MemoryPercentWithContext(ctx)
	ppid, _ := ps.Ppid()
	createTime, _ := ps.CreateTime()
	username, _ := ps.Username()
	command, _ := ps.Cmdline()
	workdir, _ := ps.Cwd()
	processes := &models.Processes{
		Pid:               ps.Pid,
		WorkDir:           workdir,
		Environ:           DeleteSliceString(environ, ""),
		UserName:          username,
		Command:           command,
		Name:              name,
		State:             1,
		Status:            status,
		CPUUsedPercent:    cpuPercentWithContext,
		MemoryUsedPercent: memPercentWithContext,
		Ppid:              ppid,
		StartTime:         time.UnixMilli(createTime).Format("2006-01-02 15:04:05"),
	}
	return processes, nil
}

// PortsByPid 获取当前虚拟机中的端口号
func PortsByPid(pid int) ([]netstat.SockTabEntry, error) {
	socks, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		return nil, err
	}
	socksv6, err := netstat.TCP6Socks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		return nil, err
	}
	socks = append(socks, socksv6...)
	entries := make([]netstat.SockTabEntry, 0)
	for _, sock := range socks {
		if sock.Process != nil && sock.Process.Pid == pid {
			entries = append(entries, sock)
		}
	}
	return entries, nil
}

// Ports 获取当前虚拟机中的端口号
func Ports() ([]netstat.SockTabEntry, error) {
	socks, err := netstat.TCPSocks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		return nil, err
	}
	socksv6, err := netstat.TCP6Socks(func(s *netstat.SockTabEntry) bool {
		return s.State == netstat.Listen
	})
	if err != nil {
		return nil, err
	}
	return append(socks, socksv6...), nil
}

func GetProcesses() ([]*models.Processes, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	ps, err := process.Processes()
	if err != nil {
		return nil, err
	}
	processes := make([]*models.Processes, 0)
	for _, p := range ps {
		name, err := p.Name()
		if err != nil {
			continue
		}
		status, err := p.Status()
		if err != nil {
			continue
		}
		cpuPercentWithContext, err := p.CPUPercentWithContext(ctx)
		if err != nil {
			continue
		}
		memPercentWithContext, err := p.MemoryPercentWithContext(ctx)
		if err != nil {
			continue
		}
		ppid, err := p.Ppid()
		if err != nil {
			continue
		}
		createTime, err := p.CreateTime()
		if err != nil {
			continue
		}
		username, err := p.Username()
		if err != nil {
			continue
		}
		command, err := p.Cmdline()
		if err != nil {
			continue
		}
		processes = append(processes, &models.Processes{
			Pid:               p.Pid,
			UserName:          username,
			Command:           command,
			Name:              name,
			State:             1,
			Status:            status,
			CPUUsedPercent:    cpuPercentWithContext,
			MemoryUsedPercent: memPercentWithContext,
			Ppid:              ppid,
			StartTime:         time.UnixMilli(createTime).Format("2006-01-02 15:04:05"),
		})

	}
	return processes, nil
}

func ReleaseInfo() (*models.Release, error) {
	version, _ := host.KernelVersion()
	fmt.Println(version)

	platform, family, version, err := host.PlatformInformation()
	if err != nil {
		return nil, err
	}
	return &models.Release{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Platform: platform,
		Family:   family,
		Version:  version,
	}, err
}
