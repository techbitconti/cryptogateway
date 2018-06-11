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

	if len(os.Args) != 4 {

		fmt.Println("not enough length : ", len(os.Args))

		return
	}

	coin := os.Args[1]
	to := os.Args[2]
	amountStr := os.Args[3]
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
			arr_deposit = append(arr_deposit, *de)
		}
	}

	fmt.Println("arr_deposit : ", arr_deposit)

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

			if i == max-1 {
				sub := balance - (total - amount)
				de.Amount = strconv.FormatFloat(sub, 'f', -1, 64)
			}
		}

		fmt.Println(de.AddressDeposit, de.Amount)

		arr_withdraw = append(arr_withdraw, de)

		if total >= amount {
			break
		}
	}

	fmt.Println("arr_withdraw : ", arr_withdraw)

	// G0-4 : check total
	if total < amount {
		fmt.Println(".........Failed..........")
		fmt.Println("Total : ", total, "  Amount Withdraw : ", amount)
		return
	}

	// GO-5 : Send Withdraw
	count := 0
	list := make([]string, 0)
	for _, de := range arr_withdraw {

		tx := sendCoin(de.Coin, de.AddressDeposit, to, de.Amount)

		if tx != "" {
			count++
			list = append(list, tx)
		}
	}

	if count >= len(arr_withdraw) {
		fmt.Println(".........Success..........")
		fmt.Println(list)
	}

}

func sendCoin(coin, from, to, amount string) (tx string) {

	max := float64(0)
	balance, _ := strconv.ParseFloat(amount, 64)

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

	value := strconv.FormatFloat(max, 'f', -1, 64)

	tx = dbScan.SendCoin(coin, from, to, value)
	fmt.Println("Tx : ", tx)

	return
}
