package app

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Router sets up the API routes
func Router() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/checkin", DailyCheckInHandler).Methods(http.MethodPost)
	return r
}