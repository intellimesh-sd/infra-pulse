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

//go:build !norunit
// +build !norunit

package collector

import (
	"github.com/prometheus-community/go-runit/runit"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/klog/v2"
)

// var runitServiceDir = flag.String("collector.runit.servicedir", "/etc/service", "Path to runit service directory.")
var runitServiceDir = "/etc/service"

type runitCollector struct {
	state          typedDesc
	stateDesired   typedDesc
	stateNormal    typedDesc
	stateTimestamp typedDesc
}

func init() {
	registerCollector("runit", defaultDisabled, NewRunitCollector)
}

// NewRunitCollector returns a new Collector exposing runit statistics.
func NewRunitCollector() (Collector, error) {
	var (
		subsystem   = "service"
		constLabels = prometheus.Labels{"supervisor": "runit"}
		labelNames  = []string{"service"}
	)

	klog.V(2).Info("msg", "This collector is deprecated and will be removed in the next major version release.")

	return &runitCollector{
		state: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "state"),
			"State of runit service.",
			labelNames, constLabels,
		), prometheus.GaugeValue},
		stateDesired: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "desired_state"),
			"Desired state of runit service.",
			labelNames, constLabels,
		), prometheus.GaugeValue},
		stateNormal: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "normal_state"),
			"Normal state of runit service.",
			labelNames, constLabels,
		), prometheus.GaugeValue},
		stateTimestamp: typedDesc{prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "state_last_change_timestamp_seconds"),
			"Unix timestamp of the last runit service state change.",
			labelNames, constLabels,
		), prometheus.GaugeValue},
	}, nil
}

func (c *runitCollector) Update(ch chan<- prometheus.Metric) error {
	services, err := runit.GetServices(runitServiceDir)
	if err != nil {
		return err
	}

	for _, service := range services {
		status, err := service.Status()
		if err != nil {
			klog.V(3).Info("msg", "Couldn't get status", "service", service.Name, "err", err)
			continue
		}

		klog.V(3).Info("msg", "duration", "service", service.Name, "status", status.State, "pid", status.Pid, "duration_seconds", status.Duration)
		ch <- c.state.mustNewConstMetric(float64(status.State), service.Name)
		ch <- c.stateDesired.mustNewConstMetric(float64(status.Want), service.Name)
		ch <- c.stateTimestamp.mustNewConstMetric(float64(status.Timestamp.Unix()), service.Name)
		if status.NormallyUp {
			ch <- c.stateNormal.mustNewConstMetric(1, service.Name)
		} else {
			ch <- c.stateNormal.mustNewConstMetric(0, service.Name)
		}
	}
	return nil
}
