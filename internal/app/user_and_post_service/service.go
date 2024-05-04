package user_and_post_service

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/khailequang334/social_network/configs"
	"github.com/khailequang334/social_network/internal/protobuf/user_and_post"
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
	// DB Mysql
	db, err := gorm.Open(mysql.New(conf.MySQL), &gorm.Config{})
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
	return &UserAndPostService{
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
