package server

import (
	"net/http"
	"parking/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(pHandler *handlers.ParckingHandler) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/cars/entry/{platenumber}", pHandler.EntryCar).Methods(http.MethodPost, http.MethodOptions)
	return r
}
