// Copyright 2020 The Prometheus Authors
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

//go:build !nofilesystem
// +build !nofilesystem

package collector

import (
	"golang.org/x/sys/unix"
)

const (
	defMountPointsExcluded = "^/(dev)($|/)"
	defFSTypesExcluded     = "^devfs$"
)

// Expose filesystem fullness.
func (c *filesystemCollector) GetStats() (stats []filesystemStats, err error) {
	var mnt []unix.Statfs_t
	size, err := unix.Getfsstat(mnt, unix.MNT_NOWAIT)
	if err != nil {
		return nil, err
	}
	mnt = make([]unix.Statfs_t, size)
	_, err = unix.Getfsstat(mnt, unix.MNT_NOWAIT)
	if err != nil {
		return nil, err
	}

	stats = []filesystemStats{}
	for _, v := range mnt {
		mountpoint := unix.ByteSliceToString(v.F_mntonname[:])
		if c.excludedMountPointsPattern.MatchString(mountpoint) {
			klog.V(3).Info("msg", "Ignoring mount point", "mountpoint", mountpoint)
			continue
		}

		device := unix.ByteSliceToString(v.F_mntfromname[:])
		fstype := unix.ByteSliceToString(v.F_fstypename[:])
		if c.excludedFSTypesPattern.MatchString(fstype) {
			klog.V(3).Info("msg", "Ignoring fs type", "type", fstype)
			continue
		}

		var ro float64
		if (v.F_flags & unix.MNT_RDONLY) != 0 {
			ro = 1
		}

		stats = append(stats, filesystemStats{
			labels: filesystemLabels{
				device:     device,
				mountPoint: mountpoint,
				fsType:     fstype,
			},
			size:      float64(v.F_blocks) * float64(v.F_bsize),
			free:      float64(v.F_bfree) * float64(v.F_bsize),
			avail:     float64(v.F_bavail) * float64(v.F_bsize),
			files:     float64(v.F_files),
			filesFree: float64(v.F_ffree),
			ro:        ro,
		})
	}
	return stats, nil
}
