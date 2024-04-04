package routes

import (
	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/controllers"
)

var RegisterEventRoutes = func(router *mux.Router) {
	router.HandleFunc("/event/create/{email}", controller.CreateEvent).Methods("POST")
	router.HandleFunc("/event/update/{eventId}", controller.UpdateEvent).Methods("POST")
	router.HandleFunc("/event/get", controller.GetAllEvents).Methods("GET")
	router.HandleFunc("/event_by_id/get/{eventId}", controller.GetEventById).Methods("GET")
	router.HandleFunc("/event/visibility/{eventId}", controller.UpdateVisibility).Methods("POST")
	router.HandleFunc("/event/delete/{eventId}", controller.DeleteEvent).Methods("DELETE")

}
