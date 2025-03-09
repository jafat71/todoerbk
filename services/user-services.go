package services

import (
	"context"
	"todoerbk/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	db *mongo.Collection
}

func NewUserService(db *mongo.Collection) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	createdUser, err := s.db.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	//i need to get the user from the database to return it
	user, err = s.GetUserByID(ctx, createdUser.InsertedID.(primitive.ObjectID).Hex())
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user models.User
	err = s.db.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := s.db.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := s.db.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) GetUserByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	err := s.db.FindOne(ctx, bson.M{"google_id": googleID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id string, user models.User) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"username":   user.Username,
			"email":      user.Email,
			"updated_at": user.UpdatedAt,
		},
	}

	_, err = s.db.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.db.DeleteOne(ctx, bson.M{"_id": objectID})
	return err
}
