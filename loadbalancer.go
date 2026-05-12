package api

import (
	"context"
	"fmt"
	"time"

	"github.com/rackcorpcloud/rackcorp-api-go/internal"
)

type LoadBalancerID int

func (id *LoadBalancerID) UnmarshalJSON(data []byte) error {
	return internal.UnmarshalJSONInt(id, data)
}

func (id LoadBalancerID) IsZero() bool {
	var zero LoadBalancerID
	return id == zero
}

type LoadBalancerScope string

const (
	LoadBalancerScopeGlobal LoadBalancerScope = "global"
	LoadBalancerScopeLocal  LoadBalancerScope = "local"
)

type LoadBalancerType string

const (
	LoadBalancerTypeCDN  LoadBalancerType = "CDN"
	LoadBalancerTypeGLB  LoadBalancerType = "GLB"
	LoadBalancerTypeHTTP LoadBalancerType = "HTTP"
	LoadBalancerTypeTCP  LoadBalancerType = "TCP"
	LoadBalancerTypeUDP  LoadBalancerType = "UDP"
)

type LoadBalancerBalanceMode string

const (
	LoadBalancerBalanceModeRoundRobin LoadBalancerBalanceMode = "roundrobin"
	LoadBalancerBalanceModeLeast      LoadBalancerBalanceMode = "least"
	LoadBalancerBalanceModeRandom     LoadBalancerBalanceMode = "random"
)

type LoadBalancerCheckMode string

const (
	LoadBalancerCheckModeHTTP LoadBalancerCheckMode = "HTTP"
	LoadBalancerCheckModeTCP  LoadBalancerCheckMode = "TCP"
)

type LoadBalancerTCPProxyMode int

const (
	LoadBalancerTCPProxyModeNone LoadBalancerTCPProxyMode = 0
	LoadBalancerTCPProxyModeV2   LoadBalancerTCPProxyMode = 2
)

type LoadBalancerBackend struct {
	Name     string
	Hostname string
	Port     int
	TLS      bool
	Timeout  time.Duration
	Weight   int
	UUID     string
	TTL      time.Duration
	TCPProxy LoadBalancerTCPProxyMode
	PortMask []int
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
	Type                 LoadBalancerType
	// Status               string // TODO LoadBalancerStatus
	DateCreated         time.Time
	DateModified        time.Time
	Version             int
	Scope               LoadBalancerScope
	ScopeNetworkID      NetworkID
	ScopeInstances      int
	MonthlyAllocationMB int
	MonthlyUsageMB      int
	TrafficRemainingMB  int

	Aliases          []string
	AutoUpgradeHTTPS bool
	Backends         []LoadBalancerBackend
	BalanceMode      LoadBalancerBalanceMode
	CheckMode        LoadBalancerCheckMode
	Ports            []int
	RefreshPatterns  []LoadBalancerRefreshPattern
	Regions          []RegionID
}

type LoadBalancerRefreshPattern struct {
	ID                         LoadBalancerRefreshPatternID
	RegularExpression          string
	MinTTL                     time.Duration
	CacheTime                  time.Duration
	MaxTTL                     time.Duration
	CheckTTL                   time.Duration
	OverrideExpire             bool
	OverrideLastModified       bool
	IgnoreSetCookie            bool
	IgnoreCacheControl         bool
	CacheAuthorizedPages       bool
	BrowserRefresh             LoadBalancerBrowserRefresh
	ForceExpireMins            time.Duration
	PseudoStreamFLV            bool
	PseudoStreamH264           bool
	CGIIgnoreParams            bool
	NoCompression              bool
	RedirectCode               int
	RedirectURL                string
	RedirectPreserveParams     bool
	RedirectForceHTTPS         bool
	IPRestrictionDefaultPolicy string
	IPRestrictions             []LoadBalancerIPRestriction
}

type LoadBalancerACLID int

func (id *LoadBalancerACLID) UnmarshalJSON(data []byte) error {
	return internal.UnmarshalJSONInt(id, data)
}

type LoadBalancerACLAction string

const (
	LoadBalancerACLActionAllow LoadBalancerACLAction = "ALLOW"
	LoadBalancerACLActionDeny  LoadBalancerACLAction = "DENY"
)

type LoadBalancerACL struct {
	ID     LoadBalancerACLID     `json:"id,omitempty"`
	Data   string                `json:"acl_data,omitempty"` // TODO ip mask
	Action LoadBalancerACLAction `json:"acl_action,omitempty"`
}

type LoadBalancerIPRestriction struct {
	IP     string `json:"ip,omitempty"`
	Action string `json:"action,omitempty"`
}

type LoadBalancerRefreshPatternID int

func (id *LoadBalancerRefreshPatternID) UnmarshalJSON(data []byte) error {
	return internal.UnmarshalJSONInt(id, data)
}

type LoadBalancerBrowserRefresh string

const (
	LoadBalancerBrowserRefreshCache           LoadBalancerBrowserRefresh = "CACHE"
	LoadBalancerBrowserRefreshIfModifiedSince LoadBalancerBrowserRefresh = "IFMODIFIEDSINCE"
	LoadBalancerBrowserRefreshRefresh         LoadBalancerBrowserRefresh = "REFRESH"
)

type loadBalancerRefreshPattern struct {
	ID                         LoadBalancerRefreshPatternID `json:"id,omitempty"`
	RegularExpression          string                       `json:"regularexpression,omitempty"`
	MinTTL                     internal.JSONInt             `json:"minttl,omitempty"`    // seconds
	CacheTime                  internal.JSONInt             `json:"cachetime,omitempty"` // seconds
	MaxTTL                     internal.JSONInt             `json:"maxttl,omitempty"`    // seconds
	CheckTTL                   internal.JSONInt             `json:"checkttl,omitempty"`  // seconds
	OverrideExpire             bool                         `json:"overrideexpire,omitempty"`
	OverrideLastModified       bool                         `json:"overridelastmodified,omitempty"`
	IgnoreSetCookie            bool                         `json:"ignoresetcookie,omitempty"`
	IgnoreCacheControl         bool                         `json:"ignorecachecontrol,omitempty"`
	CacheAuthorizedPages       bool                         `json:"cacheauthorizedpages,omitempty"`
	BrowserRefresh             LoadBalancerBrowserRefresh   `json:"browserrefresh,omitempty"`
	ForceExpireMins            internal.JSONInt             `json:"forceexpiremins,omitempty"` // minutes
	PseudoStreamFLV            bool                         `json:"pseudostreamflv,omitempty"`
	PseudoStreamH264           bool                         `json:"pseudostreammp4,omitempty"`
	CGIIgnoreParams            bool                         `json:"cgiignoreparams,omitempty"`
	NoCompression              bool                         `json:"nocompression,omitempty"`
	RedirectCode               internal.JSONInt             `json:"redirectcode,omitempty"`
	RedirectURL                string                       `json:"redirecturl,omitempty"`
	RedirectPreserveParams     bool                         `json:"redirectpreserveparams,omitempty"`
	RedirectForceHTTPS         bool                         `json:"redirectforcehttps,omitempty"`
	IPRestrictionDefaultPolicy string                       `json:"iprestrictiondefaultpolicy,omitempty"`
	IPRestrictions             []LoadBalancerIPRestriction  `json:"iprestrictions,omitempty"`
}

func (p loadBalancerRefreshPattern) ToLoadBalancerRefreshPattern() LoadBalancerRefreshPattern {
	return LoadBalancerRefreshPattern{
		ID:                         p.ID,
		RegularExpression:          p.RegularExpression,
		MinTTL:                     time.Duration(p.MinTTL) * time.Second,
		CacheTime:                  time.Duration(p.CacheTime) * time.Second,
		MaxTTL:                     time.Duration(p.MaxTTL) * time.Second,
		CheckTTL:                   time.Duration(p.CheckTTL) * time.Second,
		OverrideExpire:             p.OverrideExpire,
		OverrideLastModified:       p.OverrideLastModified,
		IgnoreSetCookie:            p.IgnoreSetCookie,
		IgnoreCacheControl:         p.IgnoreCacheControl,
		CacheAuthorizedPages:       p.CacheAuthorizedPages,
		BrowserRefresh:             p.BrowserRefresh,
		ForceExpireMins:            time.Duration(p.ForceExpireMins) * time.Minute,
		PseudoStreamFLV:            p.PseudoStreamFLV,
		PseudoStreamH264:           p.PseudoStreamH264,
		CGIIgnoreParams:            p.CGIIgnoreParams,
		NoCompression:              p.NoCompression,
		RedirectCode:               p.RedirectCode.Int(),
		RedirectURL:                p.RedirectURL,
		RedirectPreserveParams:     p.RedirectPreserveParams,
		RedirectForceHTTPS:         p.RedirectForceHTTPS,
		IPRestrictionDefaultPolicy: p.IPRestrictionDefaultPolicy,
		IPRestrictions:             p.IPRestrictions,
	}
}

type loadBalancerBackend struct {
	Name     string             `json:"name,omitempty"`
	Hostname string             `json:"hostname,omitempty"`
	Port     internal.JSONInt   `json:"port,omitempty"`
	TLS      bool               `json:"tls,omitempty"`
	Timeout  internal.JSONInt   `json:"timeout,omitempty"` // seconds
	Weight   internal.JSONInt   `json:"weight,omitempty"`
	UUID     string             `json:"uuid,omitempty"`
	TTL      internal.JSONInt   `json:"ttl,omitempty"` // seconds
	TCPProxy internal.JSONInt   `json:"tcpproxy,omitempty"`
	PortMask []internal.JSONInt `json:"portmask,omitempty"`
	Created  int64              `json:"created,omitempty"`
	Modified int64              `json:"modified,omitempty"`
}

func (b loadBalancerBackend) ToLoadBalancerBackend() LoadBalancerBackend {
	return LoadBalancerBackend{
		Name:     b.Name,
		Hostname: b.Hostname,
		Port:     b.Port.Int(),
		TLS:      b.TLS,
		Timeout:  time.Duration(b.Timeout.Int()) * time.Second,
		Weight:   b.Weight.Int(),
		UUID:     b.UUID,
		TTL:      time.Duration(b.TTL.Int()) * time.Second,
		TCPProxy: LoadBalancerTCPProxyMode(b.TCPProxy),
		PortMask: internal.JSONIntSliceInt(b.PortMask),
		// TODO Created, Modified
	}
}

func convertLoadBalancerBackend(lbBackend LoadBalancerBackend) loadBalancerBackend {
	return loadBalancerBackend{
		Name:     lbBackend.Name,
		Hostname: lbBackend.Hostname,
		Port:     internal.JSONInt(lbBackend.Port),
		TLS:      lbBackend.TLS,
		Timeout:  internal.JSONInt(int(lbBackend.Timeout.Seconds())),
		Weight:   internal.JSONInt(lbBackend.Weight),
		UUID:     lbBackend.UUID,
		TTL:      internal.JSONInt(int(lbBackend.TTL.Seconds())),
		TCPProxy: internal.JSONInt(int(lbBackend.TCPProxy)),
		PortMask: internal.IntSliceJSONInt(lbBackend.PortMask),
		// TODO Created, Modified
	}
}

type existingLoadBalancer struct {
	ID                   LoadBalancerID   `json:"id"`
	Name                 string           `json:"name"`
	CustomerID           CustomerID       `json:"customerid"`
	Hostname             string           `json:"hostname"`
	StdName              string           `json:"stdname"`
	HostSource           string           `json:"hostsource"`
	Type                 LoadBalancerType `json:"type"`
	HostSourceForceHost  string           `json:"hostsourceforcehost"`
	BackendHostnameForce string           `json:"backend_hostname_force"`
	// Status               string            `json:"status"`
	DateCreated         int64             `json:"datecreated"`
	DateModified        int64             `json:"datemodified"`
	AutoUpgradeHTTPS    bool              `json:"autoupgradehttps"`
	Version             internal.JSONInt  `json:"version"`
	Scope               LoadBalancerScope `json:"scope"`
	ScopeNetworkID      NetworkID         `json:"scope_networkid"`
	ScopeInstances      int               `json:"scope_instances"`
	MonthlyAllocationMB internal.JSONInt  `json:"monthlyallocationmb"`
	MonthlyUsageMB      int               `json:"monthlyusagemb"`
	TrafficRemainingMB  int               `json:"trafficremainingmb"`

	// other fields from Get, but not GetAll:
	ACLs            []LoadBalancerACL            `json:"acl,omitempty"`
	Aliases         []string                     `json:"aliases"`
	AllowDirectSSL  bool                         `json:"allowdirectssl"`
	Backends        []loadBalancerBackend        `json:"backends"`
	BalanceMode     LoadBalancerBalanceMode      `json:"balancemode"`
	CheckMode       LoadBalancerCheckMode        `json:"checkmode"`
	HeadersPassOn   bool                         `json:"headerspasson"`
	Ports           []internal.JSONInt           `json:"ports,omitempty"`
	RefreshPatterns []loadBalancerRefreshPattern `json:"refreshpatterns"`
	Regions         []RegionID                   `json:"regions"`
	TLS             bool                         `json:"tls"`
	// TODO "certificates", "listeners", "protocols", "extra"
}

func (e existingLoadBalancer) ToLoadBalancer() LoadBalancer {
	lb := LoadBalancer{
		ID:                   e.ID,
		Name:                 e.Name,
		CustomerID:           e.CustomerID,
		Hostname:             e.Hostname,
		StdName:              e.StdName,
		HostSource:           e.HostSource,
		Type:                 e.Type,
		HostSourceForceHost:  e.HostSourceForceHost,
		BackendHostnameForce: e.BackendHostnameForce,
		DateCreated:          time.Unix(int64(e.DateCreated), 0),
		DateModified:         time.Unix(int64(e.DateModified), 0),
		Scope:                e.Scope,
		ScopeNetworkID:       e.ScopeNetworkID,
		ScopeInstances:       e.ScopeInstances,
		MonthlyUsageMB:       e.MonthlyUsageMB,
		MonthlyAllocationMB:  e.MonthlyAllocationMB.Int(),
		TrafficRemainingMB:   e.TrafficRemainingMB,
		Version:              e.Version.Int(),

		Aliases:          e.Aliases,
		AutoUpgradeHTTPS: e.AutoUpgradeHTTPS,
		Backends:         make([]LoadBalancerBackend, len(e.Backends)),
		BalanceMode:      e.BalanceMode,
		CheckMode:        e.CheckMode,
		Ports:            internal.JSONIntSliceInt(e.Ports),
		RefreshPatterns:  make([]LoadBalancerRefreshPattern, len(e.RefreshPatterns)),
		Regions:          e.Regions,
	}

	for idx, backend := range e.Backends {
		lb.Backends[idx] = backend.ToLoadBalancerBackend()
	}

	for idx, pattern := range e.RefreshPatterns {
		lb.RefreshPatterns[idx] = pattern.ToLoadBalancerRefreshPattern()
	}

	return lb
}

type loadBalancerGetRequest struct {
	legacyRequest
	Id LoadBalancerID `json:"id"`
}

type loadBalancerGetResponse struct {
	response
	LoadBalancer existingLoadBalancer `json:"loadbalancers"`
}

type loadBalancerDeleteRequest struct {
	legacyRequest
	Id LoadBalancerID `json:"id"`
}

type loadBalancerDeleteResponse = response

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

type loadBalancerCreateRequest struct {
	legacyRequest
	Aliases              []string                     `json:"aliases,omitempty"`
	AutoUpgradeHTTPS     bool                         `json:"autoupgradehttps,omitempty"`
	BackendHostname      string                       `json:"backend_hostname,omitempty"`
	BackendHostnameForce string                       `json:"backend_hostname_force,omitempty"`
	Backends             []loadBalancerBackend        `json:"backends,omitempty"`
	BalanceMode          LoadBalancerBalanceMode      `json:"balancemode,omitempty"`
	CheckMode            LoadBalancerCheckMode        `json:"checkmode,omitempty"`
	CustomerID           CustomerID                   `json:"customerid,omitempty"`
	Name                 string                       `json:"name,omitempty"`
	Ports                []int                        `json:"ports,omitempty"`
	RefreshPatterns      []loadBalancerRefreshPattern `json:"refreshpatterns,omitempty"`
	Regions              []RegionID                   `json:"regions,omitempty"`
	Scope                LoadBalancerScope            `json:"scope,omitempty"`
	ScopeInstances       int                          `json:"scope_instances,omitempty"`
	ScopeNetworkID       NetworkID                    `json:"scope_networkid,omitempty"`
	Type                 LoadBalancerType             `json:"type,omitempty"`
	// BackendSSL           bool                         `json:"backendssl,omitempty"`
	// CheckHost            string                       `json:"checkhost,omitempty"`
	// CheckMethod          string                       `json:"checkmethod,omitempty"`
	// CheckURL             string                       `json:"checkurl,omitempty"`
	// DiscoveryDNS         string                       `json:"discoverydns,omitempty"`
	// DiscoveryMax         json.Number                  `json:"discoverymax,omitempty"`
}

type loadBalancerCreateResponse struct {
	response
	LoadBalancer existingLoadBalancer `json:"loadbalancers"`
}

type loadBalancerUpdateRequest struct {
	legacyRequest
	Backends       []loadBalancerBackend   `json:"backends"`
	BalanceMode    LoadBalancerBalanceMode `json:"balancemode,omitempty"`
	CheckMode      LoadBalancerCheckMode   `json:"checkmode,omitempty"`
	Id             LoadBalancerID          `json:"id,omitempty"`
	Name           string                  `json:"name,omitempty"`
	Ports          []int                   `json:"ports"`
	Scope          LoadBalancerScope       `json:"scope,omitempty"`
	ScopeInstances int                     `json:"scope_instances,omitempty"`
	ScopeNetworkID NetworkID               `json:"scope_networkid,omitempty"`
	Type           LoadBalancerType        `json:"type,omitempty"`
	// TODO all the other fields supported by loadbalancer.update
}

type loadBalancerUpdateResponse struct {
	response
	LoadBalancer existingLoadBalancer `json:"loadbalancers"`
}

type LoadBalancerClient interface {
	Create(ctx context.Context, lb LoadBalancer) (*LoadBalancer, error)
	Delete(ctx context.Context, id LoadBalancerID) error
	Get(ctx context.Context, id LoadBalancerID) (*LoadBalancer, error)
	GetAll(ctx context.Context, filter LoadBalancerFilter) ([]LoadBalancer, error)
	Update(ctx context.Context, lb LoadBalancer) (*LoadBalancer, error)
}

type loadBalancerClient struct {
	c *client
}

var _ LoadBalancerClient = (*loadBalancerClient)(nil)

func (lbc *loadBalancerClient) Get(ctx context.Context, id LoadBalancerID) (*LoadBalancer, error) {
	req := loadBalancerGetRequest{
		legacyRequest: legacyRequest{
			Command: "loadbalancer.get",
		},
		Id: id,
	}
	var resp loadBalancerGetResponse
	err := lbc.c.httpLegacyJson(ctx, &req, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsOK() {
		return nil, newApiError(resp.response, nil)
	}

	loadBalancer := resp.LoadBalancer.ToLoadBalancer()

	return &loadBalancer, nil
}

func (lbc *loadBalancerClient) Delete(ctx context.Context, id LoadBalancerID) error {
	req := loadBalancerDeleteRequest{
		legacyRequest: legacyRequest{
			Command: "loadbalancer.delete",
		},
		Id: id,
	}
	var resp loadBalancerDeleteResponse
	err := lbc.c.httpLegacyJson(ctx, &req, &resp)
	if err != nil {
		return err
	}
	if !resp.IsOK() {
		return newApiError(resp, nil)
	}

	return nil
}

func (lbc *loadBalancerClient) GetAll(ctx context.Context, filter LoadBalancerFilter) ([]LoadBalancer, error) {
	req := loadBalancerGetAllRequest{
		legacyRequest: legacyRequest{
			Command: "loadbalancer.getall",
		},
		LoadBalancerFilter: filter,
	}
	var resp loadBalancerGetAllResponse
	err := lbc.c.httpLegacyJson(ctx, &req, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsOK() {
		return nil, newApiError(resp.response, nil)
	}
	loadBalancers := make([]LoadBalancer, len(resp.LoadBalancers))
	for i, lb := range resp.LoadBalancers {
		loadBalancers[i] = lb.ToLoadBalancer()
	}
	return loadBalancers, nil
}

func (lbc *loadBalancerClient) Create(ctx context.Context, lb LoadBalancer) (*LoadBalancer, error) {
	if !lb.ID.IsZero() {
		return nil, fmt.Errorf("cannot create load balancer with specified ID")
	}

	req := loadBalancerCreateRequest{
		legacyRequest: legacyRequest{
			Command: "loadbalancer.create",
		},
		Aliases:              lb.Aliases,
		AutoUpgradeHTTPS:     lb.AutoUpgradeHTTPS,
		BackendHostnameForce: lb.BackendHostnameForce,
		Backends:             make([]loadBalancerBackend, len(lb.Backends)),
		BalanceMode:          lb.BalanceMode,
		CheckMode:            lb.CheckMode,
		CustomerID:           lb.CustomerID,
		Name:                 lb.Name,
		RefreshPatterns:      make([]loadBalancerRefreshPattern, len(lb.RefreshPatterns)),
		Regions:              lb.Regions,
		Scope:                lb.Scope,
		ScopeInstances:       lb.ScopeInstances,
		ScopeNetworkID:       lb.ScopeNetworkID,
		Type:                 lb.Type,
		Ports:                lb.Ports,
	}

	for idx, backend := range lb.Backends {
		req.Backends[idx] = convertLoadBalancerBackend(backend)
	}

	for idx, pattern := range lb.RefreshPatterns {
		req.RefreshPatterns[idx] = loadBalancerRefreshPattern{
			ID:                         pattern.ID,
			RegularExpression:          pattern.RegularExpression,
			MinTTL:                     internal.JSONInt(pattern.MinTTL.Seconds()),
			CacheTime:                  internal.JSONInt(pattern.CacheTime.Seconds()),
			MaxTTL:                     internal.JSONInt(pattern.MaxTTL.Seconds()),
			CheckTTL:                   internal.JSONInt(pattern.CheckTTL.Seconds()),
			OverrideExpire:             pattern.OverrideExpire,
			OverrideLastModified:       pattern.OverrideLastModified,
			IgnoreSetCookie:            pattern.IgnoreSetCookie,
			IgnoreCacheControl:         pattern.IgnoreCacheControl,
			CacheAuthorizedPages:       pattern.CacheAuthorizedPages,
			BrowserRefresh:             pattern.BrowserRefresh,
			ForceExpireMins:            internal.JSONInt(pattern.ForceExpireMins.Minutes()),
			PseudoStreamFLV:            pattern.PseudoStreamFLV,
			PseudoStreamH264:           pattern.PseudoStreamH264,
			CGIIgnoreParams:            pattern.CGIIgnoreParams,
			NoCompression:              pattern.NoCompression,
			RedirectCode:               internal.JSONInt(pattern.RedirectCode),
			RedirectURL:                pattern.RedirectURL,
			RedirectPreserveParams:     pattern.RedirectPreserveParams,
			RedirectForceHTTPS:         pattern.RedirectForceHTTPS,
			IPRestrictionDefaultPolicy: pattern.IPRestrictionDefaultPolicy,
			IPRestrictions:             pattern.IPRestrictions,
		}
	}

	var resp loadBalancerCreateResponse
	err := lbc.c.httpLegacyJson(ctx, &req, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsOK() {
		return nil, newApiError(resp.response, nil)
	}
	newLb := resp.LoadBalancer.ToLoadBalancer()
	return &newLb, nil
}

func (lbc *loadBalancerClient) Update(ctx context.Context, lb LoadBalancer) (*LoadBalancer, error) {
	if lb.ID.IsZero() {
		return nil, fmt.Errorf("cannot update load balancer without ID")
	}

	req := loadBalancerUpdateRequest{
		legacyRequest: legacyRequest{
			Command: "loadbalancer.update",
		},
		Backends:       make([]loadBalancerBackend, len(lb.Backends)),
		BalanceMode:    lb.BalanceMode,
		CheckMode:      lb.CheckMode,
		Id:             lb.ID,
		Name:           lb.Name,
		Ports:          lb.Ports,
		Scope:          lb.Scope,
		ScopeInstances: lb.ScopeInstances,
		ScopeNetworkID: lb.ScopeNetworkID,
		Type:           lb.Type,
	}

	for idx, backend := range lb.Backends {
		req.Backends[idx] = convertLoadBalancerBackend(backend)
	}

	var resp loadBalancerUpdateResponse
	err := lbc.c.httpLegacyJson(ctx, &req, &resp)
	if err != nil {
		return nil, err
	}
	if !resp.IsOK() {
		return nil, newApiError(resp.response, nil)
	}
	newLb := resp.LoadBalancer.ToLoadBalancer()
	return &newLb, nil

}
