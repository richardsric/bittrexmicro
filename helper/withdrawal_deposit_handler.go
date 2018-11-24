package helper

import (
	"encoding/json"
)

func DepositWithdrawalResponseHandler(res []byte, err error, reqType string) string {
	var bResponse BittrexJsonResponse
	eInfo := ExchangeInfo{
		ExchangeID:   1,
		ExchangeName: "Bittrex",
	}
	if err != nil {
		oResp := HistoryResponse{
			Result:   "error",
			Message:  err.Error(),
			Exchange: eInfo,
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	if err = json.Unmarshal(res, &bResponse); err != nil {
		oResp := HistoryResponse{
			Result:   "error",
			Message:  err.Error(),
			Exchange: eInfo,
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	if err = BittrexErrHandle(bResponse); err != nil {
		oResp := HistoryResponse{
			Result:   "error",
			Message:  err.Error(),
			Exchange: eInfo,
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
	switch reqType {
	case "GetWithdrawalHistory":
		var wHistory []Withdrawal
		err = json.Unmarshal(bResponse.Result, &wHistory)
		/* cStr := WithdrawalInfo{
			Withdrawals : wHistory,
		} */
		withdrawalResponse := WithdrawalResponse{
			Result:   "success",
			Exchange: eInfo,
			Details:  wHistory,
		}
		bs, _ := json.Marshal(withdrawalResponse)
		return string(bs)
	case "GetDepositHistory":
		var dHistory []Deposit
		err = json.Unmarshal(bResponse.Result, &dHistory)
		/* cStr := DepositInfo{
			Deposits : dHistory,
		} */
		depositResponse := DepositResponse{
			Result:   "success",
			Exchange: eInfo,
			Details:  dHistory,
		}
		bs, _ := json.Marshal(depositResponse)
		return string(bs)
	default:
		oResp := DepositResponse{
			Result:   "error",
			Message:  "Request not yet handled",
			Exchange: eInfo,
		}
		bs, _ := json.Marshal(oResp)
		return string(bs)
	}
}
