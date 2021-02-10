// +build linux

package loaddisks

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anfilat/final-stats/internal/common"
)

func TestReadOut(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/data1")
	require.NoError(t, err)

	isLive = true

	readOut(bytes.NewBuffer(content))
	data, err := get()
	require.NoError(t, err)
	require.Equal(t, 2, len(data))

	require.Equal(t, "sda", data[0].Name)
	require.Equal(t, 79.0, data[0].Tps)
	require.Equal(t, 1816.0, data[0].KBRead)
	require.Equal(t, 0.0, data[0].KBWrite)

	require.Equal(t, "sdb", data[1].Name)
	require.Equal(t, 17.0, data[1].Tps)
	require.Equal(t, 200.0, data[1].KBRead)
	require.Equal(t, 12.0, data[1].KBWrite)
}

func TestCounters(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/data2")
	require.NoError(t, err)

	data, err := counters(common.SplitLines(string(content)))
	require.NoError(t, err)
	require.Equal(t, 2, len(data))

	require.Equal(t, "sda", data[0].Name)
	require.Equal(t, 79.0, data[0].Tps)
	require.Equal(t, 1816.0, data[0].KBRead)
	require.Equal(t, 0.0, data[0].KBWrite)

	require.Equal(t, "sdb", data[1].Name)
	require.Equal(t, 17.0, data[1].Tps)
	require.Equal(t, 200.0, data[1].KBRead)
	require.Equal(t, 12.0, data[1].KBWrite)
}

func TestCountersWitEmptyData(t *testing.T) {
	chunk := make([]string, 0, 16)

	data, err := counters(chunk)
	require.NoError(t, err)
	require.Equal(t, 0, len(data))
}

func TestCountersFail(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/data3")
	require.NoError(t, err)

	_, err = counters(common.SplitLines(string(content)))
	require.Error(t, err)
}

func TestNoUpdateDataAfterStop(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/data1")
	require.NoError(t, err)

	mutex.Lock()
	loadDiskData = nil
	loadDiskErr = nil
	isLive = false
	mutex.Unlock()

	readOut(bytes.NewBuffer(content))
	data, err := get()
	require.NoError(t, err)
	require.Nil(t, data)
}
