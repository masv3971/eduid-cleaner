package httpserver

import (
	"context"

	"eduid-cleaner/internal/apiv1"
	"eduid-cleaner/pkg/model"
)

// Apiv1 interface
type Apiv1 interface {
	Stats(ctx context.Context) (*apiv1.StatsReply, error)
	Status(ctx context.Context) (*model.Status, error)
}
