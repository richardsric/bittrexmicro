package bittrex

import (
	"encoding/json"
	"net/http"
)

//JSON Response.

type BittrexJsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

///Get Orders..

type OrderHistoryResponse struct {
	Result  string  `json:"result"`
	Message string  `json:"message"`
	Orders  []Order `json:"orders"`
}

type GetOrderResponse struct {
	Result  string `json:"result"`
	Message string `json:"message"`
	Details Order2 `json:"order_details"`
}

// Order is the struct for getOpenOrders and getOrderHistory endpoints
type Order struct {
	OrderUuid         string  `json:"OrderUuid"`
	Exchange          string  `json:"Exchange"`
	TimeStamp         string  `json:"TimeStamp"`
	OrderType         string  `json:"OrderType"`
	Limit             float64 `json:"Limit"`
	Quantity          float64 `json:"Quantity"`
	QuantityRemaining float64 `json:"QuantityRemaining"`
	Commission        float64 `json:"Commission"`
	Price             float64 `json:"Price"`
	PricePerUnit      float64 `json:"PricePerUnit"`
}

// For getorder
type Order2 struct {
	AccountId                  string
	OrderUuid                  string `json:"OrderUuid"`
	Exchange                   string `json:"Exchange"`
	Type                       string
	Quantity                   float64 `json:"Quantity"`
	QuantityRemaining          float64 `json:"QuantityRemaining"`
	Limit                      float64 `json:"Limit"`
	Reserved                   float64
	ReserveRemaining           float64
	CommissionReserved         float64
	CommissionReserveRemaining float64
	CommissionPaid             float64
	Price                      float64 `json:"Price"`
	PricePerUnit               float64 `json:"PricePerUnit"`
	Opened                     string
	Closed                     string
	IsOpen                     bool
	Sentinel                   string
	CancelInitiated            bool
	ImmediateOrCancel          bool
	IsConditional              bool
	Condition                  string
	ConditionTarget            string
}

// Bittrex represent a bittrex client
type Bittrex struct {
	client *Client
}

type Client struct {
	ApiKey     string
	ApiSecret  string
	HttpClient *http.Client
}

type ExchangeInfo struct {
	ExchangeID   int64
	ExchangeName string
}
