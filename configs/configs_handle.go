package configs

import (
	"fmt"
	"io"
	"os"

	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
)

type UserAndPostConfig struct {
	Port  int           `yaml:"port"`
	MySQL mysql.Config  `yaml:"my_sql"`
	Redis redis.Options `yaml:"redis"`
}

type NewsfeedConfig struct {
	Port  int           `yaml:"port"`
	MySQL mysql.Config  `yaml:"my_sql"`
	Redis redis.Options `yaml:"redis"`
}

type WebConfig struct {
	Port        int `yaml:"port"`
	UserAndPost struct {
		Hosts []string `yaml:"hosts"`
	} `yaml:"user_and_post"`
	Newsfeed struct {
		Hosts []string `yaml:"hosts"`
	} `yaml:"newsfeed"`
}

type SystemConfig struct {
	MySQL             mysql.Config       `yaml:"my_sql"`
	Redis             redis.Options      `yaml:"redis"`
	UserAndPostConfig *UserAndPostConfig `yaml:"user_and_post_config"`
	NewsfeedConfig    *NewsfeedConfig    `yaml:"newsfeed_config"`
	WebConfig         *WebConfig         `yaml:"web_config"`
}

func GetSystemConfigs(path string) (*SystemConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open config (path=%s) error: %s", path, err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read config (path=%s) error: %s", path, err)
	}

	configs := &SystemConfig{}
	if err := yaml.Unmarshal(content, configs); err != nil {
		return nil, fmt.Errorf("unmarshal config (path=%s) error: %s", path, err)
	}
	return configs, nil
}

func GetUserAndPostConfig(path string) (*UserAndPostConfig, error) {
	configs, err := GetSystemConfigs(path)
	if err != nil {
		return nil, err
	}
	return configs.UserAndPostConfig, nil
}

func GetNewsfeedConfig(path string) (*NewsfeedConfig, error) {
	configs, err := GetSystemConfigs(path)
	if err != nil {
		return nil, err
	}
	return configs.NewsfeedConfig, nil
}

func GetWebConfig(path string) (*WebConfig, error) {
	configs, err := GetSystemConfigs(path)
	if err != nil {
		return nil, err
	}
	return configs.WebConfig, nil
}
