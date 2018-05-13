package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"net/http"
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
