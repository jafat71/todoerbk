package routes

import (
	"net/http"
	"todoerbk/handlers"
	"todoerbk/middlewares"

	"github.com/gorilla/mux"
)

func AuthRouter(router *mux.Router, authHandler *handlers.AuthHandler) {

	router.Handle("/register",
		middlewares.DecodeRegisterRequest(
			middlewares.ValidateRegisterRequest(
				http.HandlerFunc(authHandler.Register),
			),
		),
	).Methods("POST")

	router.Handle("/login",
		middlewares.DecodeLoginRequest(
			middlewares.ValidateLoginRequest(
				http.HandlerFunc(authHandler.Login),
			),
		),
	).Methods("POST")

}
