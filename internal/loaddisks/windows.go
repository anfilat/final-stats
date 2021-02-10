// +build windows

package loaddisks

import (
	"context"

	"github.com/anfilat/final-stats/internal/symo"
)

func Collect(_ context.Context, _ symo.MetricCommand) (symo.LoadDisksData, error) {
	return nil, nil
}
