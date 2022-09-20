package worker_manager

import (
	"context"
	"sync"
	"time"

	"eduid-cleaner/internal/storage"
	"eduid-cleaner/internal/worker_ladok"
	"eduid-cleaner/internal/worker_skv"
	"eduid-cleaner/pkg/logger"
	"eduid-cleaner/pkg/model"

	"github.com/masv3971/eduid_amapi/amapi_types"
)

type workers struct {
	register map[string]workerService
	skv      *worker_skv.Service
	ladok    *worker_ladok.Service
}

type workerService interface {
	Close(ctx context.Context) error
	IsUserChannelEmpty() bool
	AddUsersToChannel(ctx context.Context, users amapi_types.Users)
	FetchUsers(ctx context.Context) amapi_types.Users
	ReportStats(ctx context.Context) func()
	ProcessUser(ctx context.Context, user *amapi_types.User, logger *logger.Logger)
	Run(ctx context.Context, subWorkerID int)
	CheckStatus(ctx context.Context)
}

func (w *workers) start(ctx context.Context, cfg *model.Cfg, wg *sync.WaitGroup, store storage.KV, logger *logger.Logger) error {
	w.register = make(map[string]workerService)

	var err error
	w.ladok, err = worker_ladok.New(ctx, cfg, wg, store, logger.New("ladok"))
	w.register["ladok"] = w.ladok
	if err != nil {
		return err
	}

	w.skv, err = worker_skv.New(ctx, cfg, wg, store, logger.New("skv"))
	w.register["skv"] = w.skv
	if err != nil {
		return err
	}

	return nil
}

type Service struct {
	config         *model.Cfg
	logger         *logger.Logger
	workers        *workers
	quitRunChannel chan bool
	wg             *sync.WaitGroup
	kv             storage.KV
}

func New(ctx context.Context, config *model.Cfg, wg *sync.WaitGroup, store *storage.Client, logger *logger.Logger) (*Service, error) {
	s := &Service{
		config:         config,
		logger:         logger,
		workers:        &workers{},
		quitRunChannel: make(chan bool),
		wg:             wg,
		kv:             store,
	}

	if err := s.workers.start(ctx, config, wg, s.kv, logger.New("worker")); err != nil {
		s.logger.Warn("start worker", "value", err)
	}

	s.run(ctx)

	s.logger.Info("started")
	return s, nil
}

func (s *Service) run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				for _, worker := range s.workers.register {
					worker.CheckStatus(ctx)
					if worker.IsUserChannelEmpty() {
						//users := s.fetchUsers(ctx, workerType)
						users := worker.FetchUsers(ctx)
						worker.AddUsersToChannel(ctx, users)
					}
				}
			case <-s.quitRunChannel:
				ticker.Stop()
				s.logger.Info("RunChannel stopped")
				return
			}
		}
	}()
}

func (s *Service) Close(ctx context.Context) error {
	s.quitRunChannel <- true

	for workerName, workerService := range s.workers.register {
		if err := workerService.Close(ctx); err != nil {
			s.logger.Warn("close worker", "worker_name", workerName, "error", err)
		}
	}
	ctx.Done()
	s.logger.Info("Quit")
	return nil
}
