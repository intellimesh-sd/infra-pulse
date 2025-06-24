//go:build windows
// +build windows

package vm

import (
	"context"
	"github.com/clarechu/infra-pulse/src/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"k8s.io/klog/v2"
	"net"
	"os"
	"time"
)

const (
	OperationSystem = "windows"
)

func GetHostname() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return name
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
				ipv4Addr = ipv4
				continue
			}
			if ipv4 := addr.(*net.IPNet).IP.To4(); ipv4 == nil {
				ipv6Addr = addr.(*net.IPNet).IP
				continue
			}
		}
		nets[i.Name] = &models.NetworkConfig{
			Inet:         ipv4Addr.String(),
			Inet6:        ipv6Addr.String(),
			HardwareAddr: i.HardwareAddr.String(),
		}
	}
	return nets, err
}

func GetInfo() (os *models.OSInfo, err error) {
	os = &models.OSInfo{}
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
	klog.Info(v)
	return os, nil
}
