package routes

import (
	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/controllers"
)

var RegisterSocRoutes = func(router *mux.Router) {
	router.HandleFunc("/soc/register", controller.RegisterSociety).Methods("POST")
	router.HandleFunc("/soc/get", controller.GetSocietyDetails).Methods("GET")
	router.HandleFunc("/soc/update/{email}", controller.UpdateSociety).Methods("POST")
	router.HandleFunc("/soc/get/events", controller.GetSocEvents).Methods("GET")
	router.HandleFunc("/soc/dashboard/{email}", controller.GetSocDashboard).Methods("GET")
	router.HandleFunc("/soc/get/notvisible", controller.GetNotVisibleSoc).Methods("GET")
	router.HandleFunc("/soc/get/allevents/{email}", controller.GetAllSocEvents).Methods("GET")
}
