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

type TaskHandler struct {
	Service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{Service: service}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	task, ok := r.Context().Value(middlewares.TaskKey).(models.Task)

	if !ok {
		http.Error(w, "Unable to process task. Check Server", http.StatusInternalServerError)
		return
	}
	task.ID = primitive.NewObjectID()
	now := time.Now().UTC()
	task.CreatedAt = now
	task.UpdatedAt = now
	err := h.Service.CreateTask(r.Context(), &task)
	if err != nil {
		http.Error(w, "Unable to create task. Check Server", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Task created successfully",
		"task":    task,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	tasks, err := h.Service.GetTasks(r.Context())
	if err != nil {
		http.Error(w, "Unable to get tasks. Check Server", http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []models.Task{}
	}

	response := map[string]interface{}{
		"success": true,
		"message": "All tasks retrieved successfully",
		"tasks":   tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *TaskHandler) GetTaskById(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]
	taskToReturn, err := h.Service.GetTaskById(r.Context(), taskId)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Task retrieved successfully",
		"task":    taskToReturn,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(response)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskUpdateBody, ok := r.Context().Value(middlewares.TaskKey).(models.Task)
	if !ok {
		http.Error(w, "Unable to process task to update. Check Server", http.StatusInternalServerError)
		return
	}
	taskId := mux.Vars(r)["id"]
	taskToUpdate, err := h.Service.GetTaskById(r.Context(), taskId)
	if err != nil {
		http.Error(w, "Task to update not found", http.StatusNotFound)
		return
	}

	log.Println("UPDATE BODY:", taskUpdateBody)
	taskToUpdate.Title = taskUpdateBody.Title
	taskToUpdate.Doing = taskUpdateBody.Doing
	taskToUpdate.Done = taskUpdateBody.Done

	log.Println("UPDATING TASK:", taskToUpdate)
	now := time.Now().UTC()
	taskToUpdate.UpdatedAt = now

	err = h.Service.UpdateTask(r.Context(), taskId, *taskToUpdate)
	if err != nil {
		http.Error(w, "Task to update not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Task updated successfully",
		"task":    taskToUpdate,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *TaskHandler) DeleteTaskByID(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]
	taskToDelete, err := h.Service.GetTaskById(r.Context(), taskId)
	if err != nil {
		http.Error(w, "Task to delete not found", http.StatusNotFound)
		return
	}
	err = h.Service.DeleteTask(r.Context(), taskToDelete.ID.Hex())
	if err != nil {
		http.Error(w, "Task to delete not found", http.StatusNotFound)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Task with id " + taskId + " deleted successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
