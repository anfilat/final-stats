package heart

import (
	"context"
	"errors"
	"fmt"

	"github.com/anfilat/final-stats/internal/symo"
)

func loadavg(ctx context.Context, ch <-chan symo.Beat, reader symo.LoadAvg, log symo.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		case beat, ok := <-ch:
			if !ok {
				return
			}
			func() {
				workCtx, cancel := context.WithTimeout(ctx, timeToGetMetric)
				defer cancel()

				data, err := reader(workCtx)
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					return
				}
				if err != nil {
					log.Debug(fmt.Errorf("cannot get load average: %w", err))
					return
				}

				beat.Point.LoadAvg = data
			}()
		}
	}
}
