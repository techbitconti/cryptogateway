package dbScan

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"db/redis"
)

var REDIS_KEYS_DEPOSIT = "keys_deposit"

var REDIS_BTC_DEPOSIT = "btc_deposit"
var REDIS_BTC_WITHDRAW = "btc_withdraw"
var REDIS_BTC_CURRENT = "btc_current"
var REDIS_BTC_FEES = "btc_fees"

var REDIS_LTC_DEPOSIT = "ltc_deposit"
var REDIS_LTC_WITHDRAW = "ltc_withdraw"
var REDIS_LTC_CURRENT = "ltc_current"
var REDIS_LTC_FEES = "ltc_fees"

var REDIS_BCH_DEPOSIT = "bch_deposit"
var REDIS_BCH_WITHDRAW = "bch_withdraw"
var REDIS_BCH_CURRENT = "bch_current"
var REDIS_BCH_FEES = "bch_fees"

var REDIS_ETH_DEPOSIT = "eth_deposit"
var REDIS_ETH_WITHDRAW = "eth_withdraw"
var REDIS_ETH_CURRENT = "eth_current"
var REDIS_ETH_FEES = "eth_fees"

var REDIS_XLM_DEPOSIT = "xlm_deposit"
var REDIS_XLM_WITHDRAW = "xlm_withdraw"
var REDIS_XLM_CURRENT = "xlm_current"
var REDIS_XLM_FEES = "xlm_fees"

/*---------------------------------------------------------------------------*/

var BTC_DEPOSIT float64
var BTC_WITHDRAW float64
var BTC_CURRENT float64
var BTC_FEES float64

var LTC_DEPOSIT float64
var LTC_WITHDRAW float64
var LTC_CURRENT float64
var LTC_FEES float64

var BCH_DEPOSIT float64
var BCH_WITHDRAW float64
var BCH_CURRENT float64
var BCH_FEES float64

var ETH_DEPOSIT float64
var ETH_WITHDRAW float64
var ETH_CURRENT float64
var ETH_FEES float64

var XLM_DEPOSIT float64
var XLM_WITHDRAW float64
var XLM_CURRENT float64
var XLM_FEES float64

/*---------------------------------------------------------------------------*/

func Report_Deposit(coin string, de float64) {

	switch coin {
	case "BTC":
		{
			BTC_DEPOSIT += de
			SaveReport_BTC_Deposit(BTC_DEPOSIT)
		}
	case "LTC":
		{
			LTC_DEPOSIT += de
			SaveReport_LTC_Deposit(LTC_DEPOSIT)
		}
	case "BCH":
		{
			BCH_DEPOSIT += de
			SaveReport_BCH_Deposit(BCH_DEPOSIT)
		}
	case "ETH":
		{
			ETH_DEPOSIT += de
			SaveReport_ETH_Deposit(ETH_DEPOSIT)
		}
	case "XLM":
		{
			XLM_DEPOSIT += de
			SaveReport_XLM_Deposit(XLM_DEPOSIT)
		}
	}

}

func Report_Withdraw(coin string, with float64) {

	switch coin {
	case "BTC":
		{
			BTC_WITHDRAW += with
			SaveReport_BTC_Withdraw(BTC_WITHDRAW)
		}
	case "LTC":
		{
			LTC_WITHDRAW += with
			SaveReport_LTC_Withdraw(LTC_WITHDRAW)
		}
	case "BCH":
		{
			BCH_WITHDRAW += with
			SaveReport_BCH_Withdraw(BCH_WITHDRAW)
		}
	case "ETH":
		{
			ETH_WITHDRAW += with
			SaveReport_ETH_Withdraw(ETH_WITHDRAW)
		}
	case "XLM":
		{
			XLM_WITHDRAW += with
			SaveReport_XLM_Withdraw(XLM_WITHDRAW)
		}
	}

}

func Report_Fees(coin string, fees float64) {

	switch coin {
	case "BTC":
		{
			BTC_FEES += math.Abs(fees)
			SaveReport_BTC_Fees(BTC_FEES)
		}
	case "LTC":
		{
			LTC_FEES += math.Abs(fees)
			SaveReport_LTC_Fees(LTC_FEES)
		}
	case "BCH":
		{
			BCH_FEES += math.Abs(fees)
			SaveReport_BCH_Fees(BCH_FEES)
		}
	case "ETH":
		{
			ETH_FEES += fees
			SaveReport_ETH_Fees(ETH_FEES)
		}
	case "XLM":
		{
			XLM_FEES += fees
			SaveReport_XLM_Fees(XLM_FEES)
		}
	}

}

func Report_Current(coin string) {

	switch coin {
	case "BTC":
		{
			BTC_CURRENT = BTC_DEPOSIT - BTC_WITHDRAW
			SaveReport_BTC_Current(BTC_CURRENT)
		}
	case "LTC":
		{
			LTC_CURRENT = LTC_DEPOSIT - LTC_WITHDRAW
			SaveReport_LTC_Current(LTC_CURRENT)
		}
	case "BCH":
		{
			BCH_CURRENT = BCH_DEPOSIT - BCH_WITHDRAW
			SaveReport_BCH_Current(BCH_CURRENT)
		}
	case "ETH":
		{
			ETH_CURRENT = ETH_DEPOSIT - ETH_WITHDRAW
			SaveReport_ETH_Current(ETH_CURRENT)
		}
	case "XLM":
		{
			XLM_CURRENT = XLM_DEPOSIT - XLM_WITHDRAW
			SaveReport_XLM_Current(XLM_CURRENT)
		}
	}

}

/*---------------------------------------------------------------------------*/

func LoadDeposit() {

	bKeys, _ := redis.Session.Get(REDIS_KEYS_DEPOSIT).Bytes()

	var keys []string
	json.Unmarshal(bKeys, &keys)

	for _, key := range keys {

		value, _ := redis.Session.Get(key).Result()

		if value != "" {

			var de Deposit
			json.Unmarshal([]byte(value), &de)

			HMAP_DEPOSIT[key] = &de

			fmt.Println("REDIS_KEYS_DEPOSIT: ", key, value)
		}
	}
}

func SaveDeposit(de Deposit) {

	data, _ := json.Marshal(de)
	redis.Session.Set(de.AddressDeposit, data, 0).Err()

	if _, ok := HMAP_DEPOSIT[de.AddressDeposit]; !ok {

		arr, _ := redis.Session.Get(REDIS_KEYS_DEPOSIT).Result()

		var keys []string
		json.Unmarshal([]byte(arr), &keys)
		keys = append(keys, de.AddressDeposit)

		bKeys, _ := json.Marshal(keys)

		redis.Session.Set(REDIS_KEYS_DEPOSIT, bKeys, 0).Err()
	}

	fmt.Println("saveDeposit : ", de.AddressDeposit)
}

/*---------------------------------------------------------------------------*/

func LoadReport_BTC() {

	de, err1 := redis.Session.Get(REDIS_BTC_DEPOSIT).Result()
	if err1 != nil {
		fmt.Println("Error Load Report BTC Deposit")
		return
	}

	with, err2 := redis.Session.Get(REDIS_BTC_WITHDRAW).Result()
	if err2 != nil {
		fmt.Println("Error Load Report BTC Withdraw")
		return
	}

	curr, err3 := redis.Session.Get(REDIS_BTC_CURRENT).Result()
	if err3 != nil {
		fmt.Println("Error Load Report BTC Current")
		return
	}

	fees, err4 := redis.Session.Get(REDIS_BTC_FEES).Result()
	if err4 != nil {
		fmt.Println("Error Load Report BTC Fees")
		return
	}

	df, _ := strconv.ParseFloat(de, 64)
	BTC_DEPOSIT = df

	wf, _ := strconv.ParseFloat(with, 64)
	BTC_WITHDRAW = wf

	cf, _ := strconv.ParseFloat(curr, 64)
	BTC_CURRENT = cf

	bf, _ := strconv.ParseFloat(fees, 64)
	BTC_FEES = bf

	fmt.Println("Redis BTC Deposit : ", de, "  - Withdraw : ", with, "  - Current : ", curr, "  - Fees : ", fees)
}

func SaveReport_BTC_Deposit(de float64) {

	b, _ := json.Marshal(de)

	err1 := redis.Session.Set(REDIS_BTC_DEPOSIT, b, 0).Err()
	if err1 != nil {
		fmt.Println("Error Save Report BTC Deposit")
	}

	fmt.Println("SaveReport_BTC_Deposit : ", REDIS_BTC_DEPOSIT, de)
}

func SaveReport_BTC_Withdraw(with float64) {

	b, _ := json.Marshal(with)

	err2 := redis.Session.Set(REDIS_BTC_WITHDRAW, b, 0).Err()
	if err2 != nil {
		fmt.Println("Error Save Report BTC Withdraw")
	}

	fmt.Println("SaveReport_BTC_Withdraw : ", REDIS_BTC_WITHDRAW, with)
}

func SaveReport_BTC_Current(curr float64) {

	b, _ := json.Marshal(curr)

	err3 := redis.Session.Set(REDIS_BTC_CURRENT, b, 0).Err()
	if err3 != nil {
		fmt.Println("Error Save Report BTC Current")
	}

	fmt.Println("SaveReport_BTC_Current : ", REDIS_BTC_CURRENT, curr)
}

func SaveReport_BTC_Fees(fees float64) {

	b, _ := json.Marshal(fees)

	err4 := redis.Session.Set(REDIS_BTC_FEES, b, 0).Err()
	if err4 != nil {
		fmt.Println("Error Save Report BTC Fees")
	}

	fmt.Println("SaveReport_BTC_Fees : ", REDIS_BTC_FEES, fees)
}

/*---------------------------------------------------------------------------*/

func LoadReport_LTC() {

	de, err1 := redis.Session.Get(REDIS_LTC_DEPOSIT).Result()
	if err1 != nil {
		fmt.Println("Error Load Report LTC Deposit")
		return
	}

	with, err2 := redis.Session.Get(REDIS_LTC_WITHDRAW).Result()
	if err2 != nil {
		fmt.Println("Error Load Report LTC Withdraw")
		return
	}

	curr, err3 := redis.Session.Get(REDIS_LTC_CURRENT).Result()
	if err3 != nil {
		fmt.Println("Error Load Report LTC Current")
		return
	}

	fees, err4 := redis.Session.Get(REDIS_LTC_FEES).Result()
	if err4 != nil {
		fmt.Println("Error Load Report LTC Fees")
		return
	}

	df, _ := strconv.ParseFloat(de, 64)
	LTC_DEPOSIT = df

	wf, _ := strconv.ParseFloat(with, 64)
	LTC_WITHDRAW = wf

	cf, _ := strconv.ParseFloat(curr, 64)
	LTC_CURRENT = cf

	bf, _ := strconv.ParseFloat(fees, 64)
	LTC_FEES = bf

	fmt.Println("Redis LTC Deposit : ", de, "  - Withdraw : ", with, "  - Current : ", curr, "  - Fees : ", fees)
}

func SaveReport_LTC_Deposit(de float64) {

	b, _ := json.Marshal(de)

	err1 := redis.Session.Set(REDIS_LTC_DEPOSIT, b, 0).Err()
	if err1 != nil {
		fmt.Println("Error Save Report LTC Deposit")
	}

	fmt.Println("SaveReport_LTC_Deposit : ", REDIS_LTC_DEPOSIT, de)
}

func SaveReport_LTC_Withdraw(with float64) {

	b, _ := json.Marshal(with)

	err2 := redis.Session.Set(REDIS_LTC_WITHDRAW, b, 0).Err()
	if err2 != nil {
		fmt.Println("Error Save Report LTC Withdraw")
	}

	fmt.Println("SaveReport_LTC_Withdraw : ", REDIS_LTC_WITHDRAW, with)
}

func SaveReport_LTC_Current(curr float64) {

	b, _ := json.Marshal(curr)

	err3 := redis.Session.Set(REDIS_LTC_CURRENT, b, 0).Err()
	if err3 != nil {
		fmt.Println("Error Save Report LTC Current")
	}

	fmt.Println("SaveReport_LTC_Current : ", REDIS_LTC_CURRENT, curr)
}

func SaveReport_LTC_Fees(fees float64) {

	b, _ := json.Marshal(fees)

	err4 := redis.Session.Set(REDIS_LTC_FEES, b, 0).Err()
	if err4 != nil {
		fmt.Println("Error Save Report LTC Fees")
	}

	fmt.Println("SaveReport_LTC_Fees : ", REDIS_LTC_FEES, fees)
}

/*---------------------------------------------------------------------------*/

func LoadReport_BCH() {

	de, err1 := redis.Session.Get(REDIS_BCH_DEPOSIT).Result()
	if err1 != nil {
		fmt.Println("Error Load Report BCH Deposit")
		return
	}

	with, err2 := redis.Session.Get(REDIS_BCH_WITHDRAW).Result()
	if err2 != nil {
		fmt.Println("Error Load Report BCH Withdraw")
		return
	}

	curr, err3 := redis.Session.Get(REDIS_BCH_CURRENT).Result()
	if err3 != nil {
		fmt.Println("Error Load Report BCH Current")
		return
	}

	fees, err4 := redis.Session.Get(REDIS_BCH_FEES).Result()
	if err4 != nil {
		fmt.Println("Error Load Report BCH Fees")
		return
	}

	df, _ := strconv.ParseFloat(de, 64)
	BCH_DEPOSIT = df

	wf, _ := strconv.ParseFloat(with, 64)
	BCH_WITHDRAW = wf

	cf, _ := strconv.ParseFloat(curr, 64)
	BCH_CURRENT = cf

	bf, _ := strconv.ParseFloat(fees, 64)
	BCH_FEES = bf

	fmt.Println("Redis BCH Deposit : ", de, "  - Withdraw : ", with, "  - Current : ", curr, "  - Fees : ", fees)
}

func SaveReport_BCH_Deposit(de float64) {

	b, _ := json.Marshal(de)

	err1 := redis.Session.Set(REDIS_BCH_DEPOSIT, b, 0).Err()
	if err1 != nil {
		fmt.Println("Error Save Report BCH Deposit")
	}

	fmt.Println("SaveReport_BCH_Deposit : ", REDIS_BCH_DEPOSIT, de)
}

func SaveReport_BCH_Withdraw(with float64) {

	b, _ := json.Marshal(with)

	err2 := redis.Session.Set(REDIS_BCH_WITHDRAW, b, 0).Err()
	if err2 != nil {
		fmt.Println("Error Save Report BCH Withdraw")
	}

	fmt.Println("SaveReport_BCH_Withdraw : ", REDIS_BCH_WITHDRAW, with)
}

func SaveReport_BCH_Current(curr float64) {

	b, _ := json.Marshal(curr)

	err3 := redis.Session.Set(REDIS_BCH_CURRENT, b, 0).Err()
	if err3 != nil {
		fmt.Println("Error Save Report BCH Current")
	}

	fmt.Println("SaveReport_BCH_Current : ", REDIS_BCH_CURRENT, curr)
}

func SaveReport_BCH_Fees(fees float64) {

	b, _ := json.Marshal(fees)

	err4 := redis.Session.Set(REDIS_BCH_FEES, b, 0).Err()
	if err4 != nil {
		fmt.Println("Error Save Report BCH Fees")
	}

	fmt.Println("SaveReport_BCH_Fees : ", REDIS_BCH_FEES, fees)
}

/*---------------------------------------------------------------------------*/

func LoadReport_ETH() {

	de, err1 := redis.Session.Get(REDIS_ETH_DEPOSIT).Result()
	if err1 != nil {
		fmt.Println("Error Load Report ETH Deposit")
		return
	}

	with, err2 := redis.Session.Get(REDIS_ETH_WITHDRAW).Result()
	if err2 != nil {
		fmt.Println("Error Load Report ETH Withdraw")
		return
	}

	curr, err3 := redis.Session.Get(REDIS_ETH_CURRENT).Result()
	if err3 != nil {
		fmt.Println("Error Load Report ETH Current")
		return
	}

	fees, err4 := redis.Session.Get(REDIS_ETH_FEES).Result()
	if err4 != nil {
		fmt.Println("Error Load Report ETH Fees")
		return
	}

	df, _ := strconv.ParseFloat(de, 64)
	ETH_DEPOSIT = df

	wf, _ := strconv.ParseFloat(with, 64)
	ETH_WITHDRAW = wf

	cf, _ := strconv.ParseFloat(curr, 64)
	ETH_CURRENT = cf

	bf, _ := strconv.ParseFloat(fees, 64)
	ETH_FEES = bf

	fmt.Println("Redis ETH Deposit : ", de, "  - Withdraw : ", with, "  - Current : ", curr, "  - Fees : ", fees)
}

func SaveReport_ETH_Deposit(de float64) {

	b, _ := json.Marshal(de)

	err1 := redis.Session.Set(REDIS_ETH_DEPOSIT, b, 0).Err()
	if err1 != nil {
		fmt.Println("Error Save Report ETH Deposit")
	}

	fmt.Println("SaveReport_ETH_Deposit : ", REDIS_ETH_DEPOSIT, de)
}

func SaveReport_ETH_Withdraw(with float64) {

	b, _ := json.Marshal(with)

	err2 := redis.Session.Set(REDIS_ETH_WITHDRAW, b, 0).Err()
	if err2 != nil {
		fmt.Println("Error Save Report ETH Withdraw")
	}

	fmt.Println("SaveReport_ETH_Withdraw : ", REDIS_ETH_WITHDRAW, with)
}

func SaveReport_ETH_Current(curr float64) {

	b, _ := json.Marshal(curr)

	err3 := redis.Session.Set(REDIS_ETH_CURRENT, b, 0).Err()
	if err3 != nil {
		fmt.Println("Error Save Report ETH Current")
	}

	fmt.Println("SaveReport_ETH_Current : ", REDIS_ETH_CURRENT, curr)
}

func SaveReport_ETH_Fees(fees float64) {

	b, _ := json.Marshal(fees)

	err4 := redis.Session.Set(REDIS_ETH_FEES, b, 0).Err()
	if err4 != nil {
		fmt.Println("Error Save Report ETH Fees")
	}

	fmt.Println("SaveReport_ETH_Fees : ", REDIS_ETH_FEES, fees)
}

/*---------------------------------------------------------------------------*/

func LoadReport_XLM() {

	de, err1 := redis.Session.Get(REDIS_XLM_DEPOSIT).Result()
	if err1 != nil {
		fmt.Println("Error Load Report XLM Deposit")
		return
	}

	with, err2 := redis.Session.Get(REDIS_XLM_WITHDRAW).Result()
	if err2 != nil {
		fmt.Println("Error Load Report XLM Withdraw")
		return
	}

	curr, err3 := redis.Session.Get(REDIS_XLM_CURRENT).Result()
	if err3 != nil {
		fmt.Println("Error Load Report XLM Current")
		return
	}

	fees, err4 := redis.Session.Get(REDIS_XLM_FEES).Result()
	if err4 != nil {
		fmt.Println("Error Load Report XLM Fees")
		return
	}

	df, _ := strconv.ParseFloat(de, 64)
	XLM_DEPOSIT = df

	wf, _ := strconv.ParseFloat(with, 64)
	XLM_WITHDRAW = wf

	cf, _ := strconv.ParseFloat(curr, 64)
	XLM_CURRENT = cf

	bf, _ := strconv.ParseFloat(fees, 64)
	XLM_FEES = bf

	fmt.Println("Redis XLM Deposit : ", de, "  - Withdraw : ", with, "  - Current : ", curr, "  - Fees : ", fees)
}

func SaveReport_XLM_Deposit(de float64) {

	b, _ := json.Marshal(de)

	err1 := redis.Session.Set(REDIS_XLM_DEPOSIT, b, 0).Err()
	if err1 != nil {
		fmt.Println("Error Save Report XLM Deposit")
	}

	fmt.Println("SaveReport_XLM_Deposit : ", REDIS_XLM_DEPOSIT, de)
}

func SaveReport_XLM_Withdraw(with float64) {

	b, _ := json.Marshal(with)

	err2 := redis.Session.Set(REDIS_XLM_WITHDRAW, b, 0).Err()
	if err2 != nil {
		fmt.Println("Error Save Report XLM Withdraw")
	}

	fmt.Println("SaveReport_XLM_Withdraw : ", REDIS_XLM_WITHDRAW, with)
}

func SaveReport_XLM_Current(curr float64) {

	b, _ := json.Marshal(curr)

	err3 := redis.Session.Set(REDIS_XLM_CURRENT, b, 0).Err()
	if err3 != nil {
		fmt.Println("Error Save Report XLM Current")
	}

	fmt.Println("SaveReport_XLM_Current : ", REDIS_XLM_CURRENT, curr)
}

func SaveReport_XLM_Fees(fees float64) {

	b, _ := json.Marshal(fees)

	err4 := redis.Session.Set(REDIS_XLM_FEES, b, 0).Err()
	if err4 != nil {
		fmt.Println("Error Save Report XLM Fees")
	}

	fmt.Println("SaveReport_XLM_Fees : ", REDIS_XLM_FEES, fees)
}
