package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"net/http"
)

func Do_ChangeNofify(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_ChangeNofify : ", string(params)) // {"old" : "notify" , "new" : "notity"}

	resp := Writer{Api: api.CHANGE_NOTIFY}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_changeNofify(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)
	} else {

	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_changeNofify(request map[string]interface{}) bool {

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
