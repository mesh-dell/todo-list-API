# todo-list-API

This project is a Go-based RESTful API for managing todo lists, built with a clean, layered architecture. It follows the project requirements and challenges outlined by roadmap.sh.

### Inspiration

This project was inspired by the **Todo List API** project from [roadmap.sh/projects/todo-list-api](https://roadmap.sh/projects/todo-list-api).

---

### Installation and Setup

#### 1. Clone the Repository

```bash
git clone https://github.com/mesh-dell/todo-list-API.git
cd todo-list-API

```

#### 2. Install Dependencies

```bash
go mod tidy

```

#### 3. Environment Configuration

Create a `.env` file in the root directory and add the following:

```env
DB_USER=root
DB_NAME=yourDbName
DB_PASSWORD=yourPassword
DB_ADDR=localhost:3306
PORT=8080
ACCESS_SECRET=yourSecret
REFRESH_SECRET=yourRefreshSecret

```

#### 4. Run the Server

```bash
go run cmd/server/main.go

```

---

### Project Structure

* **cmd/server**: Application entry point.
* **internal/api**: API routing and rate-limiting middleware.
* **internal/auth**: User authentication, JWT management, and refresh token logic.
* **internal/todos**: Core logic for managing todo items.
* **internal/database**: GORM configuration and database connection.
* **internal/errors**: Custom error handling and standard error responses.

---

### API Endpoints

| Category | Method | Endpoint | Description |
| --- | --- | --- | --- |
| **Auth** | POST | /auth/register | Register a new user account |
| **Auth** | POST | /auth/login | Authenticate and retrieve tokens |
| **Auth** | POST | /auth/refresh | Refresh an expired access token |
| **Todos** | GET | /todos | List all todos for the authenticated user |
| **Todos** | POST | /todos | Create a new todo item |
| **Todos** | PUT | /todos/:id | Update an existing todo item |
| **Todos** | DELETE | /todos/:id | Delete a todo item |

---

### Testing

The project includes unit tests for services and repositories using mocks. To run the test suite, use the following command:

```bash
go test ./internal/...

```
