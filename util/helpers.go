package util

import (
	"backend/model"
	"context"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SanitizeInput limpia inputs peligrosos para NoSQL injection
func SanitizeInput(input string) string {
	// Trimear espacios
	input = strings.TrimSpace(input)

	// Validar longitud máxima
	if len(input) > 500 {
		input = input[:500]
	}

	return input
}

// EscapeRegex escapa caracteres especiales en expresiones regulares
func EscapeRegex(s string) string {
	return regexp.QuoteMeta(s)
}

// FindOrCreateTags busca tags existentes o crea nuevos si no existen
// Retorna los ObjectIDs de los tags - PROTEGIDO CONTRA INYECCIÓN NOSQL
func FindOrCreateTags(tagNames []string) ([]primitive.ObjectID, error) {
	if len(tagNames) == 0 {
		return []primitive.ObjectID{}, nil
	}

	var tagIDs []primitive.ObjectID

	for _, tagName := range tagNames {
		// Sanitizar entrada
		normalizedName := SanitizeInput(tagName)
		if normalizedName == "" {
			continue
		}

		// Buscar tag existente con collation (case-insensitive seguro)
		opts := options.FindOne().SetCollation(&options.Collation{
			Locale:   "en",
			Strength: 2, // case-insensitive
		})

		var tag model.Tag
		err := DB.Collection("tags").FindOne(context.TODO(),
			bson.M{"name": normalizedName},
			opts,
		).Decode(&tag)

		if err != nil {
			// Tag no existe, crear nuevo
			newTag := model.Tag{
				ID:        primitive.NewObjectID(),
				Name:      normalizedName,
				CreatedAt: time.Now(),
			}
			_, err := DB.Collection("tags").InsertOne(context.TODO(), newTag)
			if err != nil {
				return nil, err
			}
			tagIDs = append(tagIDs, newTag.ID)
		} else {
			// Tag ya existe
			tagIDs = append(tagIDs, tag.ID)
		}
	}

	return tagIDs, nil
}

// PopulateTaskTags reemplaza ObjectIDs de tags con sus nombres en una tarea
func PopulateTaskTags(task *model.Task) error {
	if len(task.Tags) == 0 {
		task.TagNames = []string{}
		return nil
	}

	cursor, err := DB.Collection("tags").Find(context.TODO(), bson.M{
		"_id": bson.M{"$in": task.Tags},
	})
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	var tags []model.Tag
	if err = cursor.All(context.TODO(), &tags); err != nil {
		return err
	}

	var tagNames []string
	for _, tag := range tags {
		tagNames = append(tagNames, tag.Name)
	}

	task.TagNames = tagNames
	return nil
}

// PopulateTasksList aplica PopulateTaskTags a una lista de tareas
func PopulateTasksList(tasks []model.Task) error {
	for i := range tasks {
		if err := PopulateTaskTags(&tasks[i]); err != nil {
			return err
		}
	}
	return nil
}

// ValidateTagIDsExist valida que todos los ObjectIDs de tags existan
func ValidateTagIDsExist(tagIDs []primitive.ObjectID) (bool, error) {
	if len(tagIDs) == 0 {
		return true, nil
	}

	count, err := DB.Collection("tags").CountDocuments(context.TODO(), bson.M{
		"_id": bson.M{"$in": tagIDs},
	})
	if err != nil {
		return false, err
	}

	return count == int64(len(tagIDs)), nil
}

// GetTagByName obtiene un tag por nombre (case-insensitive) - PROTEGIDO CONTRA INYECCIÓN NOSQL
func GetTagByName(name string) (*model.Tag, error) {
	// Sanitizar entrada
	sanitizedName := SanitizeInput(name)

	// Usar collation en lugar de regex para evitar inyección
	opts := options.FindOne().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2, // case-insensitive
	})

	var tag model.Tag
	err := DB.Collection("tags").FindOne(context.TODO(),
		bson.M{"name": sanitizedName},
		opts,
	).Decode(&tag)
	return &tag, err
}
