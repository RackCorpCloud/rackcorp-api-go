package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

type OrderContract struct {
	ContractId string `json:"contractId"`
	CustomerId string `json:"customerId"`
	DeviceId   string `json:"deviceID"`
	Status     string `json:"status"` // TODO enum
	Type       string `json:"type"`   // TODO enum
	// TODO contractInfo, created, currency, lastBilled, modified, notes, referenceID, serviceBillId
}

type orderContractGetRequest struct {
	legacyRequest
	ContractId string `json:"contractId"`
}

type orderContractGetResponse struct {
	response
	Contract *OrderContract `json:"contract"`
}

func (c *client) OrderContractGet(ctx context.Context, contractId string) (*OrderContract, error) {
	if contractId == "" {
		return nil, errors.New("contractId parameter is required")
	}

	req := &orderContractGetRequest{
		legacyRequest: legacyRequest{
			Command: "order.contract.get",
		},
		ContractId: contractId,
	}

	var resp orderContractGetResponse
	err := c.httpJson(ctx, http.MethodPost, "json.php", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get contract Id '%s': %w", contractId, err)
	}

	if resp.Code != "OK" || resp.Contract == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Contract, nil
}
