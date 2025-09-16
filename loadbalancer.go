package api

import (
	"context"
	"time"
)

type LoadBalancerID int

type LoadBalancerScope string

const (
	LoadBalancerScopeGlobal LoadBalancerScope = "global"
	LoadBalancerScopeLocal  LoadBalancerScope = "local"
)

type existingLoadBalancer struct {
	ID                   LoadBalancerID    `json:"id"`
	Name                 string            `json:"name"`
	CustomerID           CustomerID        `json:"customerid"`
	Hostname             string            `json:"hostname"`
	StdName              string            `json:"stdname"`
	HostSource           string            `json:"hostsource"`
	HostSourceForceHost  string            `json:"hostsourceforcehost"`
	BackendHostnameForce string            `json:"backend_hostname_force"`
	Status               string            `json:"status"`
	DateCreated          int64             `json:"datecreated"`
	DateModified         int64             `json:"datemodified"`
	Version              int               `json:"version"`
	Scope                LoadBalancerScope `json:"scope"`
	ScopeNetworkID       NetworkID         `json:"scope_networkid"`
	ScopeInstances       int               `json:"scope_instances"`
	MonthlyAllocationMB  int               `json:"monthlyallocationmb"`
	MonthlyUsageMB       int               `json:"monthlyusagemb"`
	TrafficRemainingMB   int               `json:"trafficremainingmb"`
}

type LoadBalancer struct {
	ID                   LoadBalancerID
	Name                 string
	CustomerID           CustomerID
	Hostname             string
	StdName              string
	HostSource           string
	HostSourceForceHost  string
	BackendHostnameForce string
	Status               string // TODO LoadBalancerStatus
	DateCreated          time.Time
	DateModified         time.Time
	Version              int
	Scope                LoadBalancerScope
	ScopeNetworkID       NetworkID
	ScopeInstances       int
	MonthlyAllocationMB  int
	MonthlyUsageMB       int
	TrafficRemainingMB   int
}

type LoadBalancerFilter struct {
	Id         LoadBalancerID `json:"id,omitempty"`
	CustomerID CustomerID     `json:"customerid,omitempty"`
	Name       string         `json:"name,omitempty"`
	StdName    string         `json:"stdname,omitempty"`
	Hostname   string         `json:"hostname,omitempty"`
	// Status     string         `json:"status,omitempty"` // TODO LoadBalancerStatus enum DELETED,INACTIVE
	// Type       []string       `json:"type,omitempty"` // TODO LoadBalancerType enum GLB,TCP,HTTP,UDP,CDN?
	// LastModified time.Time      `json:"lastmodified,omitempty"` // TODO serialize as unixepoch
	ResultStart  int `json:"resStart,omitempty"`
	ResultWindow int `json:"resWindow,omitempty"`
}

type loadBalancerGetAllRequest struct {
	legacyRequest
	LoadBalancerFilter
}

type loadBalancerGetAllResponse struct {
	response
	Count         int                    `json:"count"`
	LoadBalancers []existingLoadBalancer `json:"loadbalancers"`
}

type LoadBalancerClient interface {
	GetAll(ctx context.Context, filter LoadBalancerFilter) ([]LoadBalancer, error)
}

type loadBalancerClient struct {
	c *client
}

var _ LoadBalancerClient = (*loadBalancerClient)(nil)

func (lbc *loadBalancerClient) GetAll(ctx context.Context, filter LoadBalancerFilter) ([]LoadBalancer, error) {
	var c client
	req := loadBalancerGetAllRequest{
		legacyRequest: legacyRequest{
			Command: "loadbalancer.getall",
		},
		LoadBalancerFilter: filter,
	}
	var resp loadBalancerGetAllResponse
	err := c.httpLegacyJson(ctx, &req, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsOK() {
		return nil, newApiError(resp.response, nil)
	}
	loadBalancers := make([]LoadBalancer, len(resp.LoadBalancers))
	for i, lb := range resp.LoadBalancers {
		loadBalancers[i] = LoadBalancer{
			ID:                   lb.ID,
			Name:                 lb.Name,
			CustomerID:           lb.CustomerID,
			Hostname:             lb.Hostname,
			StdName:              lb.StdName,
			HostSource:           lb.HostSource,
			HostSourceForceHost:  lb.HostSourceForceHost,
			BackendHostnameForce: lb.BackendHostnameForce,
			Status:               lb.Status,
			DateCreated:          time.Unix(int64(lb.DateCreated), 0),
			DateModified:         time.Unix(int64(lb.DateModified), 0),
			Version:              lb.Version,
			Scope:                lb.Scope,
			ScopeNetworkID:       lb.ScopeNetworkID,
			ScopeInstances:       lb.ScopeInstances,
			MonthlyAllocationMB:  lb.MonthlyAllocationMB,
			MonthlyUsageMB:       lb.MonthlyUsageMB,
			TrafficRemainingMB:   lb.TrafficRemainingMB,
		}
	}
	return loadBalancers, nil
}
