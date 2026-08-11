package repository

import (
	"context"
	"task_manager/models"
)

type TaskRepository interface {
	AddTask(ctx context.Context, task models.Task) (models.Task, error)
	ListTasks(ctx context.Context) ([]models.Task, error)
	GetTask(ctx context.Context, id string) (models.Task, error)
	UpdateTask(ctx context.Context, id string, updatedTask models.Task) (models.Task, error)
	DeleteTask(ctx context.Context, id string) error
}
