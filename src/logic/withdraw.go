package logic

import (
	"api"
	"encoding/json"
	"fmt"
	//"math/big"
	"net/http"
	//"lib/btc"
	//"lib/eth"
	"gopkg.in/mgo.v2/bson"

	"module/dbScan"
)

func Do_Withdraw(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_Withdraw : ", string(params)) // {"coin" : "ETH/BTC/ERC20", "deposit" : "", "withdraw" : ""}

	resp := Writer{Api: api.WITHDRAW}

	request := map[string]interface{}{}
	json.Unmarshal(params, &request)

	ok := check_withdraw(request)
	if !ok {
		resp.Status = -1
		resp.Error = "Invalid input !!!"
		fmt.Println(resp.Error)

	} else {

		coin := request["coin"].(string)
		deposit_Address := request["deposit"].(string)
		withdraw_Address := request["withdraw"].(string)

		// GO-0 : check coin type
		if coin != "BTC" && coin != "ETH" { //}&& coin != "ERC20" {
			resp.Status = -2
			resp.Error = "Error Coin !!!"

			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {

			// GO-1 :  verifyAddress
			if !verifyAddress(coin, deposit_Address) || !verifyAddress(coin, withdraw_Address) {
				resp.Status = -3
				resp.Error = "Invalid verifyAddress !!!"
				fmt.Println(resp.Error)
			}
		}

		// GO-2 : check HMAP_DEPOSIT
		if resp.Status == 0 {

			mDeposit, ok := dbScan.HMAP_DEPOSIT[deposit_Address]
			if !ok {
				resp.Status = -4
				resp.Error = "Invalid Deposit Address !!!"
				fmt.Println(resp.Error)

			} else {

				if mDeposit["coin"] != coin {
					resp.Status = -5
					resp.Error = "Invalid DEPOSIT Coin !!!"
					fmt.Println(resp.Error)
				}
			}
		}

		// Go-3 : check deposit balance
		if resp.Status == 0 {

			if dbScan.HMAP_DEPOSIT[deposit_Address]["status"] != "pending" {

				resp.Status = -6
				resp.Error = "Invalid Balance Of Deposit Address = 0 !!!"
				fmt.Println(resp.Error)
			}
		}

		// Go-4 : sendTransaction
		if resp.Status == 0 {

			amount := dbScan.HMAP_DEPOSIT[deposit_Address]["amount"]
			fromAdmin := getAddressAdmin(coin)

			//Go-5: sendFrom Deposit to Admin
			mTxDeposit := map[string]string{
				"addr":     fromAdmin,
				"amount":   amount,
				"receiver": "NaN",
			}
			txFromDe := sendTransaction(coin, deposit_Address, mTxDeposit)
			fmt.Println("txFromDe : ", txFromDe)

			dbScan.HMAP_DEPOSIT[deposit_Address]["status"] = "waiting"

			//Go-6 : sendFromAdmin
			mTxAdmin := map[string]string{
				"addr":     withdraw_Address,
				"amount":   amount,
				"receiver": "NaN",
			}
			txFromAdmin := sendTransaction(coin, fromAdmin, mTxAdmin)
			fmt.Println("txFromAdmin : ", txFromAdmin)

			resp.Data = bson.M{"tx": txFromAdmin}
		}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_withdraw(request map[string]interface{}) bool {

	if len(request) != 3 {
		return false
	}

	if coin, ok := request["coin"]; !ok || !reflectString(coin) {
		return false
	}

	if deposit, ok := request["deposit"]; !ok || !reflectString(deposit) {
		return false
	}

	if withdraw, ok := request["withdraw"]; !ok || !reflectString(withdraw) {
		return false
	}

	return true
}
