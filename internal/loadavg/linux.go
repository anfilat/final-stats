// +build linux

package loadavg

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"syscall"

	"github.com/anfilat/final-stats/internal/common"
	"github.com/anfilat/final-stats/internal/symo"
)

func Read(_ context.Context) (*symo.LoadAvgData, error) {
	stat, err := fileAvg()
	if err != nil {
		return sysInfoAvg()
	}
	return stat, nil
}

func fileAvg() (*symo.LoadAvgData, error) {
	content, err := common.ReadProcFile("loadavg")
	if err != nil {
		return nil, fmt.Errorf("cannot read the loadavg file: %w", err)
	}

	values := strings.Fields(content[0])

	load1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse load1 field: %w", err)
	}
	load5, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse load5 field: %w", err)
	}
	load15, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse load15 field: %w", err)
	}

	return &symo.LoadAvgData{
		Load1:  load1,
		Load5:  load5,
		Load15: load15,
	}, nil
}

func sysInfoAvg() (*symo.LoadAvgData, error) {
	var si syscall.Sysinfo_t
	err := syscall.Sysinfo(&si)
	if err != nil {
		return nil, fmt.Errorf("cannot call syscall.Sysinfo: %w", err)
	}

	return &symo.LoadAvgData{
		Load1:  sysInfoLoadToHuman(si.Loads[0]),
		Load5:  sysInfoLoadToHuman(si.Loads[1]),
		Load15: sysInfoLoadToHuman(si.Loads[2]),
	}, nil
}

func sysInfoLoadToHuman(val uint64) float64 {
	return common.NumToFix2(float64(val) / float64(1<<16))
}
