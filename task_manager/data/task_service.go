package data

import (
	"context"
	"errors"
	"task_manager/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var collection *mongo.Collection

func InitCollection(c *mongo.Collection) {
	collection = c
}

func AddTask(task models.Task) (models.Task, error) {
	if task.ID.IsZero() {
		task.ID = bson.NewObjectID()
	}

	_, err := collection.InsertOne(context.TODO(), task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func ListTasks() ([]models.Task, error) {
	var tasks []models.Task

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	if err := cursor.All(context.TODO(), &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func GetTask(id string) (models.Task, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return models.Task{}, err
	}

	var task models.Task

	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func UpdateTask(id string, updatedTask models.Task) (models.Task, error) {
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

	_, err = collection.UpdateOne(context.TODO(), bson.M{"_id": objectID}, bson.M{"$set": set})
	if err != nil {
		return models.Task{}, err
	}

	var task models.Task
	err = collection.FindOne(context.TODO(), bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func DeleteTask(id string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil
}
