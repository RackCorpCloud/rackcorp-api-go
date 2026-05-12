package api

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationLoadBalancerGetAll(t *testing.T) {
	IntegrationTest(t)
	client, err := NewClientFromEnv()
	require.NoError(t, err, "NewClientFromEnv")
	client.SetDebugLog(func(m string) {
		t.Logf("Client Debug: %s", m)
	})

	filter := LoadBalancerFilter{}
	lbs, err := client.LoadBalancer().GetAll(t.Context(), filter)
	require.NoError(t, err, "LoadBalancer().GetAll()")
	require.NotEmpty(t, lbs, "Expected non-empty load balancer list")
	t.Logf("Got %d load balancers", len(lbs))

	for _, lb := range lbs {
		lbDetail, err := client.LoadBalancer().Get(t.Context(), lb.ID)
		require.NoError(t, err, "LoadBalancer().Get()")
		require.Equal(t, lb.ID, lbDetail.ID, "LoadBalancer ID should match")
		t.Logf("LoadBalancer %d details: %+v", lb.ID, lbDetail)
		break
	}
}

func TestIntegrationLoadBalancerCreateUpdateDelete(t *testing.T) {
	IntegrationTest(t)
	client, err := NewClientFromEnv()
	require.NoError(t, err, "NewClientFromEnv")
	client.SetDebugLog(func(m string) {
		t.Logf("Client Debug: %s", m)
	})

	desiredLB := LoadBalancer{
		BalanceMode: LoadBalancerBalanceModeRoundRobin,
		CheckMode:   LoadBalancerCheckModeTCP,

		Name:  t.Name(),
		Scope: LoadBalancerScopeGlobal,
		Type:  LoadBalancerTypeTCP,
		Ports: []int{80, 443},
		Backends: []LoadBalancerBackend{
			{
				Hostname: "192.0.2.0",
				Name:     "backendhttp",
				Port:     8080,
				PortMask: []int{80},
				Weight:   100,
				TLS:      false,
			},
			{
				Hostname: "192.0.2.0",
				Name:     "backendhttps",
				Port:     8443,
				PortMask: []int{443},
				Weight:   100,
				TLS:      false,
			},
		},
	}

	createdLB, err := client.LoadBalancer().Create(t.Context(), desiredLB)
	require.NoError(t, err, "LoadBalancer().Create()")
	require.NotNil(t, createdLB, "Created LoadBalancer should not be nil")
	t.Logf("Created LoadBalancer: %+v", createdLB)

	id := createdLB.ID
	defer func() {
		err := client.LoadBalancer().Delete(t.Context(), id)
		require.NoError(t, err, "LoadBalancer().Delete()")
	}()

	require.Equal(t, desiredLB.Name, createdLB.Name, "LoadBalancer name should match desired")
	require.Equal(t, desiredLB.Scope, createdLB.Scope, "LoadBalancer scope should match desired")
	require.Equal(t, desiredLB.Type, createdLB.Type, "LoadBalancer type should match desired")
	require.Equal(t, desiredLB.Ports, createdLB.Ports, "LoadBalancer ports should match desired")
	require.Equal(t, len(desiredLB.Backends), len(createdLB.Backends), "LoadBalancer backends count should match desired")
	for i, backend := range desiredLB.Backends {
		createdBackend := createdLB.Backends[i]
		require.Equal(t, backend.Name, createdBackend.Name, "Backend name should match desired")
		require.Equal(t, backend.Port, createdBackend.Port, "Backend port should match desired")
		require.Equal(t, backend.Weight, createdBackend.Weight, "Backend weight should match desired")
		require.Equal(t, backend.TLS, createdBackend.TLS, "Backend TLS should match desired")
		require.Equal(t, backend.PortMask, createdBackend.PortMask, "Backend port mask should match desired")
	}

	desiredChangeLB := *createdLB
	desiredChangeLB.Name = t.Name() + "-updated"

	updatedLB, err := client.LoadBalancer().Update(t.Context(), desiredChangeLB)
	require.NoError(t, err, "LoadBalancer().Update()")
	require.NotNil(t, updatedLB, "Updated LoadBalancer should not be nil")
	t.Logf("Updated LoadBalancer: %+v", updatedLB)

	require.Equal(t, desiredChangeLB.Name, updatedLB.Name, "LoadBalancer name should be updated")

	fetchedLB, err := client.LoadBalancer().Get(t.Context(), id)
	require.NoError(t, err, "LoadBalancer().Get()")
	require.NotNil(t, fetchedLB, "Fetched LoadBalancer should not be nil")
	t.Logf("Fetched LoadBalancer: %+v", fetchedLB)

	require.Equal(t, updatedLB.ID, fetchedLB.ID, "Fetched LoadBalancer ID should match updated")
	require.Equal(t, updatedLB.Name, fetchedLB.Name, "Fetched LoadBalancer name should match updated")
}
