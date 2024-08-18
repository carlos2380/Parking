package redis_test

import (
	"errors"
	"parking/internal/parking"
	register "parking/internal/register/redis"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/golang/mock/gomock"
)

func TestEntryCarWithMockRedis_ValidEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := register.NewMockRedisClient(ctrl)
	mockRedis.EXPECT().Exists("ZABC123").Return(redis.NewIntResult(0, nil)).Times(1)
	mockRedis.EXPECT().Set("ZABC123", gomock.Any(), time.Duration(0)).Return(redis.NewStatusResult("OK", nil)).Times(1)
	mockRedis.EXPECT().DBSize().Return(redis.NewIntResult(1, nil)).Times(1)

	r := &register.RegisterRedis{
		RedisDB: mockRedis,
		Parking: &parking.Parking{},
	}
	r.Parking.SetMaxCars(8)
	r.Parking.SetPriceCentsMin(4)
	r.Parking.UpdateStatus(0)

	car, err := r.EntryCar("ZABC123")
	if err != nil {
		t.Fatalf("TEST: Valid entry; EXPECTED ERROR: false; GOT ERROR: %v", err)
	}

	if car.PlateNumber != "ZABC123" {
		t.Errorf("TEST: Valid entry; EXPECTED PLATE: ZABC123; GOT: %s", car.PlateNumber)
	}
}

func TestEntryCarWithMockRedis_DuplicateEntry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := register.NewMockRedisClient(ctrl)

	mockRedis.EXPECT().Exists("ZABC123").Return(redis.NewIntResult(1, nil)).Times(1)

	r := &register.RegisterRedis{
		RedisDB: mockRedis,
		Parking: &parking.Parking{},
	}
	r.Parking.SetMaxCars(8)
	r.Parking.SetPriceCentsMin(4)
	r.Parking.UpdateStatus(0)

	_, err := r.EntryCar("ZABC123")
	if err == nil {
		t.Fatalf("TEST: Duplicate entry; EXPECTED ERROR: true; GOT ERROR: false")
	}

	expectedErr := "status 500: failed to register car entry"
	if err.Error() != expectedErr {
		t.Errorf("TEST: Duplicate entry; EXPECTED ERROR MESSAGE: %s; GOT: %s", expectedErr, err.Error())
	}
}

func TestEntryCarWithMockRedis_ParkingFull(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := register.NewMockRedisClient(ctrl)
	mockRedis.EXPECT().DBSize().Return(redis.NewIntResult(8, nil)).Times(1) // Simula parking lleno
	mockRedis.EXPECT().Exists("ZABC123").Times(0)

	r := &register.RegisterRedis{
		RedisDB: mockRedis,
		Parking: &parking.Parking{},
	}

	err := r.InitParking(r.Parking)
	if err != nil {
		t.Fatalf("TEST: InitParking; EXPECTED ERROR: false; GOT ERROR: %v", err)
	}

	_, err = r.EntryCar("ZABC123")
	if err == nil {
		t.Fatalf("TEST: Parking full; EXPECTED ERROR: true; GOT ERROR: false")
	}

	expectedErr := "status 507: Parking closed"
	if err.Error() != expectedErr {
		t.Errorf("TEST: Parking full; EXPECTED ERROR MESSAGE: %s; GOT: %s", expectedErr, err.Error())
	}
}

func TestEntryCarWithMockRedis_RedisErrorOnSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRedis := register.NewMockRedisClient(ctrl)
	mockRedis.EXPECT().DBSize().Return(redis.NewIntResult(0, nil)).AnyTimes()
	mockRedis.EXPECT().Exists("ZABC123").Return(redis.NewIntResult(0, nil)).Times(1)
	mockRedis.EXPECT().Set("ZABC123", gomock.Any(), time.Duration(0)).Return(redis.NewStatusResult("", errors.New("redis set error"))).Times(1)

	r := &register.RegisterRedis{
		RedisDB: mockRedis,
		Parking: &parking.Parking{},
	}

	r.Parking.SetMaxCars(8)
	r.Parking.UpdateStatus(0) // Explicitamente aseguramos que el parking está abierto

	_, err := r.EntryCar("ZABC123")
	if err == nil {
		t.Fatalf("TEST: Redis error on set; EXPECTED ERROR: true; GOT ERROR: false")
	}

	expectedErr := "status 500: failed to register car entry: redis set error"
	if err.Error() != expectedErr {
		t.Errorf("TEST: Redis error on set; EXPECTED ERROR MESSAGE: %s; GOT: %s", expectedErr, err.Error())
	}
}
