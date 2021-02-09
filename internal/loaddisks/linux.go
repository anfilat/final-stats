package loaddisks

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
	mutex        sync.Mutex
	isWorked     bool
	command      *exec.Cmd
	cancelWork   context.CancelFunc
	loadDiskErr  error
	loadDiskData symo.LoadDisksData
)

func Read(ctx context.Context, action symo.MetricCommand) (symo.LoadDisksData, error) {
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

	cmd := exec.CommandContext(cancelableCtx, "sh", "-c", `iostat -dky 1`)
	out, err := cmd.StdoutPipe()
	if err != nil {
		cancel()
		return fmt.Errorf("cannot get iostat pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		cancel()
		return fmt.Errorf("cannot start iostat command: %w", err)
	}

	go readOut(out)

	isWorked = true
	command = cmd
	cancelWork = cancel
	return nil
}

func stop(ctx context.Context) error {
	if !isWorked {
		return nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	isWorked = false
	cancelWork()

	stopped := make(chan interface{})
	go func() {
		_ = command.Wait()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("iostas was not stopped")
	case <-stopped:
		return nil
	}
}

func readOut(out io.ReadCloser) {
	scanner := bufio.NewScanner(out)
	chunk := make([]string, 0, 16)
	header := true
	for scanner.Scan() {
		line := scanner.Text()
		if header {
			if strings.HasPrefix(line, "Device") {
				header = false
			}
			continue
		}

		if strings.HasPrefix(line, "loop") {
			continue
		}
		if strings.TrimSpace(line) == "" {
			saveLoadDisk(counters(chunk))
			chunk = make([]string, 0, len(chunk))
			header = true
			continue
		}
		chunk = append(chunk, line)
	}
}

func counters(chunk []string) (symo.LoadDisksData, error) {
	result := make([]symo.DiskData, 0, len(chunk))
	for _, line := range chunk {
		line = strings.ReplaceAll(line, ",", ".")

		values := strings.Fields(line)
		if len(values) < 4 {
			return nil, fmt.Errorf("cannot parse iostat line: %s", line)
		}

		name := values[0]

		tps, err := strconv.ParseFloat(values[1], 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse tps field: %w", err)
		}
		kbRead, err := strconv.ParseFloat(values[2], 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse kB_read/s field: %w", err)
		}
		kbWrite, err := strconv.ParseFloat(values[3], 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse kB_wrtn/s field: %w", err)
		}

		result = append(result, symo.DiskData{
			Name:    name,
			Tps:     tps,
			KBRead:  kbRead,
			KBWrite: kbWrite,
		})
	}
	return result, nil
}

func saveLoadDisk(data symo.LoadDisksData, err error) {
	mutex.Lock()
	defer mutex.Unlock()

	if !isWorked {
		return
	}

	loadDiskData = data
	loadDiskErr = err
}

func get() (symo.LoadDisksData, error) {
	mutex.Lock()
	defer mutex.Unlock()

	return loadDiskData, loadDiskErr
}
