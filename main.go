package main

import (
	"backend/controller"
	"backend/middleware"
	"backend/util"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: Error loading .env file, using system environment variables")
	}

	// Load environment variables
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("DB_NAME")

	if mongoURI == "" || dbName == "" {
		log.Fatalf("Missing environment variables: MONGO_URI or DB_NAME")
	}

	util.ConnectDB(mongoURI, dbName)
}

func main() {
	r := mux.NewRouter()

	// Nota: envolver el router con el middleware CORS al arrancar el servidor
	// garantiza que peticiones como OPTIONS reciban los headers CORS aun
	// cuando el router respondería 405/404 antes de ejecutar middlewares.

	// Task routes
	r.HandleFunc("/tasks", controller.CreateTask).Methods("POST")
	r.HandleFunc("/tasks", controller.GetTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", controller.GetTaskByID).Methods("GET")
	r.HandleFunc("/tasks/{id}", controller.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}/complete", controller.CompleteTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", controller.DeleteTask).Methods("DELETE")

	// Tag routes
	r.HandleFunc("/tags", controller.CreateTag).Methods("POST")
	r.HandleFunc("/tags", controller.GetAllTags).Methods("GET")
	r.HandleFunc("/tags/{id}", controller.GetTag).Methods("GET")
	r.HandleFunc("/tags/{id}", controller.DeleteTag).Methods("DELETE")
	r.HandleFunc("/tags/{id}/tasks", controller.GetTasksByTag).Methods("GET")

	log.Println("Starting server on :8080")
	handler := middleware.CORSMiddleware(r)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
