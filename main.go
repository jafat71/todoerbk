package main

import (
	"log"
	"net/http"
	"todoerbk/handlers"
	"todoerbk/middlewares"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/", handlers.Root)

	router.Handle("POST /tasks",
		middlewares.DecodeTask(
			middlewares.ValidateTask(
				http.HandlerFunc(handlers.CreateTask),
			),
		),
	)
	router.Handle("GET /tasks", http.HandlerFunc(handlers.GetTasks))

	port := ":8080"
	log.Println("GO SERVER RUNNING ON PORT", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
