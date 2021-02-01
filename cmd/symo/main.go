package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	config, err := newConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger.Level)
	if err != nil {
		log.Fatal(err)
	}

	logg.Info("starting system monitor")

	<-mainCtx.Done()

	logg.Info("stopping system monitor...")
	cancel()
	shutDown(logg)
	logg.Info("system monitor is stopped")
}

func watchSignals(cancel context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals
	cancel()
}

func shutDown(_ symo.Logger) {
	// TODO
}
