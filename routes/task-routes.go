package routes

import (
	"net/http"
	"todoerbk/handlers"
	"todoerbk/middlewares"

	"github.com/gorilla/mux"
)

func TaskRouter(router *mux.Router, taskHandler *handlers.TaskHandler) {

	router.Handle("",
		http.HandlerFunc(taskHandler.GetTasks),
	).Methods("GET")

	router.Handle("",
		middlewares.DecodeTask(
			middlewares.ValidateTask(
				http.HandlerFunc(taskHandler.CreateTask),
			),
		),
	).Methods("POST")

	router.Handle("/{id}",
		middlewares.ValidateModelIdFromParams(
			http.HandlerFunc(taskHandler.GetTaskById),
		),
	).Methods("GET")

	router.Handle("/{id}",
		middlewares.DecodeTask(
			middlewares.ValidateTask(
				middlewares.ValidateModelIdFromParams(
					http.HandlerFunc(taskHandler.UpdateTask),
				),
			),
		),
	).Methods("PUT")

	router.Handle("/{id}",
		middlewares.ValidateModelIdFromParams(
			http.HandlerFunc(taskHandler.DeleteTaskByID),
		),
	).Methods("DELETE")

}
