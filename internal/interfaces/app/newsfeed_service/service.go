package newsfeed_service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/khailequang334/social_network/configs"
	"github.com/khailequang334/social_network/internal/interfaces/proto/protobuf/newsfeed"
	"github.com/khailequang334/social_network/internal/logger"
	"github.com/khailequang334/social_network/internal/model"
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

func (nfs *NewsfeedService) getNewsfeedFromCache(ctx context.Context, userId int64) (*newsfeed.GenerateNewsfeedResponse, error) {
	cacheKey := "newsfeed:" + strconv.FormatInt(userId, 10)
	cachedData, err := nfs.Redis.Get(ctx, cacheKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}

	var cachedNewsfeed newsfeed.GenerateNewsfeedResponse
	err = json.Unmarshal(cachedData, &cachedNewsfeed)
	if err != nil {
		return nil, err
	}

	return &cachedNewsfeed, nil
}

func (nfs *NewsfeedService) cacheNewsfeed(ctx context.Context, userId int64, newsfeedResp *newsfeed.GenerateNewsfeedResponse) error {
	cacheKey := "newsfeed:" + strconv.FormatInt(userId, 10)
	cacheDuration := time.Hour

	data, err := json.Marshal(newsfeedResp)
	if err != nil {
		return err
	}

	return nfs.Redis.Set(ctx, cacheKey, data, cacheDuration).Err()
}

func (nfs *NewsfeedService) GenerateNewsfeed(ctx context.Context, request *newsfeed.GenerateNewsfeedRequest) (*newsfeed.GenerateNewsfeedResponse, error) {
	err := nfs.ensureUserExist(request.UserId)
	if err != nil {
		nfs.Logger.Debug("User not found", zap.Error(err))
		return &newsfeed.GenerateNewsfeedResponse{Status: newsfeed.GenerateNewsfeedResponse_USER_NOT_FOUND}, nil
	}

	cachedNewsfeed, err := nfs.getNewsfeedFromCache(ctx, request.UserId)
	if err != nil {
		nfs.Logger.Error("failed to get newsfeed from cache", zap.Error(err), zap.Int64("UserId", request.UserId))
	}
	if cachedNewsfeed != nil {
		return cachedNewsfeed, nil
	}

	var user model.User
	err = nfs.DB.Preload("Following").Preload("Following.Posts").Find(&user, request.UserId).Error
	if err != nil {
		nfs.Logger.Error("Error retrieving user and following users", zap.Error(err))
		return nil, err
	}

	var postIDs []int64
	for _, following := range user.Following {
		for _, post := range following.Posts {
			postIDs = append(postIDs, int64(post.ID))
		}
	}

	err = nfs.cacheNewsfeed(ctx, request.UserId, &newsfeed.GenerateNewsfeedResponse{
		Status:  newsfeed.GenerateNewsfeedResponse_OK,
		PostIds: postIDs,
	})
	if err != nil {
		nfs.Logger.Error("Failed to cache generated newsfeed", zap.Error(err), zap.Int64("UserId", request.UserId))
	}

	nfs.Logger.Debug("Generated newsfeed", zap.Any("postIDs", postIDs))
	return &newsfeed.GenerateNewsfeedResponse{
		Status:  newsfeed.GenerateNewsfeedResponse_OK,
		PostIds: postIDs,
	}, nil
}

func (nfs *NewsfeedService) ensureUserExist(userId int64) error {
	var user model.User
	err := nfs.DB.Table("user").Where("id = ?", userId).First(&user).Error
	if err != nil {
		return errors.New("user not found")
	}
	return nil
}

func NewNewsfeedService(conf *configs.NewsfeedConfig) (*NewsfeedService, error) {
	db, err := gorm.Open(mysql.New(conf.MySQL), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		fmt.Println("can not connect to db ", err)
		return nil, err
	}

	rd := redis.NewClient(&conf.Redis)
	if rd == nil {
		return nil, fmt.Errorf("can not init redis client")
	}

	zapLogger, err := logger.NewLogger(nil)
	if err != nil {
		return nil, err
	}

	return &NewsfeedService{
		DB:     db,
		Redis:  rd,
		Logger: zapLogger,
	}, nil
}
