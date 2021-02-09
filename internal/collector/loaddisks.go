package collector

import (
	"context"
	"errors"
	"fmt"

	"github.com/anfilat/final-stats/internal/symo"
)

func startLoadDisks(ctx context.Context, reader symo.LoadDisks, log symo.Logger) {
	_, err := reader(ctx, symo.StartMetric)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return
	}
	if err != nil {
		log.Debug(fmt.Errorf("cannot start load disks: %w", err))
	}
}

func stopLoadDisks(ctx context.Context, reader symo.LoadDisks, log symo.Logger) {
	_, err := reader(ctx, symo.StopMetric)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return
	}
	if err != nil {
		log.Debug(fmt.Errorf("cannot stop load disks: %w", err))
	}
}

func loadDisks(ctx context.Context, ch <-chan timePoint, reader symo.LoadDisks, log symo.Logger) {
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

				data, err := reader(workCtx, symo.GetMetric)
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					return
				}
				if err != nil {
					log.Debug(fmt.Errorf("cannot get load disks: %w", err))
					return
				}
				tp.point.LoadDisks = data
			}()
		}
	}
}
