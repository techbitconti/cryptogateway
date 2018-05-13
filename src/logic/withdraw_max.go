package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"math"
	"net/http"

	"lib/eth"

	"gopkg.in/mgo.v2/bson"
)

func Do_WithdrawMax(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_WithdrawMax : ", string(params)) // {"coin" : "ETH/BTC/ERC20", "deposit" : "0x923eac92bda97a4348968a1e7d64834236319b3f"}

	resp := Writer{Api: api.WITHDRAW_MAX}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_withdraw_max(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)
	} else {

		coin := request["coin"].(string)
		addr := request["deposit"].(string)

		// GO-0 : check coin type
		if coin != "BTC" && coin != "ETH" && coin != "ERC20" {
			resp.Status = -2
			resp.Error = "Error Coin !!!"

			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {

			if coin == "ERC20" {
				coin = "ETH"
			}

			isValid := verifyAddress(coin, addr)
			if !isValid {
				resp.Status = -3
				resp.Error = "Invalid Address !!!"
				fmt.Println(resp.Error)

			} else {

				balance := getBalance(coin, addr)
				max := float64(0)

				switch coin {
				case "BTC":
					{
						fee := float64(0.001)
						max = balance - fee //BTC
					}

				case "ETH", "ERC20":
					{
						wei := balance * math.Pow10(18)

						gas := float64(21000)

						gasPriceBigI, _ := eth.SuggestGasPrice()
						gasPriceWei := gasPriceBigI.Int64()
						gasPrice := float64(gasPriceWei)

						max = wei - gas*gasPrice
						max /= math.Pow10(18) //ETH
					}
				}

				resp.Status = 0
				resp.Data = bson.M{"address": addr, "balance": balance, "max": max}
			}
		}

	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_withdraw_max(request map[string]interface{}) bool {

	if len(request) != 2 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	if deposit, ok := request["deposit"]; !ok || !reflectString(deposit) {
		return false
	}

	return true
}
