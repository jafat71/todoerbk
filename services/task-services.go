package services

import (
	"context"
	"errors"
	"todoerbk/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskService struct {
	db *mongo.Collection
}

func NewTaskService(db *mongo.Collection) *TaskService {
	return &TaskService{db: db}
}

func (s *TaskService) CreateTask(ctx context.Context, task *models.Task) error {
	_, err := s.db.InsertOne(ctx, task)
	return err
}

func (s *TaskService) DeleteTask(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task id")
	}
	_, err = s.db.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (s *TaskService) GetTaskById(ctx context.Context, id string) (*models.Task, error) {
	var task models.Task
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task id")
	}

	err = s.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&task)
	return &task, err
}

func (s *TaskService) GetTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	cursor, err := s.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id string, task models.Task) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid task id")
	}
	_, err = s.db.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": task})
	return err
}
