package redis

import (
	"fmt"
	"parking/internal/errors"
	"time"

	"github.com/go-redis/redis"
)

func InitRedis(ip string, port int, password string) error {

	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", ip, port),
		Password:     password,
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	pong, err := rdb.Ping().Result()
	if err != nil {
		return errors.Wrap(err, *errors.ErrRedisConnectionFailure)
	}
	fmt.Println("Connected to Redis:", pong)
	return nil
}
