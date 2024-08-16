package register

import (
	"parking/internal/models"
)

type Register interface {
	UpdateParking() error //Update status parking
	EntryCar(plateNumber string) (models.Car, error)
	ExitCar(plateNumber string) (models.Ticket, error)
}
