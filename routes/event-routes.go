package routes

import (
	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/controllers"
)

var RegisterEventRoutes = func(router *mux.Router) {
	router.HandleFunc("/event/create/{email}", controller.CreateEvent).Methods("POST")
	router.HandleFunc("/event/update/{eventId}", controller.UpdateEvent).Methods("POST")
	router.HandleFunc("/event/get", controller.GetAllEvents).Methods("GET")
	router.HandleFunc("/event/get/notvisible", controller.GetNotVisibleEvents).Methods("GET")
	router.HandleFunc("/event/visibility/{eventId}", controller.UpdateVisibility).Methods("POST")
	router.HandleFunc("/event/delete/{eventId}", controller.DeleteEvent).Methods("DELETE")
	router.HandleFunc("/event/upload/{eventId}", controller.UploadPhotos).Methods("POST")
	router.HandleFunc("/event/photo/delete/{eventId}", controller.DeletePhoto).Methods("DELETE")
	router.HandleFunc("/event/poster/upload/{eventId}", controller.UploadPoster).Methods("POST")
	router.HandleFunc("/event/dashboard/{eventId}", controller.GetEventDashboard).Methods("GET")
}
