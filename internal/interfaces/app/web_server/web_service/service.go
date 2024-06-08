package web_service

import (
	"github.com/khailequang334/social_network/configs"
	"github.com/khailequang334/social_network/internal/clients/newsfeed_client"
	"github.com/khailequang334/social_network/internal/clients/user_and_post_client"
	"github.com/khailequang334/social_network/internal/interfaces/proto/protobuf/newsfeed"
	"github.com/khailequang334/social_network/internal/interfaces/proto/protobuf/user_and_post"
	"github.com/khailequang334/social_network/internal/logger"
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

	zapLogger, err := logger.NewLogger(nil)
	if err != nil {
		return nil, err
	}

	return &WebService{
		UserAndPostClient: userAndPostClnt,
		NewsfeedClient:    newsfeedClnt,
		Logger:            zapLogger,
	}, nil
}
