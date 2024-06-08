package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level      string `yaml:"level" default:"info"`
	Encoding   string `yaml:"encoding" default:"json"`
	OutputPath string `yaml:"output_path" default:"stdout"`
	LogDir     string `yaml:"log_dir" default:"./logs"`
}

func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		Encoding:   "json",
		OutputPath: "stdout",
		LogDir:     "./logs",
	}
}

func NewLogger(config *Config) (*zap.Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := os.MkdirAll(config.LogDir, 0755); err != nil {
		return nil, err
	}

	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.LevelKey = "level"
	encoderConfig.MessageKey = "message"
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var core zapcore.Core

	if config.Encoding == "json" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		)
	}

	if config.OutputPath != "stdout" {
		logFile := filepath.Join(config.LogDir, config.OutputPath)
		file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			fileCore := zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(file),
				level,
			)
			core = zapcore.NewTee(core, fileCore)
		}
	}

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}
