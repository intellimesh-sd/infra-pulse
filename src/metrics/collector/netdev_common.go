// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build !nonetdev && (linux || freebsd || openbsd || dragonfly || darwin)
// +build !nonetdev
// +build linux freebsd openbsd dragonfly darwin

package collector

import (
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"net"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	/*netdevDeviceInclude = flag.String("collector.netdev.device-include", "",
		"Regexp of net devices to include (mutually exclusive to device-exclude).")
	oldNetdevDeviceInclude = flag.String("collector.netdev.device-whitelist", "",
		"DEPRECATED: Use collector.netdev.device-include")
	netdevDeviceExclude = flag.String("collector.netdev.device-exclude", "",
		"Regexp of net devices to exclude (mutually exclusive to device-include).")
	oldNetdevDeviceExclude = flag.String("collector.netdev.device-blacklist", "",
		"DEPRECATED: Use collector.netdev.device-exclude")
	netdevAddressInfo = flag.Bool("collector.netdev.address-info", false,
		"Collect address-info for every device")
	netdevDetailedMetrics = flag.Bool("collector.netdev.enable-detailed-metrics", false,
		"Use (incompatible) metric names that provide more detailed stats on Linux")
	*/
	netdevDeviceInclude    = ""
	oldNetdevDeviceInclude = ""
	netdevDeviceExclude    = ""
	oldNetdevDeviceExclude = ""
	netdevAddressInfo      = false
	netdevDetailedMetrics  = false
)

type netDevCollector struct {
	subsystem        string
	deviceFilter     deviceFilter
	metricDescsMutex sync.Mutex
	metricDescs      map[string]*prometheus.Desc
}

type netDevStats map[string]map[string]uint64

func init() {
	registerCollector("netdev", defaultEnabled, NewNetDevCollector)
}

// NewNetDevCollector returns a new Collector exposing network device stats.
func NewNetDevCollector() (Collector, error) {
	if oldNetdevDeviceInclude != "" {
		if netdevDeviceInclude == "" {
			klog.V(2).Info("msg", "--collector.netdev.device-whitelist is DEPRECATED and will be removed in 2.0.0, use --collector.netdev.device-include")
			netdevDeviceInclude = oldNetdevDeviceInclude
		} else {
			return nil, errors.New("--collector.netdev.device-whitelist and --collector.netdev.device-include are mutually exclusive")
		}
	}

	if oldNetdevDeviceExclude != "" {
		if netdevDeviceExclude == "" {
			klog.V(2).Info("msg", "--collector.netdev.device-blacklist is DEPRECATED and will be removed in 2.0.0, use --collector.netdev.device-exclude")
			netdevDeviceExclude = oldNetdevDeviceExclude
		} else {
			return nil, errors.New("--collector.netdev.device-blacklist and --collector.netdev.device-exclude are mutually exclusive")
		}
	}

	if netdevDeviceExclude != "" && netdevDeviceInclude != "" {
		return nil, errors.New("device-exclude & device-include are mutually exclusive")
	}

	if netdevDeviceExclude != "" {
		klog.Info("msg", "Parsed flag --collector.netdev.device-exclude", "flag", netdevDeviceExclude)
	}

	if netdevDeviceInclude != "" {
		klog.Info("msg", "Parsed Flag --collector.netdev.device-include", "flag", netdevDeviceInclude)
	}

	return &netDevCollector{
		subsystem:    "network",
		deviceFilter: newDeviceFilter(netdevDeviceExclude, netdevDeviceInclude),
		metricDescs:  map[string]*prometheus.Desc{},
	}, nil
}

func (c *netDevCollector) metricDesc(key string) *prometheus.Desc {
	c.metricDescsMutex.Lock()
	defer c.metricDescsMutex.Unlock()

	if _, ok := c.metricDescs[key]; !ok {
		c.metricDescs[key] = prometheus.NewDesc(
			prometheus.BuildFQName(namespace, c.subsystem, key+"_total"),
			fmt.Sprintf("Network device statistic %s.", key),
			[]string{"device"},
			nil,
		)
	}

	return c.metricDescs[key]
}

func (c *netDevCollector) Update(ch chan<- prometheus.Metric) error {
	netDev, err := getNetDevStats(&c.deviceFilter)
	if err != nil {
		return fmt.Errorf("couldn't get netstats: %w", err)
	}
	for dev, devStats := range netDev {
		if !netdevDetailedMetrics {
			legacy(devStats)
		}
		for key, value := range devStats {
			desc := c.metricDesc(key)
			ch <- prometheus.MustNewConstMetric(desc, prometheus.CounterValue, float64(value), dev)
		}
	}
	if netdevAddressInfo {
		interfaces, err := net.Interfaces()
		if err != nil {
			return fmt.Errorf("could not get network interfaces: %w", err)
		}

		desc := prometheus.NewDesc(prometheus.BuildFQName(namespace, "network_address",
			"info"), "node network address by device",
			[]string{"device", "address", "netmask", "scope"}, nil)

		for _, addr := range getAddrsInfo(interfaces) {
			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, 1,
				addr.device, addr.addr, addr.netmask, addr.scope)
		}
	}
	return nil
}

type addrInfo struct {
	device  string
	addr    string
	scope   string
	netmask string
}

func scope(ip net.IP) string {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return "link-local"
	}

	if ip.IsInterfaceLocalMulticast() {
		return "interface-local"
	}

	if ip.IsGlobalUnicast() {
		return "global"
	}

	return ""
}

// getAddrsInfo returns interface name, address, scope and netmask for all interfaces.
func getAddrsInfo(interfaces []net.Interface) []addrInfo {
	var res []addrInfo

	for _, ifs := range interfaces {
		addrs, _ := ifs.Addrs()
		for _, addr := range addrs {
			ip, ipNet, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			size, _ := ipNet.Mask.Size()

			res = append(res, addrInfo{
				device:  ifs.Name,
				addr:    ip.String(),
				scope:   scope(ip),
				netmask: strconv.Itoa(size),
			})
		}
	}

	return res
}

// https://github.com/torvalds/linux/blob/master/net/core/net-procfs.c#L75-L97
func legacy(metrics map[string]uint64) {
	if metric, ok := pop(metrics, "receive_errors"); ok {
		metrics["receive_errs"] = metric
	}
	if metric, ok := pop(metrics, "receive_dropped"); ok {
		metrics["receive_drop"] = metric + popz(metrics, "receive_missed_errors")
	}
	if metric, ok := pop(metrics, "receive_fifo_errors"); ok {
		metrics["receive_fifo"] = metric
	}
	if metric, ok := pop(metrics, "receive_frame_errors"); ok {
		metrics["receive_frame"] = metric + popz(metrics, "receive_length_errors") + popz(metrics, "receive_over_errors") + popz(metrics, "receive_crc_errors")
	}
	if metric, ok := pop(metrics, "multicast"); ok {
		metrics["receive_multicast"] = metric
	}
	if metric, ok := pop(metrics, "transmit_errors"); ok {
		metrics["transmit_errs"] = metric
	}
	if metric, ok := pop(metrics, "transmit_dropped"); ok {
		metrics["transmit_drop"] = metric
	}
	if metric, ok := pop(metrics, "transmit_fifo_errors"); ok {
		metrics["transmit_fifo"] = metric
	}
	if metric, ok := pop(metrics, "multicast"); ok {
		metrics["receive_multicast"] = metric
	}
	if metric, ok := pop(metrics, "collisions"); ok {
		metrics["transmit_colls"] = metric
	}
	if metric, ok := pop(metrics, "transmit_carrier_errors"); ok {
		metrics["transmit_carrier"] = metric + popz(metrics, "transmit_aborted_errors") + popz(metrics, "transmit_heartbeat_errors") + popz(metrics, "transmit_window_errors")
	}
}

func pop(m map[string]uint64, key string) (uint64, bool) {
	value, ok := m[key]
	delete(m, key)
	return value, ok
}

func popz(m map[string]uint64, key string) uint64 {
	if value, ok := m[key]; ok {
		delete(m, key)
		return value
	}
	return 0
}
