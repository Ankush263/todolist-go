package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Ankush263/todolist/internal/db"
	"github.com/Ankush263/todolist/internal/handler"
	"github.com/Ankush263/todolist/internal/middleware"
	"github.com/Ankush263/todolist/internal/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dsn := os.Getenv("CONNECTION_STRING")
	if dsn == "" {
		log.Fatal("CONNECTION_STRING is not set")
	}

	dbConn, err := db.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	todoRepo := repository.NewTodolistRepo(dbConn)
	userRepo := repository.NewUserRepo(dbConn)

	authHandler := handler.NewAuthHandler(userRepo)
	todoHandler := handler.NewTodolistHandler(todoRepo)

	r := mux.NewRouter()

	r.Use(middleware.Logging)
	r.Use(middleware.CORS)
	r.Use(middleware.Recover)

	r.HandleFunc("/signup", authHandler.Signup).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	todoRouter := r.PathPrefix("/").Subrouter()
	todoRouter.Use(middleware.Auth)

	todoRouter.HandleFunc("/todo", todoHandler.Create).Methods("POST")
	todoRouter.HandleFunc("/todo/{id}", todoHandler.GetById).Methods("GET")
	todoRouter.HandleFunc("/todo/{id}", todoHandler.Update).Methods("PATCH")
	todoRouter.HandleFunc("/todo/{id}", todoHandler.Delete).Methods("DELETE")

	log.Println("Server is running on the port 8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
