// +build windows

package loadavg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadAvgWindows(t *testing.T) {
	ctx := context.Background()
	avg1, err := Read(ctx)
	require.NoError(t, err)
	require.NotNil(t, avg1)
}
