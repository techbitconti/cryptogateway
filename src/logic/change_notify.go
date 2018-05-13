package logic

import (
	"api"
	"encoding/json"
	"fmt"
	"net/http"

	"config"

	"gopkg.in/mgo.v2/bson"
)

func Do_ChangeNofify(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_ChangeNofify : ", string(params)) // {"notify" : ""}

	resp := Writer{Api: api.CHANGE_NOTIFY}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_changeNofify(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)
	} else {

		notify := request["notify"].(string)

		config.NOTIFY_BALANCE = notify

		resp.Data = bson.M{"result": "success", "notify": notify}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_changeNofify(request map[string]interface{}) bool {

	if len(request) != 1 {
		return false
	}

	if notity, ok := request["notity"]; !ok || !reflectString(notity) {
		return false
	}

	return true
}
