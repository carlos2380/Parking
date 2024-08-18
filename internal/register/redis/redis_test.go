package redis_test

import (
	"log"
	"parking/internal/errors"
	"parking/internal/parking"
	register "parking/internal/register/redis"
	"testing"

	"github.com/go-redis/redis"
)

func setupTestRedis() *register.RegisterRedis {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   1,
	})
	rdb.FlushDB()

	registerRedis := &register.RegisterRedis{
		RedisDB: rdb,
		Parking: &parking.Parking{},
	}
	registerRedis.Parking.SetMaxCars(8)
	registerRedis.Parking.SetPriceCentsMin(4)
	err := registerRedis.UpdateParking()
	if err != nil {
		log.Fatalf("Error update Parking: %v", err)
	}

	return registerRedis
}

func tearDownTestRedis(r *register.RegisterRedis) {
	r.RedisDB.FlushDB()
}

func TestEntryCarIntegration(t *testing.T) {
	r := setupTestRedis()
	defer tearDownTestRedis(r)

	tests := []struct {
		name        string
		plateNumber string
		expectError bool
	}{
		{"Valid entry", "ABC123", false},
		{"Duplicate entry", "ABC123", true},
		{"Another valid entry", "DEF456", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := r.EntryCar(test.plateNumber)
			if (err != nil) != test.expectError {
				t.Errorf("TEST:%s; EXPECT ERROR:%v; GOT ERROR:%v", test.name, test.expectError, err != nil)
			}
		})
	}
}

func TestExitCarIntegration(t *testing.T) {
	r := setupTestRedis()
	defer tearDownTestRedis(r)

	_, err := r.EntryCar("XYZ789")
	if err != nil {
		t.Fatalf("Failed to setup car entry: %v", err)
	}

	tests := []struct {
		name        string
		plateNumber string
		expectError bool
	}{
		{"Valid exit", "XYZ789", false},
		{"Invalid exit", "NONEXIST", true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := r.ExitCar(test.plateNumber)
			if (err != nil) != test.expectError {
				t.Errorf("TEST:%s; EXPECT ERROR:%v; GOT ERROR:%v", test.name, test.expectError, err != nil)
			}
		})
	}
}

func TestEntryCarWhenParkingFull(t *testing.T) {
	r := setupTestRedis()
	defer tearDownTestRedis(r)

	r.Parking.SetMaxCars(2)

	_, err := r.EntryCar("CAR001")
	if err != nil {
		t.Fatalf("Failed to register first car: %v", err)
	}
	_, err = r.EntryCar("CAR002")
	if err != nil {
		t.Fatalf("Failed to register second car: %v", err)
	}

	_, err = r.EntryCar("CAR003")
	if err == nil {
		t.Errorf("Expected error when parking is full, but got no error")
	} else if err != errors.ErrParkingClosed {
		t.Errorf("Expected ErrParkingClosed, but got: %v", err)
	}
}
