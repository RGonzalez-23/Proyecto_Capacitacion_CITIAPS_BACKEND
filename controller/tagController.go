package controller

import (
	"backend/model"
	"backend/util"
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateTag crea un nuevo tag
// POST /tags
func CreateTag(w http.ResponseWriter, r *http.Request) {
	var tagReq model.TagRequest
	if err := json.NewDecoder(r.Body).Decode(&tagReq); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validar nombre
	tagReq.Name = strings.TrimSpace(tagReq.Name)
	if tagReq.Name == "" {
		http.Error(w, "Tag name cannot be empty", http.StatusBadRequest)
		return
	}

	// Verificar que el tag no exista (case-insensitive)
	existing := util.DB.Collection("tags").FindOne(context.TODO(), bson.M{
		"name": bson.M{"$regex": "^" + tagReq.Name + "$", "$options": "i"},
	})

	if existing.Err() == nil {
		http.Error(w, "Tag already exists", http.StatusConflict)
		return
	}

	tag := model.Tag{
		ID:        primitive.NewObjectID(),
		Name:      tagReq.Name,
		Color:     tagReq.Color,
		CreatedAt: time.Now(),
	}

	_, err := util.DB.Collection("tags").InsertOne(context.TODO(), tag)
	if err != nil {
		http.Error(w, "Failed to create tag", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tag)
}

// GetAllTags obtiene todos los tags disponibles
// GET /tags
func GetAllTags(w http.ResponseWriter, r *http.Request) {
	cursor, err := util.DB.Collection("tags").Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Failed to fetch tags", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var tags []model.Tag
	if err = cursor.All(context.TODO(), &tags); err != nil {
		http.Error(w, "Error decoding tags", http.StatusInternalServerError)
		return
	}

	if tags == nil {
		tags = []model.Tag{}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(tags)
}

// GetTag obtiene un tag específico por ID
// GET /tags/{id}
func GetTag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tagID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	var tag model.Tag
	err = util.DB.Collection("tags").FindOne(context.TODO(), bson.M{"_id": tagID}).Decode(&tag)
	if err != nil {
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(tag)
}

// DeleteTag elimina un tag y lo remueve de todas las tareas
// DELETE /tags/{id}
func DeleteTag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tagID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	// Eliminar tag de todas las tareas
	_, err = util.DB.Collection("tasks").UpdateMany(
		context.TODO(),
		bson.M{"tags": tagID},
		bson.M{"$pull": bson.M{"tags": tagID}},
	)
	if err != nil {
		http.Error(w, "Failed to update tasks", http.StatusInternalServerError)
		return
	}

	// Eliminar el tag
	result, err := util.DB.Collection("tags").DeleteOne(context.TODO(), bson.M{"_id": tagID})
	if err != nil {
		http.Error(w, "Failed to delete tag", http.StatusInternalServerError)
		return
	}

	if result.DeletedCount == 0 {
		http.Error(w, "Tag not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetTasksByTag obtiene todas las tareas que tienen un tag específico
// GET /tags/{id}/tasks
func GetTasksByTag(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	tagID, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid tag ID", http.StatusBadRequest)
		return
	}

	cursor, err := util.DB.Collection("tasks").Find(context.TODO(), bson.M{
		"tags": tagID,
	})
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

	// Popular nombres de tags
	util.PopulateTasksList(tasks)

	if tasks == nil {
		tasks = []model.Task{}
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(tasks)
}
