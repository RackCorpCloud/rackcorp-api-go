package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

type emptyRequest struct{}

type legacyRequest struct {
	Command string `json:"cmd"`
}

type response struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Debug   json.RawMessage `json:"debug"`
}

func (r *response) IsOK() bool {
	return strings.EqualFold(r.Code, "OK")
}

type client struct {
	baseUrl    string
	apiVersion string
	uuid       string
	secret     string
	hc         *http.Client
	userAgent  string
	debugLog   LogFunc
}

type LogFunc func(message string)

func noopLog(message string) {
	// no-op
}

type Client interface {
	OrderConfirm(ctx context.Context, orderId string) (*ConfirmedOrder, error)
	OrderCreate(ctx context.Context, productCode string, customerId string, productDetails ProductDetails) (*CreatedOrder, error)
	OrderGet(ctx context.Context, orderId string) (*Order, error)

	OrderContractGet(ctx context.Context, contractId string) (*OrderContract, error)

	DeviceGet(ctx context.Context, deviceId int) (*Device, error)
	DeviceUpdateFirewall(ctx context.Context, deviceId int, policies []FirewallPolicy) error

	TransactionCreate(ctx context.Context, transactionType string, objectType string, objectId string, confirm bool) (*Transaction, error)
	TransactionDeviceStartup(ctx context.Context, deviceId string, data TransactionStartupData) (*Transaction, error)
	TransactionGet(ctx context.Context, transactionId string) (*Transaction, error)
	TransactionGetAll(ctx context.Context, filter TransactionFilter) ([]Transaction, int, error)

	// TODO LoadBalancer() LoadBalancerClient
	Device() DeviceClient

	SetDebugLog(logFunc LogFunc)
}

var _ Client = (*client)(nil)

const (
	defaultBaseUrl    = "https://api.rackcorp.net/api/"
	defaultApiVersion = "v2.8"
)

func NewClient(uuid string, secret string) (Client, error) {
	if uuid == "" {
		return nil, errors.New("uuid argument must not be empty")
	}

	if secret == "" {
		return nil, errors.New("secret argument must not be empty")
	}

	return &client{
		baseUrl:    defaultBaseUrl,
		apiVersion: defaultApiVersion,
		uuid:       uuid,
		secret:     secret,
		hc: &http.Client{
			Timeout: 30 * time.Second,
		},
		userAgent: fmt.Sprintf("rackcorpapi/1.0 golang/%s", runtime.Version()),
		debugLog:  noopLog,
	}, nil
}

func NewClientFromEnv() (Client, error) {
	cred := newApiCredentialFromEnv()
	if cred == nil {
		return nil, errors.New("failed to load API credentials from environment")
	}
	return NewClient(cred.UUID, cred.Secret)
}

func (c *client) Device() DeviceClient {
	return &deviceClient{c: c}
}

func (c *client) LoadBalancer() LoadBalancerClient {
	return &loadBalancerClient{c: c}
}

func (c *client) SetDebugLog(logFunc LogFunc) {
	if logFunc == nil {
		c.debugLog = noopLog
	} else {
		c.debugLog = logFunc
	}
}

func (c *client) httpLegacyJson(ctx context.Context, reqObj interface{}, respObj interface{}) error {
	url, err := url.JoinPath(c.baseUrl, "rest", c.apiVersion, "json.php")
	if err != nil {
		return fmt.Errorf("failed to construct url: %w", err)
	}
	return c.httpJsonImpl(ctx, http.MethodPost, url, reqObj, respObj)

}
func (c *client) httpRestJson(ctx context.Context, method string, urlSuffix string, reqObj interface{}, respObj interface{}) error {
	url, err := url.JoinPath(c.baseUrl, c.apiVersion, urlSuffix)
	if err != nil {
		return fmt.Errorf("failed to construct url: %w", err)
	}
	return c.httpJsonImpl(ctx, method, url, reqObj, respObj)
}

func (c *client) httpJsonImpl(ctx context.Context, method string, absoluteUrl string, reqObj interface{}, respObj interface{}) error {

	c.debugLog(fmt.Sprintf("Rackcorp API HTTP request: %s %s", method, absoluteUrl))

	var bodyReader io.Reader = nil
	if reqObj != nil {
		reqBody, err := json.Marshal(reqObj)
		if err != nil {
			return fmt.Errorf("failed to JSON encode request body: %v. %w", reqObj, err)
		}
		bodyReader = bytes.NewBuffer(reqBody)
		c.debugLog(fmt.Sprintf("Rackcorp API HTTP request body: '%s'", string(reqBody)))
	}

	req, err := http.NewRequestWithContext(ctx, method, absoluteUrl, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.SetBasicAuth(c.uuid, c.secret)

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	c.debugLog(fmt.Sprintf("Rackcorp API HTTP response status: %d %s", resp.StatusCode, resp.Status))

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read HTTP response body: %w", err)
	}
	c.debugLog(fmt.Sprintf("Rackcorp API HTTP response body: '%s'", string(respBytes)))

	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		return fmt.Errorf("failed to JSON decode response body: %w", err)
	}

	return nil
}
