# Task Manager API

A complete set of REST API endpoints for managing tasks in a task management system. It covers the full lifecycle of a task including creation, retrieval, updating, and deletion. Users authenticate via **JWT Bearer tokens** to access the task endpoints, and each task is scoped to the user who created it. Task data is persisted in **MongoDB**.

The project follows a **Clean Architecture** structure (`Delivery` → `Usecases` → `repository` → `Domain`), uses the **Gin** web framework, **bcrypt** for password hashing, and the official MongoDB Go driver (v2).

**Base URL:** `http://localhost:8080`
**Authentication:** JWT Bearer token (required for all task endpoints)
**Database:** MongoDB — database `task_management_api`, collections `task` and `user`
**Token Lifetime:** 15 minutes
**Web Framework:** [Gin](https://github.com/gin-gonic/gin)

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Installation & Running the Project](#installation--running-the-project)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Authentication](#authentication)
- [Objects](#objects)
  - [User Object](#user-object)
  - [Task Object](#task-object)
- [Endpoints](#endpoints)
  - [Register User](#1-register-user)
  - [Login User](#2-login-user)
  - [Create Task](#3-create-task)
  - [List Tasks](#4-list-tasks)
  - [Retrieve Task](#5-retrieve-task)
  - [Update Task](#6-update-task)
  - [Delete Task](#7-delete-task)
- [Response Status Codes](#response-status-codes)

---

## Prerequisites

Before running the project, make sure you have the following installed:

- **Go** (v1.26 or later)
- **MongoDB** — either:
  - A local MongoDB instance running on your machine, or
  - A hosted MongoDB cluster (e.g. [MongoDB Atlas](https://www.mongodb.com/atlas))

---

## Environment Setup

The project uses environment variables to configure its MongoDB connection and the JWT signing secret. A `.env.example` file is provided in the project root as a template.

1. Duplicate the example file and rename it to `.env`:
```bash
   cp .env.example .env
```

2. Open `.env` and set the variables to match your environment:
```env
   MONGO_URI=mongodb://localhost:27017/task-manager
   SECRET_KEY=your-random-secret-key
```

   | Variable      | Required | Description                                                                                |
   |---------------|----------|----------------------------------------------------------------------------------------------|
   | `MONGO_URI`   | Yes      | The connection string used to connect to your MongoDB database (local or hosted).            |
   | `SECRET_KEY`  | Yes      | The secret used to sign and verify JSON Web Tokens. Use a long, random value.                |

3. Save the `.env` file.

> **Note:** The server will not start if the `.env` file cannot be loaded, and authentication will not work unless `SECRET_KEY` is set.


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
   Follow the steps in [Environment Setup](#environment-setup) to create your `.env` file with a valid `MONGO_URI` and `SECRET_KEY`.

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

   You can test the connection by sending a request to the [Login User](#2-login-user) or [Register User](#1-register-user) endpoint.

---

## Project Structure

The repository is organized following a layered (Clean Architecture) design:

```
task_manager/
├── Delivery/                          # HTTP layer: request/response handling
│   ├── controllers/
│   │   ├── task_controller.go         #   Task handlers (add, list, get, update, delete)
│   │   ├── task_controller_test.go    #   Unit tests for task handlers
│   │   ├── user_controller.go         #   User handlers (register, login)
│   │   └── user_controller_test.go    #   Unit tests for user handlers
│   └── routers/
│       ├── router.go                  #   Gin engine + public/protected route groups
│       ├── task_router.go             #   Task routes
│       └── user_router.go             #   User routes
├── Domain/                            # Core business types
│   ├── dto/
│   │   ├── task_dto.go                #   Task request/response structs
│   │   └── user_dto.go                #   User request/response structs
│   └── models/
│       ├── task.go                    #   Task model (persisted in MongoDB)
│       └── user.go                    #   User model (persisted in MongoDB)
├── Infrastructure/                     # External services & cross-cutting concerns
│   ├── middleware/
│   │   └── auth_middleware.go         #   JWT authentication middleware
│   └── utils/
│       ├── jwt.go                     #   JWT token generation
│       ├── password.go                #   bcrypt password hashing/verification
│       └── user_error_handler.go      #   Request validation error handler
├── Usecases/                          # Application / business logic
│   ├── task_service.go                #   Task use cases
│   ├── task_service_test.go           #   Unit tests for task service
│   ├── user_service.go                #   User use cases
│   └── user_service_test.go           #   Unit tests for user service
├── repository/                        # Data access layer (MongoDB)
│   ├── db_connenction.go              #   MongoDB connection & collection access
│   ├── errors.go                      #   Shared sentinel error (ErrNotFound)
│   ├── task_repository.go             #   Task repository interface + Mongo impl
│   └── user_repository.go             #   User repository interface + Mongo impl
├── Mocks/                             # Generated mocks for unit testing
│   ├── mock_UserRepository.go         #   Mock for the UserRepository interface
│   ├── TaskRepositoryMock.go          #   Mock for the TaskRepository interface
│   ├── TaskUsecaseMock.go             #   Mock for the TaskUsecase interface
│   └── UserUsecaseMock.go             #   Mock for the UserUsecase interface
├── docs/
│   └── api_documentation.md           #   This documentation
├── .env.example                       #   Environment variable template
├── .mockery.yaml                      #   Mockery config (generates Mocks/)
├── go.mod                             #   Go module definition & dependencies
├── go.sum                             #   Dependency checksums
└── main.go                            #   Entry point: wires everything together
```

| Layer          | Responsibility                                                                                                 |
|----------------|-----------------------------------------------------------------------------------------------------------------|
| `Delivery/`    | Handles incoming HTTP requests, binds/validates JSON bodies, and returns JSON responses (Gin controllers & routes). |
| `Domain/`      | Defines the core business entities (`models`) and the data-transfer objects (`dto`) exchanged with clients.       |
| `Usecases/`    | Implements application-specific business rules (task CRUD, user registration/login) on top of the repository.     |
| `repository/`  | Abstracts MongoDB data access through interfaces, implemented by `MongoTaskRepository` and `MongoUserRepository`. |
| `Infrastructure/` | Provides cross-cutting services: JWT auth middleware/token generation, bcrypt password utilities, validation errors. |
| `Mocks/`        | Mock implementations of the repository and usecase interfaces, used by the unit tests. |

Dependency flow: `main.go` builds the MongoDB connection, wires repositories → use cases → controllers, and registers the Gin routes in `Delivery/routers/router.go`.

---

## Getting Started

1. Ensure the Task Manager server is running locally on port `8080` and connected to MongoDB (see [Installation & Running the Project](#installation--running-the-project)).
2. Register a new user via [Register User](#1-register-user) to obtain an account.
3. Call [Login User](#2-login-user) to receive a JWT access token. This token is valid for **15 minutes**.
4. Send all task requests with the token included in the `Authorization` header as a bearer token:
```
   Authorization: Bearer <token>
```
5. Start with [Create Task](#3-create-task) to add your first task, then use the remaining endpoints to list, retrieve, update, or delete tasks.

---

## Authentication

The API uses **JWT (JSON Web Tokens)** for authentication.

- Tokens are issued by the [Login User](#2-login-user) endpoint with a validity of **15 minutes**.
- Every task endpoint requires the `Authorization` header:
```
   Authorization: Bearer <token>
```
- A request with a missing, malformed, or expired token is rejected with `401 Unauthorized`.
- Tasks are **owner-scoped**: a user can only create, view, update, and delete their own tasks. Accessing another user's task returns `401 Unauthorized`.

---

## Objects

### User Object

A registered user record (returned by the user endpoints). Field names appear **exactly as in the Go response struct** (`UserResponse` has no JSON tags), i.e. PascalCase:

| Field  | Type   | Description                          |
|--------|--------|---------------------------------------|
| `ID`   | string (ObjectID hex) | Unique identifier for the user  |
| `Name` | string | Full name of the user                 |
| `Email`| string | Email address used to log in          |
| `Role` | string | Role assigned to the user (`user`) |

> Internal fields such as the bcrypt `PasswordHash` are never returned. Note that **request** structs (`RegisterRequest`, `LoginRequest`) do carry lowercase `json` tags — responses do not.

### Task Object

Each task is represented by the following fields. As with users, field names are returned in **PascalCase**:

| Field         | Type   | Description                                         |
|---------------|--------|------------------------------------------------------|
| `ID`          | string (ObjectID hex) | Unique identifier for the task                       |
| `Title`       | string | Short name or label for the task                     |
| `Description` | string | Detailed information about the task                  |
| `DueDate`     | string | Deadline in `YYYY-MM-DD` format                      |
| `Status`      | string | Current state of the task. New tasks default to `PENDING`. |

> The `owner` field exists on the persisted `models.Task` (in MongoDB) and is used for ownership checks, but it is **not** exposed in API responses.

---

## Endpoints

### Summary

| Method   | Endpoint                  | Auth   | Description                        |
|----------|---------------------------|--------|-------------------------------------|
| `POST`   | `/api/user/register`      | No     | Register a new user                |
| `POST`   | `/api/user/login`         | No     | Log in and receive a JWT token     |
| `POST`   | `/api/tasks`              | Yes    | Create a new task                  |
| `GET`    | `/api/tasks`              | Yes    | List the current user's tasks      |
| `GET`    | `/api/tasks/{id}`         | Yes    | Retrieve a single task by ID       |
| `PUT`    | `/api/tasks/{id}`         | Yes    | Update an existing task            |
| `DELETE` | `/api/tasks/{id}`         | Yes    | Delete a task by ID                |

---

### 1. Register User

Creates a new user account.

```
POST /api/user/register
```

**Request Body**

| Field      | Type   | Required | Description                            |
|------------|--------|----------|------------------------------------------|
| `name`     | string | Yes      | Full name of the user (max 50 characters) |
| `email`    | string | Yes      | Valid email address                    |
| `password` | string | Yes      | Password (min 8 characters)            |

> Validation is enforced via struct binding tags (`required`, `max=50`, `email`, `min=8`). On failure the API returns `400` with a `details` map of per-field messages.

**Example Request**

```json
POST http://localhost:8080/api/user/register
Content-Type: application/json

{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
}
```

**Response**

Returns the created user object on success.

```json
201 Created
{
    "ID": "6a7b6f9fae2bd9a2b42b25e8",
    "Name": "John Doe",
    "Email": "john@example.com",
    "Role": "user"
}
```

Returns `400 Bad Request` if JSON binding or field validation fails (e.g. missing field, name over 50 characters, invalid email, password shorter than 8 characters) or if the email is already registered.

---

### 2. Login User

Authenticates a user and returns a JWT access token.

```
POST /api/user/login
```

**Request Body**

| Field      | Type   | Required | Description                        |
|------------|--------|----------|-------------------------------------|
| `email`    | string | Yes      | Registered email address (must be a valid email) |
| `password` | string | Yes      | Account password                   |

**Example Request**

```json
POST http://localhost:8080/api/user/login
Content-Type: application/json

{
    "email": "john@example.com",
    "password": "password123"
}
```

**Response**

Returns a token, valid for **15 minutes**.

```json
200 OK
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODY0NzUzMDAsInJvbGUiOiJ1c2VyIiwidXNlcl9pZCI6IjZhN2I2ZjlmYWUyYmQ5YTJiNDJiMjVlOCJ9.e8MsndN8FLD_j6kIwu8tCrjTm5XjyjLpV_y-RjHrA2g"
}
```

Use the returned token in the `Authorization` header for all task endpoints:

```
Authorization: Bearer <token>
```

Returns `400 Bad Request` if JSON binding or validation fails (e.g. missing or invalid email/password), and `401 Unauthorized` if the email is not registered or the password is incorrect.

---

### 3. Create Task

Creates a new task owned by the authenticated user.

```
POST /api/tasks
Authorization: Bearer <token>
```

**Request Body**

| Field         | Type   | Required | Description                    |
|---------------|--------|----------|----------------------------------|
| `title`       | string | Yes      | Short name or label for the task |
| `description` | string | No       | Details about the task           |
| `due_date`    | string | No       | Deadline in `YYYY-MM-DD` format  |

> The `status` field is set to `PENDING` automatically and cannot be provided at creation time.

**Example Request**

```json
POST http://localhost:8080/api/tasks
Authorization: Bearer <token>
Content-Type: application/json

{
    "title": "Prepare project report",
    "description": "Compile weekly progress summary",
    "due_date": "2026-11-10"
}
```

**Response**

Returns the created task object on success. Note the response keys are PascalCase (`ID`, `DueDate`, etc.) as exposed by the DTO.

```json
201 Created
{
    "ID": "6a78360b8e4cf6d96c48d993",
    "Title": "Prepare project report",
    "Description": "Compile weekly progress summary",
    "DueDate": "2026-11-10",
    "Status": "PENDING"
}
```

Returns `400 Bad Request` if the request body is invalid, and `500 Internal Server Error` if the task could not be persisted.

---

### 4. List Tasks

Retrieves the tasks belonging to the authenticated user.

```
GET /api/tasks
Authorization: Bearer <token>
```

**Request Body:** None

**Example Request**

```
GET http://localhost:8080/api/tasks
Authorization: Bearer <token>
```

**Response**

Returns an object containing an array of the current user's task objects.

```json
200 OK
{
    "tasks": [
        {
            "ID": "6a78360b8e4cf6d96c48d993",
            "Title": "Prepare project report",
            "Description": "Compile weekly progress summary",
            "DueDate": "2026-11-10",
            "Status": "PENDING"
        },
        {
            "ID": "6a78360b8e4cf6d96c48d994",
            "Title": "Task 1",
            "Description": "First task",
            "DueDate": "2026-10-11",
            "Status": "completed"
        }
    ]
}
```

---

### 5. Retrieve Task

Fetches the details of a single task by its unique identifier.

```
GET /api/tasks/{id}
Authorization: Bearer <token>
```

**Path Parameters**

| Parameter | Type   | Description                       |
|-----------|--------|------------------------------------|
| `id`      | string (ObjectID hex) | The unique identifier of the task |

**Example Request**

```
GET http://localhost:8080/api/tasks/6a78360b8e4cf6d96c48d994
Authorization: Bearer <token>
```

**Response**

Returns the full task object.

```json
200 OK
{
    "task": {
        "ID": "6a78360b8e4cf6d96c48d994",
        "Title": "Task 1",
        "Description": "First task",
        "DueDate": "2026-10-11",
        "Status": "completed"
    }
}
```

Returns `404 Not Found` if the task does not exist, and `401 Unauthorized` if the task belongs to another user.

---

### 6. Update Task

Updates the details of an existing task. Any of the fields may be omitted to update only the provided ones.

```
PUT /api/tasks/{id}
Authorization: Bearer <token>
```

**Path Parameters**

| Parameter | Type   | Description                       |
|-----------|--------|------------------------------------|
| `id`      | string (ObjectID hex) | The unique identifier of the task |

**Request Body**

| Field         | Type   | Required | Description                       |
|---------------|--------|----------|-------------------------------------|
| `title`       | string | No       | Updated task name                   |
| `description` | string | No       | Updated task details                |
| `due_date`    | string | No       | Updated deadline in `YYYY-MM-DD` format |
| `status`      | string | No       | Updated task state                  |

**Example Request**

```json
PUT http://localhost:8080/api/tasks/6a78360b8e4cf6d96c48d993
Authorization: Bearer <token>
Content-Type: application/json

{
    "title": "Prepare project report",
    "status": "completed"
}
```

**Response**

Returns the updated task object on success.

```json
200 OK
{
    "task": {
        "ID": "6a78360b8e4cf6d96c48d993",
        "Title": "Prepare project report",
        "Description": "Compile weekly progress summary",
        "DueDate": "2026-11-10",
        "Status": "completed"
    }
}
```

Returns `404 Not Found` if the task does not exist, and `401 Unauthorized` if the task belongs to another user. Returns `400 Bad Request` for a malformed request body. All other service errors (including when the body contains no updateable fields) return `404 Not Found`.

---

### 7. Delete Task

Permanently removes a task from the system.

```
DELETE /api/tasks/{id}
Authorization: Bearer <token>
```

**Path Parameters**

| Parameter | Type   | Description                       |
|-----------|--------|------------------------------------|
| `id`      | string (ObjectID hex) | The unique identifier of the task |

**Request Body:** None

**Example Request**

```
DELETE http://localhost:8080/api/tasks/6a78360b8e4cf6d96c48d993
Authorization: Bearer <token>
```

**Response**

Returns `204 No Content` on success (no response body).

Returns `404 Not Found` if the task does not exist, and `401 Unauthorized` if the task belongs to another user.

---

## Response Status Codes

| Status Code | Meaning       | Description                                  |
|-------------|---------------|-----------------------------------------------|
| `200`       | OK            | Request succeeded                            |
| `201`       | Created       | Resource created successfully                |
| `204`       | No Content    | Request succeeded, no response body returned |
| `400`       | Bad Request   | Invalid or malformed request body/validation |
| `401`       | Unauthorized  | Missing/invalid/expired token, or access to another user's task |
| `404`       | Not Found     | The requested task does not exist            |
| `500`       | Server Error  | An unexpected error occurred on the server   |

---