package main

import (
	"log"
	"net/http"
	"os"
	"strings"

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
	boardCollection := db.Collection("boards")
	boardService := services.NewBoardService(boardCollection)
	taskCollection := db.Collection("tasks")
	taskService := services.NewTaskService(taskCollection)
	boardController := handlers.NewBoardHandler(boardService, taskService)
	taskController := handlers.NewTaskHandler(taskService, boardService)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	taskRouter := apiRouter.PathPrefix("/tasks").Subrouter()
	routes.TaskRouter(taskRouter, taskController)

	boardRouter := apiRouter.PathPrefix("/boards").Subrouter()
	routes.BoardRouter(boardRouter, boardController)

	router.HandleFunc("/", handlers.Root).Methods("GET")

	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		path, err := route.GetPathTemplate()
		if err != nil {
			return err
		}

		methods, err := route.GetMethods()
		if err != nil || len(methods) == 0 {
			methods = []string{"ANY"}
		}

		log.Printf("Registered route: %s %s", strings.Join(methods, ", "), path)
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
