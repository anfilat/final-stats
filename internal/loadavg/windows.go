// +build windows

package loadavg

import (
	"context"
	"fmt"
	"log"
	"math"
	"sync"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/anfilat/final-stats/internal/symo"
)

var (
	loadErr              error
	loadAvg1M            float64
	loadAvg5M            float64
	loadAvg15M           float64
	loadAvgMutex         sync.Mutex
	loadAvgGoroutineOnce sync.Once
)

func Read(_ context.Context) (*symo.LoadAvgData, error) {
	loadAvgGoroutineOnce.Do(func() {
		go loadAvgGoroutine()
	})

	loadAvgMutex.Lock()
	defer loadAvgMutex.Unlock()

	if loadErr != nil {
		return nil, loadErr
	}

	return &symo.LoadAvgData{
		Load1:  loadAvg1M,
		Load5:  loadAvg5M,
		Load15: loadAvg15M,
	}, nil
}

func loadAvgGoroutine() {
	var (
		samplingFrequency = time.Second
		loadAvgFactor1M   = 1 / math.Exp(samplingFrequency.Seconds()/time.Minute.Seconds())
		loadAvgFactor5M   = 1 / math.Exp(samplingFrequency.Seconds()/(5*time.Minute).Seconds())
		loadAvgFactor15M  = 1 / math.Exp(samplingFrequency.Seconds()/(15*time.Minute).Seconds())
	)

	counter, err := processorQueueLengthCounter()
	if err != nil || counter == nil {
		log.Println(err)
		log.Println(counter)
		log.Println("unexpected processor queue length counter error")
		return
	}

	tick := time.NewTicker(samplingFrequency).C
	for {
		currentLoad, err := counter.GetValue()
		loadAvgMutex.Lock()
		loadErr = err
		loadAvg1M = loadAvg1M*loadAvgFactor1M + currentLoad*(1-loadAvgFactor1M)
		loadAvg5M = loadAvg5M*loadAvgFactor5M + currentLoad*(1-loadAvgFactor5M)
		loadAvg15M = loadAvg15M*loadAvgFactor15M + currentLoad*(1-loadAvgFactor15M)
		loadAvgMutex.Unlock()

		<-tick
	}
}

// copied from https://github.com/shirou/gopsutil

//nolint:golint,stylecheck
const (
	PDH_FMT_DOUBLE   = 0x00000200
	PDH_INVALID_DATA = 0xc0000bc6
	PDH_NO_DATA      = 0x800007d5
)

var (
	pdhDll = windows.NewLazySystemDLL("pdh.dll")

	pdhOpenQuery                = pdhDll.NewProc("PdhOpenQuery")
	pdhAddCounter               = pdhDll.NewProc("PdhAddEnglishCounterW")
	pdhCollectQueryData         = pdhDll.NewProc("PdhCollectQueryData")
	pdhGetFormattedCounterValue = pdhDll.NewProc("PdhGetFormattedCounterValue")
)

type win32PerformanceCounter struct {
	postName    string
	counterName string
	query       windows.Handle
	counter     windows.Handle
}

func (w *win32PerformanceCounter) GetValue() (float64, error) {
	r, _, err := pdhCollectQueryData.Call(uintptr(w.query))
	if r != 0 && err != nil {
		if r == PDH_NO_DATA {
			return 0, fmt.Errorf("this counter has not data: %w", err)
		}
		return 0, err
	}

	return getCounterValue(w.counter)
}

func processorQueueLengthCounter() (*win32PerformanceCounter, error) {
	const postName = "processor_queue_length"
	const counterName = `\System\Processor Queue Length`

	query, err := createQuery()
	if err != nil {
		return nil, err
	}
	var counter = win32PerformanceCounter{
		query:       query,
		postName:    postName,
		counterName: counterName,
	}
	r, _, err := pdhAddCounter.Call(
		uintptr(counter.query),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(counter.counterName))),
		0,
		uintptr(unsafe.Pointer(&counter.counter)),
	)
	if r != 0 {
		return nil, fmt.Errorf("pdh AddCounter: %w", err)
	}
	return &counter, nil
}

//nolint:golint,stylecheck
type PDH_FMT_COUNTERVALUE_DOUBLE struct {
	CStatus     uint32
	DoubleValue float64
}

func getCounterValue(counter windows.Handle) (float64, error) {
	var value PDH_FMT_COUNTERVALUE_DOUBLE
	r, _, err := pdhGetFormattedCounterValue.Call(uintptr(counter), PDH_FMT_DOUBLE, uintptr(0), uintptr(unsafe.Pointer(&value)))
	if r != 0 && r != PDH_INVALID_DATA {
		return 0, err
	}
	return value.DoubleValue, nil
}

func createQuery() (windows.Handle, error) {
	var query windows.Handle
	r, _, err := pdhOpenQuery.Call(0, 0, uintptr(unsafe.Pointer(&query)))
	if r != 0 {
		return 0, fmt.Errorf("create query: %w", err)
	}
	return query, nil
}
