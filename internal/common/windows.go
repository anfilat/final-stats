// +build windows

package common

import "golang.org/x/sys/windows"

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	pdhDll   = windows.NewLazySystemDLL("pdh.dll")

	GetSystemTimes = kernel32.NewProc("GetSystemTimes")

	PdhOpenQuery                = pdhDll.NewProc("PdhOpenQuery")
	PdhAddCounter               = pdhDll.NewProc("PdhAddEnglishCounterW")
	PdhCollectQueryData         = pdhDll.NewProc("PdhCollectQueryData")
	PdhGetFormattedCounterValue = pdhDll.NewProc("PdhGetFormattedCounterValue")
)

//nolint:golint,stylecheck
const (
	PDH_FMT_DOUBLE   = 0x00000200
	PDH_INVALID_DATA = 0xc0000bc6
	PDH_NO_DATA      = 0x800007d5
)

//nolint:golint,stylecheck
type PDH_FMT_COUNTERVALUE_DOUBLE struct {
	CStatus     uint32
	DoubleValue float64
}

type FILETIME struct {
	DwLowDateTime  uint32
	DwHighDateTime uint32
}
