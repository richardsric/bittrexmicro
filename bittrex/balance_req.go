package bittrex

import (
	"fmt"
	"net/http"
	//"strconv"
	"encoding/json"

	"github.com/richardsric/bittrexmicro/helper"
)

func GetBalances(w http.ResponseWriter, req *http.Request) {
	//fmt.Println("Microservice: entered GetBalances function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")
		switch {
		case apiKey == "":
			d := helper.BalancesResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.BalancesResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		//	fmt.Printf("apiKey=%v\napiSecret=%v\n", apiKey, apiSecret)
		b := New(apiKey, apiSecret)
		jsonResp := b.GetBalances()
		fmt.Fprintln(w, jsonResp)
		//	fmt.Println(jsonResp)
	} else {
		fmt.Println("Oops!!...Request method should be get")
	}
}

func GetNonZeroBalances(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("Microservice: entered GetNonZeroBalances function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")
		switch {
		case apiKey == "":
			d := helper.BalancesResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.BalancesResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		//	fmt.Printf("apiKey=%v\napiSecret=%v\n", apiKey, apiSecret)
		b := New(apiKey, apiSecret)
		jsonResp := b.GetNonZeroBalances()
		fmt.Fprintln(w, jsonResp)
		//	fmt.Println(jsonResp)
	} else {
		fmt.Println("Oops!!...Request method should be get")
	}
}

func GetBalance(w http.ResponseWriter, req *http.Request) {
	//	fmt.Println("Microservice: entered GetBalances function")
	if req.Method == "GET" {
		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		req.Header.Add("Accept", "application/json")

		apiKey := req.FormValue("apiKey")
		apiSecret := req.FormValue("secret")
		currency := req.FormValue("currency")
		switch {
		case apiKey == "":
			d := helper.BalancesResponse{
				Result:  "error",
				Message: "'apiKey' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case apiSecret == "":
			d := helper.BalancesResponse{
				Result:  "error",
				Message: "'secret' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		case currency == "":
			d := helper.BalancesResponse{
				Result:  "error",
				Message: "'currency' parameter not found.",
			}
			bs, _ := json.Marshal(d)
			fmt.Fprintln(w, string(bs))
			return
		}
		//	fmt.Printf("apiKey=%v\napiSecret=%v\n", apiKey, apiSecret)
		b := New(apiKey, apiSecret)
		jsonResp := b.GetBalance(currency)
		fmt.Fprintln(w, jsonResp)
		//fmt.Println(jsonResp)
	} else {
		fmt.Println("Oops!!...Request method should be get")
	}
}
