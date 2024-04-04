package routes

import (
	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/controllers"
)

var RegisterSocRoutes = func(router *mux.Router) {
	router.HandleFunc("/soc/register", controller.RegisterSociety).Methods("POST")
	router.HandleFunc("/soc/get", controller.GetSocietyDetails).Methods("GET")
	router.HandleFunc("/soc_by_id/get/{societyId}", controller.GetSocietyDetailsByID).Methods("GET")
	router.HandleFunc("/soc/update/{email}", controller.UpdateSociety).Methods("POST")
}
