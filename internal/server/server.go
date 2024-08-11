package server

import (
	"net/http"
	"parking/internal/handlers"

	"github.com/gorilla/mux"
)

func NewRouter(pHandler *handlers.ParckingHandler) http.Handler {
	r := mux.NewRouter()
	return r
}
