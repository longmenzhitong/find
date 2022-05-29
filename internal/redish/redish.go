// package redish implements methods for handling redis,
// which is the only option of synchronization for now.
package redish

import (
	"find/internal/config"
	"find/internal/logs"
	"github.com/go-redis/redis"
)

// Client is a pointer of redis client which other modules can use to synchronize.
var Client *redis.Client

func init() {
	if config.Conf.Backup.Redis.Address != "" {
		client := redis.NewClient(config.RedisConf())
		_, err := client.Ping().Result()
		if err != nil {
			logs.Error("ping redis error: %s\n", err.Error())
		} else {
			Client = client
		}
	}
}
