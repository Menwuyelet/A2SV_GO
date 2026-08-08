package models

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Status      string `json:"status"`
}

var Tasks = []Task{
	{ID: "1", Title: "Task 1", Description: "First task", DueDate: "2026-10-11", Status: "Pending"},
	{ID: "2", Title: "Task 2", Description: "Second task", DueDate: "2026-10-11", Status: "In Progress"},
	{ID: "3", Title: "Task 3", Description: "Third task", DueDate: "2026-10-11", Status: "Completed"},
}
