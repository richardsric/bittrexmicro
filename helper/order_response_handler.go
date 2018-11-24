package helper

import (
	"encoding/json"
	"errors"
	"fmt"
)

// BittrexErrHandle gets JSON response from Bittrex API and deal with error
func BittrexErrHandle(r BittrexJsonResponse) error {
	if !r.Success {
		return errors.New(r.Message)
	}
	return nil
}

//HandleResponse Handles Bittrex responses
func HandleResponse(res []byte, err error, reqType string) string {
	var bResponse BittrexJsonResponse
	eInfo := ExchangeInfo{
		ExchangeID:   1,
		ExchangeName: "Bittrex",
	}
	if err != nil {
		oResp := OrderResponse{
			Result:  "error",
			Message: err.Error(), //"Request to bittrex not successfull due to no connection available",//err.Error(),
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	if err = json.Unmarshal(res, &bResponse); err != nil {
		oResp := OrderResponse{
			Result:  "error",
			Message: err.Error(),
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	if err = BittrexErrHandle(bResponse); err != nil {
		oResp := OrderResponse{
			Result:  "error",
			Message: err.Error(),
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	switch reqType {
	case "SellLimit":
		var u BittrexOrderUuid
		err = json.Unmarshal(bResponse.Result, &u)
		oResp := OrderResponse{
			Result:      "success",
			OrderNumber: u.OrderNumber,
		}
		//	fmt.Println("order number for sell limit is ", u.OrderNumber)
		bs, _ := json.Marshal(oResp)
		return string(bs)
	case "BuyLimit":
		var u BittrexOrderUuid
		err = json.Unmarshal(bResponse.Result, &u)
		oResp := OrderResponse{
			Result:      "success",
			OrderNumber: u.OrderNumber,
		}
		//	fmt.Println("order number for buy limit is ", u.OrderNumber)
		bs, _ := json.Marshal(oResp)
		return string(bs)
	case "CancelOrder":
		oResp := OrderResponse{
			Result: "success",
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	case "GetOrderInfo":
		var bOrder BittrexOrderInfo
		err = json.Unmarshal(bResponse.Result, &bOrder)
		var oType string
		var oStatus string
		var oDate string
		fmt.Printf("Bittrex response for getOrder is %+v\n", bOrder)
		if bOrder.Type == "LIMIT_BUY" {
			oType = "BUY"
		} else {
			oType = "SELL"
		}
		if bOrder.IsOpen == false && bOrder.CancelInitiated == false && bOrder.Price > 0 {
			oStatus = "COMPLETED"
			oDate = bOrder.Opened
		} else if bOrder.IsOpen == false && bOrder.CancelInitiated == true && bOrder.Price == 0 {
			oStatus = "CANCELED"
			oDate = bOrder.Opened
		} else if bOrder.IsOpen == true && bOrder.Price == 0 && bOrder.CancelInitiated == false {
			oStatus = "OPEN"
			oDate = bOrder.Opened
		} else if bOrder.IsOpen == true && bOrder.CancelInitiated == false && bOrder.PricePerUnit > 0.00000000 {
			oStatus = "PARTIAL_FILL"
			oDate = bOrder.Opened
		} else if bOrder.IsOpen == false && bOrder.CancelInitiated == true && bOrder.PricePerUnit > 0.00000000 {
			oStatus = "PARTIAL_CANCELED"
			oDate = bOrder.Opened
		} else {
			oStatus = "UNKNOWN"
			oDate = bOrder.Opened
		}
		cStr := CustomOrderInfo{
			Market:            bOrder.Exchange,
			OrderType:         oType,
			ActualQuantity:    bOrder.Quantity,
			ActualRate:        bOrder.PricePerUnit,
			OrderStatus:       oStatus,
			Fee:               bOrder.CommissionPaid,
			QuantityRemaining: bOrder.QuantityRemaining,
			OrderDate:         oDate,
			Price:             bOrder.Price,
			PricePerUnit:      bOrder.PricePerUnit,
			Reserved:          bOrder.ReserveRemaining + bOrder.Price + bOrder.CommissionPaid,
			ReserveRemaining:  bOrder.ReserveRemaining,
			ReserveUsed:       bOrder.Price + bOrder.CommissionPaid,
			Exchange:          eInfo,
		}
		oInfoResponse := OrderInfoResponse{
			Result:       "success",
			OrderNumber:  bOrder.OrderUuid,
			Exchange:     eInfo,
			OrderDetails: cStr,
		}
		//	fmt.Println(oInfoResponse)
		bs, _ := json.Marshal(oInfoResponse)
		return string(bs)
	case "GetOpenOrders":
		var bOrder []BittrexOpenOrderInfo
		err = json.Unmarshal(bResponse.Result, &bOrder)
		//	fmt.Printf("Bittrex response for getOpenOrders is %v\n", bOrder)
		/* cStr := OpenOrdersInfo{
			OpenOrders: bOrder,
		} */
		oInfoResponse := OpenOrderInfoResponse{
			Result:   "success",
			Exchange: eInfo,
			Details:  bOrder,
		}
		bs, _ := json.Marshal(oInfoResponse)
		return string(bs)
	case "GetOrderHistory":
		var bOrders []BittrexOrderInfo
		err = json.Unmarshal(bResponse.Result, &bOrders)
		//	fmt.Printf("Bittrex response for getOpenOrders is %v\n", bOrder)
		oInfoResponse := OrderHistoryResponse{
			Result:       "success",
			Exchange:     eInfo,
			OrderDetails: bOrders,
		}
		bs, _ := json.Marshal(oInfoResponse)
		return string(bs)
	default:
		oResp := OrderResponse{
			Result:  "error",
			Message: "Request not yet handled",
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	oResp := OrderResponse{
		Result:  "error",
		Message: "Request not known",
	}
	bs, _ := json.Marshal(oResp)
	return string(bs)
}
