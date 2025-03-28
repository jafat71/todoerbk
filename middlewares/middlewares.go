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
const BoardKey contextKey = "board"
const UserKey contextKey = "user"
const RegisterRequestKey contextKey = "register_request"
const LoginRequestKey contextKey = "login_request"

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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in task model validation",
				"errors":  []string{"Invalid JSON. Verify the data sent"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		log.Println("TASK:", task)

		ctx := context.WithValue(r.Context(), TaskKey, task)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateTask(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("VALIDATING TASK")
		log.Println("CONTEXT Entity sent:", r.Context().Value(TaskKey))
		task, ok := r.Context().Value(TaskKey).(models.Task)
		if !ok {
			http.Error(w, "Invalid Task data", http.StatusBadRequest)
			return
		}
		log.Println("TASK:", task)

		err := validate.Struct(task)
		if err != nil {
			responseErrors := getAllValidationErrs(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in task model validation",
				"errors":  responseErrors,
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		if task.Status != "" && !task.Status.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in task model validation",
				"errors":  []string{"Invalid Task Status. < field: status, value: TODO, DOING, DONE >"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		if task.Priority != "" && !task.Priority.IsValid() {
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in task model validation",
				"errors":  []string{"Invalid Task Priority. < field: priority, value: LOW, MEDIUM, HIGH >"},
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
		_, err := primitive.ObjectIDFromHex(modelId)
		if err != nil {
			http.Error(w, "Invalid Model ID", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func DecodeBoard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var board models.Board
		err := json.NewDecoder(r.Body).Decode(&board)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in board model validation",
				"errors":  []string{"Invalid JSON. Verify the data sent"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		ctx := context.WithValue(r.Context(), BoardKey, board)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateBoard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		board, ok := r.Context().Value(BoardKey).(models.Board)
		if !ok {
			http.Error(w, "Invalid Board data", http.StatusBadRequest)
			return
		}
		err := validate.Struct(board)
		if err != nil {
			responseErrors := getAllValidationErrs(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in board model validation",
				"errors":  responseErrors,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		//Validar que fromDate sea una fecha valida
		//user send date as fromDate not FromDate
		if board.FromDate.IsZero() {
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in board model validation",
				"errors":  []string{"FromDate must be a valid date"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		//Validar que toDate sea una fecha valida
		if board.ToDate.IsZero() {
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in board model validation",
				"errors":  []string{"ToDate must be a valid date"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		//Validar que toDate sea una fecha valida y posterior a la fecha de inicio
		if board.FromDate.After(board.ToDate) {
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in board model validation",
				"errors":  []string{"FromDate must be before ToDate"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func DecodeUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in user model validation",
				"errors":  []string{"Invalid JSON. Verify the data sent"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		ctx := context.WithValue(r.Context(), UserKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(UserKey).(models.User)
		if !ok {
			http.Error(w, "Invalid User data", http.StatusBadRequest)
			return
		}
		err := validate.Struct(user)
		if err != nil {
			responseErrors := getAllValidationErrs(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error in user model validation",
				"errors":  responseErrors,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func DecodeRegisterRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var registerRequest models.RegisterRequest
		err := json.NewDecoder(r.Body).Decode(&registerRequest)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error en validación de registro",
				"errors":  []string{"JSON inválido. Verifica los datos enviados"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		ctx := context.WithValue(r.Context(), RegisterRequestKey, registerRequest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateRegisterRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		registerRequest, ok := r.Context().Value(RegisterRequestKey).(models.RegisterRequest)
		if !ok {
			http.Error(w, "Invalid Register data", http.StatusBadRequest)
			return
		}
		err := validate.Struct(registerRequest)
		if err != nil {
			responseErrors := getAllValidationErrs(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error en validación de registro",
				"errors":  responseErrors,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func DecodeLoginRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var loginRequest models.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error en validación de login",
				"errors":  []string{"JSON inválido. Verifica los datos enviados"},
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		ctx := context.WithValue(r.Context(), LoginRequestKey, loginRequest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ValidateLoginRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loginRequest, ok := r.Context().Value(LoginRequestKey).(models.LoginRequest)
		if !ok {
			http.Error(w, "Invalid Login data", http.StatusBadRequest)
			return
		}
		err := validate.Struct(loginRequest)
		if err != nil {
			responseErrors := getAllValidationErrs(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			response := map[string]interface{}{
				"success": false,
				"message": "Error en validación de login",
				"errors":  responseErrors,
			}
			json.NewEncoder(w).Encode(response)
			return
		}
		next.ServeHTTP(w, r)
	})
}
