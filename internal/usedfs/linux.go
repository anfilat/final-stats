// +build linux

package usedfs

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/anfilat/final-stats/internal/symo"
)

var (
	mutex      sync.Mutex
	isLive     bool
	command    *exec.Cmd
	cancelWork context.CancelFunc
	fsErr      error
	fsData     symo.UsedFSData
)

// Collect позволяет управлять получением информации об использовании файловых систем.
func Collect(ctx context.Context, action symo.MetricCommand) (symo.UsedFSData, error) {
	switch action {
	case symo.StartMetric:
		return nil, start(ctx)
	case symo.StopMetric:
		return nil, stop(ctx)
	default:
		return get()
	}
}

func start(ctx context.Context) error {
	mutex.Lock()
	defer mutex.Unlock()

	cancelableCtx, cancel := context.WithCancel(ctx)

	cmdLine := `while true ; do df --output="used,avail,itotal,iused,target" -x tmpfs -x squashfs; echo ---; sleep 1 ; done`
	cmd := exec.CommandContext(cancelableCtx, "sh", "-c", cmdLine)
	out, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("cannot get df pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		cancel()
		return fmt.Errorf("cannot start df command: %w", err)
	}

	go readOut(out)

	isLive = true
	command = cmd
	cancelWork = cancel
	return nil
}

func stop(ctx context.Context) error {
	if !isLive {
		return nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	isLive = false
	cancelWork()

	stopped := make(chan interface{})
	go func() {
		_ = command.Wait()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("df was not stopped")
	case <-stopped:
		return nil
	}
}

func readOut(out io.Reader) {
	scanner := bufio.NewScanner(out)
	chunk := make([]string, 0, 16)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			saveFS(usage(chunk[1:]))
			chunk = make([]string, 0, len(chunk))
			continue
		}
		chunk = append(chunk, line)
	}
}

func usage(chunk []string) (symo.UsedFSData, error) {
	result := make(symo.UsedFSData, 0, len(chunk))
	for _, line := range chunk {
		values := strings.Fields(line)
		if len(values) < 5 {
			return nil, fmt.Errorf("cannot parse df line: %s", line)
		}

		used, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse used field: %w", err)
		}
		avail, err := strconv.ParseInt(values[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse avail field: %w", err)
		}
		usedSpace := 0.0
		if used+avail != 0 {
			usedSpace = float64(used) * 100 / float64(used+avail)
		}

		iTotal, err := strconv.ParseInt(values[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse total field: %w", err)
		}
		iUsed, err := strconv.ParseInt(values[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse iused field: %w", err)
		}
		usedINode := 0.0
		if iTotal != 0 {
			usedINode = float64(iUsed) * 100 / float64(iTotal)
		}

		path := values[4]

		result = append(result, symo.FSData{
			Path:      path,
			UsedSpace: usedSpace,
			UsedInode: usedINode,
		})
	}
	return result, nil
}

func saveFS(data symo.UsedFSData, err error) {
	mutex.Lock()
	defer mutex.Unlock()

	if !isLive {
		return
	}

	fsData = data
	fsErr = err
}

func get() (symo.UsedFSData, error) {
	mutex.Lock()
	defer mutex.Unlock()

	return fsData, fsErr
}
