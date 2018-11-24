package public

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/richardsric/bittrexmicro/helper"
)

// Stat is use to track the number excution time of the programe
var stat = make([]time.Duration, 0)

//var mutexcoinmarketcap = &sync.Mutex{}

//CoinMarketCapTicker this hold the bittrex market data and can be accessed to get the data
var CoinMarketCapTicker []byte

//CoincapUnreachableServiceCount holds error count for service calls
var CoincapUnreachableServiceCount int64

//CoinMarketCapDataService this is timer for bittrex market data.
func CoinMarketCapDataService() {
	fmt.Println("Coin Market Cap Data Service Started. Refreshes every", time.Minute*10)

	for {

		CoinMarketCapService()
		time.Sleep(time.Minute * 10)

	}
}

//map[percent_change_1h:0.78 symbol:BTC price_btc:1.0 name:Bitcoin total_supply:16700825.0 max_supply:21000000.0 percent_change_7d:11.89 last_updated:1511627953 rank:1 price_usd:8662.67 24h_volume_usd:4427800000.0 market_cap_usd:144673735703 available_supply:16700825.0 percent_change_24h:5.45 id:bitcoin]

//AllCoinMarketsCap contains all the market data from coinmarketcap
var AllCoinMarketsCap map[string]CoinMarketCapDataStruct

func doCoinMarketCapDataInsert(idv string, namev string, symbolv string, rankv int64, priceusd float64, pricebtc float64, h24volumeusd float64, marketcapusd float64, availablesupply float64, totalsupply float64, percentchange1h float64, percentchange24h float64, percentchange7d float64) {
	//fmt.Println("inserting to db")
	con, err := helper.OpenConnection()
	if err != nil {
		fmt.Println("doCoinMarketCapDataInsert: connection failed:", err)
	}
	defer con.Close()
	var idkey int64
	//INSERT INTO public.coinmarketcap_data(
	//	idkey, id, name, symbol, rank, price_usd, price_btc, h24_volume_usd, market_cap_usd, available_supply, total_supply, percent_change_1h, percent_change_24h, percent_change_7d, last_updated)
	//	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	insertString := `INSERT INTO public.coinmarketcap_data
	(id, name, symbol, rank, price_usd, price_btc, h24_volume_usd, market_cap_usd, available_supply, total_supply, percent_change_1h, percent_change_24h, percent_change_7d)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT ON CONSTRAINT duplicate_symbol
		DO
		UPDATE
		SET (id, name, symbol, rank, price_usd, price_btc, h24_volume_usd, market_cap_usd, available_supply, total_supply, percent_change_1h, percent_change_24h, percent_change_7d, last_updated) 
		= (EXCLUDED.id, EXCLUDED.name, EXCLUDED.symbol, EXCLUDED.rank, EXCLUDED.price_usd, EXCLUDED.price_btc, EXCLUDED.h24_volume_usd, EXCLUDED.market_cap_usd, EXCLUDED.available_supply, EXCLUDED.total_supply, EXCLUDED.percent_change_1h, EXCLUDED.percent_change_24h, EXCLUDED.percent_change_7d, now())
		 RETURNING idkey`
	err = con.Db.QueryRow(insertString, idv, namev, symbolv, rankv, priceusd, pricebtc, h24volumeusd, marketcapusd, availablesupply, totalsupply, percentchange1h, percentchange24h, percentchange7d).Scan(&idkey)

	if err != nil {
		fmt.Println("Could not insert coin market cap data:", err)
	}
}

//CoinMarketCapService is use to return market data for Bittrex Exchange and will be used as bittrex.BittrexMarketData by other packages;.
func CoinMarketCapService() {

	var err error
	url := "https://api.coinmarketcap.com/v1/ticker/?limit=0"

	CoinMarketCapTicker, err = GetTicker(url)
	//fmt.Println(string(body))
	if err != nil {
		fmt.Println("Coin Market Cap Connection Failed. Check Network Connection. Error:", err)
		return
	}

	if len(CoinMarketCapTicker) == 0 {
		fmt.Println("CoinMarketCapService. Nil Response From Coin Market Cap", url)
		fmt.Println("Kindly Check Your Internet Connection")
		return
	}
	var inchar uint8 = 60

	if CoinMarketCapTicker[0] == inchar {

		CoincapUnreachableServiceCount++
		if CoincapUnreachableServiceCount > 99 {
			fmt.Println("CoinCap. '<' character detected many times. Seems service is not available.")
			//Send Telegram Message
			//default ChatID "430073910"

			msg := "<b>CoinCap Service Down!</b>\nCoincap Service has not been reachable for a while now!"
			SendServiceStatusIM(msg)
			//Reset Counter
			CoincapUnreachableServiceCount = 0
		}

		return
	}
	var m interface{}
	err = json.Unmarshal(CoinMarketCapTicker, &m)
	if err != nil {

		fmt.Println("CoinMarketCapService. json Unmarshal failed. ", err)
		return
	}

	if m == nil {
		fmt.Println("CoinMarketCapService. Invalid Response Received From The Request")
		return
	}

	t := m.([]interface{})
	if len(t) > 0 {

		for _, val := range t {

			MapCoinMarketCapData(val)

		}

	}

}

//MapCoinMarketCapData gets the market data from a range value of exchange market data
func MapCoinMarketCapData(val2 interface{}) {
	//fmt.Printf("Type: %T\n Value: %+v \n\n\n", val2, val2)
	//(id, name, symbol, rank, price_usd, price_btc, h24_volume_usd, market_cap_usd, available_supply, total_supply, percent_change_1h, percent_change_24h, percent_change_7d)
	var percentchange7d, percentchange24, percentchange1h, totalsupply, availablesupply, marketcapusd, volumeUsd24H, pricebtc, priceusd float64
	var rank int64
	var err error
	id := val2.(map[string]interface{})["id"].(string)
	symbol := val2.(map[string]interface{})["symbol"].(string)
	//	fmt.Println("Symbol:", val2.(map[string]interface{})["symbol"].(string))
	name := val2.(map[string]interface{})["name"].(string)
	if helper.IsNil(val2.(map[string]interface{})["rank"]) {
		rank = 0
	} else {
		rank, err = strconv.ParseInt(val2.(map[string]interface{})["rank"].(string), 10, 64)
		if err != nil {
			rank = 0
		}
	}

	if helper.IsNil(val2.(map[string]interface{})["price_usd"]) {
		priceusd = 0
	} else {
		priceusd, err = strconv.ParseFloat(val2.(map[string]interface{})["price_usd"].(string), 64)
		if err != nil {
			priceusd = 0
		}
	}

	if helper.IsNil(val2.(map[string]interface{})["price_btc"]) {
		pricebtc = 0
	} else {
		pricebtc, err = strconv.ParseFloat(val2.(map[string]interface{})["price_btc"].(string), 64)
		if err != nil {
			pricebtc = 0
		}
	}

	if val2.(map[string]interface{})["24h_volume_usd"] != nil {
		volumeUsd24H, _ = strconv.ParseFloat(val2.(map[string]interface{})["24h_volume_usd"].(string), 64)
	} else {
		volumeUsd24H = 0
	}
	if val2.(map[string]interface{})["market_cap_usd"] != nil {
		marketcapusd, _ = strconv.ParseFloat(val2.(map[string]interface{})["market_cap_usd"].(string), 64)
	} else {
		marketcapusd = 0
	}
	if val2.(map[string]interface{})["available_supply"] != nil {
		availablesupply, _ = strconv.ParseFloat(val2.(map[string]interface{})["available_supply"].(string), 64)
	} else {
		availablesupply = 0
	}
	if val2.(map[string]interface{})["available_supply"] != nil {
		totalsupply, _ = strconv.ParseFloat(val2.(map[string]interface{})["available_supply"].(string), 64)
	} else {
		totalsupply = 0
	}
	if val2.(map[string]interface{})["percent_change_1h"] != nil {
		percentchange1h, _ = strconv.ParseFloat(val2.(map[string]interface{})["percent_change_1h"].(string), 64)
	} else {
		percentchange1h = 0
	}

	if val2.(map[string]interface{})["percent_change_24h"] != nil {
		percentchange24, _ = strconv.ParseFloat(val2.(map[string]interface{})["percent_change_24h"].(string), 64)
	} else {
		percentchange24 = 0
	}
	if val2.(map[string]interface{})["percent_change_7d"] != nil {
		percentchange7d, _ = strconv.ParseFloat(val2.(map[string]interface{})["percent_change_7d"].(string), 64)
	} else {
		percentchange7d = 0
	}

	doCoinMarketCapDataInsert(id, name, symbol, rank, priceusd, pricebtc, volumeUsd24H, marketcapusd, availablesupply, totalsupply, percentchange1h, percentchange24, percentchange7d)
	//throttle the insertion to half a second per item
	time.Sleep(time.Millisecond * 400)
}
