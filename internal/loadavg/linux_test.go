// +build linux

package loadavg

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anfilat/final-stats/internal/common"
)

func TestParseLoadAvg(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/loadavg1")
	require.NoError(t, err)

	data, err := parseLoadAvg(common.SplitLines(string(content)))
	require.NoError(t, err)
	require.Equal(t, 0.60, data.Load1)
	require.Equal(t, 0.75, data.Load5)
	require.Equal(t, 0.74, data.Load15)
}

func TestParseLoadAvgFail(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/loadavg2")
	require.NoError(t, err)

	_, err = parseLoadAvg(common.SplitLines(string(content)))
	require.ErrorIs(t, err, strconv.ErrSyntax)
}
