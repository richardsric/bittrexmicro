package public

import (
	"encoding/json"
	"fmt"
	"strings"

	"net/http"
	"strconv"

	"github.com/richardsric/bittrexmicro/helper"
)

// BittrexSinglePair is the function that will return the ask and bid or error.
func BittrexSinglePair(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//fmt.Println("Request Received Here On Bittrex")

	pair := strings.ToUpper(r.FormValue("pair"))
	eID := r.FormValue("eid")
	var result AskBid
	if pair == "" || eID == "" {
		result = AskBid{
			Result:  "error",
			Message: "Empty Field Selected",
		}

		res, _ := json.Marshal(result)
		fmt.Fprint(w, string(res))
		return
	}

	eid, _ := strconv.Atoi(eID)

	if eid == helper.GetExchangeInfo().ID {
		if len(AllMarketData.AllMarkets) == 0 {
			result = AskBid{
				Result:  "error",
				Message: "No Market Data Available. Check Network",
			}
			res, _ := json.Marshal(result)
			fmt.Fprint(w, string(res))
			return
		}
		//just address the map of all market data using the pair
		//func can scan it like AllMarketData.AllMarkets.["BTC-DGB"], and it will return BittrexMarketDataStruct and boolean
		//if key exists.
		mutex.Lock()
		mdata, exists := AllMarketData.AllMarkets[pair]
		mutex.Unlock()
		if exists == false {

			result = AskBid{
				Result:  "error",
				Message: "Invalid Pair For The Selected Market",
			}
			res, _ := json.Marshal(result)
			fmt.Fprint(w, string(res))
			fmt.Println("Could not fetch price data for", pair)
			return
		}
		if exists {
			//BittrexMarketDataStruct is what is returned as mdata
			//get the data and build the result struct
			result = AskBid{
				Result:     "success",
				Market:     pair,
				Last:       mdata.Last,
				Ask:        mdata.Ask,
				Bid:        mdata.Bid,
				High:       mdata.High,
				Low:        mdata.Low,
				Volume:     mdata.Volume,
				BaseVolume: mdata.BaseVolume,
				Change:     mdata.Change24,
				OpenBuys:   mdata.OpenbuyOrders,
				OpenSells:  mdata.OpenSellOrders,
			}
			//	fmt.Printf("Result: %+v", result)
			res, _ := json.Marshal(result)
			fmt.Fprint(w, string(res))
			//	fmt.Printf("JSON Result: %+v", string(res))

		}
	} else {
		//request is not for this microexchange
		result = AskBid{
			Result:  "error",
			Message: "Request is not for this Exchange",
		}

		res, _ := json.Marshal(result)
		fmt.Fprint(w, string(res))
	}
	return
}
