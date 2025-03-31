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

	router.Handle("/forget-password",
		middlewares.DecodeForgetRequest(
			middlewares.ValidateForgetRequest(
				http.HandlerFunc(authHandler.ForgetPassword),
			),
		),
	).Methods("POST")

	router.Handle("/reset-password",
		middlewares.DecodeResetPasswordRequest(
			middlewares.ValidateResetPasswordRequest(
				http.HandlerFunc(authHandler.ResetPassword),
			),
		),
	).Methods("POST")

	router.Handle("/auth/forget-password",
		middlewares.DecodeForgetRequest(
			middlewares.ValidateForgetRequest(
				http.HandlerFunc(authHandler.ForgetPassword),
			),
		),
	).Methods("POST")

	router.Handle("/auth/reset-password",
		middlewares.DecodeResetPasswordRequest(
			middlewares.ValidateResetPasswordRequest(
				http.HandlerFunc(authHandler.ResetPassword),
			),
		),
	).Methods("POST")

}
