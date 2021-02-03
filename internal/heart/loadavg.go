package heart

import (
	"context"
	"fmt"

	"github.com/anfilat/final-stats/internal/symo"
)

func loadavg(h *heart, ch <-chan symo.Beat, reader symo.LoadAvg) {
	for {
		select {
		case <-h.ctx.Done():
			return
		case beat, ok := <-ch:
			if !ok {
				return
			}
			func() {
				ctx, cancel := context.WithTimeout(h.ctx, timeToGetMetric)
				defer cancel()

				loadAvg, err := reader(ctx)
				if err != nil {
					h.log.Debug(fmt.Errorf("cannot get load average: %w", err))
					return
				}

				h.pointsMutex.Lock()
				defer h.pointsMutex.Unlock()

				beat.Point.LoadAvg = loadAvg
			}()
		}
	}
}
