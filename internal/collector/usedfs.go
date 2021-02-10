//nolint:dupl
package collector

import (
	"context"
	"errors"
	"fmt"

	"github.com/anfilat/final-stats/internal/symo"
)

func startUsedFS(ctx context.Context, reader symo.UsedFS, log symo.Logger) {
	_, err := reader(ctx, symo.StartMetric)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return
	}
	if err != nil {
		log.Debug(fmt.Errorf("cannot start used fs: %w", err))
	}
}

func stopUsedFS(ctx context.Context, reader symo.UsedFS, log symo.Logger) {
	_, err := reader(ctx, symo.StopMetric)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return
	}
	if err != nil {
		log.Debug(fmt.Errorf("cannot stop used fs: %w", err))
	}
}

func usedFS(ctx context.Context, ch <-chan timePoint, reader symo.UsedFS, log symo.Logger) {
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
					log.Debug(fmt.Errorf("cannot get used fs: %w", err))
					return
				}
				tp.point.UsedFS = data
			}()
		}
	}
}
