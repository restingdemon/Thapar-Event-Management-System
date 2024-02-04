package routes

import (
	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/controllers"
)

var RegisterAuthRoutes = func(router *mux.Router) {
	router.HandleFunc("/create", controller.Create).Methods("POST")
}
