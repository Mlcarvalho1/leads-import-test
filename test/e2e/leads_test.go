package e2e

import (
	"testing"

	"leads-import/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHealth(t *testing.T) {
	app := testutil.SetupTestApp(t)
	defer testutil.CleanupTestApp(t)

	resp := testutil.MakeRequest(t, app, "GET", "/health", nil)
	assert.Equal(t, 200, resp.StatusCode)
}
