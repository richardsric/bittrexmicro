package bittrex

import (
	"strings"

	"github.com/richardsric/bittrexmicro/helper"
)

// GetWithdrawalHistory is used to retrieve your withdrawal history
// currency string a string literal for the currency (ie. BTC). If set to "all", will return for all currencies
func (b *Bittrex) GetWithdrawalHistoryFunc(cur string) (res string) {
	ressource := "account/getwithdrawalhistory"
	currency := strings.ToLower(cur)
	if currency != "all" {
		ressource += "?currency=" + currency
	}
	r, err := b.client.do("GET", ressource, "", true)
	res = helper.DepositWithdrawalResponseHandler(r, err, "GetWithdrawalHistory")
	return
}

// GetDepositHistory is used to retrieve your deposit history
// currency string a string literal for the currency (ie. BTC). If set to "all", will return for all currencies
func (b *Bittrex) GetDepositHistoryFunc(cur string) (res string) {
	ressource := "account/getdeposithistory"
	currency := strings.ToLower(cur)
	if currency != "all" {
		ressource += "?currency=" + currency
	}
	r, err := b.client.do("GET", ressource, "", true)
	res = helper.DepositWithdrawalResponseHandler(r, err, "GetDepositHistory")
	return
}
