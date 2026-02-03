package e2e

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"your-app/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHealthAndLeadsImport(t *testing.T) {
	app := testutil.SetupTestApp(t)
	defer testutil.CleanupTestApp(t)

	// Health
	resp := testutil.MakeRequest(t, app, "GET", "/health", nil)
	assert.Equal(t, 200, resp.StatusCode)

	// Lead import (same app/DB)
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	part, _ := w.CreateFormFile("file", "leads.csv")
	_, _ = part.Write([]byte("name,email,phone\nJoão,john@example.com,11999999999\nMaria,maria@example.com,"))
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost, "/leads/import", body)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
}
