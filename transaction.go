package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

type createdTransaction struct {
	TransactionId        int    `json:"rcTransactionId"`
	ConfirmationRequired bool   `json:"confirmationRequired"`
	ConfirmationText     string `json:"confirmationText"`
	ObjectType           string `json:"objType"`
	ObjectId             string `json:"objId"`
	Type                 string `json:"type"`
	Data                 string `json:"data"`
}

type transactionCreateRequest struct {
	ObjectType string `json:"objType"`
	ObjectId   string `json:"objId"`
	Type       string `json:"type"`
	Confirm    bool   `json:"confirmation"`
	Data       string `json:"data"`
}

type transactionCreateResponse struct {
	response
	Transaction *createdTransaction `json:"data"`
}

type existingTransaction struct {
	TransactionId string `json:"rcTransactionId"`
	ObjectType    string `json:"objType"`
	ObjectId      string `json:"objId"`
	Type          string `json:"method"`
	Data          string `json:"data"`
	Status        string `json:"status"`
	StatusInfo    string `json:"statusInfo"`
}

type transactionGetRequest struct {
	legacyRequest
	TransactionId string `json:"rcTransactionId"`
}
type transactionGetResponse struct {
	response
	Transaction *existingTransaction `json:"rcTransaction"` // json:"data" for REST
}

type TransactionFilter struct {
	ObjectType   string   `json:"objType"`
	ObjectId     []string `json:"objId,omitempty"`
	Type         []string `json:"method,omitempty"`
	Status       []string `json:"status,omitempty"`
	CustomerId   []string `json:"customerId,omitempty"`
	ResultStart  int      `json:"resStart,omitempty"`
	ResultWindow int      `json:"resWindow,omitempty"`
}

type transactionGetAllRequest struct {
	legacyRequest
	TransactionFilter
}

type transactionGetAllResponse struct {
	response
	Matches      int                   `json:"matches"`
	Transactions []existingTransaction `json:"rcTransactions"`
}

type Transaction struct {
	TransactionId        string
	ObjectType           string
	ObjectId             string
	Type                 string
	Data                 string
	ConfirmationRequired bool
	ConfirmationText     string
	Status               string
	StatusInfo           string
}

type TransactionStartupCloudInit struct {
	MetaData      string `json:"metaData,omitempty"`
	NetworkConfig string `json:"networkConfig,omitempty"`
	UserData      string `json:"userData,omitempty"`
}

type TransactionStartupData struct {
	DeployMediaImageAccessKey    string `json:"deployMediaImageAccessKey,omitempty"`
	DeployMediaImageAccessSecret string `json:"deployMediaImageAccessSecret,omitempty"`
	DeployMediaImageBucket       string `json:"deployMediaImageBucket,omitempty"`
	DeployMediaImageId           string `json:"deployMediaImageId,omitempty"`
	DeployMediaImagePath         string `json:"deployMediaImagePath,omitempty"`

	CloudInit TransactionStartupCloudInit `json:"cloudInit"`
}

const (
	TransactionObjectTypeDevice = "DEVICE"

	TransactionStatusCommenced = "COMMENCED"
	TransactionStatusCompleted = "COMPLETED"
	TransactionStatusPending   = "PENDING"

	TransactionTypeCancel        = "CANCEL"
	TransactionTypeCloseVNC      = "CLOSEVNC"
	TransactionTypeForceShutdown = "FORCESHUTDOWN"
	TransactionTypeOpenVNC       = "OPENVNC" // data parameter contains public IP that allows VNC
	TransactionTypeRefreshConfig = "REFRESHCONFIG"
	TransactionTypeSafeShutdown  = "SAFESHUTDOWN"
	TransactionTypeShutdown      = "SHUTDOWN"
	TransactionTypeStartup       = "STARTUP"
)

func (t *createdTransaction) ToTransaction() *Transaction {
	return &Transaction{
		TransactionId:        strconv.Itoa(t.TransactionId),
		ObjectType:           t.ObjectType,
		ObjectId:             t.ObjectId,
		Type:                 t.Type,
		Data:                 t.Data,
		ConfirmationRequired: t.ConfirmationRequired,
		ConfirmationText:     t.ConfirmationText,
	}
}

func (t *existingTransaction) ToTransaction() *Transaction {
	return &Transaction{
		TransactionId: t.TransactionId,
		ObjectType:    t.ObjectType,
		ObjectId:      t.ObjectId,
		Type:          t.Type,
		Data:          t.Data,
		Status:        t.Status,
		StatusInfo:    t.StatusInfo,
	}
}

func (c *client) TransactionDeviceStartup(ctx context.Context, deviceId string, data TransactionStartupData) (*Transaction, error) {
	var encodedData, err = json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON encode transaction startup data: %w", err)
	}

	return c.transactionCreateInternal(
		ctx,
		TransactionTypeStartup,
		TransactionObjectTypeDevice,
		deviceId,
		true,
		string(encodedData),
	)
}

func (c *client) TransactionCreate(ctx context.Context, transactionType string, objectType string, objectId string, confirm bool) (*Transaction, error) {
	return c.transactionCreateInternal(
		ctx,
		transactionType,
		objectType,
		objectId,
		confirm,
		"",
	)
}

func (c *client) transactionCreateInternal(ctx context.Context, transactionType string, objectType string, objectId string, confirm bool, data string) (*Transaction, error) {
	if transactionType == "" {
		return nil, errors.New("transactionType parameter is required")
	}

	if objectType == "" {
		return nil, errors.New("objectType parameter is required")
	}

	if objectId == "" {
		return nil, errors.New("objectId parameter is required")
	}

	req := &transactionCreateRequest{
		Type:       transactionType,
		ObjectType: objectType,
		ObjectId:   objectId,
		Confirm:    confirm,
		Data:       data,
	}

	var resp transactionCreateResponse
	err := c.httpRestJson(ctx, http.MethodPost, "rctransaction", req, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	if resp.Code != "OK" || resp.Transaction == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Transaction.ToTransaction(), nil
}

func (c *client) TransactionGet(ctx context.Context, transactionId string) (*Transaction, error) {
	req := &transactionGetRequest{
		legacyRequest: legacyRequest{
			Command: "rctransaction.get",
		},
		TransactionId: transactionId,
	}

	var resp transactionGetResponse
	err := c.httpLegacyJson(ctx, req, &resp)
	//	err := c.httpRestJson(ctx, http.MethodGet, fmt.Sprintf("rctransactions/%s", transactionId), emptyRequest{}, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction id '%s': %w", transactionId, err)
	}

	if resp.Code != "OK" || resp.Transaction == nil {
		return nil, newApiError(resp.response, nil)
	}

	return resp.Transaction.ToTransaction(), nil
}

func (c *client) TransactionGetAll(ctx context.Context, filter TransactionFilter) ([]Transaction, int, error) {
	if filter.ObjectType == "" {
		return nil, 0, errors.New("field ObjectType of TransactionFilter is required")
	}

	req := &transactionGetAllRequest{
		legacyRequest: legacyRequest{
			Command: "rctransaction.getall",
		},
		TransactionFilter: filter,
	}

	var resp transactionGetAllResponse
	err := c.httpLegacyJson(ctx, req, &resp)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	if resp.Code != "OK" || resp.Transactions == nil {
		return nil, 0, newApiError(resp.response, nil)
	}

	transactions := make([]Transaction, len(resp.Transactions))
	for index, transaction := range resp.Transactions {
		transactions[index] = *transaction.ToTransaction()
	}

	return transactions, resp.Matches, nil
}
