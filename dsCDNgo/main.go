package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Tkdefender88/dsCDN/server"
	"github.com/rs/cors"
)

func main() {
	router := server.NewRouter() // create routes

	//Set up CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS", "HEAD"},
	})

	//Launch the server with CORS validations
	fmt.Println("listening on port 30120")
	log.Fatal(http.ListenAndServe(":30120", c.Handler(router)))
}
