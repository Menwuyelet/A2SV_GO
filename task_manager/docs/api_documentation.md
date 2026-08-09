# Task Manager API
 
A complete set of REST API endpoints for managing tasks in a task management system. It covers the full lifecycle of a task including creation, retrieval, updating, and deletion. Task data is persisted in **MongoDB**.
 
**Base URL:** `http://localhost:8080`
**Authentication:** None required
**Database:** MongoDB (connection configured via environment variable)
 
---
 
## Table of Contents
 
- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Installation & Running the Project](#installation--running-the-project)
- [Getting Started](#getting-started)
- [Task Object](#task-object)
- [Endpoints](#endpoints)
  - [Create Task](#1-create-task)
  - [List Tasks](#2-list-tasks)
  - [Retrieve Task](#3-retrieve-task)
  - [Update Task](#4-update-task)
  - [Delete Task](#5-delete-task)
- [Response Status Codes](#response-status-codes)
---
 
## Prerequisites
 
Before running the project, make sure you have the following installed:
 
- **Go** (v1.21 or later recommended)
- **MongoDB** — either:
  - A local MongoDB instance running on your machine, or
  - A hosted MongoDB cluster (e.g. [MongoDB Atlas](https://www.mongodb.com/atlas))
---
 
## Environment Setup
 
The project uses environment variables to configure its connection to MongoDB. A `.env.example` file is provided in the project root as a template.
 
1. Duplicate the example file and rename it to `.env`:
```bash
   cp .env.example .env
```
 
2. Open `.env` and set the `MONGODB_URI` variable to your MongoDB connection string:
```env
   MONGODB_URI=mongodb://localhost:27017/task-manager
```
 
   | Variable       | Required | Description                                                                 |
   |----------------|----------|-------------------------------------------------------------------------------|
   | `MONGODB_URI`  | Yes      | The connection string used to connect to your MongoDB database (local or hosted). |
 
 
3. Save the `.env` file.
---
 
## Installation & Running the Project
 
Follow these steps to get the Task Manager API running locally:
 
1. **Clone the repository**
```bash
   git clone <repository-url>
   cd <repository-folder>
```
 
2. **Install dependencies**
```bash
   go mod download
```
 
3. **Configure environment variables**
   Follow the steps in [Environment Setup](#environment-setup) to create your `.env` file with a valid `MONGODB_URI`.
4. **Start MongoDB** (if running locally)
```bash
   mongod
```
 
   Skip this step if you're connecting to a hosted MongoDB instance.
 
5. **Run the server**
```bash
   go run main.go
```
 
6. **Verify the server is running**
   Once started, the API will be available at:
```
   http://localhost:8080
```
 
   You can test the connection by sending a request to the [List Tasks](#2-list-tasks) endpoint.
 
---
 
## Getting Started
 
1. Ensure the Task Manager server is running locally on port `8080` and connected to MongoDB (see [Installation & Running the Project](#installation--running-the-project)).
2. No authentication is required — requests can be sent directly.
3. Start with **Create Task** to add your first task, then use the remaining endpoints to list, retrieve, update, or delete tasks.
---
 
## Task Object
 
Each task is represented by the following fields:
 
| Field         | Type   | Description                                         |
|---------------|--------|------------------------------------------------------|
| `id`          | bson.ObjectID | Unique identifier for the task                       |
| `title`       | string | Short name or label for the task                     |
| `description` | string | Detailed information about the task                  |
| `due_date`    | string | Deadline in `YYYY-MM-DD` format                      |
| `status`      | string | Current state of the task (e.g. `completed`, `pending`) |
 
---
 
## Endpoints
 
### Summary
 
| Method   | Endpoint      | Description                  |
|----------|---------------|-------------------------------|
| `POST`   | `/tasks`      | Create a new task            |
| `GET`    | `/tasks`      | List all tasks               |
| `GET`    | `/tasks/{id}` | Retrieve a single task by ID |
| `PUT`    | `/tasks/{id}` | Update an existing task      |
| `DELETE` | `/tasks/{id}` | Delete a task by ID          |
 
---
 
### 1. Create Task
 
Creates a new task in the system.
 
```
POST /tasks
```
 
**Request Body**
 
| Field         | Type   | Required | Description                            |
|---------------|--------|----------|------------------------------------------|
| `title`       | string | Yes      | Task name                                |
| `description` | string | Yes      | Details about the task                   |
| `due_date`    | string | Yes      | Deadline in `YYYY-MM-DD` format          |
| `status`      | string | Yes      | Current state of the task (e.g. `completed`) |
 
**Example Request**
 
```json
POST http://localhost:8080/tasks
Content-Type: application/json
 
{
    "title": "test",
    "description": "Testst",
    "due_date": "2026-11-10",
    "status": "completed"
}
```
 
**Response**
 
Returns the created task object on success.
 
```json
200 OK
{
    "id": "6a78360b8e4cf6d96c48d993",
    "title": "test",
    "description": "Testst",
    "due_date": "2026-11-10",
    "status": "completed"
}
```
 
---
 
### 2. List Tasks
 
Retrieves a list of all tasks stored in the system.
 
```
GET /tasks
```
 
**Request Body:** None
 
**Example Request**
 
```
GET http://localhost:8080/tasks
```
 
**Response**
 
Returns an array of task objects.
 
```json
200 OK
[
    {
        "id": "6a78360b8e4cf6d96c48d993",
        "title": "test",
        "description": "Testst",
        "due_date": "2026-11-10",
        "status": "completed"
    },
    {
        "id": "6a78360b8e4cf6d96c48d994",
        "title": "Task 1",
        "description": "First task",
        "due_date": "2026-10-11",
        "status": "pending"
    }
]
```
 
---
 
### 3. Retrieve Task
 
Fetches the details of a single task by its unique identifier.
 
```
GET /tasks/{id}
```
 
**Path Parameters**
 
| Parameter | Type   | Description                       |
|-----------|--------|------------------------------------|
| `id`      | bson.ObjectID | The unique identifier of the task |
 
**Example Request**
 
```
GET http://localhost:8080/tasks/6a78360b8e4cf6d96c48d994
```
 
**Response**
 
Returns the full task object.
 
```json
200 OK
{
    "id": "6a78360b8e4cf6d96c48d994",
    "title": "Task 1",
    "description": "First task",
    "due_date": "2026-10-11",
    "status": "pending"
}
```
 
Returns a `404 Not Found` error if the task does not exist.
 
---
 
### 4. Update Task
 
Updates the details of an existing task.
 
```
PUT /tasks/{id}
```
 
**Path Parameters**
 
| Parameter | Type   | Description                       |
|-----------|--------|------------------------------------|
| `id`      | bson.ObjectID | The unique identifier of the task |
 
**Request Body**
 
| Field         | Type   | Required | Description                       |
|---------------|--------|----------|-------------------------------------|
| `title`       | string | No       | Updated task name                   |
| `description` | string | No       | Updated task details                |
| `due_date`    | string | No       | Updated deadline in `YYYY-MM-DD` format |
| `status`      | string | No       | Updated task state                  |
 
**Example Request**
 
```json
PUT http://localhost:8080/tasks/2
Content-Type: application/json
 
{
    "title": "Task 6a78360b8e4cf6d96c48d993",
    "description": "First task",
    "due_date": "2026-10-11",
    "status": "Completed"
}
```
 
**Response**
 
Returns the updated task object on success.
 
```json
200 OK
{
    "id": "6a78360b8e4cf6d96c48d993",
    "title": "Task 1",
    "description": "First task",
    "due_date": "2026-10-11",
    "status": "Completed"
}
```
 
---
 
### 5. Delete Task
 
Permanently removes a task from the system.
 
```
DELETE /tasks/{id}
```
 
**Path Parameters**
 
| Parameter | Type   | Description                       |
|-----------|--------|------------------------------------|
| `id`      | bson.ObjectID | The unique identifier of the task |
 
**Request Body:** None
 
**Example Request**
 
```
DELETE http://localhost:8080/tasks/6a78360b8e4cf6d96c48d993
```
 
**Response**
 
Returns a success confirmation.
 
```json
200 OK
{
    "message": "Task deleted successfully"
}
```
 
Returns a `404 Not Found` error if the task does not exist.
 
---
 
## Response Status Codes
 
| Status Code | Meaning       | Description                                  |
|-------------|---------------|-----------------------------------------------|
| `200`       | OK            | Request succeeded                            |
| `404`       | Not Found     | The requested task ID does not exist         |
| `500`       | Server Error  | An unexpected error occurred on the server   |
 
---
 
