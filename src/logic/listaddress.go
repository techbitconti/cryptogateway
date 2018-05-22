package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"net/http"

	"lib/btc"
	"lib/eth"
)

func Do_ListAddress(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_ListAddress : ", string(params)) // {"coin" : "ETH/BTC"}

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

		list := make([]string, 0)

		switch coin {
		case "BTC":
			listAddr := btc.ListAddress()
			for _, v := range listAddr {
				list = append(list, v.String())
			}
		case "ETH":
			list = eth.GetAccounts()
		}

		resp.Data = list

	}
}

func check_listAddress(request map[string]interface{}) bool {

	if len(request) != 4 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	return true
}
