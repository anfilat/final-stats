package clients

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/anfilat/final-stats/internal/pb"
	"github.com/anfilat/final-stats/internal/symo"
)

func makeSnapshot(data *symo.ClientsBeat, m int) *pb.Stats {
	from := data.Time.Add(time.Duration(-m) * time.Second)
	count := 0

	load1 := 0.0
	load5 := 0.0
	load15 := 0.0

	for tm, point := range data.Points {
		if tm.Before(from) {
			continue
		}
		count++

		load1 += point.LoadAvg.Load1
		load5 += point.LoadAvg.Load5
		load15 += point.LoadAvg.Load15
	}

	if count > 1 {
		load1 /= float64(count)
		load5 /= float64(count)
		load15 /= float64(count)
	}

	return &pb.Stats{
		Time: timestamppb.New(data.Time),
		LoadAvg: &pb.LoadAvg{
			Load1:  load1,
			Load5:  load5,
			Load15: load15,
		},
	}
}
