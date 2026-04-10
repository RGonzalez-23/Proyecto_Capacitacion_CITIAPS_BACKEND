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

	// Load environment variables - Permitir ambos nombres para compatibilidad
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = os.Getenv("MONGO_URI")
	}

	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		dbName = os.Getenv("DB_NAME")
	}

	if mongoURI == "" || dbName == "" {
		log.Fatalf("Missing environment variables: MONGODB_URI/MONGO_URI or MONGODB_DB_NAME/DB_NAME")
	}

	util.ConnectDB(mongoURI, dbName)
}

func main() {
	r := mux.NewRouter()

	// Nota: envolver el router con el middleware CORS al arrancar el servidor
	// garantiza que peticiones como OPTIONS reciban los headers CORS aun
	// cuando el router respondería 405/404 antes de ejecutar middlewares.

	// Task routes
	r.HandleFunc("/api/tasks", controller.CreateTask).Methods("POST")
	r.HandleFunc("/api/tasks", controller.GetTasks).Methods("GET")
	r.HandleFunc("/api/tasks/{id}", controller.GetTaskByID).Methods("GET")
	r.HandleFunc("/api/tasks/{id}", controller.UpdateTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}/complete", controller.CompleteTask).Methods("PUT")
	r.HandleFunc("/api/tasks/{id}", controller.DeleteTask).Methods("DELETE")

	// Tag routes
	r.HandleFunc("/api/tags", controller.CreateTag).Methods("POST")
	r.HandleFunc("/api/tags", controller.GetAllTags).Methods("GET")
	r.HandleFunc("/api/tags/{id}", controller.GetTag).Methods("GET")
	r.HandleFunc("/api/tags/{id}", controller.DeleteTag).Methods("DELETE")
	r.HandleFunc("/api/tags/{id}/tasks", controller.GetTasksByTag).Methods("GET")

	log.Println("Starting server on :8080")
	handler := middleware.CORSMiddleware(r)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
