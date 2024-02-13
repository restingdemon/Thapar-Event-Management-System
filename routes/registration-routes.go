package routes

import (
	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/controllers"
)

var RegisterRegistrationRoutes = func(router *mux.Router) {
	router.HandleFunc("/event/register/{eventId}", controller.CreateRegistration).Methods("POST")
}
