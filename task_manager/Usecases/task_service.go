package Usecases

import (
	"context"
	"errors"
	"task_manager/Domain/dto"
	"task_manager/Domain/models"
	"task_manager/repository"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type TaskService struct {
	repo repository.TaskRepository
}

var ErrUnauthorized = errors.New("unauthorized")

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (ts *TaskService) AddTask(ctx context.Context, task dto.TaskRequest, owner string) (dto.TaskResponse, error) {
	id := bson.NewObjectID()
	status := "PENDING"
	ownerId, err := bson.ObjectIDFromHex(owner)

	if err != nil {
		return dto.TaskResponse{}, err
	}

	newTask := models.Task{
		ID:          id,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      status,
		Owner:       ownerId,
	}

	_, err = ts.repo.Create(ctx, newTask)

	if err != nil {
		return dto.TaskResponse{}, err
	}

	result := dto.TaskResponse{
		ID:          id.Hex(),
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      status,
	}
	return result, nil
}

func (ts *TaskService) ListTasks(ctx context.Context, owner string) ([]dto.TaskResponse, error) {
	ownerID, err := bson.ObjectIDFromHex(owner)
	if err != nil {
		return nil, err
	}

	var tasks []models.Task

	cursor, err := ts.repo.List(ctx, ownerID)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	retTasks := []dto.TaskResponse{}

	for _, task := range tasks {
		newTask := dto.TaskResponse{
			ID:          task.ID.Hex(),
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Status:      task.Status,
		}
		retTasks = append(retTasks, newTask)
	}
	return retTasks, nil
}
func (ts *TaskService) GetTask(ctx context.Context, id string, owner string) (dto.TaskResponse, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	ownerID, err := bson.ObjectIDFromHex(owner)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	var task models.Task

	task, err = ts.repo.Get(ctx, objectID)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	if task.Owner != ownerID {
		return dto.TaskResponse{}, ErrUnauthorized
	}

	retTask := dto.TaskResponse{
		ID:          task.ID.Hex(),
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Status:      task.Status,
	}

	return retTask, nil
}

func (ts *TaskService) UpdateTask(ctx context.Context, id string, owner string, updatedTask dto.TaskUpdateRequest) (dto.TaskResponse, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	ownerID, err := bson.ObjectIDFromHex(owner)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	task, err := ts.repo.Get(ctx, objectID)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	if task.Owner != ownerID {
		return dto.TaskResponse{}, ErrUnauthorized
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
		return dto.TaskResponse{}, errors.New("no fields to update")
	}

	stat, err := ts.repo.Update(ctx, objectID, set)
	if err != nil {
		return dto.TaskResponse{}, err
	}

	if stat {
		task, err = ts.repo.Get(ctx, objectID)
		if err != nil {
			return dto.TaskResponse{}, err
		}

		retTask := dto.TaskResponse{
			ID:          task.ID.Hex(),
			Title:       task.Title,
			Description: task.Description,
			DueDate:     task.DueDate,
			Status:      task.Status,
		}

		return retTask, nil
	}

	return dto.TaskResponse{}, errors.New("task not found")
}

func (ts *TaskService) DeleteTask(ctx context.Context, id string, owner string) error {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	ownerID, err := bson.ObjectIDFromHex(owner)
	if err != nil {
		return err
	}

	task, err := ts.repo.Get(ctx, objectID)
	if err != nil {
		return err
	}

	if task.Owner != ownerID {
		return ErrUnauthorized
	}

	deleted, err := ts.repo.Delete(ctx, objectID)
	if err != nil {
		return err
	}

	if !deleted {
		return errors.New("task not found")
	}

	return nil
}
