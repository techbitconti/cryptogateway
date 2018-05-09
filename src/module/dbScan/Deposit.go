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
	//de.Status = STATUS_WAITING
	de.Coin = coin

	go de.run()

	return de
}

func (de *Deposit) run() {

	for {

		//		switch de.Status {
		//		case STATUS_WAITING:
		//			de.waiting()
		//		}

		de.waiting()

		time.Sleep(10 * time.Second)
	}
}

func (de *Deposit) notify() {

}

func (de *Deposit) checkBalance() (balance float64) {

	// GO-0 : check balance
	fmt.Println("GO-0 : check balance", de)

	cCoin := ""
	if de.Coin == "ERC20" {
		cCoin = "ETH"
	}
	balance = GetBalance(cCoin, de.AddressDeposit)

	// Go-1 : Notify
	de.notify()

	return
}

func (de *Deposit) waiting() {

	fmt.Println("................waiting...................")

	balance := de.checkBalance()
	//	if balance <= float64(0) {
	//		de.Status = STATUS_WAITING
	//		return
	//	}

	// GO-1 : status pending
	//fmt.Println("GO-1 : status pending")
	//de.Status = STATUS_PENDING
	de.Amount = strconv.FormatFloat(balance, 'f', -1, 64)
}
