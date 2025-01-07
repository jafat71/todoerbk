package routes

import (
	"net/http"
	"todoerbk/handlers"
	"todoerbk/middlewares"

	"github.com/gorilla/mux"
)

func BoardRouter(router *mux.Router, boardHandler *handlers.BoardHandler) {

	router.Handle("",
		http.HandlerFunc(boardHandler.GetBoards),
	).Methods("GET")

	router.Handle("",
		middlewares.DecodeBoard(
			middlewares.ValidateBoard(
				http.HandlerFunc(boardHandler.CreateBoard),
			),
		),
	).Methods("POST")

	router.Handle("/{id}",
		middlewares.ValidateModelIdFromParams(
			http.HandlerFunc(boardHandler.GetBoardById),
		),
	).Methods("GET")

	router.Handle("/{id}",
		middlewares.DecodeBoard(
			middlewares.ValidateBoard(
				middlewares.ValidateModelIdFromParams(
					http.HandlerFunc(boardHandler.UpdateBoard),
				),
			),
		),
	).Methods("PUT")

	router.Handle("/{id}",
		middlewares.ValidateModelIdFromParams(
			http.HandlerFunc(boardHandler.DeleteBoardByID),
		),
	).Methods("DELETE")

}
