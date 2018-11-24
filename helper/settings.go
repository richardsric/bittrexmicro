package helper

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

//var pair

var mutex = &sync.Mutex{}

var exchP, gateP, baseMarket map[string]string

//BaseURL variable holds the BaseURL from gateway_settings table
var ApiUrl string

//ApiMethod variable holds the api method from gateway_settings table
var ApiMethod string

//Port variable holds the service_port from gateway_settings table
var Port string

//TimeOut variable holds the request_time_out from gateway_settings table
var TimeOut time.Duration

//GetExchangeInfo returns the hardcoded exchange name and the ID
func GetExchangeInfo() exchangeInfo {
	return exchangeInfo{"Bittrex", 1}
}

// GetDefaults is used to load exapi_settings from db.
func GetDefaults() {
	//	fmt.Println("Enter To Load Settings From DB")
	con, err := OpenConnection()
	if err != nil {
		fmt.Println("SETTINGS.go ERROR in opening connection due to: ", err)
		os.Exit(3)
	}
	defer con.Close()
	err = con.Db.QueryRow("SELECT api_endpoint,api_method,request_timeout,port FROM api_settings WHERE exchange_id = $1", GetExchangeInfo().ID).Scan(&ApiUrl, &ApiMethod, &TimeOut, &Port)
	if err != nil {
		fmt.Println("SETTINGS.go: Selection of api_settings Failed Due To: ", err)
		os.Exit(3)
	}
	LoadBaseMarkets()
}

//GetValidPairs is used to load the currency pairs from the DB
func GetValidPairs() {
	//	fmt.Println("Enter To Load GetValidPairs From DB")
	con, err := OpenConnection()
	if err != nil {
		fmt.Println("SETTINGS.go ERROR in opening connection due to: ", err)
		os.Exit(3)

	}
	defer con.Close()
	rows, err := con.Db.Query("SELECT pair,exchange_format FROM currency_pairs WHERE exchange_id = $1", GetExchangeInfo().ID)
	if err != nil {
		fmt.Println("SETTINGS.go: Selection of valid pairs Failed Due To: ", err)
		os.Exit(3)

	}
	defer rows.Close()
	var p, ex string
	exPairs := make(map[string]string)
	gatePairs := make(map[string]string)
	exchP = make(map[string]string)
	for rows.Next() {

		//valid pair is d pair in our format, and Market is the exchange format
		err = rows.Scan(&p, &ex)
		if err != nil {
			fmt.Println("SETTINGS.go currency pairs rows Scan Failed Due To: ", err)
			os.Exit(3)

		}
		exPairs[p] = ex
		gatePairs[ex] = p

	}
	mutex.Lock()
	//Gateway Pair is Key
	exchP = exPairs
	//Exchange Pair is key
	gateP = gatePairs
	mutex.Unlock()

}

//Mparse2e is used to return either gateway valid pair or exchange format depending on the direction passed
//'f' direction implies that you want to get the exchange format while passing the gateway valid format as pair param
//'r' direction implies that you want to get the gateway valid format while passing the exchange format as pair param
func Mparse2e(exchangeID int, pair string, direction string) exchPair {
	var res exchPair
	if exchangeID != GetExchangeInfo().ID {
		res = exchPair{
			Msg: "Wrror!!...Passed Exchange id different from what was loaded on startup",
		}
	}

	if direction == "f" { //'f' direction implies that you want to get the exchange format while passing the gateway valid format
		//use d mapr where gateway is key so u get d exchange format
		mutex.Lock()
		mapdata, exists := exchP[pair]
		mutex.Unlock()
		if exists == true {

			res = exchPair{
				Pair: mapdata,
			}
		} else {

			res = exchPair{
				Msg: fmt.Sprintf("%s does not exist on %s", pair, GetExchangeInfo().Name),
			}

		}

	} else if direction == "r" { //'r' direction implies that you want to get the gateway valid format while passing the exchange format
		//use d mapr where gateway is key so u get d exchange format
		mutex.Lock()
		mapdata, exists := gateP[pair]
		mutex.Unlock()
		if exists == true {

			res = exchPair{
				Pair: mapdata,
			}
		} else {

			res = exchPair{
				Msg: fmt.Sprintf("%s does not exist on %s", pair, GetExchangeInfo().Name),
			}

		}

	} else {

		res = exchPair{
			Msg: fmt.Sprintf("Invalid direction of %s on %s", pair, GetExchangeInfo().Name),
		}
	}
	return res
}
func insertReturnValidPair(pair string) exchPair {
	pSlice := strings.Split(pair, "-")
	_, ok := baseMarket[pSlice[0]]
	if ok {
		vPair := pSlice[0] + "-" + pSlice[1]
		pairID, err := DBInsertReturn("INSERT INTO currency_pairs(pair,exchange_id,pcurrency,scurrency,exchange_format,status) RETURNING pair_id",
			vPair, GetExchangeInfo().ID, pSlice[0], pSlice[1], pair, 1)
		if err != nil || pairID.(int64) < 0 {
			fmt.Println("pair insert error due to", err)
			return exchPair{
				Msg: fmt.Sprintf("pair insert error due to %s", err),
			}
		}
		if pairID.(int64) > 0 {
			GetValidPairs()
			return exchPair{
				Pair: vPair,
			}
		}
	}
	return exchPair{
		Msg: fmt.Sprintf("pair insertion Failed"),
	}
}

//VetExchMarketData first calls Mpars2e(), knowing that Mparse2e().Pair returns the format if exists or empty if not
func VetExchMarketData(exchangeID int, pair string) exchPair {
	//The function first calls Mpars2e(), knowing that Mparse2e().Pair returns the format if exists or empty if not
	//fmt.Println("reaching here")
	if p := Mparse2e(exchangeID, pair, "r"); len(p.Pair) > 0 {
		//we know it is valid so we return true.
		return p
	} else {
		return insertReturnValidPair(pair)
		/*
			//we know it doesn't exist, so we insert. Insert returns the string pair or empty
			if i := insertReturnValidPair(pair); len(i.Pair) > 0 {
				//means insert was successful. And reload done automatically. We return pair.
				return i
			} else {
				//insert failed for some reason.
				return i
			}
		*/
	}
}

//LoadBaseMarkets loads the base markets configured for the exchange.
func LoadBaseMarkets() {
	con, err := OpenConnection()
	if err != nil {
		fmt.Println("SETTINGS.go ERROR in opening connection due to: ", err)
		os.Exit(3)
	}
	defer con.Close()
	rows, err := con.Db.Query("SELECT base_currency FROM exchange_base_markets WHERE exchange_id = $1", GetExchangeInfo().ID)
	if err != nil {
		fmt.Println("SETTINGS.go: Selection of valid exchange_base_markets Failed Due To: ", err)
		os.Exit(3)
	}
	defer rows.Close()
	var cur string
	baseMarket = make(map[string]string)
	for rows.Next() {
		//valid pair is d pair in our format, and Market is the exchange format
		err = rows.Scan(&cur)
		if err != nil {
			fmt.Println("SETTINGS.go exchange_base_markets rows Scan Failed Due To: ", err)
			os.Exit(3)
		}
		baseMarket[cur] = cur
	}
}
