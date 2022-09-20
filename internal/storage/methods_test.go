package storage

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
)

func mockClient(t *testing.T) (*Client, redismock.ClientMock) {
	redisClient, redisMock := redismock.NewClientMock()
	client := &Client{
		redis: redisClient,
	}
	return client, redisMock
}

func TestRedisCounter(t *testing.T) {
	ctx := context.TODO()
	client, redisMock := mockClient(t)
	redisMock.ExpectIncr("testKey").SetVal(1)

	err := client.AddToCounter(ctx, "testKey")
	assert.NoError(t, err)

	redisMock.ExpectGet("testKey").SetVal("1")

	client.GetCounter(ctx, "testKey")

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestRedisStats(t *testing.T) {
	ctx := context.TODO()
	client, redisMock := mockClient(t)

	redisMock.ExpectHSet("status", "testNode", "testStatus").SetVal(1)
	redisMock.ExpectHGetAll("status").SetVal(map[string]string{
		"testNode": "testStatus",
	})

	err := client.SetStatus(ctx, "testNode", "testStatus")
	assert.NoError(t, err)

	_, err = client.GetAllStatus(ctx)
	assert.NoError(t, err)

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
