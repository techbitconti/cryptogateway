package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"lib/eth"
	"module/dbScan"
)

func withdraw() {

	if len(os.Args) != 3 {
		return
	}

	coin := os.Args[0]
	to := os.Args[1]
	amountStr := os.Args[2]
	amount, aerr := strconv.ParseFloat(amountStr, 64)
	if aerr != nil {
		fmt.Println("Invalid amount")
		return
	}

	fmt.Println("coin", coin, "to", to, "amount", amount)

	// GO-0 : check coin type
	if coin != "BTC" && coin != "ETH" {
		fmt.Println("Invalid coin. Should be BTC or ETH")

		return
	}

	// GO-1 : load database
	dbScan.Start()

	// G0-2 : split ETH vs BTC
	arr_deposit := make([]dbScan.Deposit, 0)
	for _, de := range dbScan.HMAP_DEPOSIT {

		if coin == de.Coin {
			arr_de = append(arr_deposit, *de)
		}
	}

	// G0-3 : split array to send
	arr_withdraw := make([]dbScan.Deposit, 0)
	max := len(arr_deposit)
	total := float64(0)

	for i := 0; i < max; i++ {

		de := arr_deposit[i]

		balance := dbScan.GetBalance(coin, de.AddressDeposit)

		de.Amount = strconv.FormatFloat(balance, 'f', -1, 64)

		total += balance

		if total >= amount {

			if i < max-1 {
				sub := balance - (total - amount)
				de.Amount = strconv.FormatFloat(sub, 'f', -1, 64)
			}
		}

		arr_withdraw = append(arr_withdraw, *de)

		if total >= amount {
			break
		}
	}

	if total < amount {
		fmt.Println("Error Total : ", total, "  Amount Withdraw : ", amount)
		return
	}
}

func sendCoin(coin, from, to string, balance float64) (tx string) {

	max := float64(0)

	switch coin {
	case "BTC":
		{
			fee := float64(0.0001)
			max = balance - fee //BTC
		}

	case "ETH":
		{
			wei := balance * math.Pow10(18)

			gas := float64(21000)

			gasPriceBigI, _ := eth.SuggestGasPrice()
			gasPriceWei := gasPriceBigI.Int64()
			gasPrice := float64(gasPriceWei)

			max = wei - gas*gasPrice
			max /= math.Pow10(18) //ETH
		}
	}

	fmt.Println("withdraw input : ", balance, "withdraw max : ", max)

	amount := strconv.FormatFloat(max, 'f', -1, 64)

	tx = dbScan.SendCoin(coin, from, to, amount)
	fmt.Println("Tx : ", tx)

	return
}
