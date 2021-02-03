package clients

import (
	"time"

	"github.com/anfilat/final-stats/internal/symo"
)

func makeSnapshot(data *symo.ClientsBeat, m int) *symo.Stat {
	from := data.Time.Add(time.Duration(-m) * time.Second)
	count := 0

	stat := &symo.Stat{
		Time: data.Time,
		Stat: &symo.Point{
			LoadAvg: &symo.LoadAvgData{
				Load1:  0,
				Load5:  0,
				Load15: 0,
			},
		},
	}

	for tm, point := range data.Points {
		if tm.Before(from) {
			continue
		}
		count++

		stat.Stat.LoadAvg.Load1 += point.LoadAvg.Load1
		stat.Stat.LoadAvg.Load5 += point.LoadAvg.Load5
		stat.Stat.LoadAvg.Load15 += point.LoadAvg.Load15
	}

	if count > 1 {
		stat.Stat.LoadAvg.Load1 /= float64(count)
		stat.Stat.LoadAvg.Load5 /= float64(count)
		stat.Stat.LoadAvg.Load15 /= float64(count)
	}

	return stat
}
