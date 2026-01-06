# Bailu Backend

<p align="center">
  <strong>A lightweight, production-ready, and feature-rich backend boilerplate.</strong>
</p>
<p align="center">
  Built with Go, Gin, GORM, and Wire, Bailu is designed to help you quickly launch a secure and scalable admin panel, RESTful API, or microservice.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.21+-blue.svg" alt="Go Version">
  <img src="https://img.shields.io/badge/gin-v1.10.0-blue.svg" alt="Gin Version">
  <img src="https://img.shields.io/badge/gorm-v1.25.11-orange.svg" alt="Gorm Version">
  <img src="https://img.shields.io/badge/casbin-v2.99.0-green.svg" alt="Casbin Version">
  <img src="https://img.shields.io/badge/wire-v0.6.0-purple.svg" alt="Wire Version">
  <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License">
</p>

<div align="center">
  <strong><a href="./README.zh-CN.md">中文</a></strong> | <strong>English</strong>
</div>

---

## ⚙️ Technology Stack

-   **Web Framework**: Gin
-   **ORM**: GORM
-   **Database**: MySQL
-   **Dependency Injection**: Google Wire
-   **Permissions Management**: Casbin
-   **Authentication**: JWT
-   **Configuration Management**: Viper
-   **Logging**: Zap


## ✨ Usage Examples

Bailu is designed for developer productivity. Here are a few examples of its convenient features:

### 1. Effortless Pagination

Simply add one line in your controller to enable pagination for any query.

```go
// In your API controller function:
func (a *UserAPI) GetUserList(c *gin.Context) {
    // Automatically applies limit and offset from query parameters (e.g., ?page=1&pageSize=10)
    page.StartPage(c) 
    
    // Your data retrieval logic follows
    users, err := a.userService.ListByBuilder(c)
    if err != nil {
        resp.FailWithError(c, err)
        return
    }
    
    // The response is automatically wrapped in a pagination structure
    resp.OKWithData(c, users)
}
```

### 2. Unified Response Wrapper

Standardize your API responses with simple, expressive helpers.

```go
// For successful responses:
resp.OK(c) // Returns a standard success message
resp.OKWithData(c, data) // Returns success with a data payload

// For error responses:
resp.Fail(c) // Returns a standard failure message
resp.FailWithError(c, someError) // Returns failure with a specific error

// For more complex scenarios, you can panic with a response error.
// A global recovery middleware will catch it and format the JSON response.
if user == nil {
    panic(resp.ErrNotFound)
}
if err != nil {
    panic(resp.InternalServerErrorWithError(err))
}
```

### 3. Automatic Query Builder

Build complex GORM queries directly from your request DTOs using struct tags. This eliminates boilerplate `db.Where()` clauses.

```go
// 1. Define query parameters in your DTO with `query` tags.
//    Format: `query:"[column_name],[operator]"`
//    Supported operators: eq, neq, gt, gte, lt, lte, like, in
type UserQueryParams struct {
    dto.Pagination
    Username string `form:"username" query:"username,like"`
    Email    string `form:"email" query:"email,eq"`
    Status   int    `form:"status" query:"status,eq"`
}

// 2. Use the QueryBuilder in your repository layer.
func (r *UserRepo) FindByParams(ctx context.Context, params *dto.UserQueryParams) ([]*entity.User, error) {
    // The builder automatically constructs the WHERE clause.
    // e.g., WHERE username LIKE '%...%' AND status = ?
    builder := base.NewQueryBuilder().WithWhereStruct(params)
    
    var users []*entity.User
    err := r.FindByBuilder(ctx, builder).Find(&users).Error
    return users, err
}
```

## Frontend Project

- **Bailu Admin (Vue)**: A companion frontend project is under development. (Link to be added)
- **Live Demo**: (Link to be added)
- **Default Credentials**: `sufo` / `admin123`

## 🚀 Getting Started

Follow these steps to get your local development environment up and running.

### Prerequisites

- [Go](https://golang.org/dl/) 1.21+
- [MySQL](https://www.mysql.com/downloads/) 5.7+
- [Make](https://www.gnu.org/software/make/)
- [Wire](https://github.com/google/wire): `go install github.com/google/wire/cmd/wire@latest`
- [Swag](https://github.com/swaggo/swag): `go install github.com/swaggo/swag/cmd/swag@latest`

### Installation & Running

1.  **Clone the repository:**
    ```shell
    git clone https://github.com/sufo/bailu-backend.git
    cd bailu-backend
    ```

2.  **Configure the application:**
    -   Copy the development configuration file: `cp config/config.dev.yml config/config.yml`.
    -   Edit `config/config.yml` and update the `mysql` section with your database credentials.

3.  **Initialize the database:**
    -   Create a new database in MySQL (e.g., `bailu`).
    -   Import the initial schema and data from `sql/init_mysql.sql`.

4.  **Generate dependency injection code:**
    ```shell
    make wire
    ```

5.  **Generate API documentation:**
    ```shell
    make swagger
    ```

6.  **Run the server:**
    ```shell
    make start
    ```
    The server will start on the port specified in your config (default: `8081`).

7.  **Access API Docs:**
    Visit `http://localhost:8081/swagger/index.html` to view the interactive API documentation.

## 🐳 Docker Quick Start

1.  **Build the Docker image:**
    ```shell
    make build-image-server TAGS_OPT=latest
    ```

2.  **Run the container:**
    Make sure your `config/config.docker.yml` is correctly configured to connect to your database.
    ```shell
    docker run -d -p 8081:8081 --name bailu-server bailu-server:latest
    ```

## 🧰 Makefile Commands

This project uses `make` to simplify common tasks.

- `make start`: Start the application in development mode.
- `make build`: Build the application binary.
- `make wire`: Generate dependency injection code.
- `make swagger`: Generate Swagger API documentation.
- `make stop`: Stop the running application.
- `make build-image-server`: Build the Docker image.

## 📂 Project Structure

The project follows a modular, layered architecture to promote separation of concerns and maintainability.

```
/
├── app/                # Core application code
│   ├── api/            # API controllers and routing
│   ├── config/         # Structs for configuration
│   ├── core/           # Core components (server engine, DI)
│   ├── domain/         # Domain models (entities, DTOs, repos)
│   ├── middleware/     # Gin middleware
│   ├── service/        # Business logic layer
│   └── ...
├── config/             # Configuration files (YAML, etc.)
├── global/             # Global variables and constants
├── pkg/                # Shared utility packages
├── sql/                # SQL initialization scripts
├── utils/              # General utility functions
├── main.go             # Application entry point
├── go.mod              # Go module definitions
├── Makefile            # Makefile for common tasks
└── Dockerfile          # Docker build definition
```

## 📄 License

This project is [MIT](./LICENSE) licensed.