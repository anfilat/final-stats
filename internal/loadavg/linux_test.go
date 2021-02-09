// +build linux

package loadavg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadAvgLinuxEqualityMethods(t *testing.T) {
	avg1, err := fileAvg()
	require.NoError(t, err)
	avg2, err := sysInfoAvg()
	require.NoError(t, err)
	require.Equal(t, avg1, avg2)
}
