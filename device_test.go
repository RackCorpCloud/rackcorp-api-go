package api

import (
	"context"
	"fmt"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeviceGet(t *testing.T) {
	defer gock.OffAll()

	const deviceId = 5075
	responseBody := getTestDataString(t, "device.get.responseBody.json")

	client := getTestClient(t)

	gock.New("https://api.rackcorp.net").
		Get(fmt.Sprintf("/api/v2.8/devices/%d", deviceId)).
		Reply(200).
		BodyString(responseBody)

	device, err := client.DeviceGet(context.TODO(), deviceId)
	assertGockNoUnmatchedRequests(t)
	assert.True(t, gock.IsDone(), "gock.IsDone")

	require.NoError(t, err, "DeviceGet error")
	assert.Equal(t, 5075, device.DeviceId, "DeviceId")
}

func TestDeviceUpdateFirewall(t *testing.T) {
	defer gock.OffAll()

	const deviceId = 678
	responseBody := "{\"code\": \"OK\", \"message\": \"good to go\"}"

	client := getTestClient(t)

	gock.New("https://api.rackcorp.net").
		Put(fmt.Sprintf("/api/v2.8/devices/%d/firewall", deviceId)).
		Reply(200).
		BodyString(responseBody)

	policies := []FirewallPolicy{
		{Direction: "INPUT"},
	}
	err := client.DeviceUpdateFirewall(context.TODO(), deviceId, policies)
	assertGockNoUnmatchedRequests(t)
	assert.True(t, gock.IsDone(), "gock.IsDone")

	require.NoError(t, err, "DeviceUpdateFirewall error")
}
