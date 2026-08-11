package data

import (
	"context"
	"errors"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoTaskRepository struct {
	collection *mongo.Collection
}

func NewMongoTaskRepository(c *mongo.Collection) *MongoTaskRepository {
	return &MongoTaskRepository{
		collection: c,
	}
}

func (r *MongoTaskRepository) AddTask(ctx context.Context, task models.Task) (models.Task, error) {
	if task.ID.IsZero() {
		task.ID = bson.NewObjectID()
	}

	_, err := r.collection.InsertOne(ctx, task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r *MongoTaskRepository) ListTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *MongoTaskRepository) GetTask(ctx context.Context, id string) (models.Task, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, err
	}

	var task models.Task

	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r *MongoTaskRepository) UpdateTask(ctx context.Context, id string, updatedTask models.Task) (models.Task, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, err
	}

	set := bson.M{}
	if updatedTask.Title != "" {
		set["title"] = updatedTask.Title
	}
	if updatedTask.Description != "" {
		set["description"] = updatedTask.Description
	}
	if updatedTask.DueDate != "" {
		set["due_date"] = updatedTask.DueDate
	}
	if updatedTask.Status != "" {
		set["status"] = updatedTask.Status
	}

	if len(set) == 0 {
		return models.Task{}, errors.New("no fields to update")
	}

	_, err = r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": set})
	if err != nil {
		return models.Task{}, err
	}

	var task models.Task
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r *MongoTaskRepository) DeleteTask(ctx context.Context, id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}
