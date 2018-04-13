package logic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"api"
	"module/dbScan"
)

func Do_Deposit(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_Deposit : ", string(params)) // {"coin" : "ETH/BTC/ERC20"}

	resp := Writer{Api: api.DEPOSIT}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_genAddress(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)

	} else {

		coin := request["coin"].(string)
		deposit_Address := genAddress(coin)

		// GO-0 : check coin type
		if coin != "BTC" && coin != "ETH" && coin != "ERC20" {
			resp.Status = -2
			resp.Error = "Error Coin !!!"

			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {

			deposit_Map := map[string]string{
				"status": "waiting",
				"coin":   coin,
				"amount": "0",
			}
			dbScan.HMAP_DEPOSIT[deposit_Address] = deposit_Map

			resp.Status = 0
			resp.Data = bson.M{"coin": coin, "deposit": deposit_Address}
		}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_genAddress(request map[string]interface{}) bool {

	if len(request) != 1 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	return true
}
