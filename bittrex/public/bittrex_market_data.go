package public

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/richardsric/bittrexmicro/helper"
)

// Stat is use to track the number excution time of the programe
var Stat = make([]time.Duration, 0)
var mutex = &sync.Mutex{}

//UnreachableServiceCount keeps count of unreachable service dat has bn encountered.
var UnreachableServiceCount int64

//UnmarshalErrorCount keeps count of unreachable service dat has bn encountered.
var UnmarshalErrorCount int64

//NilRespErrorCount keeps count of unreachable service dat has bn encountered.
var NilRespErrorCount int64

// RequestNo is the number request that pass through the end point within an hr
var RequestNo = 0

// BittreMarketTicker this hold the bittrex market data and can be accessed to get the data
var BittreMarketTicker []byte

//AllMarketData is scannable by any function in the package so to get any single market data
var AllMarketData AllMarketDataStruct

var marketsData = make(map[string]BittrexMarketDataStruct)
var marketsDBData = make(map[string]BittrexMarketDataStruct)
var chanmarketsDBData = make(chan BittrexMarketDataStruct, 20000)

//SendServiceStatusIM this use to send HTML parsed  message to a telegram user.
func SendServiceStatusIM(msg string) int64 {
	//BotKey for Error Reporting
	var adminTelegram = "430073910"

	con, e := helper.OpenConnection()
	if e != nil {
		return 0
	}
	defer con.Close()
	var mid interface{}
	q := `INSERT INTO telegram_service_status_messages(msgto, message) VALUES($1, $2) RETURNING messageid`
	e = con.Db.QueryRow(q, adminTelegram, msg).Scan(&mid)
	if e != nil {
		return 0
	}
	if mid != nil {
		return mid.(int64)
	}

	return 0
}

// MarketDataService this is timer for bittrex market data.
func MarketDataService() {
	//timeInterval := helpers.GetTimerInterval("MarketDataService")
	fmt.Println("Bittrex Market Data Service Started...")
	msg := "<b>Service Status Alert!</b>\nBittrex Market Data Service has just started!"
	SendServiceStatusIM(msg)
	for {

		BittrexMarketDataService()

		//fmt.Println("Bittrex Market Service Run Successfully.... The Service Will Run Again In The Next 5 Sec")
		//	time.Sleep(time.Millisecond * timeInterval)
		time.Sleep(2000 * time.Millisecond)

	}
}

// BittrexMarketData1 is use to return market data for Bittrex Exchange and will be used as bittrex.BittrexMarketData by other packages;.
func BittrexMarketData1(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Entered BittrexMarketData Func")

	//resp, err := BittrexMarketDataService()
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	fmt.Fprint(w, string(BittreMarketTicker))

}

func doBittrexMarketDataInsert(marketpair string, highv float64, lowv float64, volumev float64, lastv float64, basevolumev float64, mktimestampv time.Time, bidv float64, askv float64, openbuyordersv float64, opensellordersv float64, prevdayv float64, createdv time.Time, change24 float64, workerChan chan int) {
	//fmt.Println("inserting to db")
	con, err := helper.OpenConnection()
	if err != nil {
		fmt.Println("doBittrexMarketDataInsert: connection failed:", err)
	}
	defer con.Close()
	var idbitmkt int64
	//INSERT INTO public.bittrex_market_data(
	//idbitmkt, market, high, low, volume, last, basevolume, mktimestamp, bid, ask, openbuyorders, opensellorders, prevday, created, exch_24change)
	//VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	insertString := `INSERT INTO bittrex_market_data(
		market, high, low, volume, last, basevolume, mktimestamp, bid, ask, openbuyorders, opensellorders, prevday, created, exch_24change)
		VALUES ($1, $2, $3, $4, $5, $6, now(), $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT ON CONSTRAINT duplicate_pair
		DO
		UPDATE
		SET (high, low, volume, last, basevolume, mktimestamp, bid, ask, openbuyorders, opensellorders, prevday, created, exch_24change) = (EXCLUDED.high, EXCLUDED.low, EXCLUDED.volume, EXCLUDED.last, EXCLUDED.basevolume, EXCLUDED.mktimestamp, EXCLUDED.bid, EXCLUDED.ask, EXCLUDED.openbuyorders, EXCLUDED.opensellorders, EXCLUDED.prevday, EXCLUDED.created, EXCLUDED.exch_24change)
		 RETURNING idbitmkt`
	if askv > 0 && bidv > 0 && prevdayv > 0 && highv > 0 && lowv > 0 {
		err = con.Db.QueryRow(insertString, marketpair, highv, lowv, volumev, lastv, basevolumev, bidv, askv, openbuyordersv, opensellordersv, prevdayv, createdv, change24).Scan(&idbitmkt)

		if err != nil {
			fmt.Println("doBittrexMarketDataInsert: Could not insert bittrex market data:", err)
		}
	}
	//work done...free the worker slot.
	<-workerChan
}

// BittrexMarketDataService is use to return market data for Bittrex Exchange and will be used as bittrex.BittrexMarketData by other packages;.
func BittrexMarketDataService() {

	//async := nasync.New(1000, 1000)
	//	defer async.Close()
	start := time.Now() // get current time
	var err error
	//var wG sync.WaitGroup
	url := fmt.Sprintf("%s/public/getmarketsummaries", helper.ApiUrl)
	//BittreMarketTicker, err = GetTicker("https://bittrex.com/api/v1.1/public/getmarketsummaries")

	BittreMarketTicker, err = GetTicker(url)

	if err != nil {
		UnreachableServiceCount++

		if UnreachableServiceCount > 99 {
			fmt.Println("BittrexMarketDataService. getticker error. Seems service is not available.")
			//Send Telegram Message
			//default ChatID "430073910"

			msg := "<b>Bittrex Service Down!</b>\nBittrex Service has not been reachable for a while now!"
			SendServiceStatusIM(msg)
			//Reset Counter
			UnreachableServiceCount = 0
		}
		//fmt.Println("BittrexMarketDataService: Bittrex Connection Failed. Check Network Connection. Error:", err)
		//	fmt.Println("BittrexMarketDataService. Bittrex Data Unmarshal:", err)
		fmt.Println("Raw BittrexMarket Ticker:", BittreMarketTicker)
		fmt.Println("String BittrexMarket Ticker:", string(BittreMarketTicker))
		return
	}

	if len(BittreMarketTicker) == 0 {
		fmt.Println("BittrexMarketDataService: Nil Response Gotten From The Request", url)
		fmt.Println("Kindly Check Your Internet Connection")
		UnreachableServiceCount++
		if UnreachableServiceCount > 99 {
			fmt.Println("BittrexMarketDataService. len of BittreMarketTicker is 0. Seems service is not available.")
			//Send Telegram Message
			//default ChatID "430073910"

			msg := "<b>Bittrex Service Down!</b>\nBittrex Service has not been reachable for a while now!"
			SendServiceStatusIM(msg)
			//Reset Counter
			UnreachableServiceCount = 0
		}
		return
	}
	if BittreMarketTicker == nil {
		fmt.Println("BittrexMarketDataService: BittreMarketTicker is Nil", url)
		UnreachableServiceCount++
		if UnreachableServiceCount > 99 {
			fmt.Println("BittrexMarketDataService. BittreMarketTicker is nil. Seems service is not available.")
			//Send Telegram Message
			//default ChatID "430073910"

			msg := "<b>Bittrex Service Down!</b>\nBittrex Service has not been reachable for a while now!"
			SendServiceStatusIM(msg)
			//Reset Counter
			UnreachableServiceCount = 0
		}
		//fmt.Println("Kindly Check Your Internet Connection")
		return
	}
	var m interface{}
	var inchar uint8 = 60

	if BittreMarketTicker[0] == inchar {

		UnreachableServiceCount++
		if UnreachableServiceCount > 99 {
			fmt.Println("BittrexMarketDataService. '<' character detected many times. Seems service is not available.")
			//Send Telegram Message
			//default ChatID "430073910"

			msg := "<b>Bittrex Service Down!</b>\nBittrex Service has not been reachable for a while now!"
			SendServiceStatusIM(msg)
			//Reset Counter
			UnreachableServiceCount = 0
		}
		fmt.Println("Raw BittrexMarket Ticker:", BittreMarketTicker)
		fmt.Println("String BittrexMarket Ticker:", string(BittreMarketTicker))
		return
	}
	UnreachableServiceCount = 0
	err = json.Unmarshal(BittreMarketTicker, &m)

	if err != nil {
		UnmarshalErrorCount++
		if UnmarshalErrorCount > 20 {
			//Send Telegram Message
			//default ChatID "430073910"

			msg := "<b>Bittrex JSON Unmarshal Error!</b>\nBittrex Service has not been able to unmarshal the JSON it gets for a while now!"
			SendServiceStatusIM(msg)
			//Reset Counter
			UnmarshalErrorCount = 0
		}

		fmt.Println("BittrexMarketDataService. Bittrex Data Unmarshal:", err)
		fmt.Println("Raw BittrexMarket Ticker:", BittreMarketTicker)
		fmt.Println("String BittrexMarket Ticker:", string(BittreMarketTicker))
		return
	}
	UnmarshalErrorCount = 0
	if m == nil {
		NilRespErrorCount++
		if NilRespErrorCount > 2 {
			//Send Telegram Message
			//default ChatID "430073910"

			msg := "<b>Bittrex Nill Response Error!</b>\nBittrex Service has been getting Nil Response for a while now!"
			SendServiceStatusIM(msg)
			//Reset Counter
			NilRespErrorCount = 0
		}
		fmt.Println("BittrexMarketDataService. Invalid Response Received From The Request")
		return
	}
	NilRespErrorCount = 0
	//sem := make(chan int, 50)
	t := m.(map[string]interface{})
	//	startLoop := time.Now()
	//tnil
	if t != nil {
		for key, val := range t {

			//fmt.Println("Got Key1 As:", key, "||", "Got Values1 As:", val)
			if key == "result" {

				for _, val2 := range val.([]interface{}) {
					//	sem <- 1
					//run concurrent
					//	wG.Add(1)
					//go MapMarketData(val2, &wG, sem)
					//go MapMarketData(val2, &wG)
					if val2 != nil {
						MapBittrexMarketData(val2)
						//end run concurrent
					}
				}
				//fmt.Println("Finished Processing Batch...")
				//makrket data is fully mapped...now assign it to d global body to be scanned by singlemarket
				//func can scan it like AllMarketData.AllMarkets.["BTC-DGB"], and it will return BittrexMarketDataStruct and boolean
				//if key exists.
				//	wG.Wait()
				mutex.Lock()
				AllMarketData = AllMarketDataStruct{
					LastChanged: time.Now(),
					AllMarkets:  marketsData,
				}
				//	marketsDBData = marketsData
				mutex.Unlock()
				//end of assignment. functions can now scan it
			}

		}
	} else {
		fmt.Println("BittrexMarketDataService. t is nil Invalid Response Received From The Request")
		return
	}
	//t nil
	//fmt.Println("Time taken to Range Through All Market Data is:", time.Since(startLoop))
	//fmt.Println("...........................................END Of MarketData Range............................................................................")

	RequestNo = RequestNo + 1
	elapsed := time.Since(start)
	Stat = append(Stat, elapsed)
	//fmt.Println("Total Time taken by MarketData function process is:", elapsed)
	//fmt.Println("...........................................END Of MarketData Function............................................................................")

	//return body, nil

}

//MapBittrexMarketData gets the market data from a range value of exchange market data
func MapBittrexMarketData(val2 interface{}) {
	//	fmt.Printf("%+v \n\n\n", val2)
	//time.Sleep(10 * time.Second)
	//	defer w.Done()
	//<-sem
	//	startDataExtract := time.Now()
	//fmt.Println("starting Loop")
	//fmt.Println("Got Key2 As:", key2, "||", "Got Values2 As:", val2)
	//market, high, low, volume, last, basevolume, mktimestamp, bid, ask, openbuyorders, opensellorders, prevday, created, exch_24change
	if val2 == nil {
		return
	}
	var pair string
	if val2.(map[string]interface{})["MarketName"] != nil {
		pair = val2.(map[string]interface{})["MarketName"].(string)
	}
	var ask float64
	if val2.(map[string]interface{})["Ask"] != nil {
		ask = val2.(map[string]interface{})["Ask"].(float64)
	}
	var bid float64
	if val2.(map[string]interface{})["Bid"] != nil {
		bid = val2.(map[string]interface{})["Bid"].(float64)
	}
	var last float64
	if val2.(map[string]interface{})["Last"] != nil {
		last = val2.(map[string]interface{})["Last"].(float64)
	}
	var high24hr float64
	if val2.(map[string]interface{})["High"] != nil {
		high24hr = val2.(map[string]interface{})["High"].(float64)
	}
	var low24hr float64
	if val2.(map[string]interface{})["Low"] != nil {
		low24hr = val2.(map[string]interface{})["Low"].(float64)
	}
	var vol float64
	if val2.(map[string]interface{})["Volume"] != nil {
		vol = val2.(map[string]interface{})["Volume"].(float64)
	}
	var baseVol float64
	if val2.(map[string]interface{})["BaseVolume"] != nil {
		baseVol = val2.(map[string]interface{})["BaseVolume"].(float64)
	}
	var mktimeStamp time.Time
	if val2.(map[string]interface{})["TimeStamp"] != nil {
		mktimeStamp, _ = time.Parse("2006-01-02T15:04:05.99", val2.(map[string]interface{})["TimeStamp"].(string))
	}
	var openbuyOrders float64
	if val2.(map[string]interface{})["OpenBuyOrders"] != nil {
		openbuyOrders = val2.(map[string]interface{})["OpenBuyOrders"].(float64)
	}
	var opensellOrders float64
	if val2.(map[string]interface{})["OpenSellOrders"] != nil {
		opensellOrders = val2.(map[string]interface{})["OpenSellOrders"].(float64)
	}
	var creaTed time.Time
	if val2.(map[string]interface{})["Created"] != nil {
		creaTed, _ = time.Parse("2006-01-02T15:04:05.99", val2.(map[string]interface{})["Created"].(string))
	}
	var prevDay float64
	if val2.(map[string]interface{})["PrevDay"] != nil {
		prevDay = val2.(map[string]interface{})["PrevDay"].(float64)
	}
	if ask == 0 || bid == 0 || prevDay == 0 {
		return
	}
	exchangeID := helper.GetExchangeInfo().ID

	changE24 := ((last - prevDay) / prevDay) * 100
	//		fmt.Println("exchange pair is ", pair)
	vPair := helper.VetExchMarketData(exchangeID, pair)
	//		fmt.Printf("pair-%s, ask-%v,bid-%v,last-%v,high-%v,low-%v,volume-%v,baseVol-%v,exchangeId-%v\n",
	//			vPair.Pair, ask, bid, last, high24hr, low24hr, vol, baseVol, exchangeID)
	//we can generate multirow insert values.
	//(first row set), (second row set)... (nth row set)
	//async.Do(doMarketDatainsert, vPair.Pair, ask, bid, last, high24hr, low24hr, vol, baseVol, exchangeID)
	//doMarketDatainsert(vPair.Pair, ask, bid, last, high24hr, low24hr, vol, baseVol, exchangeID)
	//market, high, low, volume, last, basevolume, mktimestamp, bid, ask, openbuyorders, opensellorders, prevday, created, exch_24change
	result := BittrexMarketDataStruct{
		Ask:            ask,
		Bid:            bid,
		Last:           last,
		High:           high24hr,
		Low:            low24hr,
		Volume:         vol,
		BaseVolume:     baseVol,
		MktimeStamp:    mktimeStamp,
		Created:        creaTed,
		OpenbuyOrders:  openbuyOrders,
		OpenSellOrders: opensellOrders,
		PrevDay:        prevDay,
		Change24:       changE24,
		ExchangeID:     exchangeID,
		Market:         vPair.Pair,
	}
	//Insert The map keys
	mutex.Lock()
	marketsData[vPair.Pair] = result
	mutex.Unlock()

	//concurrently send through the channel
	if ask > 0 && bid > 0 && prevDay > 0 {
		go func() {
			select {
			case chanmarketsDBData <- result:
			default:
			}
		}()
	}
	//	doBittrexMarketDataInsert(vPair.Pair, high24hr, low24hr, vol, last, baseVol, mktimeStamp, bid, ask, openbuyOrders, opensellOrders, prevDay, creaTed, changE24)

	//fmt.Println("Time To Map Market data for ", vPair.Pair, " is: ", time.Since(startDataExtract))

}

// Statz use t0 show statz
func Statz(w http.ResponseWriter, r *http.Request) {
	fmt.Println("...........................................Entered Stat Function............................................................................")
	var n, smallest, biggest time.Duration
	x := Stat

	for _, v := range x {
		if v > n {
			fmt.Println(v, ">", n)
			n = v
			biggest = n
		} else {
			fmt.Println(v, "<", n)
		}
	}

	fmt.Println("The biggest number is ", biggest)
	for _, v := range x {
		if v > n {
			fmt.Println(v, ">", n)
		} else {
			fmt.Println(v, "<", n)
			n = v
			smallest = n
		}
	}
	fmt.Println("The smallest number is ", smallest)

	fmt.Fprint(w, "iTradeCoin Bittrex Market Update Service Running", "\n\n")
	fmt.Fprint(w, "Number Of Request Recevied Within 1Hr: ", RequestNo, "\n")
	fmt.Fprint(w, "Minimum Execution Time Within 1Hr: ", smallest, "\n")
	fmt.Fprint(w, "Maximum Execution Time Within 1Hr: ", biggest, "\n")
}

func isNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}

//DBMarketsInsertService inserts market data on DB
func DBMarketsInsertService() {
	//create workers control chan

	numW := 12
	workerChan := make(chan int, numW)

	for i := 1; i < numW; i++ {
		workerChan <- 1
		val := <-chanmarketsDBData
		fmt.Println("Starting DBDataInsert Worker:", i)

		go doBittrexMarketDataInsert(val.Market, val.High, val.Low, val.Volume, val.Last, val.BaseVolume, val.MktimeStamp, val.Bid, val.Ask, val.OpenbuyOrders, val.OpenSellOrders, val.PrevDay, val.Created, val.Change24, workerChan)

	}
	msg := "<b>Service Status Alert!</b>\nBittrex Market Data DB Insert Service has just started!\n" + fmt.Sprint(numW) + " worker processes provisioned."
	SendServiceStatusIM(msg)
	fmt.Println("Starting Worker Controller... And remaining Workers.", numW, "Workers Total")
	for {
		//provision new worker
		workerChan <- 1
		//get work to worker
		val := <-chanmarketsDBData
		//	fmt.Println("starting worker to process write", val)
		go doBittrexMarketDataInsert(val.Market, val.High, val.Low, val.Volume, val.Last, val.BaseVolume, val.MktimeStamp, val.Bid, val.Ask, val.OpenbuyOrders, val.OpenSellOrders, val.PrevDay, val.Created, val.Change24, workerChan)

	}

}
