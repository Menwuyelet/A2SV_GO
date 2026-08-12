package dto

type TaskRequest struct {
	Title       string
	Description string
	DueDate     string
}

type TaskUpdateRequest struct {
	Title       string
	Description string
	DueDate     string
	Status      string
}

type TaskResponse struct {
	ID          string
	Title       string
	Description string
	DueDate     string
	Status      string
}
