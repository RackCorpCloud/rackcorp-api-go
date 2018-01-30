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
	"time"
)

type emptyRequest struct{}

type legacyRequest struct {
	Command string `json:"cmd"`
}

type response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Debug   string `json:"debug"`
}

type client struct {
	baseUrl string
	uuid    string
	secret  string
	hc      *http.Client
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
}

const (
	defaultBaseUrl = "https://api.rackcorp.net/api/v2.8/"
)

func NewClient(uuid string, secret string) (Client, error) {
	if uuid == "" {
		return nil, errors.New("uuid argument must not be empty")
	}

	if secret == "" {
		return nil, errors.New("secret argument must not be empty")
	}

	return &client{
		baseUrl: defaultBaseUrl,
		uuid:    uuid,
		secret:  secret,
		hc: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c client) httpJson(ctx context.Context, method string, urlSuffix string, reqObj interface{}, respObj interface{}) error {
	var bodyReader io.Reader = nil
	if reqObj != nil {
		reqBody, err := json.Marshal(reqObj)
		if err != nil {
			return fmt.Errorf("failed to JSON encode request body: %v. %w", reqObj, err)
		}
		bodyReader = bytes.NewBuffer(reqBody)
	}

	url, err := url.JoinPath(c.baseUrl, urlSuffix)
	if err != nil {
		return fmt.Errorf("failed to construct url: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}
	if bodyReader != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.uuid, c.secret)

	resp, err := c.hc.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	err = json.NewDecoder(resp.Body).Decode(&respObj)
	if err != nil {
		return fmt.Errorf("failed to JSON decode response body: %w", err)
	}

	return nil
}
