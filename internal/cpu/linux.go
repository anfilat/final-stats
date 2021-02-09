// +build linux

package cpu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anfilat/final-stats/internal/common"
)

func getCPU() (*cpuData, error) {
	content, err := common.ReadProcFile("stat")
	if err != nil {
		return nil, fmt.Errorf("cannot read the stat file: %w", err)
	}

	values := strings.Fields(content[0])[1:]

	user, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse user field: %w", err)
	}
	nice, err := strconv.ParseInt(values[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse nice field: %w", err)
	}
	system, err := strconv.ParseInt(values[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse system field: %w", err)
	}
	idle, err := strconv.ParseInt(values[3], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse idle field: %w", err)
	}
	iowait, err := strconv.ParseInt(values[4], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse iowait field: %w", err)
	}
	irq, err := strconv.ParseInt(values[5], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse irq field: %w", err)
	}
	softirq, err := strconv.ParseInt(values[6], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("cannot parse softirq field: %w", err)
	}
	steal := int64(0)
	if len(values) > 7 { // Linux >= 2.6.11
		steal, err = strconv.ParseInt(values[7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse steal field: %w", err)
		}
	}

	return &cpuData{
		total:  float64(user + nice + system + idle + iowait + irq + softirq + steal),
		user:   float64(user),
		system: float64(system),
		idle:   float64(idle),
	}, nil
}
