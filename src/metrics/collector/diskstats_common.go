// Copyright 2019 The Prometheus Authors
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

//go:build !nodiskstats && (openbsd || linux || darwin)
// +build !nodiskstats
// +build openbsd linux darwin

package collector

import (
	"errors"
	"k8s.io/klog/v2"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	diskSubsystem = "disk"
)

var (
	diskLabelNames = []string{"device"}

	diskstatsDeviceExcludeSet bool
	/*	diskstatsDeviceExclude    = flag.String(
			"collector.diskstats.device-exclude", diskstatsDefaultIgnoredDevices,
			"Regexp of diskstats devices to exclude (mutually exclusive to device-include).",
		)
		oldDiskstatsDeviceExclude = flag.String(
			"collector.diskstats.ignored-devices", "",
			"DEPRECATED: Use collector.diskstats.device-exclude")

		diskstatsDeviceInclude = flag.String("collector.diskstats.device-include", "",
			"Regexp of diskstats devices to include (mutually exclusive to device-exclude).")
	*/

	diskstatsDeviceExclude    = diskstatsDefaultIgnoredDevices
	oldDiskstatsDeviceExclude = ""

	diskstatsDeviceInclude = ""

	readsCompletedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "reads_completed_total"),
		"The total number of reads completed successfully.",
		diskLabelNames, nil,
	)

	readBytesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "read_bytes_total"),
		"The total number of bytes read successfully.",
		diskLabelNames, nil,
	)

	writesCompletedDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "writes_completed_total"),
		"The total number of writes completed successfully.",
		diskLabelNames, nil,
	)

	writtenBytesDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "written_bytes_total"),
		"The total number of bytes written successfully.",
		diskLabelNames, nil,
	)

	ioTimeSecondsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "io_time_seconds_total"),
		"Total seconds spent doing I/Os.",
		diskLabelNames, nil,
	)

	readTimeSecondsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "read_time_seconds_total"),
		"The total number of seconds spent by all reads.",
		diskLabelNames,
		nil,
	)

	writeTimeSecondsDesc = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, diskSubsystem, "write_time_seconds_total"),
		"This is the total number of seconds spent by all writes.",
		diskLabelNames,
		nil,
	)
)

func newDiskstatsDeviceFilter() (deviceFilter, error) {
	if oldDiskstatsDeviceExclude != "" {
		if !diskstatsDeviceExcludeSet {
			klog.V(2).Info("msg", "--collector.diskstats.ignored-devices is DEPRECATED and will be removed in 2.0.0, use --collector.diskstats.device-exclude")
			diskstatsDeviceExclude = oldDiskstatsDeviceExclude
		} else {
			return deviceFilter{}, errors.New("--collector.diskstats.ignored-devices and --collector.diskstats.device-exclude are mutually exclusive")
		}
	}

	if diskstatsDeviceExclude != "" && diskstatsDeviceInclude != "" {
		return deviceFilter{}, errors.New("device-exclude & device-include are mutually exclusive")
	}

	if diskstatsDeviceExclude != "" {
		klog.Info("msg", "Parsed flag --collector.diskstats.device-exclude", "flag", diskstatsDeviceExclude)
	}

	if diskstatsDeviceInclude != "" {
		klog.Info("msg", "Parsed Flag --collector.diskstats.device-include", "flag", diskstatsDeviceInclude)
	}

	return newDeviceFilter(diskstatsDeviceExclude, diskstatsDeviceInclude), nil
}
