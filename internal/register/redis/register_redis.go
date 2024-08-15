package redis

import (
	"fmt"
	"log"
	"parking/internal/errors"
	"parking/internal/models"
	"time"
)

func (RClient *RegisterRedis) EntryCar(plateNumber string) (models.Car, error) {

	now := time.Now()

	exists, err := RClient.RedisDB.Exists(plateNumber).Result()
	if err != nil {
		log.Printf("Failed to register car entry. Plate: %s, Timestamp: %d, Error: %v", plateNumber, now, err)
		return models.Car{}, errors.Wrap(err, *errors.ErrRedisCarEntryFailed)
	}

	if exists > 0 {
		log.Printf("Car on the redis DB. Plate: %s, Timestamp: %d, Error: %v", plateNumber, now, err)
		return models.Car{}, errors.Wrap(err, *errors.ErrRedisCarEntryFailed)
	} else {
		err = RClient.RedisDB.Set(plateNumber, now, 0).Err()
		if err != nil {
			log.Printf("Failed to register car entry. Plate: %s, Timestamp: %d, Error: %v", plateNumber, now, err)
			return models.Car{}, errors.Wrap(err, *errors.ErrRedisCarEntryFailed)
		}

		fmt.Println("New car:", plateNumber, now)
		return models.Car{
			PlateNumber: plateNumber,
			EntryDate:   now.Format(time.RFC822),
			ExitDate:    "",
		}, nil
	}
}
