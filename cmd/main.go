package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"eduid-cleaner/internal/apiv1"
	"eduid-cleaner/internal/httpserver"
	"eduid-cleaner/internal/storage"
	"eduid-cleaner/internal/worker_manager"
	"eduid-cleaner/pkg/configuration"
	"eduid-cleaner/pkg/logger"
)

type service interface {
	Close(ctx context.Context) error
}

func main() {
	wg := &sync.WaitGroup{}
	ctx := context.Background()

	var (
		log      *logger.Logger
		mainLog  *logger.Logger
		services = make(map[string]service)
	)

	cfg, err := configuration.Parse(logger.NewSimple("Configuration"))
	if err != nil {
		panic(err)
	}

	mainLog = logger.New("main", cfg.Production)
	log = logger.New("eduid-cleaner", cfg.Production)

	store, err := storage.New(cfg)
	if err != nil {
		panic(err)
	}

	workerManager, err := worker_manager.New(ctx, cfg, wg, store, log.New("worker_manager"))
	services["workerManager"] = workerManager
	if err != nil {
		panic(err)
	}

	apiv1Client, err := apiv1.New(cfg, store, log.New("apiv1"))
	if err != nil {
		panic(err)
	}

	httpserverService, err := httpserver.New(ctx, cfg, apiv1Client, log.New("httpserver"))
	services["httpserverService"] = httpserverService
	if err != nil {
		panic(err)
	}

	// Handle sigterm and await termChan signal
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)

	<-termChan // Blocks here until interrupted

	mainLog.Info("HALTING SIGNAL!")

	for serviceName, service := range services {
		err := service.Close(ctx)
		if err != nil {
			mainLog.Warn("Service:", serviceName, "error", err)
		}
	}

	wg.Wait() // Block here until all workers are done

	mainLog.Info("Stopped")
}
