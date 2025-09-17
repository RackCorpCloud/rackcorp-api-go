package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeviceGetAll(t *testing.T) {
	IntegrationTest(t)
	client, err := NewClientFromEnv()
	client.SetDebugLog(func(m string) {
		t.Logf("Client Debug: %s", m)
	})
	require.NoError(t, err, "NewClientFromEnv")

	filter := DeviceGetAllFilter{}
	devices, err := client.Device().GetAll(t.Context(), filter)
	require.NoError(t, err, "DeviceGetAll")
	require.NotEmpty(t, devices, "Expected non-empty device list")
	t.Logf("Got %d devices", len(devices))
}
