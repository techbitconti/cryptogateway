package logic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"api"
)

func Do_GetBalance(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_GetBalance : ", string(params)) // {"coin" : "ETH/BTC/ERC20", "address" : "", "contract" : ""}

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

		// GO-0 : check coin type
		if coin != "BTC" && coin != "ETH" && coin != "ERC20" {
			resp.Status = -2
			resp.Error = "Error Coin !!!"

			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {

			switch coin {
			case "BTC", "ETH":
				{
					isValid := verifyAddress(coin, addr)
					if !isValid {
						resp.Status = -3
						resp.Error = "Invalid Address !!!"
						fmt.Println(resp.Error)

					} else {

						balance := getBalance(coin, addr)
						resp.Status = 0
						resp.Data = bson.M{"address": addr, "balance": balance}
					}
				}

			case "ERC20":
				{
					contract := request["contract"].(string)

					isValid := verifyAddress("ETH", addr)
					isValid = verifyAddress(coin, contract)

					if !isValid {
						resp.Status = -3
						resp.Error = "Invalid Address !!!"
						fmt.Println(resp.Error)

					} else {

						balance := getBalanceOf(contract, addr)
						resp.Status = 0
						resp.Data = bson.M{"address": addr, "balance": balance}
					}
				}
			}

		}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_getBalance(request map[string]interface{}) bool {

	if len(request) != 4 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	if addr, ok := request["address"]; !ok || !reflectString(addr) {
		return false
	}

	if contract, ok := request["contract"]; !ok || !reflectString(contract) {
		return false
	}
	return true
}
