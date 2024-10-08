# Todolist API

## Project Description

This project is a simple **Todolist API** that provides basic task management features. It allows users to register, log in, and perform CRUD (Create, Read, Update, Delete) operations on tasks. It also supports token-based authentication using JWTs for user management and security.

The API is built following the REST architecture and returns responses in JSON format.

### Key Features
- **User Registration and Authentication**: Users can register and log in with secure password management using JWT tokens.
- **Task Management**: Users can add, update, fetch, and delete their tasks.
- **Pagination**: Supports pagination for retrieving tasks.
- **JWT-Based Authentication**: Secure routes with token authentication for task management.

## Getting Started

### Prerequisites
To run this project, you’ll need:
- **Go** (version 1.XX or later)
- **Database** (e.g., PostgreSQL or MySQL) set up for task persistence
- **Postman** or **cURL** (optional) for testing the API

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/NeGat1FF/todolist-api.git
   cd todolist-api
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables for your database and JWT secret. You can use a `.env` file or directly set them up in your environment:
   ```bash
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=your_db_user
   export DB_PASSWORD=your_db_password
   export DB_NAME=todolist
   export SECRET_KEY=your_secret_key
   ```

4. Start the server:
   ```bash
   go run ./cmd/todolist-api/main.go
   ```

### API Documentation

The API is available at `http://localhost:8080/` and exposes the following endpoints:

#### Authentication
- **POST /users/register**: Register a new user.
- **POST /users/login**: Log in an existing user.
- **POST /users/refresh**: Refresh access token.

#### Tasks
- **GET /tasks**: Retrieve all tasks (supports pagination).
- **POST /tasks**: Create a new task.
- **PUT /tasks/{id}**: Update a specific task by ID.
- **DELETE /tasks/{id}**: Delete a specific task by ID.

### Example API Requests

1. **User Registration**:
   ```bash
   curl -X POST http://localhost:8080/users/register \
   -H "Content-Type: application/json" \
   -d '{"username":"johndoe", "email":"johndoe@example.com", "password":"mypassword"}'
   ```

2. **Create a Task** (requires JWT Token in Authorization Header):
   ```bash
   curl -X POST http://localhost:8080/tasks \
   -H "Authorization: Bearer YOUR_JWT_TOKEN" \
   -H "Content-Type: application/json" \
   -d '{"title":"Buy Groceries", "description":"Buy eggs, milk, and bread"}'
   ```

3. **Get All Tasks**:
   ```bash
   curl -X GET http://localhost:8080/tasks?page=1&limit=10 \
   -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

4. **Delete a Task**:
   ```bash
   curl -X DELETE http://localhost:8080/tasks/1 \
   -H "Authorization: Bearer YOUR_JWT_TOKEN"
   ```

### Response Format

All responses are in JSON format. Example response for getting tasks:
```json
{
  "data": [
    {
      "id": 1,
      "title": "Buy Groceries",
      "description": "Buy eggs, milk, and bread"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 1
}
```

## Authentication

JWT (JSON Web Token) is used for securing the API. After successful login, a JWT token is provided that should be included in the `Authorization` header for all protected routes. Example:
```
Authorization: Bearer <JWT-TOKEN>
```

## Future Improvements
- Add more advanced filtering for tasks (e.g., by date or status).
- Implement user-specific task management so tasks are tied to specific users.
- Add role-based access control (RBAC) for different user roles.

https://roadmap.sh/projects/todo-list-api
