package logic

import (
	"api"
	"encoding/json"
	"fmt"
	//"math/big"
	"net/http"
	"strconv"

	"config"
	"module/dbScan"
	//"lib/btc"
	//"lib/eth"

	"gopkg.in/mgo.v2/bson"
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
			//			if resp.Status == 0 {
			//				if mDeposit.Status != dbScan.STATUS_PENDING {
			//					resp.Status = -6
			//					resp.Error = "Invalid Balance Of Deposit Address = 0 !!!"
			//					fmt.Println(resp.Error)
			//				}
			//			}

			// Go-4 : sendTransaction
			if resp.Status == 0 {

				switch coin {
				case "BTC", "ETH":
					{
						amountDespsit := mDeposit.Amount
						aMDe, _ := strconv.ParseFloat(amountDespsit, 64)
						aMWith, _ := strconv.ParseFloat(amountWithdraw, 64)

						balance := aMDe - aMWith

						if balance < float64(0) {
							resp.Status = -7
							resp.Error = "Amount deposit less then withdraw !!!"
							fmt.Println(resp.Error)
						} else {

							fromAdmin := ""
							switch coin {
							case "BTC":
								fromAdmin = config.BTC_SIM.Address
							case "ETH":
								fromAdmin = config.ETH_SIM.Address
							}

							if resp.Status == 0 {
								//Go-5 : sendFrom Deposit to Admin
								from := deposit_Address
								to := fromAdmin
								amount := amountWithdraw

								txToAdmin := sendCoin(coin, from, to, amount)
								fmt.Println("txToAdmin : ", txToAdmin)
								if txToAdmin == "" {
									resp.Status = -8
									resp.Error = "ETH not enough txToAdmin !!!"
									fmt.Println(resp.Error)
								}
							}

							if resp.Status == 0 {
								//Go-6 : sendFromAdmin to Receipt

								txFromAdmin := sendCoin(coin, fromAdmin, withdraw_Address, amountWithdraw)
								fmt.Println("txFromAdmin : ", txFromAdmin)

								//Go-7 : update HMAP_DEPOSIT
								mDeposit.Amount = strconv.FormatFloat(balance, 'f', -1, 64)
								//							if balance <= 0 {
								//								mDeposit.Status = dbScan.STATUS_WAITING
								//							}

								resp.Data = bson.M{"tx": txFromAdmin}
							}
						}
					}
				case "ERC20":
					{
						rating := getRatingFromEtherScan(mDeposit.AddressContract, "test")
						if rating <= float64(0) {
							resp.Status = -8
							resp.Error = "Contract Not On EtherScan !!!" + mDeposit.AddressContract
							fmt.Println(resp.Error)
						}

						if resp.Status == 0 {
							ethDeposit, _ := strconv.ParseFloat(mDeposit.Amount, 64)
							tokenDeposit := int64(ethDeposit / rating)
							fmt.Println("ethDeposit : ", ethDeposit, "tokenDeposit :", tokenDeposit)

							tokenWidthDraw, _ := strconv.ParseInt(amountWithdraw, 0, 64)
							ethWidthDraw := float64(tokenWidthDraw) * rating
							fmt.Println("ethWidthDraw : ", ethWidthDraw, "tokenWidthDraw :", tokenWidthDraw)

							ethBalance := ethDeposit - ethWidthDraw
							tokenBalance := tokenDeposit - tokenWidthDraw

							fmt.Println("ethBalance : ", ethBalance)
							fmt.Println("tokenBalance : ", tokenBalance)

							if tokenBalance <= int64(0) || ethBalance < float64(0) {
								resp.Status = -9
								resp.Error = "Token deposit less then Token Withdraw !!!"
								fmt.Println(resp.Error)
							}

							if resp.Status == 0 {
								//Go-5 : sendFrom Deposit to Admin
								from := deposit_Address
								to := config.ETH_SIM.Address
								amount := strconv.FormatFloat(ethWidthDraw, 'f', -1, 64)

								txToAdmin := sendCoin("ETH", from, to, amount)
								fmt.Println("txToAdmin : ", txToAdmin)
								if txToAdmin == "" {
									resp.Status = -10
									resp.Error = "ETH not enough txToAdmin !!!"
									fmt.Println(resp.Error)
								}
							}

							if resp.Status == 0 {
								//Go-6 : sendFromAdmin to Receipt
								tokens := strconv.FormatInt(tokenWidthDraw, 10)
								txFromAdmin := sendERC20(mDeposit.AddressContract, withdraw_Address, tokens)
								fmt.Println("txFromAdmin : ", txFromAdmin)
								if txFromAdmin == "" {
									resp.Status = -10
									resp.Error = "ERC20 not enough txFromAdmin !!!"
									fmt.Println(resp.Error)
								}
								resp.Data = bson.M{"tx": txFromAdmin}
							}

							if resp.Status == 0 {
								//Go-7 : update HMAP_DEPOSIT
								mDeposit.Amount = strconv.FormatFloat(ethBalance, 'f', -1, 64)
								//								if ethBalance <= 0 {
								//									mDeposit.Status = dbScan.STATUS_WAITING
								//								}
							}

						}
					}
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
