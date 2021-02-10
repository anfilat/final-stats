// +build linux

package cpu

import (
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anfilat/final-stats/internal/common"
)

func TestParseCPU(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/stat1")
	require.NoError(t, err)

	data, err := parseCPU(common.SplitLines(string(content)))
	require.NoError(t, err)
	require.Equal(t, float64(631341+1284+109025+3744304+11237+1685), data.total)
	require.Equal(t, 631341.0, data.user)
	require.Equal(t, 109025.0, data.system)
	require.Equal(t, 3744304.0, data.idle)
}

func TestParseCPUFail(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/stat2")
	require.NoError(t, err)

	_, err = parseCPU(common.SplitLines(string(content)))
	require.ErrorIs(t, err, strconv.ErrSyntax)
}
