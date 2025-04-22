package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type Device struct {
	DeviceId         int              `json:"deviceId"`
	Name             string           `json:"name"`
	CustomerId       int              `json:"customerId"`
	PrimaryIP        string           `json:"primaryIP"`
	Status           string           `json:"status"`
	DataCenterId     int              `json:"dcid"`
	FirewallPolicies []FirewallPolicy `json:"firewallPolicies"`
	StdName          string           `json:"stdName"`
	DateCreated      int64            `json:"dateCreated"`
	DateModified     int64            `json:"dateModified"`
	TrafficShared    bool             `json:"trafficShared,omitempty"`
	TrafficCurrent   string           `json:"trafficCurrent"`
	TrafficEstimated float64          `json:"trafficEstimated"`
	TrafficMB        int64            `json:"trafficMB"`
	DCName           string           `json:"dcName"`
	// TODO assets, dcDescription, ips, networkRoutes, ports,

	Extra map[string]interface{} `json:"extra"`
}

type deviceGetResponse struct {
	response
	Device *Device `json:"data"`
}

type deviceUpdateRequest struct {
	FirewallPolicies []FirewallPolicy `json:"firewallPolicies"`
}

type deviceUpdateResponse struct {
	response
}

func (c *client) DeviceGet(ctx context.Context, deviceId int) (*Device, error) {
	if deviceId == 0 {
		return nil, errors.New("deviceId parameter is required")
	}

	var resp deviceGetResponse
	err := c.httpRestJson(ctx, http.MethodGet, fmt.Sprintf("devices/%d", deviceId), emptyRequest{}, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get device for device Id '%d': %w", deviceId, err)
	}

	if resp.Code != "OK" || resp.Device == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Device, nil
}

//	Note that if you want to delete an existing policy, you need to have it's policy set to DELETED
//
// (instead of ALLOW/REJECT/DISABLED) in the firewallPolicies array
func (c *client) DeviceUpdateFirewall(ctx context.Context, deviceId int, firewallPolicies []FirewallPolicy) error {
	if deviceId == 0 {
		return errors.New("deviceId parameter is required")
	}
	if len(firewallPolicies) == 0 {
		return errors.New("must update with Firewall Policies")
	}

	req := &deviceUpdateRequest{
		FirewallPolicies: firewallPolicies,
	}

	var resp deviceUpdateResponse
	err := c.httpRestJson(ctx, http.MethodPut, fmt.Sprintf("devices/%d/firewall", deviceId), req, &resp)
	if err != nil {
		return fmt.Errorf("failed to update firewall for device Id '%d': %w", deviceId, err)
	}

	if resp.Code != "OK" {
		return newApiError(resp.response, nil)
	}

	return nil
}
