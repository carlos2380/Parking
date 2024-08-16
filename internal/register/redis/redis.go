package redis

import (
	"fmt"
	"parking/internal/errors"
	"parking/internal/parking"
	"time"

	"github.com/go-redis/redis"
)

type RegisterRedis struct {
	RedisDB *redis.Client
	Parking *parking.Parking
}

func (RClient *RegisterRedis) InitRedis(ip string, port int, password string) error {

	RClient.RedisDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", ip, port),
		Password:     password,
		DB:           0,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	pong, err := RClient.RedisDB.Ping().Result()
	if err != nil {
		return errors.Wrap(err, *errors.ErrRedisConnectionFailure)
	}
	fmt.Println("Connected to Redis:", pong)
	return nil
}

func (RClient *RegisterRedis) InitParking(parking *parking.Parking) error {
	RClient.Parking = parking
	if err := RClient.UpdateParking(); err != nil {
		return err
	}
	return nil
}
