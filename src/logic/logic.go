package logic

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"config"
	"lib/btc"
	"lib/eth"
	"lib/ltc"
	"module/dbScan"

	"github.com/PuerkitoBio/goquery"
)

func verifyAddress(coin string, addr string) bool {

	switch coin {
	case "BTC":
		_, err := btc.ValidateAddress(addr)
		if err != nil {
			return false
		}
	case "LTC":
		_, err := ltc.ValidateAddress(addr)
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

	case "LTC":
		{
			result, err := ltc.GetTransaction(tx)
			if err != nil {
				return nil, false
			}

			if result.BlockHash == "" {
				return nil, false
			}

			data := make(map[string]interface{})
			data["token"] = "LTC"
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

	case "ETH", "ERC20":
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
	case "LTC":
		address = genAddressLTC()
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

func genAddressLTC() string {

	utc := time.Now().Unix()
	decode := strconv.FormatInt(utc, 10)
	encode := base64.StdEncoding.EncodeToString([]byte(decode))

	address, _ := ltc.GetNewAddress(encode)

	return address.String()
}

func genAddressETH() string {

	keyHex, address, _ := eth.NewAccount()
	eth.StoreAccount(keyHex, "123456", config.PATH_ETH)

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

	return dbScan.SendCoin(coin, from, to, amount)
}

func sendERC20(contract, receiver, amount string) (tx string) {

	return dbScan.SendERC20(contract, receiver, amount)
}
