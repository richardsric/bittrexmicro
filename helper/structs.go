package helper

import "encoding/json"

type ExchangeInfo struct {
	ExchangeID   int64
	ExchangeName string
}

type exchangePairs struct {
	//ValidPair is Gateway Format
	ValidPair string

	//Market is Exchange Format
	Market string
}

type gateWayPairs struct {
	//ValidPair is Gateway Format
	gatePair string

	//ePair is Exchange Format
	ePair string
}

type exchPair struct {
	Pair string
	Msg  string
}

type exchangeInfo struct {
	Name string
	ID   int
}

type ModifyDb struct {
	AffectedRows int64
	ErrorMsg     string
}

type RowSelect struct {
	Columns  map[string]interface{}
	ErrorMsg string
}

type Withdrawal struct {
	PaymentUuid    string  `json:"PaymentUuid"`
	Currency       string  `json:"Currency"`
	Amount         float64 `json:"Amount"`
	Address        string  `json:"Address"`
	Opened         string  `json:"Opened"`
	Authorized     bool    `json:"Authorized"`
	PendingPayment bool    `json:"PendingPayment"`
	TxCost         float64 `json:"TxCost"`
	TxId           string  `json:"TxId"`
	Canceled       bool    `json:"Canceled"`
}

type HistoryResponse struct {
	Result   string       `json:"result"`
	Message  string       `json:"message"`
	Exchange ExchangeInfo `json:"exchange_info"`
}

type Deposit struct {
	Id            int64   `json:"Id"`
	Amount        float64 `json:"Amount"`
	Currency      string  `json:"Currency"`
	Confirmations int     `json:"Confirmations"`
	LastUpdated   string  `json:"LastUpdated"`
	TxId          string  `json:"TxId"`
	CryptoAddress string  `json:"CryptoAddress"`
}

type DepositResponse struct {
	Result   string       `json:"result"`
	Message  string       `json:"message"`
	Exchange ExchangeInfo `json:"exchange_info"`
	Details  []Deposit    `json:"details"`
}

type DepositInfo struct {
	Deposits []Deposit
}

type WithdrawalResponse struct {
	Result   string       `json:"result"`
	Message  string       `json:"message"`
	Exchange ExchangeInfo `json:"exchange_info"`
	Details  []Withdrawal `json:"details"`
}

type WithdrawalInfo struct {
	Withdrawals []Withdrawal
}

type BittrexOpenOrderInfo struct {
	OrderUuid         string  `json:"OrderUuid"`
	Exchange          string  `json:"Exchange"`
	OrderType         string  `json:"OrderType"`
	Limit             float64 `json:"Limit"`
	Reserved          float64 `json:"Reserved"`
	ReserveRemaining  float64 `json:"ReserveRemaining"`
	Quantity          float64 `json:"Quantity"`
	QuantityRemaining float64 `json:"QuantityRemaining"`
	Commission        float64 `json:"CommissionPaid"`
	Price             float64 `json:"Price"`
	PricePerUnit      float64 `json:"PricePerUnit"`
	Opened            string  `json:"Opened"`
	Closed            string  `json:"Closed"`
	CancelInitiated   bool    `json:"CancelInitiated"`
	ImmediateOrCancel bool    `json:"ImmediateOrCancel"`
	IsConditional     bool    `json:"IsConditional"`
	Condition         string  `json:"Condition"`
	ConditionTarget   string  `json:"ConditionTarget"`
}
type BittrexOrderInfo struct {
	AccountId                  string
	OrderUuid                  string  `json:"OrderUuid"`
	Exchange                   string  `json:"Exchange"`
	Type                       string  `json:"OrderType"`
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
	CancelInitiated            bool `json:"CancelInitiated"`
	ImmediateOrCancel          bool `json:"ImmediateOrCancel"`
	IsConditional              bool
	Condition                  string
	ConditionTarget            string
}
type OpenOrdersInfo struct {
	OpenOrders []BittrexOpenOrderInfo
}
type CustomOrderInfo struct {
	Market            string  `json:"market"`
	OrderType         string  `json:"order_type"`
	ActualQuantity    float64 `json:"actual_quantity"`
	QuantityRemaining float64 `json:"QuantityRemaining"`
	ActualRate        float64 `json:"actual_rate"`
	OrderStatus       string  `json:"order_status"`
	Fee               float64 `json:"fee"`
	OrderDate         string  `json:"order_date"`
	Price             float64 `json:"price"`
	PricePerUnit      float64 `json:"pricePerUnit"`
	Reserved          float64 `json:"reserved"`
	ReserveRemaining  float64 `json:"ReserveRemaining"`
	ReserveUsed       float64 `json:"ReserveUsed"`
	Exchange          ExchangeInfo
}
type BittrexJsonResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Result  json.RawMessage `json:"result"`
}

type BittrexOrderUuid struct {
	OrderNumber string `json:"uuid"`
}

type OrderResponse struct {
	Result      string `json:"result"`
	Message     string `json:"message"`
	OrderNumber string `json:"order_number"`
}
type OrderInfoResponse struct {
	Result       string          `json:"result"`
	Message      string          `json:"message"`
	OrderNumber  string          `json:"order_number"`
	Exchange     ExchangeInfo    `json:"exchange_info"`
	OrderDetails CustomOrderInfo `json:"order_details"`
}
type OrderHistoryResponse struct {
	Result       string             `json:"result"`
	Message      string             `json:"message"`
	Exchange     ExchangeInfo       `json:"exchange_info"`
	OrderDetails []BittrexOrderInfo `json:"details"`
}
type OpenOrderInfoResponse struct {
	Result   string                 `json:"result"`
	Message  string                 `json:"message"`
	Exchange ExchangeInfo           `json:"exchange_info"`
	Details  []BittrexOpenOrderInfo `json:"details"`
}

//about balances

type BalancesResponse struct {
	Result  string             `json:"result"`
	Message string             `json:"message"`
	Details map[string]Balance `json:"details"`
}
type BalanceResponse struct {
	Result  string  `json:"result"`
	Message string  `json:"message"`
	Details Balance `json:"details"`
}
type Address struct {
	Currency string `json:"Currency"`
	Address  string `json:"Address"`
}
type AllBalance struct {
	Values []Balance
}
type Balance struct {
	Currency      string  `json:"Currency"`
	Balance       float64 `json:"Balance"`
	Available     float64 `json:"Available"`
	Pending       float64 `json:"Pending"`
	CryptoAddress string  `json:"CryptoAddress"`
}
