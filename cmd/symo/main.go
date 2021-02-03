package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/anfilat/final-stats/internal/clients"
	"github.com/anfilat/final-stats/internal/grpc"
	"github.com/anfilat/final-stats/internal/heart"
	"github.com/anfilat/final-stats/internal/loadavg"
	"github.com/anfilat/final-stats/internal/logger"
	"github.com/anfilat/final-stats/internal/symo"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "", "Path to configuration file")
}

func main() {
	flag.Parse()

	if isVersionCommand() {
		printVersion()
		os.Exit(0)
	}

	mainCtx, cancel := context.WithCancel(context.Background())

	go watchSignals(cancel)

	config, err := symo.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Log.Level)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("starting system monitor")

	wg := &sync.WaitGroup{}

	readers := symo.MetricReaders{
		LoadAvg: loadavg.Read,
	}

	toHeartChan := make(symo.ClientsToHeartChan, 1)
	toClientsChan := make(symo.HeartToClientsChan, 1)

	clientsService := clients.NewClients(mainCtx, logg, toHeartChan, toClientsChan)
	clientsService.Start(wg)

	heartService := heart.NewHeart(mainCtx, logg, config.Metric, readers, toHeartChan, toClientsChan)
	heartService.Start(wg)

	grpcServer := grpc.NewServer(logg, clientsService)
	go func() {
		err := grpcServer.Start(":" + config.Server.Port)
		if err != nil {
			logg.Error(err)
			cancel()
		}
	}()

	logg.Info("system monitor is running...")

	<-mainCtx.Done()

	logg.Info("stopping system monitor...")
	shutDown(logg, wg, grpcServer)
	logg.Info("system monitor is stopped")
}

func watchSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}

func shutDown(logg symo.Logger, wg *sync.WaitGroup, grpcServer symo.GRPCServer) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Stop(ctx); err != nil {
			logg.Error(err)
		}
	}()

	wg.Wait()
}
