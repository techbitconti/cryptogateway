package logic

import (
	"api"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"config"
	"lib/btc"
	"lib/eth"
	"module/dbScan"

	"gopkg.in/mgo.v2/bson"
)

func Do_GenAddress(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_GenAddress : ", string(params)) // {"coin" : "ETH/BTC/ERC20"}

	resp := Writer{Api: api.ADDRESS}

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

		deposit_Map := map[string]string{
			"status": "waiting",
			"coin":   coin,
			"amount": "0",
		}
		dbScan.HMAP_DEPOSIT[deposit_Address] = deposit_Map

		resp.Status = 0
		resp.Data = bson.M{"coin": coin, "deposit": deposit_Address}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func genAddress(coin string) (address string) {

	switch coin {
	case "BTC":
		address = genAddressBTC()
	case "ETH":
		address = genAddressETH()
	}

	fmt.Println("genAddress --- ", address)

	return
}

func genAddressBTC() string {

	utc := time.Now().Unix()
	decode := strconv.FormatInt(utc, 10)
	encode := base64.StdEncoding.EncodeToString([]byte(decode))

	address, _ := btc.GetNewAddress(encode)

	return address.String()
}

func genAddressETH() string {

	address, _ := eth.NewAccount(config.Path_ETH, "123456")

	return address.Hex()
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
