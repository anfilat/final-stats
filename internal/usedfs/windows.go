// +build windows

// based on https://github.com/shirou/gopsutil
package usedfs

import (
	"context"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/anfilat/final-stats/internal/common"
	"github.com/anfilat/final-stats/internal/symo"
)

func Collect(_ context.Context, action symo.MetricCommand) (symo.UsedFSData, error) {
	switch action {
	case symo.StartMetric:
		return nil, nil
	case symo.StopMetric:
		return nil, nil
	default:
		return get()
	}
}

func get() (symo.UsedFSData, error) {
	disks, err := partitions()
	if err != nil {
		return nil, err
	}

	result := make(symo.UsedFSData, 0, len(disks))
	for _, path := range disks {
		data, err := usage(path)
		if err != nil {
			continue
		}
		result = append(result, *data)
	}

	return result, nil
}

func partitions() ([]string, error) {
	result := make([]string, 0, 5)

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
		if typeRet != windows.DRIVE_REMOVABLE && typeRet != windows.DRIVE_FIXED &&
			typeRet != windows.DRIVE_REMOTE && typeRet != windows.DRIVE_CDROM {
			continue
		}

		lpVolumeNameBuffer := make([]byte, 256)
		lpVolumeSerialNumber := int64(0)
		lpMaximumComponentLength := int64(0)
		lpFileSystemFlags := int64(0)
		lpFileSystemNameBuffer := make([]byte, 256)
		volPath, _ := windows.UTF16PtrFromString(string(rune(v)) + ":/")
		driveRet, _, err := common.GetVolumeInformation.Call(
			uintptr(unsafe.Pointer(volPath)),
			uintptr(unsafe.Pointer(&lpVolumeNameBuffer[0])),
			uintptr(len(lpVolumeNameBuffer)),
			uintptr(unsafe.Pointer(&lpVolumeSerialNumber)),
			uintptr(unsafe.Pointer(&lpMaximumComponentLength)),
			uintptr(unsafe.Pointer(&lpFileSystemFlags)),
			uintptr(unsafe.Pointer(&lpFileSystemNameBuffer[0])),
			uintptr(len(lpFileSystemNameBuffer)))
		if driveRet == 0 {
			if typeRet == 5 || typeRet == 2 {
				continue // device is not ready will happen if there is no disk in the drive
			}
			return result, fmt.Errorf("call GetVolumeInformation error: %w", err)
		}

		result = append(result, path)
	}

	return result, nil
}

func usage(path string) (*symo.FSData, error) {
	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	diskRet, _, err := common.GetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(path))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)))
	if diskRet == 0 {
		return nil, fmt.Errorf("call GetDiskFreeSpaceExW error: %w", err)
	}
	result := &symo.FSData{
		Path:      path,
		UsedSpace: (float64(lpTotalNumberOfBytes) - float64(lpTotalNumberOfFreeBytes)) / float64(lpTotalNumberOfBytes) * 100,
		UsedInode: 0,
	}
	return result, nil
}
