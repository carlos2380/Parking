package redis

import (
	"fmt"
	"log"
	"parking/internal/errors"
	"parking/internal/models"
	"parking/internal/utils"
	"time"

	"github.com/go-redis/redis"
)

func (RClient *RegisterRedis) UpdateParking() error {
	numCars, err := RClient.RedisDB.DBSize().Result()
	if err != nil {
		log.Fatal("Error getting number cars in Redis:", err)
	}
	RClient.Parking.UpdateStatus(numCars)
	return nil
}

func (RClient *RegisterRedis) EntryCar(plateNumber string) (models.Car, error) {

	if !RClient.Parking.IsOpen() {
		return models.Car{}, errors.ErrParkingClosed
	}

	now := utils.DateToStr(time.Now())
	exists, err := RClient.RedisDB.Exists(plateNumber).Result()
	if err != nil {
		log.Printf("Failed to register car entry. Plate: %s, Timestamp: %s, Error: %v", plateNumber, now, err)
		return models.Car{}, errors.Wrap(err, *errors.ErrRedisCarEntryFailed)
	}

	if exists > 0 {
		log.Printf("Car on the redis DB. Plate: %s, Timestamp: %s, Error: %v", plateNumber, now, err)
		return models.Car{}, errors.Wrap(err, *errors.ErrRedisCarEntryFailed)
	} else {
		err = RClient.RedisDB.Set(plateNumber, now, 0).Err()
		if err != nil {
			log.Printf("Failed to register car entry. Plate: %s, Timestamp: %s, Error: %v", plateNumber, now, err)
			return models.Car{}, errors.Wrap(err, *errors.ErrRedisCarEntryFailed)
		}

		fmt.Println("New car:", plateNumber, now)
		if err := RClient.UpdateParking(); err != nil {
			return models.Car{}, err
		}

		return models.Car{
			PlateNumber: plateNumber,
			EntryDate:   now,
		}, nil
	}
}

func (RClient *RegisterRedis) ExitCar(plateNumber string) (models.Ticket, error) {
	fmt.Println("EXIT CARR")
	now := time.Now()

	entryDateStr, err := RClient.RedisDB.Get(plateNumber).Result()
	if err == redis.Nil {
		log.Printf("Failed to on Redis. Plate: %s, Timestamp: %d, Error: %v", plateNumber, now.Unix(), err)
		return models.Ticket{}, errors.Wrap(err, *errors.ErrRedisCarEntryFailed)
	} else if err != nil {
		log.Printf("Failed to register car Plate: %s, Timestamp: %d, Error: %v", plateNumber, now.Unix(), err)
		return models.Ticket{}, errors.Wrap(err, *errors.ErrRedisCarEntryFailed)
	} else {
		_, err := RClient.RedisDB.Del(plateNumber).Result()
		if err != nil {
			log.Fatalf("Error on Plate: %s, Timestamp: %d, Error: %v", plateNumber, now.Unix(), err)
			return models.Ticket{}, err
		}
	}

	entryDate, err := utils.StrDateToDate(entryDateStr)
	if err != nil {
		return models.Ticket{}, err
	}

	minutes := utils.GetMinutes(&entryDate, &now)
	price, err := RClient.Parking.GetPriceDollar(minutes)
	if err != nil {
		log.Fatalf("Error on parse entry date %s : %v", entryDateStr, err)
		return models.Ticket{}, err
	}

	if err := RClient.UpdateParking(); err != nil {
		return models.Ticket{}, err
	}

	return models.Ticket{
		PlateNumber: plateNumber,
		EntryDate:   utils.DateToStr(entryDate),
		ExitDate:    utils.DateToStr(now),
		Price:       fmt.Sprintf("%s €", utils.CentsIntToEurStr(price)),
	}, nil
}
