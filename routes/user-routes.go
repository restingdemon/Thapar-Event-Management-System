package routes

import (
	"github.com/gorilla/mux"
	controller "github.com/restingdemon/thaparEvents/controllers"
)

var RegisterUserRoutes = func(router *mux.Router) {
	router.HandleFunc("/users/get", controller.GetUserByEmail).Methods("GET")
	router.HandleFunc("/users/update/{email}", controller.UpdateUser).Methods("POST")
}
