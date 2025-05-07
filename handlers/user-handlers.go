package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"todoerbk/middlewares"
	"todoerbk/models"
	"todoerbk/services"
)

type UserHandler struct {
	Service      *services.UserService
	BoardService *services.BoardService
	TaskService  *services.TaskService
}

func NewUserHandler(service *services.UserService, boardService *services.BoardService, taskService *services.TaskService) *UserHandler {
	return &UserHandler{Service: service, BoardService: boardService, TaskService: taskService}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	user, err := h.Service.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Unable to get user. Check Server", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "User retrieved successfully",
		"user":    user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	user, err := h.Service.GetUserByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Unable to get user. Check Server", http.StatusInternalServerError)
		return
	}

	var updatedUser models.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Unable to decode user. Check Server", http.StatusInternalServerError)
		return
	}

	updatedUser.ID = user.ID
	updatedUser.CreatedAt = user.CreatedAt
	updatedUser.UpdatedAt = time.Now()

	err = h.Service.UpdateUser(r.Context(), userID, updatedUser)
	if err != nil {
		http.Error(w, "Unable to update user. Check Server", http.StatusInternalServerError)
		return
	}

}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	err := h.Service.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Unable to delete user. Check Server", http.StatusInternalServerError)
		return
	}

	//delete all boards iterating over the boards
	boards, err := h.BoardService.GetBoardsByOwnerID(r.Context(), userID)
	if err != nil {
		http.Error(w, "Unable to delete boards. Check Server", http.StatusInternalServerError)
		return
	}
	for _, board := range boards {
		//delete all tasks associated with the board
		tasks, err := h.TaskService.GetTasksByBoardId(r.Context(), board.ID.Hex())
		if err != nil {
			http.Error(w, "Unable to delete tasks. Check Server", http.StatusInternalServerError)
			return
		}
		for _, task := range tasks {
			err = h.TaskService.DeleteTask(r.Context(), task.ID.Hex())
			if err != nil {
				http.Error(w, "Unable to delete tasks. Check Server", http.StatusInternalServerError)
				return
			}
		}
		err = h.BoardService.DeleteBoard(r.Context(), board.ID.Hex())
		if err != nil {
			http.Error(w, "Unable to delete boards. Check Server", http.StatusInternalServerError)
			return
		}
	}

	response := map[string]interface{}{
		"success": true,
		"message": "User deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) InactivateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	err := h.Service.InactivateUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Unable to inactivate user. Check Server", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"success": true,
		"message": "User inactivated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middlewares.UserIDKey).(string)
	err := h.Service.ActivateUser(r.Context(), userID)
	if err != nil {
		http.Error(w, "Unable to activate user. Check Server", http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"success": true,
		"message": "User activated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
