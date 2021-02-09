package clients

import (
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

func makeSnapshot(data *symo.MetricsData, m int) *symo.Stats {
	from := data.Time.Add(time.Duration(-m) * time.Second)
	count := 0

	load1 := 0.0
	load5 := 0.0
	load15 := 0.0
	cpuUser := 0.0
	cpuSystem := 0.0
	cpiIdle := 0.0

	for tm, point := range data.Points {
		if tm.Before(from) {
			continue
		}
		count++

		load1 += point.LoadAvg.Load1
		load5 += point.LoadAvg.Load5
		load15 += point.LoadAvg.Load15
		cpuUser += point.CPU.User
		cpuSystem += point.CPU.System
		cpiIdle += point.CPU.Idle
	}

	if count > 1 {
		load1 /= float64(count)
		load5 /= float64(count)
		load15 /= float64(count)
		cpuUser /= float64(count)
		cpuSystem /= float64(count)
		cpiIdle /= float64(count)
	}

	return &symo.Stats{
		Time: data.Time,
		LoadAvg: &symo.LoadAvgData{
			Load1:  load1,
			Load5:  load5,
			Load15: load15,
		},
		CPU: &symo.CPUData{
			User:   cpuUser,
			System: cpuSystem,
			Idle:   cpiIdle,
		},
	}
}
