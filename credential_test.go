package api

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApiCredentialFromIni(t *testing.T) {
	contents := "[general]\napiuuid=dummy-uuid\napisecret=dummy-secret\n"
	r := strings.NewReader(contents)
	cred, err := newApiCredentialFromIni(r)
	require.NoError(t, err, "newApiCredentialFromIni")
	require.NotNil(t, cred, "cred")
	assert.Equal(t, "dummy-uuid", cred.UUID, "cred.UUID")
	assert.Equal(t, "dummy-secret", cred.Secret, "cred.Secret")
}
