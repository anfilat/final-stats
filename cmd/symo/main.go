package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anfilat/final-stats/internal/clients"
	"github.com/anfilat/final-stats/internal/cpu"
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

	stopper := newServiceStopper()

	readers := symo.MetricReaders{
		LoadAvg: loadavg.Read,
		CPU:     cpu.Read,
	}

	toHeartChan := make(symo.ClientsToHeartChan, 1)
	toClientsChan := make(symo.HeartToClientsChan, 1)

	clientsService := clients.NewClients(logg)
	clientsService.Start(mainCtx, toHeartChan, toClientsChan)
	stopper.add(clientsService.Stop)

	heartService := heart.NewHeart(logg)
	heartService.Start(mainCtx, config.Metric, readers, toHeartChan, toClientsChan)
	stopper.add(heartService.Stop)

	grpcServer := grpc.NewServer(logg)
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

func watchSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}
