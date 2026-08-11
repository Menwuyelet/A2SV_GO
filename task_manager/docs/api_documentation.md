# Task Manager API

A complete set of REST API endpoints for managing tasks in a task management system. It covers the full lifecycle of a task including creation, retrieval, updating, and deletion. Users authenticate via **JWT Bearer tokens** to access the task endpoints, and each task is scoped to the user who created it. Task data is persisted in **MongoDB**.

**Base URL:** `http://localhost:8080`
**Authentication:** JWT Bearer token (required for all task endpoints)
**Database:** MongoDB (connection configured via environment variable)
**Token Lifetime:** 15 minutes

---

## Table of Contents

- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Installation & Running the Project](#installation--running-the-project)
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
   SECRETE_KEY=your-random-secret-key
```

   | Variable       | Required | Description                                                                                |
   |----------------|----------|----------------------------------------------------------------------------------------------|
   | `MONGO_URI`    | Yes      | The connection string used to connect to your MongoDB database (local or hosted).            |
   | `SECRETE_KEY`  | Yes      | The secret used to sign and verify JSON Web Tokens. Use a long, random value.                |

3. Save the `.env` file.

> **Note:** The server will not start if the `.env` file cannot be loaded, and authentication will not work unless `SECRETE_KEY` is set.

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
   Follow the steps in [Environment Setup](#environment-setup) to create your `.env` file with a valid `MONGO_URI` and `SECRETE_KEY`.

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

A registered user record (returned by the user endpoints):

| Field  | Type   | Description                          |
|--------|--------|---------------------------------------|
| `id`   | bson.ObjectID | Unique identifier for the user  |
| `name` | string | Full name of the user                 |
| `email`| string | Email address used to log in          |
| `role` | string | Role assigned to the user (e.g. `user`) |

### Task Object

Each task is represented by the following fields:

| Field         | Type   | Description                                         |
|---------------|--------|------------------------------------------------------|
| `id`          | bson.ObjectID | Unique identifier for the task                       |
| `title`       | string | Short name or label for the task                     |
| `description` | string | Detailed information about the task                  |
| `due_date`    | string | Deadline in `YYYY-MM-DD` format                      |
| `status`      | string | Current state of the task. New tasks default to `PENDING`. |
| `owner`          | bson.ObjectID | Unique identifier for the owner of the task                  |

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
    "id": "6a7b6f9fae2bd9a2b42b25e8",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "user"
}
```

Returns a `400 Bad Request` error if validation fails or the email is already registered.

---

### 2. Login User

Authenticates a user and returns a JWT access token.

```
POST /api/user/login
```

**Request Body**

| Field      | Type   | Required | Description                 |
|------------|--------|----------|-------------------------------|
| `email`    | string | Yes      | Registered email address    |
| `password` | string | Yes      | Account password            |

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

Returns `401 Unauthorized` for invalid credentials.

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

Returns the created task object on success.

```json
201 Created
{
    "id": "6a78360b8e4cf6d96c48d993",
    "title": "Prepare project report",
    "description": "Compile weekly progress summary",
    "due_date": "2026-11-10",
    "status": "PENDING"
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
            "id": "6a78360b8e4cf6d96c48d993",
            "title": "Prepare project report",
            "description": "Compile weekly progress summary",
            "due_date": "2026-11-10",
            "status": "PENDING"
        },
        {
            "id": "6a78360b8e4cf6d96c48d994",
            "title": "Task 1",
            "description": "First task",
            "due_date": "2026-10-11",
            "status": "completed"
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
| `id`      | bson.ObjectID | The unique identifier of the task |

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
        "id": "6a78360b8e4cf6d96c48d994",
        "title": "Task 1",
        "description": "First task",
        "due_date": "2026-10-11",
        "status": "completed"
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
        "id": "6a78360b8e4cf6d96c48d993",
        "title": "Prepare project report",
        "description": "Compile weekly progress summary",
        "due_date": "2026-11-10",
        "status": "completed"
    }
}
```

Returns `404 Not Found` if the task does not exist, and `401 Unauthorized` if the task belongs to another user.

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
| `id`      | bson.ObjectID | The unique identifier of the task |

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