package clients

import (
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

func makeSnapshot(data *symo.MetricsData, m int) *symo.Stats {
	result := &symo.Stats{
		Time: data.Time,
	}

	from := data.Time.Add(time.Duration(-m) * time.Second)
	points := make([]*symo.Point, 0, len(data.Points))

	for tm, point := range data.Points {
		if tm.Before(from) {
			continue
		}
		points = append(points, point)
	}

	fillLoadAvg(result, points)
	fillCPU(result, points)
	fillLoadDisks(result, points)
	fillUsedFS(result, points)

	return result
}

func fillLoadAvg(result *symo.Stats, points []*symo.Point) {
	countLoad := 0
	load1 := 0.0
	load5 := 0.0
	load15 := 0.0

	for _, point := range points {
		if point.LoadAvg != nil {
			countLoad++
			load1 += point.LoadAvg.Load1
			load5 += point.LoadAvg.Load5
			load15 += point.LoadAvg.Load15
		}
	}

	if countLoad > 0 {
		result.LoadAvg = &symo.LoadAvgData{
			Load1:  load1 / float64(countLoad),
			Load5:  load5 / float64(countLoad),
			Load15: load15 / float64(countLoad),
		}
	}
}

func fillCPU(result *symo.Stats, points []*symo.Point) {
	countCPU := 0
	cpuUser := 0.0
	cpuSystem := 0.0
	cpiIdle := 0.0

	for _, point := range points {
		if point.CPU != nil {
			countCPU++
			cpuUser += point.CPU.User
			cpuSystem += point.CPU.System
			cpiIdle += point.CPU.Idle
		}
	}

	if countCPU > 0 {
		result.CPU = &symo.CPUData{
			User:   cpuUser / float64(countCPU),
			System: cpuSystem / float64(countCPU),
			Idle:   cpiIdle / float64(countCPU),
		}
	}
}

func fillLoadDisks(result *symo.Stats, points []*symo.Point) {
	type loadDisk struct {
		count   int
		tps     float64
		kbRead  float64
		kbWrite float64
	}
	disks := make(map[string]*loadDisk, len(points))

	for _, point := range points {
		if point.LoadDisks != nil {
			for _, diskData := range point.LoadDisks {
				data := disks[diskData.Name]
				if data == nil {
					data = &loadDisk{}
					disks[diskData.Name] = data
				}
				data.count++
				data.tps += diskData.Tps
				data.kbRead += diskData.KBRead
				data.kbWrite += diskData.KBWrite
			}
		}
	}

	if len(disks) > 0 {
		for name, data := range disks {
			result.LoadDisks = append(result.LoadDisks, symo.DiskData{
				Name:    name,
				Tps:     data.tps / float64(data.count),
				KBRead:  data.kbRead / float64(data.count),
				KBWrite: data.kbWrite / float64(data.count),
			})
		}
	}
}

func fillUsedFS(result *symo.Stats, points []*symo.Point) {
	type usedFS struct {
		count     int
		usedSpace float64
		usedInode float64
	}
	fss := make(map[string]*usedFS, len(points))

	for _, point := range points {
		if point.UsedFS != nil {
			for _, fsData := range point.UsedFS {
				data := fss[fsData.Path]
				if data == nil {
					data = &usedFS{}
					fss[fsData.Path] = data
				}
				data.count++
				data.usedSpace += fsData.UsedSpace
				data.usedInode += fsData.UsedInode
			}
		}
	}

	if len(fss) > 0 {
		for path, data := range fss {
			result.UsedFS = append(result.UsedFS, symo.FSData{
				Path:      path,
				UsedSpace: data.usedSpace / float64(data.count),
				UsedInode: data.usedInode / float64(data.count),
			})
		}
	}
}
