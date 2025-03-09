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
	Service     *services.BoardService
	TaskService *services.TaskService
}

func NewBoardHandler(service *services.BoardService, taskService *services.TaskService) *BoardHandler {
	return &BoardHandler{Service: service, TaskService: taskService}
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

	// userID := r.Context().Value(middlewares.UserIDKey).(string)
	// board.OwnerID, err = primitive.ObjectIDFromHex(userID)
	// if err != nil {
	// 	http.Error(w, "Invalid user ID", http.StatusBadRequest)
	// 	return
	// }

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

	log.Println("BOARD TO RETURN:", boardToReturn)
	log.Println("BOARD ID:", boardToReturn.ID.Hex())
	//Find all tasks associated with the board
	tasks, err := h.TaskService.GetTasksByBoardId(r.Context(), boardToReturn.ID.Hex())
	if err != nil {
		http.Error(w, "Unable to get tasks. Check Server", http.StatusInternalServerError)
		return
	}
	log.Println("TASKS:", tasks)

	if tasks == nil {
		tasks = []models.Task{}
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Board retrieved successfully",
		"board":   boardToReturn,
		"tasks":   tasks,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BoardHandler) GetBoardByUserId(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value(middlewares.UserIDKey).(string)
	boards, err := h.Service.GetBoardsByOwnerID(r.Context(), userId)
	if err != nil {
		http.Error(w, "Unable to get boards. Check Server", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Boards retrieved successfully",
		"boards":  boards,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	boardToUpdate.FromDate = boardUpdateBody.FromDate
	boardToUpdate.ToDate = boardUpdateBody.ToDate

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

	//Delete all tasks associated with the board
	err = h.TaskService.DeleteTasksByBoardId(r.Context(), boardToDelete.ID.Hex())
	if err != nil {
		http.Error(w, "Unable to delete tasks. Check Server", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Board with id " + boardId + " and all associated tasks deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BoardHandler) GetBoardsByOwnerID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	boards, err := h.Service.GetBoardsByOwnerID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Unable to get boards. Check Server", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Boards retrieved successfully",
		"boards":  boards,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
