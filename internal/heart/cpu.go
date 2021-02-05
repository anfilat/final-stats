package heart

import (
	"context"
	"errors"
	"fmt"

	"github.com/anfilat/final-stats/internal/symo"
)

func initCPU(ctx context.Context, reader symo.CPU, log symo.Logger) {
	workCtx, cancel := context.WithTimeout(ctx, timeToGetMetric)
	defer cancel()

	_, err := reader(workCtx, true)
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return
	}
	if err != nil {
		log.Debug(fmt.Errorf("cannot init cpu: %w", err))
	}
}

func cpu(ctx context.Context, ch <-chan symo.Beat, reader symo.CPU, log symo.Logger) {
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

				data, err := reader(workCtx, false)
				if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
					return
				}
				if err != nil {
					log.Debug(fmt.Errorf("cannot get cpu: %w", err))
					return
				}

				beat.Point.CPU = data
			}()
		}
	}
}
