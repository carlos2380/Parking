package register

import (
	"parking/internal/models"
)

type Register interface {
	EntryCar(plateNumber string) (models.Car, error)
	//ExitCar(plateNumber string)
}
