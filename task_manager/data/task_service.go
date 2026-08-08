package data

import (
	"errors"
	"task_manager/models"
)

func AddTask(task models.Task) string {
	newTask := task

	models.Tasks = append(models.Tasks, newTask)

	return "Success"
}

func ListTasks() []models.Task {
	return models.Tasks
}

func GetTask(id string) (models.Task, error) {
	for _, task := range models.Tasks {
		if task.ID == id {
			return task, nil
		}
	}
	err := errors.New("Not found.")
	return models.Task{}, err
}

func UpdateTask(id string, updatedTask models.Task) (models.Task, error) {
	for idx, task := range models.Tasks {

		if task.ID == id {
			if updatedTask.Title != "" {

				models.Tasks[idx].Title = updatedTask.Title
			}
			if updatedTask.Description != "" {

				models.Tasks[idx].Description = updatedTask.Description
			}
			if updatedTask.DueDate != "" {

				models.Tasks[idx].DueDate = updatedTask.DueDate
			}
			if updatedTask.Status != "" {

				models.Tasks[idx].Status = updatedTask.Status
			}

			return models.Tasks[idx], nil
		}

	}

	err := errors.New("Task not found")

	return models.Task{}, err
}

func DeleteTask(id string) string {

	for i, val := range models.Tasks {
		if val.ID == id {
			models.Tasks = append(models.Tasks[:i], models.Tasks[i+1:]...)
			return "Success"
		}
	}

	return "Failed"

}
