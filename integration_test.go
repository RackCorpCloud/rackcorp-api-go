package api

import (
	"os"
	"testing"
)

func IntegrationTest(t testing.TB) {
	t.Helper()
	if len(os.Getenv("INTEGRATION_TEST")) == 0 {
		t.Skip("Skipping integration test, set INTEGRATION_TEST environment variable to enable.")
	}
}
