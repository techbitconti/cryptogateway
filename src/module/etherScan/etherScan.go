package etherScan

import (
	"bytes"
	"encoding/json"
	//"fmt"
	"lib/eth"
	"math"
	"strconv"
	"time"

	"module/dbScan"
)

var MapBlock = map[string]interface{}{}

func Start() {

	update()
}

func update() {

	go func() {
		for {

			time.Sleep(1000 * time.Millisecond)

			getBlock()
		}
	}()
}

func getBlock() {

	numHash := eth.GetBlockNumber()
	//numInt, _ := strconv.ParseInt(numHash, 0, 64)

	block := eth.GetBlockByNumber(numHash)
	//bBytes, _ := json.MarshalIndent(block, " ", "")

	//	fmt.Println("..............................................................")
	//	fmt.Println("blockNumber: ", numHash, numInt)
	//	fmt.Println(string(bBytes))
	//	fmt.Println("..............................................................")

	blockOld, exist := MapBlock[numHash]
	bytesOld, _ := json.Marshal(blockOld)
	bytesNew, _ := json.Marshal(block)

	if !exist || !bytes.Equal(bytesOld, bytesNew) {

		MapBlock[numHash] = block

		parse(block)
	}
}

func parse(block map[string]interface{}) {

	var trans []map[string]interface{}

	tranB, _ := json.Marshal(block["transactions"])
	json.Unmarshal(tranB, &trans)

	for _, txObj := range trans {

		data := make(map[string]interface{})
		data["token"] = "ETH"
		data["transaction_id"] = txObj["hash"]
		data["from_address"] = txObj["from"]
		data["to_address"] = txObj["to"]

		value, _ := strconv.ParseInt(txObj["value"].(string), 0, 64)
		data["amount"] = float64(value) / math.Pow10(18)

		gas, _ := strconv.ParseInt(txObj["gas"].(string), 0, 64)
		gasPrice, _ := strconv.ParseInt(txObj["gasPrice"].(string), 0, 64)
		fee := float64(gas*gasPrice) / math.Pow10(18)
		data["transaction_fee"] = fee

		if txObj["to"] != nil {
			if de, ok := dbScan.HMAP_DEPOSIT[txObj["to"].(string)]; ok {
				de.Notify(data)

				dbScan.Report_Fees("ETH", fee)
				dbScan.Report_Deposit("ETH", data["amount"].(float64))
				dbScan.Report_Current("ETH")
			}
		}

	}
}
