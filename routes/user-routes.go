package routes

import (
	"github.com/gorilla/mux"
	controller "github.com/restingdemon/thaparEvents/controllers"
)

var RegisterUserRoutes = func(router *mux.Router) {
	router.HandleFunc("/users/get/{email}", controller.GetUserByEmail).Methods("GET")
	router.HandleFunc("/users/getall", controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/users/update/{email}", controller.UpdateUser).Methods("POST") 
}
