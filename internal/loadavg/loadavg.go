package loadavg

import (
	"context"
	"time"
)

type Data struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

func GetStream(ctx context.Context) <-chan *Data {
	result := make(chan *Data)

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				avg, _ := Avg(ctx)
				result <- avg
			}
		}
	}()

	return result
}
