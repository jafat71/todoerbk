package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"todoerbk/database"
	"todoerbk/handlers"
	"todoerbk/middlewares"
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
	taskCollection := db.Collection("tasks")
	userCollection := db.Collection("users")

	boardService := services.NewBoardService(boardCollection)
	taskService := services.NewTaskService(taskCollection)
	userService := services.NewUserService(userCollection)
	authService := services.NewAuthService(userService)

	boardController := handlers.NewBoardHandler(boardService, taskService)
	taskController := handlers.NewTaskHandler(taskService, boardService)
	authController := handlers.NewAuthHandler(authService, userService)
	userController := handlers.NewUserHandler(userService, boardService, taskService)

	authMiddleware := middlewares.NewAuthMiddleware(authService)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	taskRouter := apiRouter.PathPrefix("/tasks").Subrouter()
	routes.TaskRouter(taskRouter, taskController, authMiddleware)

	boardRouter := apiRouter.PathPrefix("/boards").Subrouter()
	routes.BoardRouter(boardRouter, boardController, authMiddleware)

	userRouter := apiRouter.PathPrefix("/users").Subrouter()
	routes.UserRouter(userRouter, userController, authMiddleware)

	authRouter := apiRouter.PathPrefix("/auth").Subrouter()
	routes.AuthRouter(authRouter, authController, authMiddleware)

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
		gorillaHandlers.AllowCredentials(),
		gorillaHandlers.ExposedHeaders([]string{"Set-Cookie"}),
	)

	if err := http.ListenAndServe(port, corsOptions(router)); err != nil {
		log.Fatal(err)
	}
}
