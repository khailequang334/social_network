package newsfeed_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/khailequang334/social_network/configs"
	"github.com/khailequang334/social_network/internal/protobuf/newsfeed"
	"github.com/khailequang334/social_network/internal/types"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type NewsfeedService struct {
	newsfeed.UnimplementedNewsfeedServer
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.Logger
}

func (nfs *NewsfeedService) GenerateNewsfeed(ctx context.Context, request *newsfeed.GenerateNewsfeedRequest) (*newsfeed.GenerateNewsfeedResponse, error) {
	// Ensure the user exists
	err := nfs.ensureUserExist(request.UserId)
	if err != nil {
		nfs.Logger.Debug("User not found", zap.Error(err))
		return &newsfeed.GenerateNewsfeedResponse{Status: newsfeed.GenerateNewsfeedResponse_USER_NOT_FOUND}, nil
	}

	// Query the user and their following users with their posts
	var user types.User
	err = nfs.DB.Preload("Following").Preload("Following.Posts").Find(&user, request.UserId).Error
	if err != nil {
		nfs.Logger.Error("Error retrieving user and following users", zap.Error(err))
		return nil, err
	}

	// Collect the IDs of posts from following users
	var postIDs []int64
	for _, following := range user.Following {
		for _, post := range following.Posts {
			postIDs = append(postIDs, int64(post.ID))
		}
	}

	// Return the generated newsfeed
	nfs.Logger.Debug("Generated newsfeed", zap.Any("postIDs", postIDs))
	return &newsfeed.GenerateNewsfeedResponse{
		Status:  newsfeed.GenerateNewsfeedResponse_OK,
		PostIds: postIDs,
	}, nil
}

func (nfs *NewsfeedService) ensureUserExist(userId int64) error {
	var user types.User
	err := nfs.DB.Table("user").Where("id = ?", userId).First(&user).Error
	if err != nil {
		return errors.New("user not found")
	}
	return nil
}

func NewNewsfeedService(conf *configs.NewsfeedConfig) (*NewsfeedService, error) {
	// DB Mysql
	db, err := gorm.Open(mysql.New(conf.MySQL), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println("can not connect to db ", err)
		return nil, err
	}

	// Redis
	rd := redis.NewClient(&conf.Redis)
	if rd == nil {
		return nil, fmt.Errorf("can not init redis client")
	}

	// Zap Logger
	logger, err := NewLogger()
	if err != nil {
		return nil, err
	}

	return &NewsfeedService{
		DB:     db,
		Redis:  rd,
		Logger: logger,
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
