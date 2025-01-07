package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"todoerbk/middlewares"
	"todoerbk/models"
	"todoerbk/services"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BoardHandler struct {
	Service *services.BoardService
}

func NewBoardHandler(service *services.BoardService) *BoardHandler {
	return &BoardHandler{Service: service}
}

func (h *BoardHandler) CreateBoard(w http.ResponseWriter, r *http.Request) {
	board, ok := r.Context().Value(middlewares.BoardKey).(models.Board)

	if !ok {
		http.Error(w, "Unable to process board. Check Server", http.StatusInternalServerError)
		return
	}
	board.ID = primitive.NewObjectID()
	now := time.Now().UTC()
	board.CreatedAt = now
	board.UpdatedAt = now
	board.Completed = false
	board.Tasks = []models.Task{}

	err := h.Service.CreateBoard(r.Context(), &board)
	if err != nil {
		http.Error(w, "Unable to create board. Check Server", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Board created successfully",
		"board":   board,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *BoardHandler) GetBoards(w http.ResponseWriter, r *http.Request) {
	boards, err := h.Service.GetBoards(r.Context())
	if err != nil {
		http.Error(w, "Unable to get boards. Check Server", http.StatusInternalServerError)
		return
	}

	if boards == nil {
		boards = []models.Board{}
	}

	response := map[string]interface{}{
		"success": true,
		"message": "All boards retrieved successfully",
		"boards":  boards,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BoardHandler) GetBoardById(w http.ResponseWriter, r *http.Request) {
	boardId := mux.Vars(r)["id"]
	boardToReturn, err := h.Service.GetBoardById(r.Context(), boardId)
	if err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Board retrieved successfully",
		"board":   boardToReturn,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(response)
}

func (h *BoardHandler) UpdateBoard(w http.ResponseWriter, r *http.Request) {
	boardUpdateBody, ok := r.Context().Value(middlewares.BoardKey).(models.Board)
	if !ok {
		http.Error(w, "Unable to process board to update. Check Server", http.StatusInternalServerError)
		return
	}
	boardId := mux.Vars(r)["id"]
	boardToUpdate, err := h.Service.GetBoardById(r.Context(), boardId)
	if err != nil {
		http.Error(w, "Board to update not found", http.StatusNotFound)
		return
	}

	boardToUpdate.Title = boardUpdateBody.Title
	log.Println("UPDATING BOARD:", boardToUpdate)
	now := time.Now().UTC()
	boardToUpdate.UpdatedAt = now

	err = h.Service.UpdateBoard(r.Context(), boardId, *boardToUpdate)
	if err != nil {
		http.Error(w, "Board to update not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Board updated successfully",
		"board":   boardToUpdate,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BoardHandler) DeleteBoardByID(w http.ResponseWriter, r *http.Request) {
	boardId := mux.Vars(r)["id"]
	boardToDelete, err := h.Service.GetBoardById(r.Context(), boardId)
	if err != nil {
		http.Error(w, "Board to delete not found", http.StatusNotFound)
		return
	}
	err = h.Service.DeleteBoard(r.Context(), boardToDelete.ID.Hex())
	if err != nil {
		http.Error(w, "Board to delete not found", http.StatusNotFound)
		return
	}

	//TODO:Delete all tasks associated with the board

	response := map[string]interface{}{
		"success": true,
		"message": "Board with id " + boardId + " deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
