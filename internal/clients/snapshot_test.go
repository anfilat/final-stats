package clients

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/anfilat/final-stats/internal/symo"
)

var (
	la1 = symo.LoadAvgData{
		Load1:  0.1,
		Load5:  0.2,
		Load15: 0.3,
	}
	la2 = symo.LoadAvgData{
		Load1:  0.2,
		Load5:  0.3,
		Load15: 0.4,
	}
	laSum12 = symo.LoadAvgData{
		Load1:  0.15,
		Load5:  0.25,
		Load15: 0.35,
	}
)

var (
	cpu1 = symo.CPUData{
		User:   10,
		System: 20,
		Idle:   30,
	}
	cpu2 = symo.CPUData{
		User:   20,
		System: 30,
		Idle:   40,
	}
	cpuSum12 = symo.CPUData{
		User:   15,
		System: 25,
		Idle:   35,
	}
)

var (
	ld1 = symo.LoadDisksData{{
		Name:    "sda",
		Tps:     100,
		KBRead:  200,
		KBWrite: 300,
	}, {
		Name:    "sdb",
		Tps:     5,
		KBRead:  6,
		KBWrite: 7,
	}}
	ld2 = symo.LoadDisksData{{
		Name:    "sda",
		Tps:     200,
		KBRead:  300,
		KBWrite: 400,
	}, {
		Name:    "sdb",
		Tps:     6,
		KBRead:  7,
		KBWrite: 8,
	}}
	ld3 = symo.LoadDisksData{{
		Name:    "sda",
		Tps:     200,
		KBRead:  300,
		KBWrite: 400,
	}, {
		Name:    "sdb",
		Tps:     6,
		KBRead:  7,
		KBWrite: 8,
	}, {
		Name:    "sdc",
		Tps:     7,
		KBRead:  7,
		KBWrite: 7,
	}}
	ldSum12 = symo.LoadDisksData{{
		Name:    "sda",
		Tps:     150,
		KBRead:  250,
		KBWrite: 350,
	}, {
		Name:    "sdb",
		Tps:     5.5,
		KBRead:  6.5,
		KBWrite: 7.5,
	}}
	ldSum13 = symo.LoadDisksData{{
		Name:    "sda",
		Tps:     150,
		KBRead:  250,
		KBWrite: 350,
	}, {
		Name:    "sdb",
		Tps:     5.5,
		KBRead:  6.5,
		KBWrite: 7.5,
	}, {
		Name:    "sdc",
		Tps:     7,
		KBRead:  7,
		KBWrite: 7,
	}}
)

var (
	fs1 = symo.UsedFSData{{
		Path:      "/",
		UsedSpace: 20,
		UsedInode: 70,
	}, {
		Path:      "/data",
		UsedSpace: 3,
		UsedInode: 6,
	}}
	fs2 = symo.UsedFSData{{
		Path:      "/",
		UsedSpace: 30,
		UsedInode: 80,
	}, {
		Path:      "/data",
		UsedSpace: 4,
		UsedInode: 7,
	}}
	fs3 = symo.UsedFSData{{
		Path:      "/",
		UsedSpace: 30,
		UsedInode: 80,
	}, {
		Path:      "/data",
		UsedSpace: 4,
		UsedInode: 7,
	}, {
		Path:      "/mount/c",
		UsedSpace: 5,
		UsedInode: 13,
	}}
	fsSum12 = symo.UsedFSData{{
		Path:      "/",
		UsedSpace: 25,
		UsedInode: 75,
	}, {
		Path:      "/data",
		UsedSpace: 3.5,
		UsedInode: 6.5,
	}}
	fsSum13 = symo.UsedFSData{{
		Path:      "/",
		UsedSpace: 25,
		UsedInode: 75,
	}, {
		Path:      "/data",
		UsedSpace: 3.5,
		UsedInode: 6.5,
	}, {
		Path:      "/mount/c",
		UsedSpace: 5,
		UsedInode: 13,
	}}
)

//nolint:funlen
func TestSnapshot(t *testing.T) {
	now := time.Now().Truncate(time.Second)
	tests := []struct {
		name     string
		data     *symo.MetricsData
		m        int
		expected *symo.Stats
	}{
		{
			name: "with data",
			data: &symo.MetricsData{
				Time: now,
				Points: symo.Points{
					now.Add(-time.Second): {
						LoadAvg:   &la1,
						CPU:       &cpu1,
						LoadDisks: ld1,
						UsedFS:    fs1,
					},
					now.Add(-2 * time.Second): {
						LoadAvg:   &la2,
						CPU:       &cpu2,
						LoadDisks: ld2,
						UsedFS:    fs2,
					},
				},
			},
			m: 2,
			expected: &symo.Stats{
				Time:      now,
				LoadAvg:   &laSum12,
				CPU:       &cpuSum12,
				LoadDisks: ldSum12,
				UsedFS:    fsSum12,
			},
		},
		{
			name: "with old points",
			data: &symo.MetricsData{
				Time: now,
				Points: symo.Points{
					now.Add(-time.Second): {
						LoadAvg:   &la1,
						CPU:       &cpu1,
						LoadDisks: ld1,
						UsedFS:    fs1,
					},
					now.Add(-2 * time.Second): {
						LoadAvg:   &la2,
						CPU:       &cpu2,
						LoadDisks: ld2,
						UsedFS:    fs2,
					},
					// эта секунда не попадает в интервал усреднения
					now.Add(-3 * time.Second): {
						LoadAvg:   &la2,
						CPU:       &cpu2,
						LoadDisks: ld2,
						UsedFS:    fs2,
					},
				},
			},
			m: 2,
			expected: &symo.Stats{
				Time:      now,
				LoadAvg:   &laSum12,
				CPU:       &cpuSum12,
				LoadDisks: ldSum12,
				UsedFS:    fsSum12,
			},
		},
		{
			name: "with only old points",
			data: &symo.MetricsData{
				Time: now,
				Points: symo.Points{
					now.Add(-10 * time.Second): {
						LoadAvg:   &la1,
						CPU:       &cpu1,
						LoadDisks: ld1,
						UsedFS:    fs1,
					},
					now.Add(-11 * time.Second): {
						LoadAvg:   &la2,
						CPU:       &cpu2,
						LoadDisks: ld2,
						UsedFS:    fs2,
					},
				},
			},
			m: 2,
			expected: &symo.Stats{
				Time:      now,
				LoadAvg:   nil,
				CPU:       nil,
				LoadDisks: nil,
				UsedFS:    nil,
			},
		},
		{
			name: "with changed set of disks",
			data: &symo.MetricsData{
				Time: now,
				Points: symo.Points{
					now.Add(-time.Second): {
						LoadAvg:   &la1,
						CPU:       &cpu1,
						LoadDisks: ld1,
						UsedFS:    fs1,
					},
					now.Add(-2 * time.Second): {
						LoadAvg:   &la2,
						CPU:       &cpu2,
						LoadDisks: ld3,
						UsedFS:    fs3,
					},
				},
			},
			m: 5,
			expected: &symo.Stats{
				Time:      now,
				LoadAvg:   &laSum12,
				CPU:       &cpuSum12,
				LoadDisks: ldSum13,
				UsedFS:    fsSum13,
			},
		},
		{
			name: "without data",
			data: &symo.MetricsData{
				Time:   now,
				Points: nil,
			},
			m: 5,
			expected: &symo.Stats{
				Time:      now,
				LoadAvg:   nil,
				CPU:       nil,
				LoadDisks: nil,
				UsedFS:    nil,
			},
		},
		{
			name: "with empty point",
			data: &symo.MetricsData{
				Time: now,
				Points: symo.Points{
					now.Add(-time.Second): {
						LoadAvg:   nil,
						CPU:       nil,
						LoadDisks: nil,
						UsedFS:    nil,
					},
				},
			},
			m: 5,
			expected: &symo.Stats{
				Time:      now,
				LoadAvg:   nil,
				CPU:       nil,
				LoadDisks: nil,
				UsedFS:    nil,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			stats := makeSnapshot(tt.data, tt.m)

			require.True(t, tt.expected.Time.Equal(stats.Time))

			if tt.expected.LoadAvg == nil {
				require.Nil(t, stats.LoadAvg)
			} else {
				require.NotNil(t, stats.LoadAvg)
				require.InEpsilon(t, tt.expected.LoadAvg.Load1, stats.LoadAvg.Load1, 0.001)
				require.InEpsilon(t, tt.expected.LoadAvg.Load5, stats.LoadAvg.Load5, 0.001)
				require.InEpsilon(t, tt.expected.LoadAvg.Load15, stats.LoadAvg.Load15, 0.001)
			}

			if tt.expected.CPU == nil {
				require.Nil(t, stats.CPU)
			} else {
				require.NotNil(t, stats.CPU)
				require.InEpsilon(t, tt.expected.CPU.User, stats.CPU.User, 0.001)
				require.InEpsilon(t, tt.expected.CPU.System, stats.CPU.System, 0.001)
				require.InEpsilon(t, tt.expected.CPU.Idle, stats.CPU.Idle, 0.001)
			}

			require.Len(t, stats.LoadDisks, len(tt.expected.LoadDisks))
			for i := 0; i < len(tt.expected.LoadDisks); i++ {
				expectedLd := tt.expected.LoadDisks[i]
				ld := findLoadDisk(expectedLd.Name, stats.LoadDisks)
				require.NotNilf(t, ld, "disk %s not found", expectedLd.Name)
				require.InEpsilon(t, expectedLd.Tps, ld.Tps, 0.001)
				require.InEpsilon(t, expectedLd.KBRead, ld.KBRead, 0.001)
				require.InEpsilon(t, expectedLd.KBWrite, ld.KBWrite, 0.001)
			}

			require.Len(t, stats.UsedFS, len(tt.expected.UsedFS))
			for i := 0; i < len(tt.expected.UsedFS); i++ {
				expectedFs := tt.expected.UsedFS[i]
				fs := findUsedFS(expectedFs.Path, stats.UsedFS)
				require.NotNilf(t, fs, "fs %s not found", expectedFs.Path)
				require.InEpsilon(t, expectedFs.UsedSpace, fs.UsedSpace, 0.001)
				require.InEpsilon(t, expectedFs.UsedInode, fs.UsedInode, 0.001)
			}
		})
	}
}

func findLoadDisk(name string, ld symo.LoadDisksData) *symo.DiskData {
	for i := 0; i < len(ld); i++ {
		if ld[i].Name == name {
			return &ld[i]
		}
	}
	return nil
}

func findUsedFS(path string, fs symo.UsedFSData) *symo.FSData {
	for i := 0; i < len(fs); i++ {
		if fs[i].Path == path {
			return &fs[i]
		}
	}
	return nil
}
