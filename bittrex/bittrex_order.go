package bittrex

import (
	"encoding/json"
	"strconv"
	//"fmt"
	"strings"

	"github.com/richardsric/bittrexmicro/helper"
)

// SellLimit is used to place a limited sell order in a specific market.
func (b *Bittrex) SellLimit(market string, quantity, rate float64) (res string) {

	r, err := b.client.do("GET", "market/selllimit?market="+market+"&quantity="+strconv.FormatFloat(quantity, 'f', 8, 64)+"&rate="+strconv.FormatFloat(rate, 'f', 8, 64), "", true)
	//fmt.Println(r,err)
	res = helper.HandleResponse(r, err, "SellLimit")
	return
}

// BuyLimit is used to place a limited buy order in a specific market.
func (b *Bittrex) BuyLimit(market string, quantity, rate float64) (res string) {
	r, err := b.client.do("GET", "market/buylimit?market="+market+"&quantity="+strconv.FormatFloat(quantity, 'f', 8, 64)+"&rate="+strconv.FormatFloat(rate, 'f', 8, 64), "", true)
	res = helper.HandleResponse(r, err, "BuyLimit")
	return
}

// CancelOrder is used to cancel a buy or sell order.
func (b *Bittrex) CancelOrder(orderID string) (res string) {
	if orderID == "" {
		re := helper.OrderResponse{
			Result:  "error",
			Message: "Order_number not found!",
		}
		bs, _ := json.Marshal(re)
		return string(bs)
	}
	r, err := b.client.do("GET", "market/cancel?uuid="+orderID, "", true)
	res = helper.HandleResponse(r, err, "CancelOrder")
	return
}

// GetOpenOrders returns orders that you currently have opened.
// If market is set to "all", GetOpenOrders return all orders
// If market is set to a specific order, GetOpenOrders return orders for this market
func (b *Bittrex) GetOpenOrders(market string) (res string) {
	ressource := "market/getopenorders"
	if market != "all" {
		ressource += "?market=" + strings.ToUpper(market)
	}
	r, err := b.client.do("GET", ressource, "", true)
	res = helper.HandleResponse(r, err, "GetOpenOrders")
	return
}

// GetOrderHistory used to retrieve your order history.
// market string literal for the market (ie. BTC-LTC). If set to "all", will return for all market
func (b *Bittrex) GetOrderHistory(market string) (res string) {
	ressource := "account/getorderhistory"
	if market != "all" {
		ressource += "?market=" + market
	}
	r, err := b.client.do("GET", ressource, "", true)
	res = helper.HandleResponse(r, err, "GetOrderHistory")
	return
}

//GetOrderInfo is used to retrieve a particular order details/information
func (b *Bittrex) GetOrderInfo(order_uuid string) (res string) {
	ressource := "account/getorder?uuid=" + order_uuid

	r, err := b.client.do("GET", ressource, "", true)
	res = helper.HandleResponse(r, err, "GetOrderInfo")
	return
}
