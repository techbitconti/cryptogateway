package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"net/http"
)

func Do_ChangeIP(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_ChangeIP : ", string(params)) // {"old" : "ip:port", "new" : "ip_port"}

	resp := Writer{Api: api.CHANGE_IP}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_changeIP(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)
	} else {

	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_changeIP(request map[string]interface{}) bool {

	if len(request) != 2 {
		return false
	}

	if old, ok := request["old"]; !ok || !reflectString(old) {
		return false
	}

	if nnew, ok := request["new"]; !ok || !reflectString(nnew) {
		return false
	}

	return true
}
