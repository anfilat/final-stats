package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/anfilat/final-stats/internal/pb"
)

var nFrom int
var nTo int
var mFrom int
var mTo int
var count int
var loadTime int

func init() {
	flag.IntVar(&nFrom, "nf", 1, "Send stats every N seconds. Low limit")
	flag.IntVar(&nTo, "nt", 5, "Send stats every N seconds. High limit")
	flag.IntVar(&mFrom, "mf", 5, "Send stats for last M seconds. Low limit")
	flag.IntVar(&mTo, "mt", 30, "Send stats for last M seconds. High limit")
	flag.IntVar(&count, "count", 100_000, "Count of clients")
	flag.IntVar(&loadTime, "time", 600, "Wait stats for time seconds")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "./bin/symo")
	cmd.Env = append(os.Environ(),
		"LOG_LEVEL=INFO",
	)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Print("server: ", scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(3 * time.Second)

	load()

	cancel()
	_ = cmd.Wait()
}

func load() {
	conn, err := grpc.Dial(":8000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(loadTime)*time.Second)
	defer cancel()

	wg := &sync.WaitGroup{}
	runClient(ctx, wg, conn, true, nFrom, mFrom)
	for i := 1; i < count; i++ {
		n := nFrom
		if nTo > nFrom {
			//nolint:gosec
			n += rand.Intn(nTo - nFrom)
		}
		m := mFrom
		if mTo > mFrom {
			//nolint:gosec
			m += rand.Intn(mTo - mFrom)
		}
		runClient(ctx, wg, conn, false, n, m)
	}
	wg.Wait()
}

func runClient(timeoutCtx context.Context, wg *sync.WaitGroup, conn grpc.ClientConnInterface, first bool, n, m int) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		ctx, cancel := context.WithCancel(timeoutCtx)
		defer cancel()

		client := pb.NewSymoClient(conn)
		req := &pb.StatsRequest{
			N: int32(n),
			M: int32(m),
		}
		reqClient, err := client.GetStats(ctx, req)
		if err != nil {
			log.Print(fmt.Errorf("client request fail: %w", err))
			return
		}

		for {
			stats, err := reqClient.Recv()
			if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) ||
				errors.Is(ctx.Err(), context.Canceled) ||
				errors.Is(err, io.EOF) {
				return
			}
			if err != nil {
				log.Print(fmt.Errorf("error: %w", err))
				return
			}
			// результаты первого клиента показываются, чтобы в консоли что-то менялось
			if first {
				log.Print(fmt.Sprintf("%s %5.2f %5.2f %5.2f", stats.Time.AsTime(), stats.Cpu.User, stats.Cpu.System, stats.Cpu.Idle))
			}
			// переподключение клиента
			//nolint:gosec
			if !first && rand.Float64() < 0.1 {
				cancel()
				runClient(timeoutCtx, wg, conn, first, n, m)
			}
		}
	}()
}
