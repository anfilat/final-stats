package loadavg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadAvg(t *testing.T) {
	ctx := context.Background()

	data, err := Collect(ctx)
	require.NoError(t, err)
	require.NotNil(t, data)
	require.GreaterOrEqual(t, data.Load1, 0.0)
	require.GreaterOrEqual(t, data.Load5, 0.0)
	require.GreaterOrEqual(t, data.Load15, 0.0)
}
