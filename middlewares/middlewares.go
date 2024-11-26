package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"todoerbk/data"
	"todoerbk/models"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

type contextKey string

const TaskKey contextKey = "task"
const TaskUpdateKey contextKey = "taskUpdate"

func getAllValidationErrs(err error) []map[string]string {
	var validationErrors validator.ValidationErrors
	var responseErrors []map[string]string
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			responseErrors = append(responseErrors, map[string]string{
				"field":   fieldError.Field(),
				"tag":     fieldError.Tag(),
				"value":   fieldError.Param(),
				"message": getValidationMessage(fieldError),
			})
		}
	}
	return responseErrors
}

func getValidationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "The field '" + fe.Field() + "' is required."
	case "min":
		return "The field '" + fe.Field() + "' must have at least " + fe.Param() + " characters."
	default:
		return "Validation failed on field '" + fe.Field() + "'"
	}
}

func DecodeTask(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var task models.Task
		err := json.NewDecoder(r.Body).Decode(&task)
		if err != nil {
			log.Println("ERROR IN DECODING JSON:", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), TaskKey, task)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func DecodeTaskUpdate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var taskUpdate models.TaskUpdate
		err := json.NewDecoder(r.Body).Decode(&taskUpdate)
		if err != nil {
			log.Println("ERROR IN DECODING JSON:", err)
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}
		log.Println("TASK UPDATE DECODED:", taskUpdate)
		ctx := context.WithValue(r.Context(), TaskUpdateKey, taskUpdate)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateTask(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		task, ok := r.Context().Value(TaskKey).(models.Task)
		if !ok {
			http.Error(w, "Invalid Task data", http.StatusBadRequest)
			return
		}

		err := validate.Struct(task)
		if err != nil {
			responseErrors := getAllValidationErrs(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": responseErrors,
			})
			return
		}

		if !isIdUnique(task.Id) {
			http.Error(w, "Task ID already exists", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateTaskUpdate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskUpdate, ok := r.Context().Value(TaskUpdateKey).(models.TaskUpdate)
		if !ok {
			http.Error(w, "Invalid Task data", http.StatusBadRequest)
			return
		}

		err := validate.Struct(taskUpdate)
		if err != nil {
			responseErrors := getAllValidationErrs(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": responseErrors,
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateTaskIdFromParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		taskId := mux.Vars(r)["id"]
		if taskId == "" {
			http.Error(w, "Task ID is required", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isIdUnique(taskId string) bool {
	for _, task := range data.Tasks {
		if task.Id == taskId {
			return false
		}
	}
	return true
}
