package dbScan

import (
	"encoding/json"
	"fmt"
	"strconv"

	"db/redis"
)

var REDIS_KEYS_DEPOSIT = "keys_deposit"

var REDIS_BTC_DEPOSIT = "btc_deposit"
var REDIS_BTC_WITHDRAW = "btc_withdraw"
var REDIS_BTC_CURRENT = "btc_current"
var REDIS_BTC_FEES = "btc_fees"

var REDIS_ETH_DEPOSIT = "eth_deposit"
var REDIS_ETH_WITHDRAW = "eth_withdraw"
var REDIS_ETH_CURRENT = "eth_current"
var REDIS_ETH_FEES = "eth_fees"

var BTC_DEPOSIT float64
var BTC_WITHDRAW float64
var BTC_CURRENT float64
var BTC_FEES float64

var ETH_DEPOSIT float64
var ETH_WITHDRAW float64
var ETH_CURRENT float64
var ETH_FEES float64

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
