package dbScan

import (
	"encoding/json"
	"fmt"
	"lib/btc"
	"lib/eth"
	"math"
	"math/big"

	"config"
	"db/redis"

	"github.com/ethereum/go-ethereum/common"
)

var HMAP_DEPOSIT = map[string]*Deposit{}

var REDIS_KEYS_DEPOSIT = "keys_deposit"

func Start() {

	loadDeposit()
}

func loadDeposit() {

	bKeys, _ := redis.Session.Get(REDIS_KEYS_DEPOSIT).Bytes()

	var keys []string
	json.Unmarshal(bKeys, &keys)

	for _, key := range keys {

		deStr := redis.Session.Get(key).Val()

		if deStr != "" {

			b, _ := json.Marshal(deStr)
			var de Deposit
			json.Unmarshal(b, &de)

			HMAP_DEPOSIT[key] = &de

			fmt.Println("REDIS_KEYS_DEPOSIT: ", key, string(b))
		}
	}
}

func saveDeposit(de Deposit) {

	data, _ := json.Marshal(de)
	redis.Session.Set(de.AddressDeposit, data, 0).Err()

	arr, _ := redis.Session.Get(REDIS_KEYS_DEPOSIT).Result()
	var keys []string
	json.Unmarshal([]byte(arr), &keys)
	keys = append(keys, de.AddressDeposit)

	bKeys, _ := json.Marshal(keys)

	redis.Session.Set(REDIS_KEYS_DEPOSIT, bKeys, 0).Err()

	fmt.Println("saveDeposit : ", de.AddressDeposit)
}

func GetBalance(coin, addr string) float64 {

	switch coin {
	case "BTC":
		amount := btc.GetBalance(addr)
		fmt.Println("getBalance BTC : ", amount)
		return amount

	case "ETH":
		bigInt := eth.GetBalance(addr)
		bigFloat := new(big.Float).SetInt(bigInt)
		wei, _ := bigFloat.Float64()

		ether := wei / math.Pow10(18)

		fmt.Println("getBalance :", ether, "ether", " -- wei", wei)
		return ether
	}

	return 0
}

func GetBalanceOf(ercAddr, toAddr string) int64 {

	hex, err := eth.SolidityCallRaw(config.ETH_ADDR, ercAddr, `balanceOf(address)`, toAddr)
	if err != nil {
		fmt.Println("ERC20 Token not enough !!!")
		return int64(0)
	}
	//fmt.Println("balanceOf", common.BytesToHash(hex).Big(), new(big.Int).SetBytes(hex))

	return common.BytesToHash(hex).Big().Int64()
}
