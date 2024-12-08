package main

import (
	"log"
	"net/http"
	"os"

	"todoerbk/database"
	"todoerbk/handlers"
	"todoerbk/routes"
	"todoerbk/services"

	gorillaHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

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

	db, client, ctx, cancel := database.SetupMongoDB(mongoURL)
	defer database.CloseConnection(client, ctx, cancel)

	taskCollection := db.Collection("tasks")
	taskService := services.NewTaskService(taskCollection)
	taskController := handlers.NewTaskHandler(taskService)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()
	taskRouter := apiRouter.PathPrefix("/tasks").Subrouter()
	routes.TaskRouter(taskRouter, taskController)

	router.HandleFunc("/", handlers.Root).Methods("GET")

	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		log.Println("Registered route:", path)
		return nil
	})
	if err != nil {
		log.Println("Error walking routes:", err)
	}

	log.Println("GO SERVER RUNNING ON PORT", port)

	corsOptions := gorillaHandlers.CORS(
		gorillaHandlers.AllowedOrigins([]string{"http://localhost:5173"}),
		gorillaHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		gorillaHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	if err := http.ListenAndServe(port, corsOptions(router)); err != nil {
		log.Fatal(err)
	}
}
