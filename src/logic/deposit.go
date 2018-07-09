package logic

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"api"
	"module/dbScan"
)

func Do_Deposit(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_Deposit : ", string(params)) // {"coin" : "ETH/BTC/ERC20/LTC", "contract" : ""}

	resp := Writer{Api: api.DEPOSIT}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_deposit(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)

	} else {

		coin := request["coin"].(string)

		// GO-0 : check coin type
		if coin != "BTC" && coin != "LTC" && coin != "BCH" && coin != "ETH" && coin != "ERC20" && coin != "XLM" {
			resp.Status = -2
			resp.Error = "Error Coin !!!"

			fmt.Println(resp.Error)
		}

		// Go-1 : new deposit address
		if resp.Status == 0 {

			deposit_Address, seed := genAddress(coin)

			switch coin {
			case "BTC", "LTC", "BCH", "ETH", "XLM":
				{
					deposit_Map := dbScan.NewDepositCoin(seed, deposit_Address, coin)
					dbScan.HMAP_DEPOSIT[deposit_Address] = deposit_Map

					resp.Data = bson.M{"coin": coin, "deposit": deposit_Address}
				}
			case "ERC20":
				{
					contract := request["contract"].(string)

					// GO : verifyAddress Contract
					if !verifyAddress(coin, contract) {
						resp.Status = -3
						resp.Error = "Error Address Not Contract !!!"

						fmt.Println(resp.Error)
					}

					// Go : get rating contract on etherscan
					if resp.Status == 0 {

						ratio := getRatingFromEtherScan(contract, "test")

						if ratio <= float64(0) {
							resp.Status = -4
							resp.Error = "Contract Not On EtherScan !!!" + contract

							fmt.Println(resp.Error)
						}
					}

					// Go : new depost address
					if resp.Status == 0 {

						deposit_Map := dbScan.NewDepositERC20(deposit_Address, contract, coin)
						dbScan.HMAP_DEPOSIT[deposit_Address] = deposit_Map

						resp.Data = bson.M{"coin": coin, "deposit": deposit_Address, "contract": contract}
					}

				}
			}

			fmt.Println(resp.Data)
		}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_deposit(request map[string]interface{}) bool {

	if len(request) != 2 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	if contract, ok := request["contract"]; !ok || !reflectString(contract) {
		return false
	}

	return true
}
