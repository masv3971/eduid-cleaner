package configuration

import (
	"fmt"
	"os"
	"testing"

	"eduid-cleaner/pkg/logger"
	"eduid-cleaner/pkg/model"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"gopkg.in/yaml.v3"
)

var mockConfig = []byte(`
---
eduid:
  worker:
    common:
      debug: yes
    worker_ladok-x:
      api_server:
        host: :8080
    worker_ladok:
      api_server:
        host: :8080
      production: false
      tracing:
        kind: jaeger
        endpoint: "http://localhost:14268/api/traces"
      schools:
        kf: 
          saml_name: student.konstfack.se
        lnu:
          saml_name: lnu.se 
      worker_ladok:
        url: https://api.integrationstest.ladok.se
        certificate:
          folder: cert
        atom:
          periodicity: 60 
        permissions:
          90019: "rattighetsniva.las"
          51001: "rattighetsniva.las"
          61001: "rattighetsniva.las"
          11004: "rattighetsniva.las"
          860131: "rattighetsniva.las"
      eduid:
        iam:
          url: https://api.dev.eduid.se/scim/test 
      sunet:
        auth:
          url: https://auth-test.sunet.se 
      redis:
        db: 3
        host: localhost:6379
        sentinel_hosts:
        #  - localhost:1231
        #  - localhost:12313
        sentinel_service_name: redis-cluster
    x_service:
      api_server:
        host: 8080
`)

func TestParse(t *testing.T) {
	tempDir := t.TempDir()

	tts := []struct {
		name           string
		setEnvVariable bool
	}{
		{
			name:           "OK",
			setEnvVariable: true,
		},
	}

	for _, tt := range tts {
		path := fmt.Sprintf("%s/test.cfg", tempDir)
		if err := os.WriteFile(path, mockConfig, 0666); err != nil {
			assert.NoError(t, err)
		}
		if tt.setEnvVariable {
			os.Setenv("EDUID_CONFIG_YAML", path)
		}

		want := &model.Config{}
		err := yaml.Unmarshal(mockConfig, want)
		assert.NoError(t, err)

		testLog := logger.Logger{
			Logger: *zaptest.NewLogger(t, zaptest.Level(zap.PanicLevel)),
		}

		t.Run(tt.name, func(t *testing.T) {
			cfg, err := Parse(&testLog)
			assert.NoError(t, err)

			assert.Equal(t, &want.EduID.Worker.Cleaner, cfg)

		})
	}

}
