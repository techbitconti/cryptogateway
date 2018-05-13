package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"config"

	"gopkg.in/mgo.v2/bson"
)

func Do_RegisApp(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_RegisApp : ", string(params)) // {"url" : "ip:port", "notify_balance" : "notify", "pass_wallet" : ""}

	resp := Writer{Api: api.REGIS_APP}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_regis_app(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)
	} else {

		url := request["url"].(string)
		notify_balance := request["notify_balance"].(string)
		pass_wallet := request["pass_wallet"].(string)

		host, port, err := net.SplitHostPort(url)
		if err != nil {
			resp.Status = -2
			resp.Error = "Error SplitHostPort !!!"
			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {
			config.IP_ALLOW = host
			config.PORT_ALLOW = port
			config.NOTIFY_BALANCE = notify_balance
			config.PASS_WALLET = pass_wallet

			resp.Data = bson.M{"IP_ALLOW": config.IP_ALLOW}
		}

	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_regis_app(request map[string]interface{}) bool {

	if len(request) != 3 {
		return false
	}

	if url, ok := request["url"]; !ok || !reflectString(url) {
		return false
	}

	if notify_balance, ok := request["notify_balance"]; !ok || !reflectString(notify_balance) {
		return false
	}

	if pass_wallet, ok := request["pass_wallet"]; !ok || !reflectString(pass_wallet) {
		return false
	}

	return true
}
