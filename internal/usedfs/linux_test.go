// +build linux

package usedfs

import (
	"bytes"
	"io/ioutil"
	"strconv"
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
	require.Equal(t, 5, len(data))

	require.Equal(t, "/", data[0].Path)
	require.Equal(t, 75.76559030572227, data[0].UsedSpace)
	require.Equal(t, 23.453521728515625, data[0].UsedInode)

	require.Equal(t, "/mnt/c", data[1].Path)
	require.Equal(t, 78.84516661716783, data[1].UsedSpace)
	require.Equal(t, 4.584831119069495, data[1].UsedInode)
}

func TestUsage(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/data2")
	require.NoError(t, err)

	data, err := usage(common.SplitLines(string(content)))
	require.NoError(t, err)
	require.Equal(t, 5, len(data))

	require.Equal(t, "/", data[0].Path)
	require.Equal(t, 75.76559030572227, data[0].UsedSpace)
	require.Equal(t, 23.453521728515625, data[0].UsedInode)

	require.Equal(t, "/mnt/c", data[1].Path)
	require.Equal(t, 78.84516661716783, data[1].UsedSpace)
	require.Equal(t, 4.584831119069495, data[1].UsedInode)
}

func TestUsageWitEmptyData(t *testing.T) {
	chunk := make([]string, 0, 16)

	data, err := usage(chunk)
	require.NoError(t, err)
	require.Equal(t, 0, len(data))
}

func TestUsageFail(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/data3")
	require.NoError(t, err)

	_, err = usage(common.SplitLines(string(content)))
	require.ErrorIs(t, err, strconv.ErrSyntax)
}

func TestNoUpdateDataAfterStop(t *testing.T) {
	content, err := ioutil.ReadFile("./testdata/data1")
	require.NoError(t, err)

	mutex.Lock()
	fsData = nil
	fsErr = nil
	isLive = false
	mutex.Unlock()

	readOut(bytes.NewBuffer(content))
	data, err := get()
	require.NoError(t, err)
	require.Nil(t, data)
}
