package logic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"api"
)

func Do_VerifyTx(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_VerifyTx : ", string(params)) // {"coin" : "ETH/BTC", "tx" : ""}

	resp := Writer{Api: api.VERIFY}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_verify(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)

	} else {

		coin := request["coin"].(string)
		tx := request["tx"].(string)

		// GO-0 : check coin type
		if coin != "BTC" && coin != "ETH" {
			resp.Status = -2
			resp.Error = "Error Coin !!!"

			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {

			receipt, ok := verityTX(coin, tx)
			if !ok {
				resp.Status = -3
				resp.Error = "Unkown Transaction !!!"

				fmt.Println(resp.Error)
			} else {

				resp.Data = receipt
			}
		}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_verify(request map[string]interface{}) bool {

	if len(request) != 2 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	if tx, ok := request["tx"]; !ok || !reflectString(tx) {
		return false
	}

	return true
}
