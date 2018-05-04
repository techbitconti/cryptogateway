package logic

import (
	"api"
	"encoding/json"
	"fmt"
	//"math/big"
	"net/http"
	"strconv"

	//"lib/btc"
	//"lib/eth"
	"gopkg.in/mgo.v2/bson"

	"module/dbScan"
)

func Do_Withdraw(ip string, w http.ResponseWriter, params []byte) {
	fmt.Println("Do_Withdraw : ", string(params)) // {"coin" : "ETH/BTC/ERC20", "deposit" : "", "withdraw" : "", "amount" : ""}

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
		amountWithdraw := request["amount"].(string)

		// GO-0 : check coin type
		if coin != "BTC" && coin != "ETH" && coin != "ERC20" {
			resp.Status = -2
			resp.Error = "Error Coin !!!"

			fmt.Println(resp.Error)
		}

		if resp.Status == 0 {

			var cCoin string
			if coin == "ERC20" {
				cCoin = "ETH"
			}

			// GO-1 :  verifyAddress
			if !verifyAddress(cCoin, deposit_Address) || !verifyAddress(cCoin, withdraw_Address) {
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
			}

			if resp.Status == 0 {
				if mDeposit.Coin != coin {
					resp.Status = -5
					resp.Error = "Invalid DEPOSIT Coin !!!"
					fmt.Println(resp.Error)
				}
			}

			// Go-3 : check deposit balance
			if resp.Status == 0 {

				if mDeposit.Status != dbScan.STATUS_PENDING {

					resp.Status = -6
					resp.Error = "Invalid Balance Of Deposit Address = 0 !!!"
					fmt.Println(resp.Error)
				}
			}

			// Go-4 : sendTransaction
			if resp.Status == 0 {

				amountDespsit := mDeposit.Amount

				aMDe, _ := strconv.ParseFloat(amountDespsit, 64)
				aMWith, _ := strconv.ParseFloat(amountWithdraw, 64)

				balance := aMWith - aMDe

				if balance < float64(0) {
					resp.Status = -7
					resp.Error = "Amount deposit less then withdraw !!!"
					fmt.Println(resp.Error)

				} else {

					fromAdmin := ""
					switch coin {
					case "BTC":
						fromAdmin = config.BTC_SIM.Address
					case "ETH", "ERC20":
						fromAdmin = config.ETH_SIM.Address
					}

					//Go-5 : sendFrom Deposit to Admin
					mTxDeposit := map[string]string{
						"addr":     fromAdmin,
						"amount":   amountWithdraw,
						"receiver": "NaN",
					}
					txFromDe := sendCoin(coin, deposit_Address, mTxDeposit)
					fmt.Println("txFromDe : ", txFromDe)

					//Go-6 : sendFromAdmin toreceipt
					var txFromAdmin string
					switch coin {
					case "BTC", "ETH":
						{
							mTxAdmin := map[string]string{
								"addr":     withdraw_Address,
								"amount":   amountWithdraw,
								"receiver": "NaN",
							}
							txFromAdmin = sendCoin(coin, fromAdmin, mTxAdmin)
							fmt.Println("txFromAdmin : ", txFromAdmin)
						}
					case "ERC20":
						{
							tokens := amountWithdraw
							txFromAdmin = sendERC20(mDeposit.AddressContract, withdraw_Address, tokens)
						}
					}

					//Go-7 : update HMAP_DEPOSIT
					mDeposit.Amount = strconv.FormatFloat(balance, 'f', -1, 64)
					if balance <= 0 {
						mDeposit.Status = dbScan.STATUS_WAITING
					}

					resp.Data = bson.M{"tx": txFromAdmin}
				}
			}
		}
	}

	data, _ := json.Marshal(resp)
	w.Write(data)
}

func check_withdraw(request map[string]interface{}) bool {

	if len(request) != 4 {
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

	if amount, ok := request["amount"]; !ok || !reflectString(amount) {
		return false
	}

	return true
}
