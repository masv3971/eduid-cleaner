package worker_skv

import (
	"context"
	"sync"

	"eduid-cleaner/internal/storage"
	"eduid-cleaner/pkg/logger"
	"eduid-cleaner/pkg/model"
	"github.com/masv3971/eduid_amapi"
	"github.com/masv3971/eduid_amapi/amapi_types"
)

type Service struct {
	cleaningType string
	config       *model.Cfg
	logger       *logger.Logger
	UserChannel  chan *amapi_types.User
	subWorkers   chan int
	wg           *sync.WaitGroup
	amapiClient  *eduid_amapi.Client
	kv           storage.KV
}

func New(ctx context.Context, config *model.Cfg, wg *sync.WaitGroup, kv storage.KV, logger *logger.Logger) (*Service, error) {
	s := &Service{
		cleaningType: "skv",
		config:       config,
		logger:       logger,
		UserChannel:  make(chan *amapi_types.User, 100000),
		subWorkers:   make(chan int),
		wg:           wg,
		kv:           kv,
	}

	var err error
	s.amapiClient, err = eduid_amapi.New(&eduid_amapi.Config{
		URL: s.config.Sunet.AmAPI.URL,
	})
	if err != nil {
		return nil, err
	}

	for subWorkerID := 0; subWorkerID <= s.config.Workers.SKV.SubWorkerAmount; subWorkerID++ {
		s.wg.Add(1)
		go s.Run(ctx, subWorkerID)
	}

	s.logger.Info("Started")
	return s, nil
}

func (s *Service) FetchUsers(ctx context.Context) amapi_types.Users {
	reply, httpResp, err := s.amapiClient.Sampler.Get(ctx, &eduid_amapi.SamplerRequest{
		Periodicity:         s.config.Workers.SKV.Periodicity,
		DurationOfExecution: 1.0,
		CleanedType:         "skv",
	})
	s.logger.Warn("FetchUsers", "value", httpResp, "reply", reply, "url", s.config.Sunet.AmAPI.URL)
	if err != nil {
		s.logger.Warn("FetchUsers", "value", err.Error())
		return nil
	}
	if reply.Data.Status {
		return reply.Data.Users
	}

	s.logger.Warn("FetchUsers", "msg", "return status from EduID AmAPI is false")
	return nil
}

func (s *Service) AddUsersToChannel(ctx context.Context, users amapi_types.Users) {
	for _, user := range users {
		select {
		case s.UserChannel <- user:
			s.logger.Info("AddUsersToChannel", "value", user.Eppn)
		default:
			s.logger.Info("AddUsersToChannel", "value", "UserChannel is full")
		}
	}
}

func (s *Service) IsUserChannelEmpty() bool {
	if len(s.UserChannel) == 0 {
		s.logger.Info("worker channel is empty")
		return true
	}
	return false
}

func (s *Service) CheckStatus(ctx context.Context) {

}

func (s *Service) Close(ctx context.Context) error {
	ctx.Done()
	for subWorkerID := 0; subWorkerID <= s.config.Workers.SKV.SubWorkerAmount; subWorkerID++ {
		s.logger.New("sub_worker").New("Quit")
		s.wg.Done()
	}
	s.logger.Info("Quit")

	return nil
}
