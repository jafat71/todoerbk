package main

import (
	"log"
	"net/http"
	"todoerbk/handlers"
	"todoerbk/middlewares"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.Root).Methods("GET")

	router.Handle("/tasks",
		middlewares.DecodeTask(
			middlewares.ValidateTask(
				http.HandlerFunc(handlers.CreateTask),
			),
		),
	).Methods("POST")

	router.Handle("/tasks", http.HandlerFunc(handlers.GetTasks)).Methods("GET")
	router.Handle("/tasks/{id}",
		middlewares.ValidateTaskIdFromParams(
			http.HandlerFunc(handlers.GetTaskById),
		),
	).Methods("GET")

	router.Handle("/tasks/{id}",
		middlewares.DecodeTaskUpdate(
			middlewares.ValidateTaskUpdate(
				middlewares.ValidateTaskIdFromParams(
					http.HandlerFunc(handlers.UpdateTask),
				),
			),
		),
	).Methods("PUT")

	router.Handle("/tasks/{id}",
		middlewares.ValidateTaskIdFromParams(
			http.HandlerFunc(handlers.DeleteTaskByID),
		),
	).Methods("DELETE")

	//TODo: DELETE + ORM - MONGO
	port := ":8080"
	log.Println("GO SERVER RUNNING ON PORT", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
