package repository

import (
	"context"
	"task_manager/Domain/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TaskRepository interface {
	Create(ctx context.Context, task models.Task) (bool, error)
	List(ctx context.Context, ownerID bson.ObjectID) (*mongo.Cursor, error)
	Get(ctx context.Context, objectID bson.ObjectID) (models.Task, error)
	Update(ctx context.Context, objectID bson.ObjectID, set bson.M) (bool, error)
	Delete(ctx context.Context, objectID bson.ObjectID) (bool, error)
}

type MongoTaskRepository struct {
	collection *mongo.Collection
}

func NewMongoTaskRepository(c *mongo.Collection) *MongoTaskRepository {
	return &MongoTaskRepository{
		collection: c,
	}
}

func (t *MongoTaskRepository) Create(ctx context.Context, task models.Task) (bool, error) {
	res, err := t.collection.InsertOne(ctx, task)
	if err != nil {
		return false, err
	}
	return res.Acknowledged, nil
}

func (t *MongoTaskRepository) List(ctx context.Context, ownerID bson.ObjectID) (*mongo.Cursor, error) {
	filter := bson.M{"owner": ownerID}

	cursor, err := t.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func (t *MongoTaskRepository) Get(ctx context.Context, objectID bson.ObjectID) (models.Task, error) {
	var task models.Task

	err := t.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (t *MongoTaskRepository) Update(ctx context.Context, objectID bson.ObjectID, set bson.M) (bool, error) {
	req, err := t.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": set})
	if err != nil {
		return false, err
	}

	return req.MatchedCount > 0, nil
}

func (t *MongoTaskRepository) Delete(ctx context.Context, objectID bson.ObjectID) (bool, error) {
	req, err := t.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return false, err
	}
	return req.DeletedCount > 0, nil
}
