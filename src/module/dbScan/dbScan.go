package dbScan

import (
	"fmt"
	"lib/btc"
	"lib/eth"
	"math"
	"math/big"

	"config"

	"github.com/ethereum/go-ethereum/common"
)

var HMAP_DEPOSIT = map[string]*Deposit{}

func getBalance(coin, addr string) float64 {

	switch coin {
	case "BTC":
		amount := btc.GetBalance(addr)
		fmt.Println("getBalance BTC : ", amount)
		return amount

	case "ETH":
		bigInt := eth.GetBalance(addr)
		bigFloat := new(big.Float).SetInt(bigInt)
		wei, _ := bigFloat.Float64()

		ether := wei / math.Pow10(18)

		fmt.Println("getBalance ETH : ", ether, "wei", wei)
		return ether
	}

	return 0
}

func getBalanceOf(ercAddr, toAddr string) int64 {

	hex, err := eth.SolidityCallRaw(config.ETH_SIM.Address, ercAddr, `balanceOf(address)`, toAddr)
	if err != nil {
		fmt.Println("ERC20 Token not enough !!!")
		return int64(0)
	}
	//fmt.Println("balanceOf", common.BytesToHash(hex).Big(), new(big.Int).SetBytes(hex))

	return common.BytesToHash(hex).Big().Int64()
}
