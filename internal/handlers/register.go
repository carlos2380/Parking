package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"parking/internal/errors"

	"github.com/gorilla/mux"
)

func (phandler ParckingHandler) EntryCar(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		params := mux.Vars(r)
		plateNumber := params["platenumber"]
		entryCar, err := phandler.Register.EntryCar(plateNumber)
		if err != nil {
			log.Println(err)
			errors.Wrap(err, *errors.ErrInternalServerError).Respond(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(entryCar)
	default:
		errors.ErrMethodNotAllowed.Respond(w)
	}
}

func (phandler *ParckingHandler) ExitCar(w http.ResponseWriter, r *http.Request) {
	setHeaders(w)
	fmt.Println("EXIT HANDLER")
	switch r.Method {
	case http.MethodOptions:
		w.WriteHeader(http.StatusOK)
	case http.MethodDelete:
		fmt.Println("EXIT CAR DELETEE")
		params := mux.Vars(r)
		plateNumber := params["platenumber"]
		ticket, err := phandler.Register.ExitCar(plateNumber)
		if err != nil {
			log.Println(err)
			errors.Wrap(err, *errors.ErrInternalServerError).Respond(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(ticket)
	default:
		errors.ErrMethodNotAllowed.Respond(w)
	}
}
