package dbScan

import (
	"fmt"
	"lib/bch"
	"lib/btc"
	"lib/eth"
	"lib/ltc"
	"lib/xlm"
	"math"
	"math/big"
	"strconv"

	"config"

	"github.com/ethereum/go-ethereum/common"
)

var HMAP_DEPOSIT = map[string]*Deposit{}

func Start() {

	LoadDeposit()
	LoadReport_ETH()
	LoadReport_BTC()
	LoadReport_LTC()
	LoadReport_BCH()
	LoadReport_XLM()
}

func GetBalance(coin, addr string) float64 {

	switch coin {
	case "BTC":
		{
			amount := btc.GetBalance(addr)
			fmt.Println("getBalance BTC : ", amount)
			return amount
		}

	case "BCH":
		{
			amount := bch.GetBalance(addr)
			fmt.Println("getBalance BCH : ", amount)
			return amount
		}

	case "LTC":
		{
			amount := ltc.GetBalance(addr)
			fmt.Println("getBalance LTC : ", amount)
			return amount
		}

	case "ETH":
		{
			bigInt := eth.GetBalance(addr)
			bigFloat := new(big.Float).SetInt(bigInt)
			wei, _ := bigFloat.Float64()

			ether := wei / math.Pow10(18)

			fmt.Println("getBalance :", ether, "ether", " -- wei", wei)
			return ether
		}

	case "XLM":
		{
			return xlm.GetBalance(config.XLM_NET, addr, "native", "")
		}
	}

	return 0
}

func GetBalanceOf(ercAddr, toAddr string) int64 {

	hex, err := eth.SolidityCallRaw(config.ETH_ADDR, ercAddr, `balanceOf(address)`, toAddr)
	if err != nil {
		fmt.Println("ERC20 Token not enough !!!")
		return int64(0)
	}
	//fmt.Println("balanceOf", common.BytesToHash(hex).Big(), new(big.Int).SetBytes(hex))

	return common.BytesToHash(hex).Big().Int64()
}

func SendCoin(coin, from, to, amount string) (tx string) {

	switch coin {
	case "BTC":
		{
			aMountBTC := btc.ToBTC(amount)
			//satoshi := btc.ToSatoshi(amount)

			fund := aMountBTC + config.BTC_FEE

			if GetBalance(coin, from) < fund {
				fmt.Println("BTC not enough !!!", fund)
				return ""
			}

			//btc.WalletPassphrase("123456", 10)
			//txHash, err := btc.SendFrom(from, to, satoshi)
			txHash, err := btc.SendFrom(from, to, aMountBTC)

			if err != nil {
				return ""
			}
			tx = txHash.String()
			fmt.Println("tx BTC : ", tx)
		}

	case "BCH":
		{
			aMountBCH := bch.ToBCH(amount)
			//satoshi := bch.ToSatoshi(amount)

			fund := aMountBCH + config.BTC_FEE

			if GetBalance(coin, from) < fund {
				fmt.Println("BCH not enough !!!", fund)
				return ""
			}

			//bch.WalletPassphrase("123456", 10)
			//txHash, err := bch.SendFrom(from, to, satoshi)
			txHash, err := bch.SendFrom(from, to, aMountBCH)

			if err != nil {
				return ""
			}
			tx = txHash.String()
			fmt.Println("tx BTC : ", tx)
		}

	case "LTC":
		{
			aMountLTC := ltc.ToBTC(amount)
			//satoshi := ltc.ToSatoshi(amount)

			fund := aMountLTC + config.BTC_FEE

			if GetBalance(coin, from) < fund {
				fmt.Println("BTC not enough !!!", fund)
				return ""
			}

			//ltc.WalletPassphrase("123456", 10)
			//txHash, err := ltc.SendFrom(from, to, satoshi)
			txHash, err := ltc.SendFrom(from, to, aMountLTC)

			if err != nil {
				return ""
			}
			tx = txHash.String()
			fmt.Println("tx LTC : ", tx)
		}

	case "ETH":
		{
			amETH, _ := strconv.ParseFloat(amount, 64)
			valueWei := eth.ToWei(amETH, "ether") // amETH * math.Pow10(18)

			gasWei := config.ETH_GAS
			gasPriceBigI, _ := eth.SuggestGasPrice()
			gasPriceWei := gasPriceBigI.Int64()

			fundWei := gasWei*float64(gasPriceWei) + valueWei
			fundETH := eth.FromWei(fundWei, "ether") // fundWei / math.Pow10(18)

			fmt.Println("gasWei : ", gasWei, "gasPriceWei : ", gasPriceWei, "fundWei : ", fundWei)

			balanceETH := GetBalance(coin, from)
			fmt.Println("balanceETH : ", balanceETH, "fundETH : ", fundETH)
			if balanceETH < fundETH {
				fmt.Println("ETH not enough balance !!!", balanceETH, "funds : ", fundETH)
				return ""
			}

			valueAM := eth.ToBigNumber(uint64(valueWei))
			valueGAS := eth.ToBigNumber(uint64(gasWei))
			valueGASPr := eth.ToBigNumber(uint64(gasPriceWei))

			eth.UnlockAccount(from, "123456", uint64(10))
			msg := map[string]interface{}{
				"from":     from,
				"to":       to,
				"value":    valueAM,
				"gas":      valueGAS,
				"gasPrice": valueGASPr,
			}
			tx = eth.SendTransaction(msg)

			fmt.Println("tx ETH : ", tx)
		}

	case "XLM":
		{
			kp, errKp := xlm.KeyPairParse(from)
			if errKp != nil {
				return ""
			}

			fund := xlm.ToLumens(amount) + config.XLM_FEE

			if GetBalance(coin, kp.Address()) < fund {
				fmt.Println("XLM not enough !!!", strconv.FormatFloat(fund, 'f', -1, 64))
				return ""
			}

			txBuilder, err := xlm.TxBuilder(config.XLM_NET, from, to, amount)
			if err != nil {
				fmt.Println("XLM send fail !!!")
				return ""
			}

			txeB64 := xlm.TxSign(txBuilder, from)
			tx = xlm.TxSubmit("test", txeB64)
		}

	}

	return tx
}

func SendERC20(contract, receiver, amount string) (tx string) {

	// Go : get balance of sender
	amETH := GetBalance("ETH", config.ETH_ADDR)
	valueWei := int64(eth.ToWei(amETH, "ether")) // int64(amETH * math.Pow10(18))

	// GO : convert
	amountBigI, _ := strconv.ParseInt(amount, 0, 64)

	//GO : get bytecode of contract function
	byteCode := eth.GetByteCode(contract, "transfer", big.NewInt(amountBigI))
	gasUsed, _ := eth.EstimateGas(contract, nil, byteCode)
	fmt.Println("gasUsed : ", gasUsed, "byteCode : ", byteCode)

	gasWei := int64(gasUsed)
	gasPriceBigI, _ := eth.SuggestGasPrice()
	gasPriceWei := gasPriceBigI.Int64()
	// Cost returns amount + gasprice * gaslimit.
	fundWei := gasWei*gasPriceWei + valueWei

	fmt.Println("gasWei : ", gasWei, "gasPriceWei : ", gasPriceWei, "fundWei : ", fundWei)

	if valueWei < fundWei {
		fmt.Println("ERC20 not enough gas!!!")
		return ""
	}

	token, _ := strconv.ParseInt(amount, 0, 64)
	balance := GetBalanceOf(contract, config.ETH_ADDR)

	if balance < token {
		fmt.Println("ERC20 Token not enough !!!", balance, token)
		return ""
	}

	tx = eth.SolidityTransactRaw(config.ETH_PRIV, contract, `transfer(address,uint256)`, nil, receiver, big.NewInt(token))
	fmt.Println("tx ERC20 : ", tx)

	return tx
}
