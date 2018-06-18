package dbScan

import (
	"fmt"
	"lib/btc"
	"lib/eth"
	"math"
	"math/big"
	"strconv"

	"config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var HMAP_DEPOSIT = map[string]*Deposit{}

func Start() {

	LoadDeposit()
	LoadReport_BTC()
	LoadReport_ETH()
}

func GetBalance(coin, addr string) float64 {

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

		fmt.Println("getBalance :", ether, "ether", " -- wei", wei)
		return ether
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
				fmt.Println("BTC not enough !!!")
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

	case "ETH":
		{
			amETH, _ := strconv.ParseFloat(amount, 64)
			valueWei := amETH * math.Pow10(18)

			gasWei := config.ETH_GAS
			gasPriceBigI, _ := eth.SuggestGasPrice()
			gasPriceWei := gasPriceBigI.Int64()

			// Cost returns amount + gasprice * gaslimit.
			fundWei := gasWei*float64(gasPriceWei) + valueWei
			fundETH := fundWei / math.Pow10(18)

			fmt.Println("gasWei : ", gasWei, "gasPriceWei : ", gasPriceWei, "fundWei : ", fundWei)

			balanceETH := GetBalance(coin, from)
			fmt.Println("balanceETH : ", balanceETH, "fundETH : ", fundETH)
			if balanceETH < fundETH {
				fmt.Println("ETH not enough balance !!!", balanceETH, "funds : ", fundETH)
				return ""
			}

			weiBig := big.NewInt(int64(valueWei))
			b1 := hexutil.Big(*weiBig)
			valueAM := b1.String()

			gasBig := big.NewInt(int64(gasWei))
			b2 := hexutil.Big(*gasBig)
			valueGAS := b2.String()

			gasPrBig := big.NewInt(int64(gasPriceWei))
			b3 := hexutil.Big(*gasPrBig)
			valueGASPr := b3.String()

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
	}

	return tx
}

func SendERC20(contract, receiver, amount string) (tx string) {

	// Go : get balance of sender
	amETH := GetBalance("ETH", config.ETH_ADDR)
	valueWei := int64(amETH * math.Pow10(18))

	// GO : convert
	//weiBigI := big.NewInt(valueWei)
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

func Report_Deposit(coin string, de float64) {

	switch coin {
	case "BTC":
		{
			BTC_DEPOSIT += de
			SaveReport_BTC_Deposit(BTC_DEPOSIT)
		}
	case "ETH":
		{
			ETH_DEPOSIT += de
			SaveReport_ETH_Deposit(ETH_DEPOSIT)
		}
	}

}

func Report_Withdraw(coin string, with float64) {

	switch coin {
	case "BTC":
		{
			BTC_WITHDRAW += with
			SaveReport_BTC_Withdraw(BTC_WITHDRAW)
		}
	case "ETH":
		{
			ETH_WITHDRAW += with
			SaveReport_ETH_Withdraw(ETH_WITHDRAW)
		}
	}

}

func Report_Fees(coin string, fees float64) {

	switch coin {
	case "BTC":
		{
			BTC_FEES += math.Abs(fees)
			SaveReport_BTC_Fees(BTC_FEES)
		}
	case "ETH":
		{
			ETH_FEES += fees
			SaveReport_ETH_Fees(ETH_FEES)
		}
	}

}

func Report_Current(coin string) {

	switch coin {
	case "BTC":
		{
			BTC_CURRENT = BTC_DEPOSIT - BTC_WITHDRAW
			SaveReport_BTC_Current(BTC_CURRENT)
		}
	case "ETH":
		{
			ETH_CURRENT = ETH_DEPOSIT - ETH_WITHDRAW
			SaveReport_ETH_Current(ETH_CURRENT)
		}
	}

}
