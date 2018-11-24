package public

import "time"

//CoinMarketCapDataStruct is a struct use to return ask and bid of request pair.
type CoinMarketCapDataStruct struct {
	/*"
		"id": "bitcoin",
	        "name": "Bitcoin",
	        "symbol": "BTC",
	        "rank": "1",
	        "price_usd": "573.137",
	        "price_btc": "1.0",
	        "24h_volume_usd": "72855700.0",
	        "market_cap_usd": "9080883500.0",
	        "available_supply": "15844176.0",
	        "total_supply": "15844176.0",
	        "percent_change_1h": "0.04",
	        "percent_change_24h": "-0.3",
	        "percent_change_7d": "-0.57",
	        "last_updated": "1472762067"
	*/
	Symbol           string
	ID               string
	Name             string
	Rank             int64
	PriceUsd         float64
	PriceBtc         float64
	VolumeUsd24H     float64
	MarketCapUsd     float64
	AvailableSupply  float64
	TotalSupply      float64
	PercentChange1H  float64
	PercentChange24H float64
	PercentChange7D  float64
}

//BittrexMarketDataStruct is a struct use to return ask and bid of request pair.
type BittrexMarketDataStruct struct {
	Ask            float64
	Bid            float64
	Last           float64
	High           float64
	Low            float64
	Volume         float64
	BaseVolume     float64
	MktimeStamp    time.Time
	Created        time.Time
	OpenbuyOrders  float64
	OpenSellOrders float64
	PrevDay        float64
	Change24       float64
	ExchangeID     int
	Market         string
}

// BittrexMarketDataService is use to return market data for Bittrex Exchange and will be used as bittrex.BittrexMarketData by other packages;.

//AllMarketDataStruct is Used to hold the vetted market data, scannable by single market data
type AllMarketDataStruct struct {
	LastChanged time.Time
	AllMarkets  map[string]BittrexMarketDataStruct
}

// AskBid is a struct use to return ask and bid of request pair.
type AskBid struct {
	Result     string  `json:"result"`
	Message    string  `json:"message"`
	Market     string  `json:"market"`
	Last       float64 `json:"last"`
	Ask        float64 `json:"ask"`
	Bid        float64 `json:"bid"`
	High       float64 `json:"high"`
	Low        float64 `json:"low"`
	Volume     float64 `json:"volume"`
	BaseVolume float64 `json:"basevolume"`
	Change     float64 `json:"change"`
	OpenBuys   float64 `json:"openbuys"`
	OpenSells  float64 `json:"opensells"`
}

// MainAskBid1 this is use to get single request
type MainAskBid1 struct {
	Values []AskBid
}
