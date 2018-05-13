package dbScan

import (
	"fmt"
	"strconv"
	"time"
)

type Deposit struct {
	AddressDeposit  string `json:"deposit"`
	AddressContract string `json:"contract"`
	//Status          string `json:"status"`
	Coin   string `json:"coin"`
	Amount string `json:"amount"`
}

type Receipt struct {
	AddressTo       string `json:"addr"`
	AddressReceiver string `json:"receiver"`
	Amount          string `json:"amount"`
	Tx              string `json:"tx"`
	Ok              string `json:"ok"`
}

var STATUS_WAITING = "waiting"
var STATUS_PENDING = "pending"
var STATUS_SUCCESS = "success"

func NewDepositCoin(deposit, coin string) *Deposit {

	de := &Deposit{}
	de.AddressDeposit = deposit
	//de.Status = STATUS_WAITING
	de.Coin = coin

	go de.run()

	return de
}

func NewDepositERC20(depositAddr, contractAddr, coin string) *Deposit {

	de := &Deposit{}
	de.AddressDeposit = depositAddr
	de.AddressContract = contractAddr
	de.Coin = coin

	go de.run()

	return de
}

func (de *Deposit) run() {

	for {

		de.waiting()

		time.Sleep(10 * time.Second)
	}
}

func (de *Deposit) notify(balance float64) {

}

func (de *Deposit) checkBalance() (balance float64) {

	// GO-0 : check balance
	fmt.Println("GO-0 : check balance", de)

	cCoin := de.Coin
	if cCoin == "ERC20" {
		cCoin = "ETH"
	}
	balance = GetBalance(cCoin, de.AddressDeposit)

	return
}

func (de *Deposit) waiting() {

	fmt.Println("................waiting...................")

	// Go-0 : getBalance
	balance := de.checkBalance()

	// GO-1 : update new balance for deposit address
	amount := strconv.FormatFloat(balance, 'f', -1, 64)
	de.Amount = amount

	// Go-2 : Notify
	if amount != de.Amount {
		de.notify(balance)
	}

}
