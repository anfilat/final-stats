package cpu

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/anfilat/final-stats/internal/symo"
)

func TestCPU(t *testing.T) {
	ctx := context.Background()

	_, err := Read(ctx, symo.StartMetric)
	require.NoError(t, err)

	time.Sleep(time.Second)

	data, err := Read(ctx, symo.GetMetric)
	require.NoError(t, err)
	require.NotNil(t, data)
	require.GreaterOrEqual(t, data.User, 0.0)
	require.LessOrEqual(t, data.User, 100.0)
	require.GreaterOrEqual(t, data.System, 0.0)
	require.LessOrEqual(t, data.System, 100.0)
	require.GreaterOrEqual(t, data.Idle, 0.0)
	require.LessOrEqual(t, data.Idle, 100.0)
}
