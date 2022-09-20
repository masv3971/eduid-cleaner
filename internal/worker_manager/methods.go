package worker_manager

import (
	"context"
	"fmt"
	"time"

	"github.com/masv3971/eduid_amapi/amapi_types"
)

func (s *Service) fetchUsers(ctx context.Context, workerType string) amapi_types.Users {
	users := amapi_types.Users{}
	for i := 1; i < 1000; i++ {
		users = append(users, &amapi_types.User{Eppn: fmt.Sprintf("hubba-%d", i)})

	}
	time.Sleep(2 * time.Second)

	return users
}

//func (s *Service) manageWorkerChannels(ctx context.Context) error {
//	fmt.Println("run")
//	for workerName, worker := range s.workers.register {
//		if worker.IsUserChannelEmpty() {
//			s.logger.Info("worker channel is empty", "value", workerName)
//			users, err := s.fetchUsers(ctx)
//			if err != nil {
//				return err
//			}
//			worker.AddUsersToChannel(ctx, users)
//		}
//	}
//
//	return nil
//}
