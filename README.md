# Go Fiber API Template

A clean and well-structured Go API template following Go best practices and conventions. Built with [Fiber](https://github.com/gofiber/fiber) web framework and [GORM](https://gorm.io/) for database operations.

## 🚀 Features

- 🏗️ Clean architecture following Go conventions
- 📁 Standard Go project layout (`cmd/`, `internal/`, `handlers/`, etc.)
- 🔐 Environment variable configuration
- 🗄️ PostgreSQL database integration with GORM
- 🛣️ Organized routing system
- 🛡️ Middleware support
- 📦 Go modules for dependency management
- ✅ End-to-end testing with SQLite

## 📋 Prerequisites

- Go 1.23.5 or higher
- PostgreSQL database
- Basic understanding of Go programming

## 🛠️ Project Structure

```
.
├── cmd/
│   └── api/
│       └── main.go         # Application entry point
├── database/
│   └── database.go         # Database configuration
├── handlers/               # HTTP request handlers
│   └── dummy_handler.go
├── internal/
│   └── testutil/           # Test utilities (internal package)
│       └── setup.go
├── middlewares/            # Custom middleware functions
│   └── dummy_middleware.go
├── models/                 # Data models
│   └── dummy_user.go
├── routes/                 # Route definitions
│   ├── routes.go
│   └── dummy_routes.go
├── services/               # Business logic
│   └── dummy_service.go
├── test/
│   └── e2e/                # End-to-end tests
│       └── dummy_test.go
├── go.mod
├── go.sum
└── .env                    # Environment variables (not in repo)
```

### Directory Conventions

- **cmd/**: Main applications for this project
- **internal/**: Private application code (not importable by other projects)
- **handlers/**: HTTP handlers (Go community prefers "handlers" over "controllers")
- **services/**: Business logic layer
- **models/**: Data structures and database schemas
- **routes/**: Route definitions and groupings
- **middlewares/**: HTTP middleware functions
- **test/e2e/**: End-to-end integration tests

## 🚀 Getting Started

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd golang-fiber-base
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Create a `.env` file in the root directory:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_username
   DB_PASSWORD=your_password
   DB_NAME=your_database
   ```

4. Run the application:
   ```bash
   go run ./cmd/api
   ```

The server will start on `http://localhost:3000`

## 🧪 Testing

The project includes end-to-end tests using SQLite for isolated testing.

### Running Tests

Run all tests:
```bash
go test ./... -v
```

Run e2e tests only:
```bash
go test ./test/e2e/... -v
```

Run a specific test:
```bash
go test ./test/e2e/... -v -run TestDummyCRUD
```

### Test Utilities

The `internal/testutil` package provides helper functions:

- `testutil.SetupTestApp(t)`: Creates a new Fiber app with SQLite database
- `testutil.MakeRequest(t, app, method, path, body)`: Makes HTTP requests
- `testutil.ParseResponseBody(t, resp, v)`: Parses response body into a struct
- `testutil.CleanupTestApp(t)`: Cleans up after tests

Example:
```go
func TestYourFeature(t *testing.T) {
    app := testutil.SetupTestApp(t)
    defer testutil.CleanupTestApp(t)

    resp := testutil.MakeRequest(t, app, "GET", "/your-endpoint", nil)
    assert.Equal(t, 200, resp.StatusCode)
}
```

## 📚 Go Best Practices Applied

### File Naming
- All files use `snake_case.go` (e.g., `dummy_handler.go`, not `DummyHandler.go`)

### Package Naming
- Short, lowercase, single-word names
- No underscores or mixedCaps

### Project Layout
- Follows [Standard Go Project Layout](https://github.com/golang-standards/project-layout)
- `cmd/` for main applications
- `internal/` for private packages

### Testing
- Test utilities in `internal/testutil/`
- E2E tests separated in `test/e2e/`
- Uses `t.Helper()` in test helper functions

## 🛠️ Dependencies

- [Fiber v2](https://github.com/gofiber/fiber) - Fast HTTP framework
- [GORM](https://gorm.io/) - ORM for database operations
- [godotenv](https://github.com/joho/godotenv) - Environment variable loader
- [testify](https://github.com/stretchr/testify) - Testing utilities

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is open source and available under the [MIT License](LICENSE).
