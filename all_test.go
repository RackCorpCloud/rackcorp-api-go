package api

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestClient(t *testing.T) Client {
	const uuid = "dummy-uuid"
	const secret = "dummy-secret"
	client, err := NewClient(uuid, secret)
	assert.Nil(t, err, "NewClient error")
	return client
}

func getTestDataString(t *testing.T, filename string) string {
	bytes, err := os.ReadFile(path.Join("testdata", filename)) //nolint:gosec
	assert.Nil(t, err, "ReadFile(%s) error", filename)
	return string(bytes)
}
