package routes

import (
	"net/http"
	"todoerbk/handlers"
	"todoerbk/middlewares"

	"github.com/gorilla/mux"
)

func AuthRouter(router *mux.Router, authHandler *handlers.AuthHandler, authMiddleware *middlewares.AuthMiddleware) {

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

	router.Handle("/logout",
		middlewares.DecodeLogoutRequest(
			middlewares.ValidateLogoutRequest(
				http.HandlerFunc(authHandler.Logout),
			),
		),
	).Methods("POST")

	router.Handle("/check",
		authMiddleware.CheckAuth(
			http.HandlerFunc(authHandler.CheckAuthStatus),
		),
	).Methods("GET")

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

}
