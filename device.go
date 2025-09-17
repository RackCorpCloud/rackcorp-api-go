package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/netip"
	"strings"

	"github.com/rackcorpcloud/rackcorp-api-go/apiv2"
)

type deviceGetAllPort struct {
	DeviceInterface int         `json:"deviceInterface"`
	ActivityInMbit  json.Number `json:"activityInMbit"`
	ActivityOutMbit json.Number `json:"activityOutMbit"`
	Status          string      `json:"status"` // "UP", "DOWN"
}

type deviceGetAllDevice struct {
	ID           json.Number `json:"id,omitempty"`
	Name         string      `json:"name,omitempty"`
	StdName      string      `json:"stdName,omitempty"`
	Status       string      `json:"status,omitempty"`
	TrafficMB    json.Number `json:"trafficMB,omitempty"`
	DateCreated  json.Number `json:"dateCreated,omitempty"`
	DateModified json.Number `json:"dateModified,omitempty"`
	// TODO Assets       []any       `json:"assets,omitempty"`
	// TODO IPs []any      `json:"ips,omitempty"`
	PrimaryIP            string             `json:"primaryIP,omitempty"`
	TrafficShared        bool               `json:"trafficShared,omitempty"`
	DCName               string             `json:"dcName,omitempty"`
	Type                 string             `json:"type,omitempty"`
	ProcessorUtilisation any                `json:"processorUtilisation,omitempty"`
	Ports                []deviceGetAllPort `json:"ports,omitempty"`
	OSGuess              string             `json:"osGuess,omitempty"`
	OSState              string             `json:"osState,omitempty"`
	OnlineStatus         string             `json:"onlineStatus,omitempty"`
	VMHostName           string             `json:"vmhostName,omitempty"`
	VMHostID             json.Number        `json:"vmhostId,omitempty"`
	// TODO VMHostProcessorUtilisation any    `json:"vmhostProcessorUtilisation,omitempty"`
	TrafficCurrent   json.Number `json:"trafficCurrent,omitempty"`
	TrafficEstimated json.Number `json:"trafficEstimated,omitempty"`
	DCID             json.Number `json:"dcId,omitempty"`
	DeviceID         json.Number `json:"deviceId,omitempty"`
	CustomerID       json.Number `json:"customerId,omitempty"`
	Extra            any         `json:"extra,omitempty"`
}

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

type deviceGetAllRequest struct {
	legacyRequest
	DeviceGetAllFilter
}

type deviceGetAllResponse struct {
	response
	Devices []deviceGetAllDevice `json:"devices"`
}

type DeviceGetAllFilter struct {
	ID         apiv2.DeviceID `json:"id,omitempty"`
	CustomerID CustomerID     `json:"customerID,omitempty"`
	Name       string         `json:"name,omitempty"`
	StdName    string         `json:"stdName,omitempty"`
	// TODO ipAddress
	// TODO dcID
	// TODO dcName
	// TODO status
	// TODO creationFromDate
	// TODO creationToDate
	// TODO assetID
	// TODO assetCode
	// TODO deviceType
	// TODO trafficShared
	// TODO hostDeviceID
	// TODO hostDeviceName
	// TODO transactionsPending
	ResultStart  int `json:"resStart,omitempty"`
	ResultWindow int `json:"resWindow,omitempty"`
	// TODO ordering
}

type DeviceClient interface {
	GetAll(ctx context.Context, filter DeviceGetAllFilter) ([]apiv2.Device, error)
}

type deviceClient struct {
	c *client
}

var _ DeviceClient = (*deviceClient)(nil)

func (dc *deviceClient) GetAll(ctx context.Context, filter DeviceGetAllFilter) ([]apiv2.Device, error) {
	req := deviceGetAllRequest{
		legacyRequest: legacyRequest{
			Command: "device.getall",
		},
		DeviceGetAllFilter: filter,
	}
	var resp deviceGetAllResponse
	err := dc.c.httpLegacyJson(ctx, &req, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsOK() {
		return nil, newApiError(resp.response, nil)
	}
	devices := make([]apiv2.Device, len(resp.Devices))
	for i, d := range resp.Devices {
		devices[i] = apiv2.Device{
			Name:    d.Name,
			StdName: d.StdName,
		}
		num, err := d.DeviceID.Int64()
		if err != nil {
			return nil, fmt.Errorf("failed to parse device ID %q as int: %w", d.DeviceID.String(), err)
		}
		devices[i].DeviceID = apiv2.DeviceID(num)
		num, err = d.CustomerID.Int64()
		if err != nil {
			return nil, fmt.Errorf("failed to parse customer ID %q as int: %w", d.CustomerID.String(), err)
		}
		devices[i].CustomerID = apiv2.CustomerID(num)
		pip := d.PrimaryIP
		if len(pip) > 0 {
			if strings.Contains(pip, "/") {
				// Strip off any CIDR suffix if present
				parts := strings.SplitN(pip, "/", 2)
				pip = parts[0]
			}
			addr, err := netip.ParseAddr(pip)
			if err != nil {
				return nil, fmt.Errorf("failed to parse primary IP %q as ip: %w", d.PrimaryIP, err)
			}
			devices[i].PrimaryIP = addr
		}
	}
	return devices, nil
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
