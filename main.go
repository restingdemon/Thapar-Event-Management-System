package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"

	"github.com/restingdemon/thaparEvents/middleware"
	"github.com/restingdemon/thaparEvents/routes"
)

func main() {
	r := mux.NewRouter()
	
	r.Use(middleware.Authenticate)
	
	routes.RegisterAuthRoutes(r)
	routes.RegisterUserRoutes(r)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})

	handler := c.Handler(r)
	http.Handle("/", handler)
	// Retrieve the PORT environment variable, default to 9010 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "9010"
	}

	addr := ":" + port
	log.Printf("Server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}
