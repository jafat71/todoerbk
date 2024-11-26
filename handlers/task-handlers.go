package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"todoerbk/data"
	"todoerbk/middlewares"
	"todoerbk/models"

	"github.com/gorilla/mux"
)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	task, ok := r.Context().Value(middlewares.TaskKey).(models.Task)
	if !ok {
		http.Error(w, "Unable to process task. Check Server", http.StatusInternalServerError)
		return
	}

	data.Tasks = append(data.Tasks, task)
	log.Println("TASK CREATED:", task)

	response := map[string]interface{}{
		"success": true,
		"message": "Task created successfully",
		"task":    task,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {

	response := map[string]interface{}{
		"success": true,
		"message": "All tasks retrieved successfully",
		"tasks":   data.Tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetTaskById(w http.ResponseWriter, r *http.Request) {
	taskId := mux.Vars(r)["id"]
	taskToReturn, found, _ := findTaskById(taskId)
	if !found {
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

func findTaskById(taskId string) (models.Task, bool, int) {
	for index, task := range data.Tasks {
		if task.Id == taskId {
			return task, true, index
		}
	}
	return models.Task{}, false, -1
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskUpdateBody, ok := r.Context().Value(middlewares.TaskUpdateKey).(models.TaskUpdate)
	if !ok {
		http.Error(w, "Unable to process task to update. Check Server", http.StatusInternalServerError)
		return
	}
	taskId := mux.Vars(r)["id"]
	taskToUpdate, found, index := findTaskById(taskId)
	if !found {
		http.Error(w, "Task to update not found", http.StatusNotFound)
		return
	}

	taskToUpdate.Title = taskUpdateBody.Title
	taskToUpdate.Doing = taskUpdateBody.Doing
	taskToUpdate.Done = taskUpdateBody.Done

	data.Tasks[index] = taskToUpdate

	response := map[string]interface{}{
		"success": true,
		"message": "Task updated successfully",
		"task":    taskToUpdate,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
