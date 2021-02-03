package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os/exec"
	"sync"
	"time"

	"google.golang.org/grpc"

	grpcClient "github.com/anfilat/final-stats/internal/grpc"
)

var paramN int
var paramM int
var count int
var loadTime int

func init() {
	flag.IntVar(&paramN, "n", 1, "Send stats every N seconds")
	flag.IntVar(&paramM, "m", 5, "Send stats for last M seconds")
	flag.IntVar(&count, "count", 10000, "Count of clients")
	flag.IntVar(&loadTime, "time", 30, "Wait stats for time seconds")
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(ctx, "./bin/symo")
	err := cmd.Start()
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
	for i := 0; i < count; i++ {
		runClient(ctx, wg, conn, i)
	}
	wg.Wait()
}

func runClient(ctx context.Context, wg *sync.WaitGroup, conn grpc.ClientConnInterface, index int) {
	client := grpcClient.NewSymoClient(conn)
	req := &grpcClient.StatsRequest{
		N: int32(paramN),
		M: int32(paramM),
	}
	reqClient, err := client.GetStats(ctx, req)
	if err != nil {
		log.Print(fmt.Errorf("client request fail: %w", err))
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			stats, err := reqClient.Recv()
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				return
			}
			if errors.Is(err, io.EOF) {
				return
			} else if err != nil {
				log.Print(fmt.Errorf("error: %w", err))
				return
			}
			if index == 0 {
				log.Print(stats)
			}
		}
	}()
}
