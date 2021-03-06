package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/benbjohnson/clock"

	"github.com/anfilat/final-stats/internal/clients"
	"github.com/anfilat/final-stats/internal/collector"
	"github.com/anfilat/final-stats/internal/cpu"
	"github.com/anfilat/final-stats/internal/grpc"
	"github.com/anfilat/final-stats/internal/loadavg"
	"github.com/anfilat/final-stats/internal/loaddisks"
	"github.com/anfilat/final-stats/internal/logger"
	"github.com/anfilat/final-stats/internal/symo"
	"github.com/anfilat/final-stats/internal/usedfs"
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

	go watchSignals(mainCtx, cancel)

	config, err := symo.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Log.Level)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("starting system monitor")
	logg.Info("time to keep metrics: ", config.App.MaxSeconds, " seconds")

	stopper := newServiceStopper()

	collectors := symo.MetricCollectors{
		LoadAvg:   loadavg.Collect,
		CPU:       cpu.Collect,
		LoadDisks: loaddisks.Collect,
		UsedFS:    usedfs.Collect,
	}

	toClientsCh := make(symo.CollectorToClientsCh, 1)

	clientsService := clients.NewClients(logg, clock.New())
	clientsService.Start(mainCtx, toClientsCh)
	stopper.add(clientsService.Stop)

	collectorService := collector.NewCollector(logg, config)
	collectorService.Start(mainCtx, collectors, toClientsCh)
	stopper.add(collectorService.Stop)

	grpcServer := grpc.NewServer(logg, config)
	go func() {
		err := grpcServer.Start(":"+config.Server.Port, clientsService)
		if err != nil {
			logg.Error(err)
			cancel()
			return
		}
	}()
	stopper.add(grpcServer.Stop)

	logg.Info("system monitor is running...")

	<-mainCtx.Done()

	logg.Info("stopping system monitor...")
	stopper.stop()
	logg.Info("system monitor is stopped")
}

func watchSignals(mainCtx context.Context, cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-mainCtx.Done():
	case <-signals:
	}
	cancel()
}
