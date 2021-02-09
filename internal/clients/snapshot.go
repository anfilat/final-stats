package clients

import (
	"fmt"
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

func makeSnapshot(data *symo.MetricsData, m int) *symo.Stats {
	from := data.Time.Add(time.Duration(-m) * time.Second)

	countLoad := 0
	load1 := 0.0
	load5 := 0.0
	load15 := 0.0

	countCPU := 0
	cpuUser := 0.0
	cpuSystem := 0.0
	cpiIdle := 0.0

	type LoadDisk struct {
		count   int
		tps     float64
		kbRead  float64
		kbWrite float64
	}
	loadDisks := make(map[string]*LoadDisk, len(data.Points))

	for tm, point := range data.Points {
		if tm.Before(from) {
			continue
		}

		if point.LoadAvg != nil {
			countLoad++
			load1 += point.LoadAvg.Load1
			load5 += point.LoadAvg.Load5
			load15 += point.LoadAvg.Load15
		}
		if point.CPU != nil {
			countCPU++
			cpuUser += point.CPU.User
			cpuSystem += point.CPU.System
			cpiIdle += point.CPU.Idle
		}
		if point.LoadDisks != nil {
			for _, diskData := range point.LoadDisks {
				fmt.Println(diskData)
				data := loadDisks[diskData.Name]
				if data == nil {
					data = &LoadDisk{}
					loadDisks[diskData.Name] = data
				}
				data.count++
				data.tps += diskData.Tps
				data.kbRead += diskData.KBRead
				data.kbWrite += diskData.KBWrite
			}
		}
	}

	result := &symo.Stats{
		Time: data.Time,
	}

	if countLoad > 1 {
		result.LoadAvg = &symo.LoadAvgData{
			Load1:  load1 / float64(countLoad),
			Load5:  load5 / float64(countLoad),
			Load15: load15 / float64(countLoad),
		}
	}
	if countCPU > 1 {
		result.CPU = &symo.CPUData{
			User:   cpuUser / float64(countCPU),
			System: cpuSystem / float64(countCPU),
			Idle:   cpiIdle / float64(countCPU),
		}
	}
	if len(loadDisks) > 0 {
		for name, data := range loadDisks {
			result.LoadDisks = append(result.LoadDisks, symo.DiskData{
				Name:    name,
				Tps:     data.tps / float64(data.count),
				KBRead:  data.kbRead / float64(data.count),
				KBWrite: data.kbWrite / float64(data.count),
			})
		}
	}

	return result
}
