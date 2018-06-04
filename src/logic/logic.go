package logic

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"config"
	"lib/btc"
	"lib/eth"
	"module/dbScan"

	"github.com/PuerkitoBio/goquery"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

func verityTX(coin, tx string) (interface{}, bool) {

	var receipt interface{}

	switch coin {
	case "BTC":
		{
			result, err := btc.GetTransaction(tx)
			if err != nil {
				return nil, false
			}

			if result.BlockHash == "" {
				return nil, false
			}

			data := make(map[string]interface{})
			data["token"] = "BTC"
			data["transaction_id"] = result.TxID
			data["transaction_fee"] = result.Fee

			for _, obj := range result.Details {
				if obj.Category == "send" {
					data["from_address"] = obj.Address
				} else if obj.Category == "receive" {
					data["to_address"] = obj.Address
					data["amount"] = obj.Amount
				}
			}

			receipt = data
			fmt.Println("receipt : ", receipt)
		}

	case "ETH":
		{
			//receipt := eth.GetTransactionReceipt(tx)
			parsed := eth.GetTransactionByHash(tx)
			if parsed["result"] == nil {
				return nil, false
			}
			result := parsed["result"].(map[string]interface{})

			blockHash := result["blockHash"].(string)
			blockNum, _ := strconv.ParseInt(blockHash, 0, 64)

			if blockNum == int64(0) {
				return nil, false
			}

			data := make(map[string]interface{})
			data["token"] = "ETH"
			data["transaction_id"] = result["hash"]
			data["from_address"] = result["from"]
			data["to_address"] = result["to"]

			value, _ := strconv.ParseInt(result["value"].(string), 0, 64)
			data["amount"] = float64(value) / math.Pow10(18)

			gas, _ := strconv.ParseInt(result["gas"].(string), 0, 64)
			gasPrice, _ := strconv.ParseInt(result["gasPrice"].(string), 0, 64)
			fee := float64(gas*gasPrice) / math.Pow10(18)
			data["transaction_fee"] = fee

			receipt = data
			fmt.Println("receipt : ", receipt)

		}

	}

	return receipt, true
}

func getBalance(coin string, addr string) float64 {

	return dbScan.GetBalance(coin, addr)
}

func getBalanceOf(ercAddr, toAddr string) int64 {
	return dbScan.GetBalanceOf(ercAddr, toAddr)
}

func genAddress(coin string) (address string) {

	switch coin {
	case "BTC":
		address = genAddressBTC()
	case "ETH", "ERC20":
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

func getRatingFromEtherScan(addr, net string) float64 {

	if net != "main" {
		return float64(1)
	}

	res, err := http.Get("https://etherscan.io/token/" + addr)
	if err != nil {
		fmt.Println(err)
		return float64(0)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		fmt.Println("status code error: %d %s", res.StatusCode, res.Status)
		return float64(0)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return float64(0)
	}

	rating := "0"
	doc.Find(".container #ContentPlaceHolder1_divSummary #ContentPlaceHolder1_tr_valuepertoken").Each(func(i int, s *goquery.Selection) {

		data := s.Find("td").Text()
		fmt.Println("Html td : ", data)

		in := string(data)
		r := csv.NewReader(strings.NewReader(in))

		records, err := r.ReadAll()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Html records : ", records)

			record := records[1][0]
			fmt.Println("Html record", record)

			rating = strings.Split(record, " ")[2]
			fmt.Println("Html rating", rating)
		}
	})

	f, _ := strconv.ParseFloat(rating, 64)
	fmt.Println(f, "Eth")

	return f
}

func sendCoin(coin, from, to, amount string) (tx string) {

	switch coin {
	case "BTC":
		{
			aMountBTC := btc.ToBTC(amount)
			//satoshi := btc.ToSatoshi(amount)

			fee := float64(0.001)
			fund := aMountBTC + fee

			if getBalance(coin, from) < fund {
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

			balanceETH := getBalance(coin, from)
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

func sendERC20(contract, receiver, amount string) (tx string) {

	// Go : get balance of sender
	amETH := getBalance("ETH", config.ETH_ADDR)
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
	balance := getBalanceOf(contract, config.ETH_ADDR)

	if balance < token {
		fmt.Println("ERC20 Token not enough !!!", balance, token)
		return ""
	}

	tx = eth.SolidityTransactRaw(config.ETH_PRIV, contract, `transfer(address,uint256)`, nil, receiver, big.NewInt(token))
	fmt.Println("tx ERC20 : ", tx)

	return tx
}
