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
	Port int `yaml:"port"`
}

type NewsfeedConfig struct {
	Port int `yaml:"port"`
}

type WebConfig struct {
	Port int `yaml:"port"`
}

type SystemConfig struct {
	MySQL                mysql.Config       `yaml:"my_sql"`
	Redis                redis.Options      `yaml:"redis"`
	AuthenticationConfig *UserAndPostConfig `yaml:"user_and_post_config"`
	NewsfeedConfig       *NewsfeedConfig    `yaml:"newsfeed_config"`
	WebConfig            *WebConfig         `yaml:"web_config"`
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
