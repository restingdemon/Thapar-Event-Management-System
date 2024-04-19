package routes

import (
	"github.com/gorilla/mux"
	"github.com/restingdemon/thaparEvents/controllers"
)

var RegisterFeedbackRoutes = func(router *mux.Router) {
	router.HandleFunc("/feedback", controller.Feedback).Methods("POST")
}
