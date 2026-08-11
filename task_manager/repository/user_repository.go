package repository

import (
	"context"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserRepository interface {
	Create(ctx context.Context, user models.User) (bool, error)
	FindByEmail(ctx context.Context, email string) (models.User, error)
	FindByID(ctx context.Context, id string) (models.User, error)
}

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(c *mongo.Collection) *MongoUserRepository {
	return &MongoUserRepository{collection: c}
}

func (u *MongoUserRepository) Create(ctx context.Context, user models.User) (bool, error) {
	res, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return false, err
	}
	return res.Acknowledged, nil
}

func (u *MongoUserRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User

	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *MongoUserRepository) FindByID(ctx context.Context, id string) (models.User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.User{}, err
	}

	var user models.User

	err = u.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
