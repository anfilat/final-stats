// +build linux

package loadavg

import (
	"context"
	"strconv"
	"strings"
	"syscall"

	"github.com/anfilat/final-stats/internal/common"
	"github.com/anfilat/final-stats/internal/symo"
)

func Avg(ctx context.Context) (*symo.LoadAvgData, error) {
	stat, err := fileAvg(ctx)
	if err != nil {
		stat, err = sysInfoAvg(ctx)
	}
	return stat, err
}

func fileAvg(_ context.Context) (*symo.LoadAvgData, error) {
	content, err := common.ReadProcFile("loadavg")
	if err != nil {
		return nil, err
	}

	values := strings.Fields(content[0])

	Load1, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return nil, err
	}
	Load5, err := strconv.ParseFloat(values[1], 64)
	if err != nil {
		return nil, err
	}
	Load15, err := strconv.ParseFloat(values[2], 64)
	if err != nil {
		return nil, err
	}

	return &symo.LoadAvgData{
		Load1,
		Load5,
		Load15,
	}, nil
}

func sysInfoAvg(_ context.Context) (*symo.LoadAvgData, error) {
	var si syscall.Sysinfo_t
	err := syscall.Sysinfo(&si)
	if err != nil {
		return nil, err
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
