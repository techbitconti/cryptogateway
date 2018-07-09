package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func Do_Transfer(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_Transfer : ", string(params)) // {"coin" : "ETH/BTC/BCH/LTC", "from" : "", "to" : "", "amount" : ""}

	resp := Writer{Api: api.TRANSFER}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_transfer(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)

	} else {

		coin := request["coin"].(string)
		from := request["from"].(string)
		to := request["to"].(string)
		amount := request["amount"].(string)

		// GO-0 : check coin type
		if coin != "BTC" && coin != "LTC" && coin != "BCH" && coin != "ETH" && coin != "XLM" {
			resp.Status = -2
			resp.Error = "Error Coin !!!"

			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {

			tx := sendCoin(coin, from, to, amount)
			fmt.Println("tx : ", tx)

			if tx == "" {
				resp.Status = -3
				resp.Error = "Can not transfer !!!" + coin
				fmt.Println(resp.Error)

			} else {

				resp.Data = bson.M{"tx": tx}
			}
		}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_transfer(request map[string]interface{}) bool {

	if len(request) != 4 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	if from, ok := request["from"]; !ok || !reflectString(from) {
		return false
	}

	if to, ok := request["to"]; !ok || !reflectString(to) {
		return false
	}

	if amount, ok := request["amount"]; !ok || !reflectString(amount) {
		return false
	}

	return true
}
