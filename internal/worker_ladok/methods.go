package worker_ladok

import (
	"context"
	"fmt"
	"time"

	"eduid-cleaner/pkg/logger"
	"github.com/masv3971/eduid_amapi/amapi_types"
)

func (s *Service) ReportStats(ctx context.Context) func() {
	start := time.Now()
	s.kv.AddToCounter(ctx, "ladok")
	fn := func() {
		fmt.Printf("Took %v\n", time.Since(start))
	}
	return fn
}

func (s *Service) ProcessUser(ctx context.Context, user *amapi_types.User, logger *logger.Logger) {
	time.Sleep(1 * time.Second)
	logger.Info("processUser", "value", user.Eppn)
}

func (s *Service) Run(ctx context.Context, subWorkerID int) {
	subWorkerLogger := s.logger.New(fmt.Sprintf("%d", subWorkerID))
	for {
		select {
		case user, ok := <-s.UserChannel:
			if ok {
				subWorkerLogger.Info("consume message")
				s.ProcessUser(ctx, user, subWorkerLogger)
			} else {
				subWorkerLogger.Info("Channel has been closed")
			}
		default:
		}
	}
}
