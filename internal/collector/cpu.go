package collector

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/anfilat/final-stats/internal/symo"
)

func cpuCollect(ctx context.Context, mutex sync.Locker, ch <-chan timePoint, collector symo.CPU, log symo.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		case tp := <-ch:
			func() {
				workCtx, cancel := context.WithTimeout(ctx, timeToGetMetric)
				defer cancel()

				data, err := collector(workCtx, symo.GetMetric)
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					return
				}
				if err != nil {
					log.Debug(fmt.Errorf("cannot get cpu: %w", err))
					return
				}

				mutex.Lock()
				defer mutex.Unlock()

				tp.point.CPU = data
			}()
		}
	}
}
