package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/NeGat1FF/todolist-api/cmd/todolist-api/docs"
	"github.com/NeGat1FF/todolist-api/internal/database"
	"github.com/NeGat1FF/todolist-api/internal/handlers"
	"github.com/NeGat1FF/todolist-api/internal/middleware"
	"github.com/NeGat1FF/todolist-api/internal/repository"
	"github.com/NeGat1FF/todolist-api/internal/service"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			todolist-API
//	@version		0.1
//	@description	This is a simple todo list API
//	@host			localhost:8080
//	@BasePath		/
//	@schemes		http
//	@produce		json
//	@consumes		json

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mux := http.NewServeMux()

	db := database.InitDB()

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	servLogger := log.New(os.Stderr, "[SYSTEM] ", log.Ldate|log.Ltime|log.Lshortfile)

	serv := service.NewService(userRepo, taskRepo, servLogger)

	taskHandler := handlers.NewTaskHandler(serv)
	userHandler := handlers.NewUserHandler(serv)

	rateLimiter := middleware.NewRateLimiter(50, time.Minute)

	mux.HandleFunc("POST /register", rateLimiter.Middleware(middleware.ValidateRegistration(userHandler.RegisterUser)))
	mux.HandleFunc("POST /login", rateLimiter.Middleware(middleware.ValidateLogin(userHandler.LoginUser)))
	mux.HandleFunc("POST /refresh", rateLimiter.Middleware(userHandler.RefreshToken))

	mux.HandleFunc("POST /todos", rateLimiter.Middleware(middleware.AuthUserMiddleware(middleware.ValidateAddTask(taskHandler.AddTask))))
	mux.HandleFunc("GET /todos", rateLimiter.Middleware(middleware.AuthUserMiddleware(taskHandler.GetTasks)))
	mux.HandleFunc("PUT /todos/{id}", rateLimiter.Middleware(middleware.AuthUserMiddleware(middleware.ValidateUpdateTask(taskHandler.UpdateTask))))
	mux.HandleFunc("DELETE /todos/{id}", rateLimiter.Middleware(middleware.AuthUserMiddleware(taskHandler.DeleteTask)))

	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	http.ListenAndServe("localhost:8080", mux)
}
