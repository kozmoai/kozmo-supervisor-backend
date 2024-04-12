package redis

import (
	"github.com/kozmoai/kozmo-supervisor-backend/src/utils/config"
	redis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const RETRY_TIMES = 6

type RedisConfig struct {
	Addr     string `env:"KOZMO_REDIS_ADDR" envDefault:"localhost"`
	Port     string `env:"KOZMO_REDIS_PORT" envDefault:"6379"`
	Password string `env:"KOZMO_REDIS_PASSWORD" envDefault:"kozmo2022"`
	Database int    `env:"KOZMO_REDIS_DATABASE" envDefault:"0"`
}

func NewRedisConnectionByGlobalConfig(config *config.Config, logger *zap.SugaredLogger) (*redis.Client, error) {
	redisConfig := &RedisConfig{
		Addr:     config.GetRedisAddr(),
		Port:     config.GetRedisPort(),
		Password: config.GetRedisPassword(),
		Database: config.GetRedisDatabase(),
	}
	return NewRedisConnection(redisConfig, logger)
}

func NewRedisConnection(config *RedisConfig, logger *zap.SugaredLogger) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr + ":" + config.Port,
		Password: config.Password,
		DB:       config.Database,
	})

	logger.Infow("connected with redis", "redis", config)

	return rdb, nil
}
