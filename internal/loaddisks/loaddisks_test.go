package loaddisks

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/anfilat/final-stats/internal/symo"
)

func TestLoadDisk(t *testing.T) {
	ctx := context.Background()

	_, err := Read(ctx, symo.StartMetric)
	require.NoError(t, err)

	time.Sleep(1500 * time.Millisecond)

	data, err := Read(ctx, symo.GetMetric)
	require.NoError(t, err)
	require.NotNil(t, data)
	require.Greater(t, len(data), 0)
	for _, disk := range data {
		require.GreaterOrEqual(t, disk.Tps, 0.0)
		require.GreaterOrEqual(t, disk.KBRead, 0.0)
		require.GreaterOrEqual(t, disk.KBWrite, 0.0)
	}

	_, err = Read(ctx, symo.StopMetric)
	require.NoError(t, err)
}
