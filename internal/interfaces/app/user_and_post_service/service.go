package user_and_post_service

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/khailequang334/social_network/configs"
	"github.com/khailequang334/social_network/internal/interfaces/proto/protobuf/user_and_post"
	"github.com/khailequang334/social_network/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserAndPostService struct {
	user_and_post.UnimplementedUserAndPostServer
	DB     *gorm.DB
	Redis  *redis.Client
	Logger *zap.Logger
}

func NewUserAndPostService(conf *configs.UserAndPostConfig) (*UserAndPostService, error) {
	db, err := gorm.Open(mysql.New(conf.MySQL), &gorm.Config{})
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
	return &UserAndPostService{
		DB:     db,
		Redis:  rd,
		Logger: zapLogger,
	}, nil
}
