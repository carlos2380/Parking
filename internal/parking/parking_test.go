package parking

import (
	"testing"
)

func TestParking_SetMaxCars(t *testing.T) {
	parking := &Parking{}
	maxCars := int64(50)
	parking.SetMaxCars(maxCars)

	if parking.maxCars != maxCars {
		t.Errorf("Expected maxCars to be %d, but got %d", maxCars, parking.maxCars)
	}
}

func TestParking_SetPriceCentsMin(t *testing.T) {
	parking := &Parking{}
	centsMin := 10
	parking.SetPriceCentsMin(centsMin)

	if parking.centsMin != centsMin {
		t.Errorf("Expected centsMin to be %d, but got %d", centsMin, parking.centsMin)
	}
}

func TestParking_UpdateStatus(t *testing.T) {
	tests := []struct {
		name     string
		maxCars  int64
		numCars  int64
		expected bool
	}{
		{"Parking open when numCars < maxCars", 50, 30, true},
		{"Parking closed when numCars = maxCars", 50, 50, false},
		{"Parking closed when numCars > maxCars", 50, 60, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parking := &Parking{}
			parking.SetMaxCars(test.maxCars)
			parking.UpdateStatus(test.numCars)

			if parking.IsOpen() != test.expected {
				t.Errorf("Expected parking open status to be %v, but got %v", test.expected, parking.IsOpen())
			}
		})
	}
}

func TestParking_IsOpen(t *testing.T) {
	parking := &Parking{}
	parking.SetMaxCars(50)
	parking.UpdateStatus(30)

	if !parking.IsOpen() {
		t.Errorf("Expected parking to be open, but it was closed")
	}

	parking.UpdateStatus(50)

	if parking.IsOpen() {
		t.Errorf("Expected parking to be closed, but it was open")
	}
}

func TestParking_GetPriceDollar(t *testing.T) {
	tests := []struct {
		name      string
		centsMin  int
		minutes   int
		expected  int
		expectErr bool
	}{
		{"Valid minutes and price", 10, 60, 600, false},
		{"Zero minutes", 10, 0, 0, false},
		{"Negative minutes", 10, -5, 0, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parking := &Parking{}
			parking.SetPriceCentsMin(test.centsMin)
			price, err := parking.GetPriceDollar(test.minutes)
			if (err != nil) != test.expectErr {
				t.Errorf("Expected error = %v, got %v", test.expectErr, err)
				return
			}
			if price != test.expected {
				t.Errorf("Expected price = %d, got %d", test.expected, price)
			}
		})
	}
}
