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

func Do_ChangeIP(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_ChangeIP : ", string(params)) // {"url" : "ip_port"}

	resp := Writer{Api: api.CHANGE_IP}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_changeIP(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)
	} else {

		url := request["url"].(string)

		host, port, err := net.SplitHostPort(url)
		if err != nil {
			resp.Status = -2
			resp.Error = "Error SplitHostPort !!!"
			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {

			if ip != config.IP_ALLOW {
				resp.Status = -3
				resp.Error = "Error LAST IP Not Same !!!"
				fmt.Println(resp.Error)
			}
		}

		if resp.Status == 0 {

			config.IP_ALLOW = host
			config.PORT_ALLOW = port

			resp.Data = bson.M{"IP_ALLOW": config.IP_ALLOW, "PORT_ALLOW": config.PORT_ALLOW}
		}

	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_changeIP(request map[string]interface{}) bool {

	if len(request) != 1 {
		return false
	}

	if url, ok := request["url"]; !ok || !reflectString(url) {
		return false
	}

	return true
}
