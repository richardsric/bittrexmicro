package bittrex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/richardsric/bittrexmicro/helper"
)

func BuyLimit(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("Microservice: entered Buy Limit order function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := strings.Trim(req.FormValue("apiKey"), " ")
		apiSecret := strings.Trim(req.FormValue("secret"), " ")

		market := strings.Trim(req.FormValue("market"), " ")
		quantityStr := strings.Trim(req.FormValue("quantity"), " ")
		rateStr := strings.Trim(req.FormValue("rate"), " ")
		switch {
		case apiKey == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case market == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'market' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case quantityStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'quantity' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case rateStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'rate' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		quantity, _ := strconv.ParseFloat(quantityStr, 64)
		rate, _ := strconv.ParseFloat(rateStr, 64)
		//		fmt.Printf("apiKey=%v\napiSecret=%v\nmarket=%v\nquantity=%v\nrate=%v\n", apiKey, apiSecret, market, quantity, rate)

		b := New(apiKey, apiSecret)
		jsonResp := b.BuyLimit(market, quantity, rate)
		fmt.Fprintln(w, jsonResp)
		//	fmt.Println(jsonResp)
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}

func SellLimit(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("******************************************************************************************************************")
	//	fmt.Println("Microservice: entered Sell Limit order function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")
		apiKey := strings.Trim(req.FormValue("apiKey"), " ")
		apiSecret := strings.Trim(req.FormValue("secret"), " ")

		market := strings.Trim(req.FormValue("market"), " ")
		quantityStr := strings.Trim(req.FormValue("quantity"), " ")
		rateStr := strings.Trim(req.FormValue("rate"), " ")

		switch {
		case apiKey == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case market == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'market' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case quantityStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'quantity' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case rateStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'rate' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		quantity, err := strconv.ParseFloat(quantityStr, 64)
		if err != nil {
			fmt.Println(err)
		}
		rate, err := strconv.ParseFloat(rateStr, 64)
		if err != nil {
			fmt.Println(err)
		}
		//		fmt.Printf("apiKey=%v\napiSecret=%v\nmarket=%v\nquantity=%v\nrate=%v\n", apiKey, apiSecret, market, quantity, rateStr)

		b := New(apiKey, apiSecret)
		jsonResp := b.SellLimit(market, quantity, rate)
		//		fmt.Println(jsonResp)
		fmt.Fprintln(w, jsonResp)
		//		fmt.Println("******************************************************************************************************************")
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}

func CancelOrder(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("******************************************************************************************************************")
	//	fmt.Println("Microservice: entered cancel order function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")
		apiKey := strings.Trim(req.FormValue("apiKey"), " ")
		apiSecret := strings.Trim(req.FormValue("secret"), " ")

		orderId := strings.Trim(req.FormValue("uuid"), " ")
		switch {
		case apiKey == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		//		fmt.Printf("apiKey=%v\napiSecret=%v\norder number=%v\n", apiKey, apiSecret, orderId)
		b := New(apiKey, apiSecret)
		res := b.CancelOrder(orderId)
		//		fmt.Println(res)
		fmt.Fprintln(w, res)
		//		fmt.Println("******************************************************************************************************************")
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}

func GetOrderHistory(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("******************************************************************************************************************")
	//	fmt.Println("entered get orders function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")
		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")
		//accountIdStr := req.FormValue("aid") //required from the worker calling the gateway

		market := req.FormValue("market")
		//eidStr := req.FormValue("eid")
		//eid, _ := strconv.Atoi(eidStr)
		//accountId, _ := strconv.Atoi(accountIdStr)
		//	fmt.Printf("apiKey=%v\napiSecret=%v\n", apiKey, apiSecret)
		switch {
		case apiKey == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		/* case eidStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'eid' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case accountIdStr == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'aid' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return */
		case market == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'market' parameter can not be empty (input either 'all' or a market).",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		//	instanciats the particulat exchange to send the request to
		b := New(apiKey, apiSecret)
		jsonResp := b.GetOrderHistory(market)
		//	fmt.Println(jsonResp)
		fmt.Fprintln(w, jsonResp)

		//	fmt.Println("******************************************************************************************************************")
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}

func GetOpenOrders(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("******************************************************************************************************************")
	//	fmt.Println("entered getOpenOrders function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")
		apiKey := strings.Trim(req.FormValue("apiKey"), " ")
		apiSecret := strings.Trim(req.FormValue("secret"), " ")
		market := strings.Trim(req.FormValue("market"), " ")
		if apiKey == "" {
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		if apiSecret == "" {
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		if market == "" {
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'market' parameter can not be empty (input either 'all' or a market).",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		//	fmt.Printf("apiKey=%v\napiSecret=%v\nmarket=%v\n", apiKey, apiSecret, market)

		//	instanciats the particulat exchange to send the request to
		b := New(apiKey, apiSecret)
		jsonResp := b.GetOpenOrders(market)
		//	fmt.Println(jsonResp)
		fmt.Fprintln(w, jsonResp)

		//	fmt.Println("******************************************************************************************************************")
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}

func GetOrderInfo(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("******************************************************************************************************************")
	//	fmt.Println("Microservice: entered get order info function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")
		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")
		uuid := req.FormValue("uuid")
		switch {
		case apiKey == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.OrderResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		//	fmt.Printf("apiKey=%v\napiSecret=%v\nuuid=%v\n", apiKey, apiSecret, uuid)
		b := New(apiKey, apiSecret)
		jsonResp := b.GetOrderInfo(uuid)
		//	fmt.Println(jsonResp)
		fmt.Fprintln(w, jsonResp)
		//	fmt.Println("******************************************************************************************************************")
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}
