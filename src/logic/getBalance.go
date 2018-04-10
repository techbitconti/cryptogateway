package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"lib/btc"
	"lib/eth"

	"gopkg.in/mgo.v2/bson"
)

func Do_GetBalance(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_GetBalance : ", string(params)) // {"coin" : "ETH/BTC/ERC20", "address" : ""}

	resp := Writer{Api: api.BALANCE}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_getBalance(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)

	} else {

		coin := request["coin"].(string)
		addr := request["address"].(string)

		isValid := verifyAddress(coin, addr)
		if !isValid {
			resp.Status = -2
			resp.Error = "Invalid Address !!!"
			fmt.Println(resp.Error)

		} else {
			balance := getBalance(coin, addr)
			resp.Status = 0
			resp.Data = bson.M{"address": addr, "balance": balance}
		}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func verifyAddress(coin string, addr string) bool {

	switch coin {
	case "BTC":
		address, err := btc.ValidateAddress(addr)
		if err != nil || !address.IsValid {
			return false
		}

	case "ETH":
		if !eth.IsHexAddress(addr) {
			return false
		}
	}

	return true
}

func getBalance(coin string, addr string) float64 {

	switch coin {
	case "BTC":
		amount := btc.GetBalance(addr)
		fmt.Println("getBalance BTC : ", amount)
		return amount

	case "ETH":
		bigInt := eth.GetBalance(addr)
		bigFloat := new(big.Float).SetInt(bigInt)
		f, _ := bigFloat.Float64()

		fmt.Println("getBalance ETH : ", f)
		return f
	}

	return float64(0)
}

func check_getBalance(request map[string]interface{}) bool {

	if len(request) != 2 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	if addr, ok := request["address"]; !ok || !reflectString(addr) {
		return false
	}

	return true
}
