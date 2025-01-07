package services

import (
	"context"
	"errors"
	"todoerbk/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BoardService struct {
	db *mongo.Collection
}

func NewBoardService(db *mongo.Collection) *BoardService {
	return &BoardService{db: db}
}

func (s *BoardService) CreateBoard(ctx context.Context, board *models.Board) error {
	_, err := s.db.InsertOne(ctx, board)
	return err
}

func (s *BoardService) DeleteBoard(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid board id")
	}
	_, err = s.db.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

func (s *BoardService) GetBoardById(ctx context.Context, id string) (*models.Board, error) {
	var board models.Board
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid task id")
	}

	err = s.db.FindOne(ctx, bson.M{"_id": objID}).Decode(&board)
	return &board, err
}

func (s *BoardService) GetBoards(ctx context.Context) ([]models.Board, error) {
	var boards []models.Board
	cursor, err := s.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var board models.Board
		if err := cursor.Decode(&board); err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}
	return boards, nil
}

func (s *BoardService) UpdateBoard(ctx context.Context, id string, board models.Board) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid board id")
	}
	_, err = s.db.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": board})
	return err
}
