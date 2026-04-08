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

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task model.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	task.ID = primitive.NewObjectID()
	task.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	_, err := util.DB.Collection("tasks").InsertOne(context.TODO(), task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	cursor, err := util.DB.Collection("tasks").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var tasks []model.Task
	for cursor.Next(context.TODO()) {
		var task model.Task
		if err := cursor.Decode(&task); err != nil {
			http.Error(w, "Error decoding task", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var task model.Task
	err = util.DB.Collection("tasks").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&task)
	if err != nil {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

func CompleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	update := bson.M{"$set": bson.M{"completed": true}}
	_, err = util.DB.Collection("tasks").UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task marked as completed"))
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	_, err = util.DB.Collection("tasks").DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Task deleted successfully"))
}
