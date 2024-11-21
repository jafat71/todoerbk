package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"todoerbk/dtos"
	"todoerbk/models"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()
var tasks = make([]models.Task, 0)

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
	case "max":
		return "The field '" + fe.Field() + "' must have at most " + fe.Param() + " characters."
	case "uuid4":
		return "The field '" + fe.Field() + "' must be a valid UUID."
	default:
		return "Validation failed on field '" + fe.Field() + "'."
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("TODOER API"))
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTask(w http.ResponseWriter, r *http.Request, dto dtos.TaskDTO) {
	newTask := models.Task{
		Id:    dto.Id,
		Title: dto.Title,
		Doing: dto.Doing,
		Done:  dto.Done,
	}

	tasks = append(tasks, newTask)
	log.Println("TASK CREATED", newTask)
	log.Println(tasks)
	response := map[string]interface{}{
		"success": true,
		"message": "Task created successfully",
		"task": map[string]interface{}{
			"id":    newTask.Id,
			"title": newTask.Title,
			"doing": newTask.Doing,
			"done":  newTask.Done,
		},
	}

	b, err := json.Marshal(response)
	if err != nil {
		log.Println("ERROR MARSHALING RESPONSE", err)
		http.Error(w, "ERROR CREATING RESPONSE", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

func decodeTaskDTO(w http.ResponseWriter, r *http.Request) (*dtos.TaskDTO, error) {
	var dto dtos.TaskDTO
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		log.Println("ERROR IN DECODING JSON", err)
		http.Error(w, "INVALID JSON", http.StatusBadRequest)
		return nil, err
	}
	return &dto, nil
}

func HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	//TODO: EXTRACT MIDDLEWARES
	dto, err := decodeTaskDTO(w, r)
	if err != nil {
		return
	}
	err = validate.Struct(dto)

	if err != nil {
		responseErrors := getAllValidationErrs(err)
		if len(responseErrors) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"errors": responseErrors,
			})
			return
		}
		http.Error(w, "INVALID DATA PROVIDED, ", http.StatusBadRequest)
	}

	CreateTask(w, r, *dto)
}
