package routes

import (
	"net/http"
	"todoerbk/handlers"
	"todoerbk/middlewares"

	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router, userHandler *handlers.UserHandler, authMiddleware *middlewares.AuthMiddleware) {

	router.Handle("/{id}",
		authMiddleware.RequireAuth(
			middlewares.ValidateModelIdFromParams(
				http.HandlerFunc(userHandler.GetUserByID),
			),
		),
	).Methods("GET")

	router.Handle("/{id}",
		authMiddleware.RequireAuth(
			middlewares.ValidateModelIdFromParams(
				http.HandlerFunc(userHandler.DeleteUser),
			),
		),
	).Methods("DELETE")

	router.Handle("/{id}",
		authMiddleware.RequireAuth(
			middlewares.ValidateModelIdFromParams(
				http.HandlerFunc(userHandler.UpdateUser),
			),
		),
	).Methods("PUT")

}
