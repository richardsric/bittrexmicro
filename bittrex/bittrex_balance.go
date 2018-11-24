package bittrex

import (
	//"encoding/json"
	"strings"

	"github.com/richardsric/bittrexmicro/helper"
)

//GetBalances is used to retrieve all balances from your account
func (b *Bittrex) GetBalances() (res string) {
	r, err := b.client.do("GET", "account/getbalances", "", true)
	res = helper.BalanceResponseHandler(r, err, "GetBalances")
	return
}

//GetNonZeroBalances is used to retrieve all balances that are not zero from your account
func (b *Bittrex) GetNonZeroBalances() (res string) {
	r, err := b.client.do("GET", "account/getbalances", "", true)
	res = helper.BalanceResponseHandler(r, err, "GetNonZeroBalances")
	return
}

// Getbalance is used to retrieve the balance from your account for a specific currency.
// currency: a string literal for the currency (ex: LTC)
func (b *Bittrex) GetBalance(currency string) (res string) { //balance helper.Balance, err error) {
	r, err := b.client.do("GET", "account/getbalance?currency="+strings.ToUpper(currency), "", true)
	res = helper.BalanceResponseHandler(r, err, "GetBalance")
	return
}
