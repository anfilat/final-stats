package collector

import (
	"context"
	"errors"
	"fmt"

	"github.com/anfilat/final-stats/internal/symo"
)

func cpu(ctx context.Context, ch <-chan timePoint, collector symo.CPU, log symo.Logger) {
	for {
		select {
		case <-ctx.Done():
			return
		case tp, ok := <-ch:
			if !ok {
				return
			}
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

				tp.point.CPU = data
			}()
		}
	}
}
