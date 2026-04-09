package controller

import (
	"backend/model"
	"backend/util"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTask crea una nueva tarea con tags
// POST /tasks
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var taskReq model.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validar título
	if taskReq.Title == "" {
		http.Error(w, "Task title cannot be empty", http.StatusBadRequest)
		return
	}

	// Convertir nombres de tags a ObjectIDs
	tagIDs, err := util.FindOrCreateTags(taskReq.Tags)
	if err != nil {
		http.Error(w, "Failed to process tags", http.StatusInternalServerError)
		return
	}

	// Crear tarea
	task := model.Task{
		ID:          primitive.NewObjectID(),
		Title:       taskReq.Title,
		Description: taskReq.Description,
		Completed:   false,
		Tags:        tagIDs,
		CreatedAt:   primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err = util.DB.Collection("tasks").InsertOne(context.TODO(), task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Popular tagNames en la respuesta
	util.PopulateTaskTags(&task)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetTasks obtiene todas las tareas con tags populados
// GET /tasks
func GetTasks(w http.ResponseWriter, r *http.Request) {
	cursor, err := util.DB.Collection("tasks").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var tasks []model.Task
	if err = cursor.All(context.TODO(), &tasks); err != nil {
		http.Error(w, "Error decoding tasks", http.StatusInternalServerError)
		return
	}

	// Popular nombres de tags para todas las tareas
	util.PopulateTasksList(tasks)

	if tasks == nil {
		tasks = []model.Task{}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByID obtiene una tarea específica por ID
// GET /tasks/{id}
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task model.Task
	err = util.DB.Collection("tasks").FindOne(context.TODO(), bson.M{"_id": taskID}).Decode(&task)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Popular nombres de tags
	util.PopulateTaskTags(&task)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(task)
}

// UpdateTask actualiza una tarea incluyendo tags
// PUT /tasks/{id}
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var taskReq model.TaskRequest
	if err := json.NewDecoder(r.Body).Decode(&taskReq); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Convertir nombres de tags a ObjectIDs
	tagIDs, err := util.FindOrCreateTags(taskReq.Tags)
	if err != nil {
		http.Error(w, "Failed to process tags", http.StatusInternalServerError)
		return
	}

	// Actualizar tarea
	update := bson.M{
		"$set": bson.M{
			"title":       taskReq.Title,
			"description": taskReq.Description,
			"tags":        tagIDs,
		},
	}

	_, err = util.DB.Collection("tasks").UpdateOne(context.TODO(), bson.M{"_id": taskID}, update)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Obtener y retornar tarea actualizada
	var updatedTask model.Task
	err = util.DB.Collection("tasks").FindOne(context.TODO(), bson.M{"_id": taskID}).Decode(&updatedTask)
	if err != nil {
		http.Error(w, "Failed to fetch updated task", http.StatusInternalServerError)
		return
	}

	// Popular nombres de tags
	util.PopulateTaskTags(&updatedTask)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(updatedTask)
}

// CompleteTask marca una tarea como completada
// PUT /tasks/{id}/complete
func CompleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": bson.M{"completed": true}}
	_, err = util.DB.Collection("tasks").UpdateOne(context.TODO(), bson.M{"_id": taskID}, update)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteTask elimina una tarea
// DELETE /tasks/{id}
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	_, err = util.DB.Collection("tasks").DeleteOne(context.TODO(), bson.M{"_id": taskID})
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
