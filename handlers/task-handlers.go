package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"todoerbk/middlewares"
	"todoerbk/models"
)

var tasks = make([]models.Task, 0)

func CreateTask(w http.ResponseWriter, r *http.Request) {
	task, ok := r.Context().Value(middlewares.TaskKey).(models.Task)
	if !ok {
		http.Error(w, "Unable to process task. Check Server", http.StatusInternalServerError)
		return
	}

	tasks = append(tasks, task)
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
		"task":    tasks,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
