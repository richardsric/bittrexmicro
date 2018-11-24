package bittrex

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/richardsric/bittrexmicro/helper"
)

//GetWithdrawalHistory handles withdrawal history request
func GetWithdrawalHistory(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("Microservice: entered GetWithdrawalHistory Handler func")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := strings.Trim(req.FormValue("apiKey"), " ")
		apiSecret := strings.Trim(req.FormValue("secret"), " ")

		currency := strings.Trim(req.FormValue("currency"), " ")
		switch {
		case apiKey == "":
			d := helper.HistoryResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.HistoryResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case currency == "":
			d := helper.HistoryResponse{
				Result:  "error",
				Message: "'currency' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		//fmt.Printf("apiKey=%v\napiSecret=%v\ncurrency=%v\n", apiKey, apiSecret, currency)
		b := New(apiKey, apiSecret)
		jsonResp := b.GetWithdrawalHistoryFunc(currency)
		fmt.Fprintln(w, jsonResp)
		//	fmt.Println(jsonResp)
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}

//GetDepositHistory handles deposit history requests
func GetDepositHistory(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("Microservice: entered GetDepositHistory Handler func")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := strings.Trim(req.FormValue("apiKey"), " ")
		apiSecret := strings.Trim(req.FormValue("secret"), " ")

		currency := strings.Trim(req.FormValue("currency"), " ")
		switch {
		case apiKey == "":
			d := helper.HistoryResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.HistoryResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case currency == "":
			d := helper.HistoryResponse{
				Result:  "error",
				Message: "'currency' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		//fmt.Printf("apiKey=%v\napiSecret=%v\ncurrency=%v\n", apiKey, apiSecret, currency)
		b := New(apiKey, apiSecret)
		jsonResp := b.GetDepositHistoryFunc(currency)
		fmt.Fprintln(w, jsonResp)
		//	fmt.Println(jsonResp)
	} else {
		fmt.Println("The method shouldn't be POST method")
	}
}
