package etherScan

import (
	"encoding/json"
	"fmt"
	"lib/eth"
	//"math/big"
	//"strconv"
	"time"
)

func Start() {

	update()
}

func update() {

	go func() {
		for {

			getBlock()

			time.Sleep(1 * time.Second)
		}
	}()
}

func getTrans(number interface{}) {

	trans := eth.GetTransactions(number)

	fmt.Println("Transaction : ", trans)

}

func getBlock() {

	numHash := eth.GetBlockNumber()

	/*
		numInt, _ := strconv.ParseInt(numHash, 0, 64)
		numBig := big.NewInt(numInt)

		block := eth.GetBlockByNumber(numHash)
		bBytes, _ := json.MarshalIndent(block, " ", "")

		fmt.Println("..............................................................")
		fmt.Println("blockNumber: ", numBig, numHash)
		fmt.Println(string(bBytes))
		fmt.Println("..............................................................")

		parse(block)

	*/

	getTrans(numHash)
}

func parse(block map[string]interface{}) {

	tranB, _ := json.Marshal(block["transactions"])
	var trans []map[string]interface{}
	json.Unmarshal(tranB, &trans)

	for _, txObj := range trans {

		tx := txObj["hash"].(string)

		//RECEIPT
		receipt := eth.GetTransactionReceipt(tx)
		if receipt["result"] == nil {
			return
		}
		result := receipt["result"].(map[string]interface{})

		//LOGS
		var logs []map[string]interface{}
		b, _ := json.Marshal(result["logs"])
		json.Unmarshal(b, &logs)

		if len(logs) > 0 {
			txObj["amount"] = logs[0]["data"].(string)
		}
	}
}
