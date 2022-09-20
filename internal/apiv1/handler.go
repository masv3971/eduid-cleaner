package apiv1

import (
	"context"
	"fmt"

	"eduid-cleaner/pkg/model"
)

type StatsReply struct {
	SKVCount   int `json:"skv_count"`
	LadokCount int `json:"ladok_count"`
}

func (c *Client) Stats(ctx context.Context) (*StatsReply, error) {
	skvCount, err := c.kv.GetCounter(ctx, "skv")
	fmt.Println("skv", skvCount, err)
	if err != nil {
		return nil, err
	}
	ladokCount, err := c.kv.GetCounter(ctx, "ladok")
	fmt.Println("ladok", skvCount, err)
	if err != nil {
		return nil, err
	}

	reply := &StatsReply{
		SKVCount:   skvCount,
		LadokCount: ladokCount,
	}

	return reply, nil
}

// Status return status for each worker_ladok instance
func (c *Client) Status(ctx context.Context) (*model.Status, error) {
	manyStatus := model.ManyStatus{}

	status := manyStatus.Check()

	return status, nil
}
