package redis

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
)

func InitRedis(ip string, port int, password string) error {

	rdb := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", ip, port),
		Password:     password,
		DB:           0,
		DialTimeout:  10 * time.Second, // Tiempo de espera conexión
		ReadTimeout:  10 * time.Second, // Tiempo de espera lectura
		WriteTimeout: 10 * time.Second, // Tiempo de espera escritura
	})

	pong, err := rdb.Ping().Result()
	if err != nil {
		log.Fatalf("Unable to connect to Redis: %v", err)
		return err
	}
	fmt.Println("Connected to Redis:", pong)
	return nil
}
