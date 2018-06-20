package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"net/http"

	"lib/btc"
	"lib/eth"
	"lib/ltc"
)

func Do_ListAddress(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_ListAddress : ", string(params)) // {"coin" : "ETH/BTC/LTC"}

	resp := Writer{Api: api.LIST_ADDRESS}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_listAddress(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)

	} else {

		coin := request["coin"].(string)

		// GO-0 : check coin type
		if coin != "BTC" && coin != "LTC" && coin != "ETH" {
			resp.Status = -2
			resp.Error = "Error Coin !!!"

			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {
			list := make([]string, 0)

			switch coin {
			case "BTC":
				listAddr := btc.ListAddress()
				for _, v := range listAddr {
					list = append(list, v.String())
				}
			case "LTC":
				listAddr := ltc.ListAddress()
				for _, v := range listAddr {
					list = append(list, v.String())
				}
			case "ETH":
				arr := eth.GetAccounts()
				list = append(list, arr...)
			}

			resp.Data = list
		}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_listAddress(request map[string]interface{}) bool {

	if len(request) != 1 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	return true
}
