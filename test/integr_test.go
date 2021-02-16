package test

import (
	"context"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	grpcClient "github.com/anfilat/final-stats/internal/grpc"
)

// интеграционный тест.
func TestSymo(t *testing.T) {
	err := compile()
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "./bin/symo")
	cmd.Dir = ".."

	err = cmd.Start()
	require.NoError(t, err)
	time.Sleep(3 * time.Second)

	// проверяем, что что-то приходит
	stats := getMetrics(t)
	require.NotNil(t, stats.LoadAvg)
	require.NotNil(t, stats.Cpu)
	require.NotNil(t, stats.LoadDisks)
	require.NotNil(t, stats.UsedFs)

	// проверяем, что метрика - "нагрузка процессора" действительно увеличивается при загрузке процессора вычислениями
	// уровень нагрузки процессора в ненагруженной системе
	cpuUser1 := getCPU(t)

	// уровень нагрузки процессора в нагруженной системе
	go load(3 * time.Second)
	go load(3 * time.Second)
	cpuUser2 := getCPU(t)

	require.Greater(t, cpuUser2, cpuUser1)

	cancel()
	_ = cmd.Wait()
}

func compile() error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("./build.bat")
	} else {
		cmd = exec.Command("make", "build")
	}
	cmd.Dir = ".."
	return cmd.Run()
}

func getMetrics(t *testing.T) *grpcClient.Stats {
	conn, err := grpc.Dial(":8000", grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := grpcClient.NewSymoClient(conn)
	req := &grpcClient.StatsRequest{
		N: int32(1),
		M: int32(2),
	}
	reqClient, err := client.GetStats(ctx, req)
	require.NoError(t, err)

	stats, err := reqClient.Recv()
	require.NoError(t, err)

	return stats
}

func getCPU(t *testing.T) float64 {
	conn, err := grpc.Dial(":8000", grpc.WithInsecure())
	require.NoError(t, err)
	defer conn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := grpcClient.NewSymoClient(conn)
	req := &grpcClient.StatsRequest{
		N: int32(1),
		M: int32(3),
	}
	reqClient, err := client.GetStats(ctx, req)
	require.NoError(t, err)

	stats, err := reqClient.Recv()
	require.NoError(t, err)

	return stats.Cpu.User
}

func load(dur time.Duration) {
	begin := time.Now()

	for {
		if time.Since(begin) > dur {
			break
		}
	}
}
