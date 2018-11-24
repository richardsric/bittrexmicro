package main

import (
	"fmt"
	"net/http"

	"github.com/richardsric/bittrexmicro/bittrex"
	"github.com/richardsric/bittrexmicro/bittrex/public"
	"github.com/richardsric/bittrexmicro/helper"
)

func main() {
	//*******HANDLERS FOR DEPOSIT AND WITHDRAWAL HISTORY***********//
	http.HandleFunc("/returnDepositHistory", bittrex.GetDepositHistory)
	http.HandleFunc("/returnWithdrawalHistory", bittrex.GetWithdrawalHistory)

	//*******HANDLERS FOR ORDERS***********//
	http.HandleFunc("/returnOrderBuy", bittrex.BuyLimit)
	http.HandleFunc("/returnOrderSell", bittrex.SellLimit)
	http.HandleFunc("/returnOrderCancel", bittrex.CancelOrder)
	http.HandleFunc("/returnOrderInfo", bittrex.GetOrderInfo)
	http.HandleFunc("/returnOpenOrders", bittrex.GetOpenOrders)
	http.HandleFunc("/returnOrderHistory", bittrex.GetOrderHistory)

	//*******HANDLERS FOR BALANCES***********//
	http.HandleFunc("/returnBalances", bittrex.GetBalances)
	http.HandleFunc("/returnNonZeroBalances", bittrex.GetNonZeroBalances)
	http.HandleFunc("/returnBalance", bittrex.GetBalance)

	//*******HANDLERS FOR MARKET DATA***********//
	http.HandleFunc("/stat", public.Statz)
	http.HandleFunc("/pair/price", public.BittrexSinglePair)
	http.HandleFunc("/ticker", public.BittrexMarketData1)

	http.ListenAndServe(fmt.Sprintf(":%s", helper.Port), nil)
}

func init() {
	helper.GetDefaults()
	helper.GetValidPairs()

	var name = "iTradeCoin Bitt-MicroService"
	var version = "0.001 DEVEL"
	var developer = "iYochu Nig LTD"

	fmt.Println("App Name: ", name)
	fmt.Println("App Version: ", version)
	fmt.Println("Developer Name: ", developer)
	fmt.Println("Service Port: ", helper.Port)

	fmt.Printf("iTradeCoin Bittrex Microservice: Running on url 'http://localhost:%s/'\n", helper.Port)
	go public.MarketDataService()
	go public.CoindeskDataService()
	//go public.CoinMarketCapDataService()
	go public.DBMarketsInsertService()
	fmt.Println("Instruction to start All Pseudo Services Completed")

}
