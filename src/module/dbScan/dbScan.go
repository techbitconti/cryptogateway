package dbScan

import (
	"encoding/json"
	"fmt"
	"lib/btc"
	"lib/eth"
	"math"
	"math/big"
	"strconv"

	"config"
	"db/redis"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var HMAP_DEPOSIT = map[string]*Deposit{}

var REDIS_KEYS_DEPOSIT = "keys_deposit"

func Start() {

	loadDeposit()
}

func loadDeposit() {

	bKeys, _ := redis.Session.Get(REDIS_KEYS_DEPOSIT).Bytes()

	var keys []string
	json.Unmarshal(bKeys, &keys)

	for _, key := range keys {

		deStr := redis.Session.Get(key).Val()

		if deStr != "" {

			b, _ := json.Marshal(deStr)
			var de Deposit
			json.Unmarshal(b, &de)

			HMAP_DEPOSIT[key] = &de

			fmt.Println("REDIS_KEYS_DEPOSIT: ", key, string(b))
		}
	}
}

func saveDeposit(de Deposit) {

	data, _ := json.Marshal(de)
	redis.Session.Set(de.AddressDeposit, data, 0).Err()

	arr, _ := redis.Session.Get(REDIS_KEYS_DEPOSIT).Result()
	var keys []string
	json.Unmarshal([]byte(arr), &keys)
	keys = append(keys, de.AddressDeposit)

	bKeys, _ := json.Marshal(keys)

	redis.Session.Set(REDIS_KEYS_DEPOSIT, bKeys, 0).Err()

	fmt.Println("saveDeposit : ", de.AddressDeposit)
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

			fee := float64(0.0001)
			fund := aMountBTC + fee

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

			gasWei := float64(21000)
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
