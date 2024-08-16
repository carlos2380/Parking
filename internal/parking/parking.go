package parking

import "fmt"

type Parking struct {
	maxCars  int64
	open     bool
	centsMin int
}

func (parking *Parking) SetMaxCars(maxCars int64) {
	parking.maxCars = maxCars
}

func (parking *Parking) SetPriceCentsMin(centsMin int) {
	parking.centsMin = centsMin
}

func (parking *Parking) UpdateStatus(numCars int64) {
	if numCars < parking.maxCars {
		parking.open = true
	} else {
		parking.open = false
	}
}

func (parking *Parking) IsOpen() bool {
	return parking.open
}

func (parking *Parking) GetPriceDollar(minutes int) (int, error) {
	if minutes < 0 {
		return 0, fmt.Errorf("minutes cannot be negative")
	}

	return parking.centsMin * minutes, nil
}

// func (parking *Parking) convertCentsToDollars(cents int) float64 {
// 	const centsInDollar = 100.0
// 	return float64(cents) / centsInDollar
// }
