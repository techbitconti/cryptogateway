package dbScan

import (
	"fmt"
	"lib/btc"
	"lib/eth"
	"math/big"
	"strconv"
	"time"
)

var HMAP_DEPOSIT = map[string]map[string]string{}

/*
{
  "depositA" :
  {
    "status" : "waiting - peding - success",
    "coin" : "BTC/ETH/ERC20",
    "amount" : "0",
  }
}
*/

func initial() {

	//ETH
	HMAP_DEPOSIT["0x1760b401d3bf1b092e625aa6e5914c2123646923"] = map[string]string{

		"status": "waiting",
		"coin":   "ETH",
		"amount": "0",
	}

	//ERC20
	HMAP_DEPOSIT["0x53fd925caa530338cb2e21db445b2b0543ec09ce"] = map[string]string{

		"status": "waiting",
		"coin":   "ERC20",
		"amount": "0",
	}

	//BTC
	HMAP_DEPOSIT["2N4X6cHZ7My19oHFWiod2ftYr3dAGSUtTHC"] = map[string]string{

		"status": "waiting",
		"coin":   "BTC",
		"amount": "0",
	}

}

func Start() {

	//initial()

	go func() {
		for {
			update()
			time.Sleep(5 * time.Second)
		}
	}()
}

func update() {

	for k, v := range HMAP_DEPOSIT {
		switch v["status"] {
		case "waiting":
			go waiting(k, v)
		}
	}
}

func getBalance(coin, addr string) float64 {

	switch coin {
	case "BTC":
		amount := btc.GetBalance(addr)
		fmt.Println("getBalance BTC : ", amount)
		return amount

	case "ETH":
		bigInt := eth.GetBalance(addr)
		bigFloat := new(big.Float).SetInt(bigInt)
		f, _ := bigFloat.Float64()

		fmt.Println("getBalance ETH : ", f)
		return f
	}

	return 0
}

func waiting(addr_deposit string, data map[string]string) {

	fmt.Println("................waiting...................")

	// GO-0 : check balance
	fmt.Println("GO-0 : check balance")

	balance := getBalance(data["coin"], addr_deposit)
	if balance <= float64(0) {
		return
	}

	// GO-4 : status pending
	fmt.Println("GO-4 : status pending")

	HMAP_DEPOSIT[addr_deposit]["status"] = "pending"
	HMAP_DEPOSIT[addr_deposit]["amount"] = strconv.FormatFloat(balance, 'f', -1, 64)

	fmt.Println("HMAP_DEPOSIT : ", HMAP_DEPOSIT[addr_deposit])
}
