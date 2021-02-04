// +build linux

package cpu

import (
	"strconv"
	"strings"

	"github.com/anfilat/final-stats/internal/common"
)

func getCPU() (*cpuData, error) {
	content, err := common.ReadProcFile("stat")
	if err != nil {
		return nil, err
	}

	values := strings.Fields(content[0])[1:]

	user, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return nil, err
	}
	nice, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		return nil, err
	}
	system, err := strconv.ParseInt(values[2], 10, 64)
	if err != nil {
		return nil, err
	}
	idle, err := strconv.ParseInt(values[3], 10, 64)
	if err != nil {
		return nil, err
	}
	iowait, err := strconv.ParseInt(values[4], 10, 64)
	if err != nil {
		return nil, err
	}
	irq, err := strconv.ParseInt(values[5], 10, 64)
	if err != nil {
		return nil, err
	}
	softirq, err := strconv.ParseInt(values[6], 10, 64)
	if err != nil {
		return nil, err
	}
	steal, err := strconv.ParseInt(values[7], 10, 64)
	if err != nil {
		return nil, err
	}

	return &cpuData{
		total:  float64(user + nice + system + idle + iowait + irq + softirq + steal),
		user:   float64(user),
		system: float64(system),
		idle:   float64(idle),
	}, nil
}
