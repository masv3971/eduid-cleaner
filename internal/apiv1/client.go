package apiv1

import (
	"eduid-cleaner/internal/storage"
	"eduid-cleaner/pkg/logger"
	"eduid-cleaner/pkg/model"
)

// Client holds the publicapi object
type Client struct {
	config *model.Cfg
	logger *logger.Logger
	kv     storage.KV
}

// New creates a new instance of publicapi
func New(config *model.Cfg, kv *storage.Client, logger *logger.Logger) (*Client, error) {
	c := &Client{
		config: config,
		logger: logger,
		kv:     kv,
	}

	c.logger.Info("Started")

	return c, nil
}
