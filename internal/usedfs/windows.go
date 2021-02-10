// +build windows

package usedfs

import (
	"context"

	"github.com/anfilat/final-stats/internal/symo"
)

func Collect(_ context.Context, _ symo.MetricCommand) (symo.UsedFSData, error) {
	return nil, nil
}
