package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/anfilat/final-stats/internal/clients"
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

	logg, err := logger.New(config.Logger.Level)
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

	clients.
		NewClients(mainCtx, logg, toHeartChan, toClientsChan).
		Start(wg)

	heart.
		NewHeart(mainCtx, logg, config.Metric, readers, toHeartChan, toClientsChan).
		Start(wg)

	<-mainCtx.Done()

	logg.Info("stopping system monitor...")
	wg.Wait()
	logg.Info("system monitor is stopped")
}

func watchSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}
