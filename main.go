package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"todoerbk/handlers"
	"todoerbk/middlewares"
	"todoerbk/services"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	//TODO: REVISAR MODELO + INTERACTION CON MONGO
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error al cargar el archivo .env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT no está configurado en el archivo .env")
	}

	mongoURL := os.Getenv("MONGO_URL")
	if mongoURL == "" {
		log.Fatal("MONGO_URL no está configurado en el archivo .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURL)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Error al conectar a MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("No se pudo hacer ping al servidor MongoDB: %v", err)
	}

	log.Println("-*-*-*-*-*-Conexión exitosa a MongoDB-*-*-*-*-")

	db := client.Database("todoer")
	taskCollection := db.Collection("tasks")

	taskService := services.NewTaskService(taskCollection)
	taskHandler := handlers.TaskHandler{Service: *taskService}

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.Root).Methods("GET")

	router.Handle("/tasks",
		middlewares.DecodeTask(
			middlewares.ValidateTask(
				http.HandlerFunc(taskHandler.CreateTask),
			),
		),
	).Methods("POST")

	router.Handle("/tasks",
		http.HandlerFunc(taskHandler.GetTasks),
	).Methods("GET")

	router.Handle("/tasks/{id}",
		middlewares.ValidateTaskIdFromParams(
			http.HandlerFunc(taskHandler.GetTaskById),
		),
	).Methods("GET")

	router.Handle("/tasks/{id}",
		middlewares.DecodeTaskUpdate(
			middlewares.ValidateTaskUpdate(
				middlewares.ValidateTaskIdFromParams(
					http.HandlerFunc(taskHandler.UpdateTask),
				),
			),
		),
	).Methods("PUT")

	router.Handle("/tasks/{id}",
		middlewares.ValidateTaskIdFromParams(
			http.HandlerFunc(taskHandler.DeleteTaskByID),
		),
	).Methods("DELETE")

	log.Println("GO SERVER RUNNING ON PORT", port)

	if err := http.ListenAndServe(port, router); err != nil {
		log.Fatal(err)
	}
}
