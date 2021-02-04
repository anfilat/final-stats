package heart

import (
	"context"
	"fmt"

	"github.com/anfilat/final-stats/internal/symo"
)

func initCPU(h *heart) {
	ctx, cancel := context.WithTimeout(h.ctx, timeToGetMetric)
	defer cancel()

	_, err := h.readers.CPU(ctx, true)
	if err != nil {
		h.log.Debug(fmt.Errorf("cannot init cpu: %w", err))
	}
}

func cpu(h *heart, ch <-chan symo.Beat, reader symo.CPU) {
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

				cpu, err := reader(ctx, false)
				if err != nil {
					h.log.Debug(fmt.Errorf("cannot get cpu: %w", err))
					return
				}

				h.pointsMutex.Lock()
				defer h.pointsMutex.Unlock()

				beat.Point.CPU = cpu
			}()
		}
	}
}
