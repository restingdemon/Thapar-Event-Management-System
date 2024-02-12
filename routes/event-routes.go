package routes

import (
	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/controllers"
)

var RegisterEventRoutes = func(router *mux.Router) {
	router.HandleFunc("/event/create/{email}", controller.CreateEvent).Methods("POST")
	router.HandleFunc("/event/update/{eventId}", controller.UpdateEvent).Methods("POST")
}
