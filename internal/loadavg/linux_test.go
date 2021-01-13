package loadavg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEqualityLinuxMethods(t *testing.T) {
	ctx := context.Background()
	avg1, err := fileAvg(ctx)
	require.NoError(t, err)
	avg2, err := sysInfoAvg(ctx)
	require.NoError(t, err)
	require.Equal(t, &avg1, &avg2)
}
