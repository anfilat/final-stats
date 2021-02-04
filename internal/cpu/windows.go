// +build windows

package cpu

func getCPU() (*cpuData, error) {
	return &cpuData{
		total:  0,
		user:   0,
		system: 0,
		idle:   0,
	}, nil
}
