package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"todoerbk/models"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

type contextKey string

const TaskKey contextKey = "task"

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
			w.Header().Set("Content-Type", "application/json") // Asegúrate de establecer el tipo de contenido
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in task creation",
				"errors":  responseErrors,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateModelIdFromParams(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		modelId := mux.Vars(r)["id"]
		if modelId == "" {
			http.Error(w, "Model ID is required", http.StatusBadRequest)
			return
		}
		if !primitive.IsValidObjectID(modelId) {
			http.Error(w, "Invalid Model ID", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
