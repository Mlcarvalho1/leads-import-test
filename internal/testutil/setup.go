package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"your-app/database"
	"your-app/routes"

	"github.com/gofiber/fiber/v2"
)

// SetupTestEnv sets up the test environment variables (SQLite)
func SetupTestEnv(t *testing.T) {
	t.Helper()
	os.Setenv("DB_PATH", "test.db")
}

// CleanupTestEnv cleans up after tests
func CleanupTestEnv(t *testing.T) {
	t.Helper()
	if err := os.Remove("test.db"); err != nil {
		t.Logf("Failed to remove test database: %v", err)
	}
}

// SetupTestApp creates a new Fiber app for testing
func SetupTestApp(t *testing.T) *fiber.App {
	t.Helper()
	SetupTestEnv(t)

	app := fiber.New()
	database.ConnectDb()
	routes.SetupRoutes(app)

	return app
}

// CleanupTestApp cleans up after tests
func CleanupTestApp(t *testing.T) {
	t.Helper()
	CleanupTestEnv(t)
}

// MakeRequest makes an HTTP request to the test app
func MakeRequest(t *testing.T, app *fiber.App, method, path string, body interface{}) *http.Response {
	t.Helper()

	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	return resp
}

// ParseResponseBody parses the response body into the given interface
func ParseResponseBody(t *testing.T, resp *http.Response, v interface{}) {
	t.Helper()
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}
}
