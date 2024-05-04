package web_service

import (
	"encoding/json"

	"github.com/khailequang334/social_network/configs"
	"github.com/khailequang334/social_network/internal/clients/newsfeed_client"
	"github.com/khailequang334/social_network/internal/clients/user_and_post_client"
	"github.com/khailequang334/social_network/internal/protobuf/newsfeed"
	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
	"go.uber.org/zap"
)

type WebService struct {
	UserAndPostClient user_and_post.UserAndPostClient
	NewsfeedClient    newsfeed.NewsfeedClient
	Logger            *zap.Logger
}

func NewWebService(conf *configs.WebConfig) (*WebService, error) {
	userAndPostClnt, err := user_and_post_client.NewClient(conf.UserAndPost.Hosts)
	if err != nil {
		return nil, err
	}

	newsfeedClnt, err := newsfeed_client.NewClient(conf.Newsfeed.Hosts)
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger()
	if err != nil {
		return nil, err
	}

	return &WebService{
		UserAndPostClient: userAndPostClnt,
		NewsfeedClient:    newsfeedClnt,
		Logger:            logger,
	}, nil
}

func NewLogger() (*zap.Logger, error) {
	configJson := []byte(`{
		"level": "debug",
		"encoding": "json",
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		},
		"outputPaths": ["stdout", "/tmp/logs"],
		"errorOutputPaths": ["stderr"]
	}`)

	var cfg zap.Config
	if err := json.Unmarshal(configJson, &cfg); err != nil {
		return nil, err
	}
	logger := zap.Must(cfg.Build())
	return logger, nil
}
