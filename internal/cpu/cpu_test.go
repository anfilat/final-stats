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
	cpu, err := Read(ctx, symo.GetMetric)
	require.NoError(t, err)
	require.GreaterOrEqual(t, cpu.User, 0.0)
	require.LessOrEqual(t, cpu.User, 100.0)
	require.GreaterOrEqual(t, cpu.System, 0.0)
	require.LessOrEqual(t, cpu.System, 100.0)
	require.GreaterOrEqual(t, cpu.Idle, 0.0)
	require.LessOrEqual(t, cpu.Idle, 100.0)
}
