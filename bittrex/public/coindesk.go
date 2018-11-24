package public

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/richardsric/bittrexmicro/helper"
)

//CoinDeskResponse stores coindesk curency conversion
var CoinDeskResponse []byte

//CoindeskUnreachableServiceCount holds error count for service calls
var CoindeskUnreachableServiceCount int64

func doCoindeskDataInsert(currencyv string, btcusdv float64, btclocalv float64) {
	//fmt.Println("inserting to db")
	con, err := helper.OpenConnection()
	if err != nil {
		fmt.Println("doCoindeskDataInsert. connection failed:", err)
	}
	defer con.Close()
	var idconvv int64
	//INSERT INTO public.coindesk_conv(
	//	idconv, currency, btc_usd, btc_local, last_update)
	//	VALUES (?, ?, ?, ?, ?);
	insertString := `INSERT INTO public.coindesk_conv(
		currency, btc_usd, btc_local)
		VALUES ($1, $2, $3) ON CONFLICT ON CONSTRAINT unique_currency
		DO
		UPDATE
		SET(btc_usd, btc_local, last_update) = (EXCLUDED.btc_usd, EXCLUDED.btc_local, now())
		RETURNING idconv`
	err = con.Db.QueryRow(insertString, currencyv, btcusdv, btclocalv).Scan(&idconvv)

	if err != nil {
		fmt.Println("doCoindeskDataInsert. Could not insert Coindesk:", err)
	}
}

//CoindeskDataService is use to return BTC price in USD and local currencies.
func CoindeskDataService() {
	fmt.Println("Starting Coindesk Local Currency Conversion Service")
	msg := "<b>Service Status Alert!</b>\nCoindesk Data Service has just started!"
	SendServiceStatusIM(msg)
	var cur string
	var err error
	//	var ec int64
	con, err := helper.OpenConnection()
	if err != nil {
		fmt.Println("coindesk.go CoindeskDataService ERROR in opening connection due to: ", err)
		//	ec = 1
	}
	defer con.Close()
	for {

		rows, err := con.Db.Query("select currency from local_currencies")
		if err != nil {
			fmt.Println("coindesk.go:CoindeskDataService Selection of country currencies Failed Due To: ", err)

			//ec = 1
			continue
		}
		//defer rows.Close()

		for rows.Next() {

			//valid pair is d pair in our format, and Market is the exchange format
			err = rows.Scan(&cur)
			if err != nil {
				fmt.Println("coindesk.go:CoindeskDataService Selection of country currencies Failed Due To: ", err)

				//	ec = 1
				continue
			}
			//start db scan
			if cur == "" {
				//sets NGN as default currency
				cur = "NGN"
			}
			//if ec not 1

			url := fmt.Sprintf("https://api.coindesk.com/v1/bpi/currentprice/%s.json", cur)

			CoinDeskResponse, err = GetTicker(url)

			if err != nil {
				fmt.Println("Coindesk. CoindeskDataService Connection Failed. Check Network Connection. Error:", err)

				//		ec = 1
				continue
			}

			if len(CoinDeskResponse) == 0 {
				fmt.Println("CoindeskDataService. Nil Response Gotten From Coindesk", url)
				fmt.Println("Kindly Check Your Internet Connection")
				//	ec = 1
				continue
			}
			var inchar uint8 = 60

			if CoinDeskResponse[0] == inchar {

				CoindeskUnreachableServiceCount++
				if CoindeskUnreachableServiceCount > 99 {
					fmt.Println("CoinDesk. '<' character detected many times. Seems service is not available.")
					//Send Telegram Message
					//default ChatID "430073910"

					msg := "<b>Coindesk Service Down!</b>\nCoindesk Service has not been reachable for a while now!"
					SendServiceStatusIM(msg)
					//Reset Counter
					CoindeskUnreachableServiceCount = 0
				}

				continue
			}
			var m interface{}
			err = json.Unmarshal(CoinDeskResponse, &m)
			if err != nil {
				fmt.Println("CoindeskDataService: jsonUmarshal of m failed: ", err)
				//	ec = 1
				continue
			}

			if m == nil {
				fmt.Println("CoindeskDataService: Invalid Response Received From Coindesk")
				//	ec = 1
				continue
			}

			usdval := m.(map[string]interface{})["bpi"].(map[string]interface{})["USD"].(map[string]interface{})["rate_float"].(float64)
			localval := m.(map[string]interface{})["bpi"].(map[string]interface{})[cur].(map[string]interface{})["rate_float"].(float64)

			doCoindeskDataInsert(cur, usdval, localval)
			//end db scan
			//ec not 1

		}
		//close d rows and wait for some seconds
		rows.Close()
		time.Sleep(time.Second * 48)
	}

}
