package main

import (
	"log"
	"net/http"
	"todoerbk/handlers"
)

func main() {

	router := http.NewServeMux()

	router.HandleFunc("GET /", handlers.Root)

	router.HandleFunc("POST /task", handlers.HandleCreateTask)

	port := ":8080"
	log.Println("GO SERVER RUNNING ON PORT ", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
