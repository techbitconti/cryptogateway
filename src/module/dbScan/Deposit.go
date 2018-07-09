package dbScan

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"config"
	"lib/cryptocy"
)

type Deposit struct {
	Seed            string //`json:"seed"`
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

func NewDepositCoin(seed, deposit, coin string) *Deposit {

	de := &Deposit{}
	de.AddressDeposit = deposit
	//de.Status = STATUS_WAITING
	de.Coin = coin
	de.Encrypt(seed)

	//DB
	SaveDeposit(*de)

	go de.run()

	return de
}

func NewDepositERC20(depositAddr, contractAddr, coin string) *Deposit {

	de := &Deposit{}
	de.AddressDeposit = depositAddr
	de.AddressContract = contractAddr
	de.Coin = coin

	//DB
	SaveDeposit(*de)

	go de.run()

	return de
}

func (de *Deposit) run() {

	for {

		de.waiting()

		time.Sleep(5 * time.Second)
	}
}

func (de *Deposit) Notify(data map[string]interface{}) {

	addr_deposit := data["to_address"].(string)

	if de.Coin == "ETH" || de.Coin == "ERC20" {
		addr_deposit = strings.ToLower(addr_deposit)
	}

	fmt.Println("notify deposit : ", addr_deposit, de.AddressDeposit)

	if addr_deposit != de.AddressDeposit {
		return
	}

	if config.IP_ALLOW == "" || config.PORT_ALLOW == "" || config.NOTIFY_BALANCE == "" {
		return
	}

	url := "http://" + config.IP_ALLOW + ":" + config.PORT_ALLOW + config.NOTIFY_BALANCE
	b, _ := json.Marshal(data)

	res, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(".......notify........", url, data, string(body))
}

func (de *Deposit) Encrypt(msg string) {

	CIPHER_KEY := cryptocy.GenKey([]byte(de.AddressDeposit))

	encrypted, err := cryptocy.Encrypt(CIPHER_KEY, msg)
	if err != nil {
		de.Seed = msg
	} else {
		de.Seed = encrypted
	}
}

func (de *Deposit) Decrypt() string {

	CIPHER_KEY := cryptocy.GenKey([]byte(de.AddressDeposit))
	encrypted := de.Seed

	decrypted, err := cryptocy.Decrypt(CIPHER_KEY, encrypted)
	if err != nil {
		return ""
	}

	return decrypted
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
	if amount != de.Amount {
		// Go-2 : Notify
		SaveDeposit(*de)
	}
	de.Amount = amount
}
