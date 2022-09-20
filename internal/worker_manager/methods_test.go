package worker_manager

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"eduid-cleaner/pkg/logger"
	"eduid-cleaner/pkg/model"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/masv3971/eduid_amapi/amapi_mocks"
	"github.com/masv3971/eduid_amapi/amapi_types"

	"github.com/stretchr/testify/assert"
)

type mockKV struct {
	redis *redis.Client
}

func (m *mockKV) AddToCounter(ctx context.Context, key string) error {
	return m.redis.Incr(ctx, key).Err()
}

func (m *mockKV) GetCounter(ctx context.Context, key string) (int, error) {
	return m.redis.Get(ctx, key).Int()
}

func (m *mockKV) SetStatus(ctx context.Context, nodeName string, statusMSG interface{}) error {
	return m.redis.HSet(ctx, "status", nodeName, statusMSG).Err()
}

func (m *mockKV) GetAllStatus(ctx context.Context) (map[string]string, error) {
	return m.redis.HGetAll(ctx, "status").Result()
}

func mockStore() (*mockKV, redismock.ClientMock) {
	redisClient, redisMock := redismock.NewClientMock()
	m := &mockKV{
		redis: redisClient,
	}
	return m, redisMock
}

func mockService(t *testing.T) (*Service, redismock.ClientMock) {
	ctx := context.TODO()
	wg := &sync.WaitGroup{}

	config := &model.Cfg{
		Production: false,
		Workers: model.Workers{
			SKV: model.Worker{
				Periodicity:     1,
				SubWorkerAmount: 5,
			},
			Ladok: model.Worker{
				Periodicity:     1,
				SubWorkerAmount: 5,
			},
		},
		Sunet: model.Sunet{
			Auth: model.RemoteAPI{},
			AmAPI: model.RemoteAPI{
				URL: mockAmAPIHttpServer(t), // Starts the mock EduID AmAPI http server and return its url
			},
		},
	}

	mockStore, redisMock := mockStore()

	s := &Service{
		config:  config,
		logger:  logger.New("test", false).New("sample"),
		workers: &workers{},
		kv:      mockStore,
	}

	err := s.workers.start(ctx, config, wg, mockStore, s.logger.New("test-worker"))
	assert.NoError(t, err)

	return s, redisMock
}

func mockGenericEndpointServer(t *testing.T, mux *http.ServeMux, method, url string, reply []byte, statusCode int) {
	mux.HandleFunc(url,
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			testMethod(t, r, method)
			testURL(t, r, url)
			w.Write(reply)
		},
	)
}

func testMethod(t *testing.T, r *http.Request, want string) {
	assert.Equal(t, want, r.Method)
}

func testURL(t *testing.T, r *http.Request, want string) {
	assert.Equal(t, want, r.RequestURI)
}

func mockAmAPIHttpServer(t *testing.T) string {
	mux := http.NewServeMux()

	server := httptest.NewServer(mux)

	args := []struct {
		name       string
		url        string
		method     string
		reply      []byte
		statusCode int
	}{
		{
			name:       "Update Name",
			url:        "%s/users/hubba-1/name",
			method:     http.MethodPut,
			reply:      amapi_mocks.UpdateNameReplyOKJSON,
			statusCode: http.StatusOK,
		},
		{
			name:       "Update Email",
			url:        "%s/users/hubba-1/email",
			method:     http.MethodPut,
			reply:      amapi_mocks.UpdateEmailReplyOKJSON,
			statusCode: http.StatusOK,
		},
		{
			name:       "Update Meta",
			url:        "%s/users/hubba-1/meta",
			method:     http.MethodPut,
			reply:      amapi_mocks.UpdateMetaReplyOKJSON,
			statusCode: http.StatusOK,
		},
		{
			name:       "Update Language",
			url:        "%s/users/hubba-1/language",
			method:     http.MethodPut,
			reply:      amapi_mocks.UpdateLanguageReplyOKJSON,
			statusCode: http.StatusOK,
		},
		{
			name:       "Update Phone",
			url:        "%s/users/hubba-1/phone",
			method:     http.MethodPut,
			reply:      amapi_mocks.UpdatePhoneReplyOKJSON,
			statusCode: http.StatusOK,
		},
		{
			name:       "Terminate user",
			url:        "%s/users/hubba-1",
			method:     http.MethodDelete,
			reply:      amapi_mocks.TerminateUserReplyOKJSON,
			statusCode: http.StatusOK,
		},
		{
			name:       "Sample users",
			url:        "%s/sampler/",
			method:     http.MethodPost,
			reply:      amapi_mocks.SamplerReplyOKJSON(t),
			statusCode: http.StatusOK,
		},
	}

	for _, arg := range args {
		url := fmt.Sprintf(arg.url, server.URL)
		t.Logf("starting endpoint %s", url)
		mockGenericEndpointServer(t, mux, arg.method, url, arg.reply, arg.statusCode)
	}

	return server.URL
}

func TestRun(t *testing.T) {
	s, _ := mockService(t)

	s.run(context.TODO())

	time.Sleep(3 * time.Second)

	//assert.Equal(t, false, s.workers.skv.IsUserChannelEmpty())

}

func TestChannelEmpty(t *testing.T) {
	s, redisMock := mockService(t)

	s.workers.skv.UserChannel <- &amapi_types.User{Eppn: "hubba-1"}
	s.workers.skv.UserChannel <- &amapi_types.User{Eppn: "hubba-2"}
	assert.Equal(t, false, s.workers.skv.IsUserChannelEmpty())
	assert.Equal(t, true, s.workers.ladok.IsUserChannelEmpty())

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestChannelAddUser(t *testing.T) {
	ctx := context.TODO()
	s, redisMock := mockService(t)
	assert.Equal(t, true, s.workers.skv.IsUserChannelEmpty())

	users := amapi_types.Users{
		{
			Eppn: "hubba-1",
		},
		{
			Eppn: "hubba-2",
		},
	}

	s.workers.skv.AddUsersToChannel(ctx, users)
	assert.Equal(t, 2, len(s.workers.skv.UserChannel))

	if err := redisMock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}
