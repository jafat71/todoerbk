package routes

import (
	"net/http"
	"todoerbk/handlers"
	"todoerbk/middlewares"

	"github.com/gorilla/mux"
)

func BoardRouter(router *mux.Router, boardHandler *handlers.BoardHandler, authMiddleware *middlewares.AuthMiddleware) {

	router.Handle("",
		authMiddleware.RequireAuth(
			http.HandlerFunc(boardHandler.GetBoards),
		),
	).Methods("GET")

	router.Handle("",
		authMiddleware.RequireAuth(
			middlewares.DecodeBoard(
				middlewares.ValidateBoard(
					http.HandlerFunc(boardHandler.CreateBoard),
				),
			),
		),
	).Methods("POST")

	router.Handle("/{id}",
		authMiddleware.RequireAuth(
			middlewares.ValidateModelIdFromParams(
				http.HandlerFunc(boardHandler.GetBoardById),
			),
		),
	).Methods("GET")

	router.Handle("/user/{id}",
		authMiddleware.RequireAuth(
			middlewares.ValidateModelIdFromParams(
				http.HandlerFunc(boardHandler.GetBoardsByUserId),
			),
		),
	).Methods("GET")

	router.Handle("/{id}",
		authMiddleware.RequireAuth(
			middlewares.DecodeBoard(
				middlewares.ValidateBoard(
					middlewares.ValidateModelIdFromParams(
						http.HandlerFunc(boardHandler.UpdateBoard),
					),
				),
			),
		),
	).Methods("PUT")

	router.Handle("/{id}",
		authMiddleware.RequireAuth(
			middlewares.ValidateModelIdFromParams(
				http.HandlerFunc(boardHandler.DeleteBoardByID),
			),
		),
	).Methods("DELETE")

}
