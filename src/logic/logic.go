package logic

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"config"
	"lib/btc"
	"lib/eth"
)

func verifyAddress(coin string, addr string) bool {

	switch coin {
	case "BTC":
		_, err := btc.ValidateAddress(addr)
		if err != nil {
			return false
		}

	case "ETH":
		if !eth.IsHexAddress(addr) {
			return false
		}

	case "ERC20":
		if !eth.IsHexContract(addr) {
			return false
		}
	}

	return true
}

func verityTX(coin, tx string) bool {

	switch coin {
	case "BTC":
		receipt, err := btc.GetRawTransactionVerbose(tx)
		if err != nil {
			return false
		}
		fmt.Println("receipt tx : ", tx, "  result : ", receipt)

	case "ETH", "ERC20":
		receipt := eth.GetTransactionReceipt(tx)
		if receipt["result"] == nil {
			return false
		}
		fmt.Println("receipt tx : ", tx, "  result : ", receipt)

	}

	return true
}

func getBalance(coin string, addr string) float64 {

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

	return float64(0)
}

func getAddressAdmin(coin string) string {
	switch coin {
	case "BTC":
		return config.BTC_SIM.Address
	case "ETH":
		return config.ETH_SIM.Address
	case "ERC20":
		return config.ERC20_SIM.Address
	}

	return ""
}

func genAddress(coin string) (address string) {

	switch coin {
	case "BTC":
		address = genAddressBTC()
	case "ETH":
		address = genAddressETH()
	}

	fmt.Println("genAddress --- ", address)

	return
}

func genAddressBTC() string {

	utc := time.Now().Unix()
	decode := strconv.FormatInt(utc, 10)
	encode := base64.StdEncoding.EncodeToString([]byte(decode))

	address, _ := btc.GetNewAddress(encode)

	return address.String()
}

func genAddressETH() string {

	address, _ := eth.NewAccount(config.PATH_ETH, "123456")

	return address
}

func sendTransaction(coin, from string, obj map[string]string) (tx string) {

	to := obj["addr"]
	amount := obj["amount"]
	receiver := obj["receiver"]

	aMountBTC := btc.ToBTC(amount)

	amETH, _ := strconv.ParseInt(amount, 0, 64)
	aMountETH := big.NewInt(amETH)

	switch coin {
	case "BTC":
		{
			//if getBalance(coin, from)+float64(0.001) < aMountBTC {
			if getBalance(coin, from) < aMountBTC {
				fmt.Println("BTC not enough !!!")
				return ""
			}

			//btc.WalletPassphrase("123456", 10)
			txHash, err := btc.SendFrom(from, to, aMountBTC)
			if err != nil {
				return ""
			}
			tx = txHash.String()
			fmt.Println("tx BTC : ", tx)
		}

	case "ETH":
		{
			//if getBalance(coin, from)+float64(21000) < float64(aMountETH.Int64()) {
			if getBalance(coin, from) < float64(aMountETH.Int64()) {
				fmt.Println("ETH not enough !!!")
				return ""
			}

			//tx = eth.SendTransactionRaw(config.ETH_SIM.PrivKey, to, aMountETH, []byte{})

			b := hexutil.Big(*aMountETH)
			value := b.String()

			eth.UnlockAccount(from, "123456", uint64(10))
			msg := map[string]interface{}{
				"from":  from,
				"to":    to,
				"value": value,
			}
			tx = eth.SendTransaction(msg)

			fmt.Println("tx ETH : ", tx)
		}

	case "ERC20":
		{
			if getBalance(coin, config.ETH_SIM.Address) < float64(21000) {
				fmt.Println("ERC20 not enough !!!")
				return ""
			}

			hex, err := eth.SolidityCallRaw(config.ETH_SIM.Address, to, `balanceOf(address)`, config.ETH_SIM.Address)
			if err != nil {
				fmt.Println("ERC20 Token not enough !!!")
				return ""
			}
			fmt.Println(common.BytesToHash(hex).Big(), new(big.Int).SetBytes(hex))

			tokens := aMountETH
			tx = eth.SolidityTransactRaw(config.ETH_SIM.PrivKey, to, `transfer(address,uint256)`, nil, common.HexToAddress(receiver), tokens)
			fmt.Println("tx ERC20 : ", tx)
		}

	}

	return tx
}
