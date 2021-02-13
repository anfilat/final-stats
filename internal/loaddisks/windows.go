// +build windows

// based on https://github.com/shirou/gopsutil
package loaddisks

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/anfilat/final-stats/internal/symo"
)

var (
	mutex    sync.Mutex
	prevData []diskStat
)

func Collect(_ context.Context, action symo.MetricCommand) (symo.LoadDisksData, error) {
	switch action {
	case symo.StartMetric:
		return nil, start()
	case symo.StopMetric:
		return nil, nil
	default:
		return get()
	}
}

func start() error {
	data, err := getLoadDisks()
	if err != nil {
		return err
	}

	mutex.Lock()
	defer mutex.Unlock()

	prevData = data

	return nil
}

func get() (symo.LoadDisksData, error) {
	data, err := getLoadDisks()
	if err != nil {
		return nil, err
	}

	mutex.Lock()
	defer mutex.Unlock()

	if prevData == nil {
		prevData = data
		return nil, nil
	}

	result := make(symo.LoadDisksData, 0, len(data))
	for _, diskData := range data {
		prev := findLoadDisk(diskData.name, prevData)
		if prev != nil {
			diskData.bytesRead -= prev.bytesRead
			diskData.bytesWrite -= prev.bytesWrite
		}
		result = append(result, symo.DiskData{
			Name:    diskData.name,
			Tps:     0, // непонятно, как считать
			KBRead:  float64(diskData.bytesRead) / 1024,
			KBWrite: float64(diskData.bytesWrite) / 1024,
		})
	}
	prevData = data
	return result, nil
}

func findLoadDisk(name string, ld []diskStat) *diskStat {
	for i := 0; i < len(ld); i++ {
		if ld[i].name == name {
			return &ld[i]
		}
	}
	return nil
}

func getLoadDisks() ([]diskStat, error) {
	result := make([]diskStat, 0, 5)

	lpBuffer := make([]uint16, 254)
	lpBufferLen, err := windows.GetLogicalDriveStrings(uint32(len(lpBuffer)), &lpBuffer[0])
	if err != nil {
		return result, fmt.Errorf("GetLogicalDriveStrings error: %w", err)
	}

	for _, v := range lpBuffer[:lpBufferLen] {
		if v < 'A' || v > 'Z' {
			continue
		}
		path := string(rune(v)) + ":"
		typePath, _ := windows.UTF16PtrFromString(path)
		typeRet := windows.GetDriveType(typePath)
		if typeRet == 0 {
			return result, fmt.Errorf("call GetDriveType error: %w", windows.GetLastError())
		}
		if typeRet != windows.DRIVE_FIXED {
			continue
		}

		data, err := driveStat(path)
		if err != nil {
			return result, err
		}
		if data == nil {
			continue
		}
		result = append(result, *data)
	}

	return result, nil
}

func driveStat(path string) (*diskStat, error) {
	szDevice := fmt.Sprintf(`\\.\%s`, path)
	name, err := syscall.UTF16PtrFromString(szDevice)
	if err != nil {
		return nil, fmt.Errorf("call syscall.UTF16PtrFromString error: %w", err)
	}
	h, err := windows.CreateFile(name, 0, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil, windows.OPEN_EXISTING, 0, 0)
	if err != nil {
		if errors.Is(err, windows.ERROR_FILE_NOT_FOUND) {
			return nil, nil
		}
		return nil, fmt.Errorf("call windows.CreateFile error: %w", err)
	}
	defer func() {
		_ = windows.CloseHandle(h)
	}()

	//nolint:golint,stylecheck
	const IOCTL_DISK_PERFORMANCE = 0x70020
	var dp diskPerformance
	var dpSize uint32
	err = windows.DeviceIoControl(h, IOCTL_DISK_PERFORMANCE, nil, 0, (*byte)(unsafe.Pointer(&dp)), uint32(unsafe.Sizeof(dp)), &dpSize, nil)
	if err != nil {
		return nil, fmt.Errorf("call windows.DeviceIoControl error: %w", err)
	}

	return &diskStat{
		name:       path,
		bytesRead:  dp.BytesRead,
		bytesWrite: dp.BytesWritten,
	}, nil
}

type diskPerformance struct {
	BytesRead           int64
	BytesWritten        int64
	ReadTime            int64
	WriteTime           int64
	IdleTime            int64
	ReadCount           uint32
	WriteCount          uint32
	QueueDepth          uint32
	SplitCount          uint32
	QueryTime           int64
	StorageDeviceNumber uint32
	StorageManagerName  [8]uint16
	//nolint:structcheck
	alignmentPadding uint32 // necessary for 32bit support, see https://github.com/elastic/beats/pull/16553
}

type diskStat struct {
	name       string
	bytesRead  int64
	bytesWrite int64
}
