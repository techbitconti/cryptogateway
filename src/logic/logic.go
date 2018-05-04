package logic

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"time"

	"config"
	"lib/btc"
	"lib/eth"
	"module/dbScan"

	"github.com/PuerkitoBio/goquery"
	"github.com/ethereum/go-ethereum/common"
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

	return dbScan.getBalance(coin, addr)
}

func getBalanceOf(ercAddr, toAddr string) int64 {
	return dbScan.getBalanceOf(ercAddr, toAddr)
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
			return float64(0)
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

func sendCoin(coin, from string, obj map[string]string) (tx string) {

	to := obj["addr"]
	amount := obj["amount"]
	//receiver := obj["receiver"]

	switch coin {
	case "BTC":
		{
			aMountBTC := btc.ToBTC(amount)
			satoshi := btc.ToSatoshi(amount)

			if getBalance(coin, from) < aMountBTC {
				fmt.Println("BTC not enough !!!")
				return ""
			}

			//btc.WalletPassphrase("123456", 10)
			txHash, err := btc.SendFrom(from, to, satoshi)
			if err != nil {
				return ""
			}
			tx = txHash.String()
			fmt.Println("tx BTC : ", tx)
		}

	case "ETH":
		{
			amETH, _ := strconv.ParseFloat(amount, 64)
			wei := amETH * math.Pow10(18)

			if getBalance(coin, from) < float64(wei) {
				fmt.Println("ETH not enough !!!")
				return ""
			}

			//tx = eth.SendTransactionRaw(config.ETH_SIM.PrivKey, to, aMountETH, []byte{})

			weiBig := big.NewInt(int64(wei))
			b := hexutil.Big(weiBig)
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
	}

	return tx
}

func sendERC20(contract, to, amount string) (tx string) {

	if getBalance(coin, config.ETH_SIM.Address) < float64(21000) {
		fmt.Println("ERC20 not enough !!!")
		return ""
	}

	tokens, _ := strconv.ParseInt(amount, 0, 64)
	balance := getBalanceOf(contract, config.ETH_SIM.Address)

	if balance < tokens {
		fmt.Println("ERC20 Token not enough !!!")
		return ""
	}

	tx = eth.SolidityTransactRaw(config.ETH_SIM.PrivKey, contract, `transfer(address,uint256)`, nil, common.HexToAddress(to), tokens)
	fmt.Println("tx ERC20 : ", tx)
}
