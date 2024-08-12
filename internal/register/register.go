package register

type Register interface {
	EntryCar(plateNumber string)
	ExitCar(plateNumber string)
}
