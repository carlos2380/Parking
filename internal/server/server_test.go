package server_test

import (
	"net/http"
	"net/http/httptest"
	"parking/internal/handlers"
	"parking/internal/models"
	"parking/internal/register"
	"parking/internal/server"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestEntryCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegister := register.NewMockRegister(ctrl)
	pHandler := &handlers.ParckingHandler{Register: mockRegister}
	router := server.NewRouter(pHandler)

	mockCar := models.Car{
		PlateNumber: "1234ABC",
		EntryDate:   "01 Jan 21 15:04 UTC",
	}

	mockRegister.EXPECT().
		EntryCar("1234ABC").
		Return(mockCar, nil).
		Times(1)

	req, err := http.NewRequest("POST", "/api/cars/entry/1234ABC", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: GOT %v EXPECT %v",
			status, http.StatusOK)
	}

	expected := `"plate_number":"1234ABC"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: GOT %v EXPECT %v",
			rr.Body.String(), expected)
	}
}

func TestExitCar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegister := register.NewMockRegister(ctrl)
	pHandler := &handlers.ParckingHandler{Register: mockRegister}
	router := server.NewRouter(pHandler)

	mockTicket := models.Ticket{
		PlateNumber: "1234ABC",
		EntryDate:   "01 Jan 21 15:04 UTC",
		ExitDate:    "01 Jan 21 16:04 UTC",
		Price:       "4.00 €",
	}

	mockRegister.EXPECT().
		ExitCar("1234ABC").
		Return(mockTicket, nil).
		Times(1)

	req, err := http.NewRequest("DELETE", "/api/cars/exit/1234ABC", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: GOT %v EXPECT %v",
			status, http.StatusOK)
	}

	expected := `"plate_number":"1234ABC"`
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("handler returned unexpected body: GOT %v EXPECT %v",
			rr.Body.String(), expected)
	}

	expectedPrice := `"price":"4.00 €"`
	if !strings.Contains(rr.Body.String(), expectedPrice) {
		t.Errorf("handler returned unexpected body: GOT %v EXPECT %v",
			rr.Body.String(), expectedPrice)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRegister := register.NewMockRegister(ctrl)
	pHandler := &handlers.ParckingHandler{Register: mockRegister}
	router := server.NewRouter(pHandler)

	req, err := http.NewRequest("PUT", "/api/cars/entry/1234ABC", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: GOT %v EXPECT %v",
			status, http.StatusMethodNotAllowed)
	}
}
