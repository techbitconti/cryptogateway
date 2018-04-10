package etherScan

import (
	"encoding/json"
	"fmt"
	"lib/eth"
	"math/big"
	"strconv"
	"time"

	"db/redis"
)

var depth int64 = 5

var Pool_Queuing map[string]interface{}
var Pool_Pending []map[string]interface{}
var Pool_TX map[string]interface{}

func Start() {

	Pool_Queuing = make(map[string]interface{})
	Pool_Pending = make([]map[string]interface{}, 0)
	Pool_TX = make(map[string]interface{})

	update()
}

func update() {

	go func() {
		for {
			getBlock()
			time.Sleep(1 * time.Millisecond)
		}
	}()
}

func getBlock() {

	numHash := eth.GetBlockNumber()

	numInt, _ := strconv.ParseInt(numHash, 0, 64)
	numBig := big.NewInt(numInt)

	//	if _, ok := Pool_Queuing[numHash]; ok {
	//		return
	//	}

	val := redis.Session.Get(numHash).Val()
	if val != "" {
		return
	}

	block := eth.GetBlockByNumber(numHash)
	bBytes, _ := json.MarshalIndent(block, " ", "")

	fmt.Println("..............................................................")
	fmt.Println("blockNumber: ", numBig, numHash)
	fmt.Println(string(bBytes))
	fmt.Println("..............................................................")

	//Pool_Queuing[numHash] = block
	Pool_Pending = append(Pool_Pending, block)
	redis.Session.Set(numHash, bBytes, 0)

	parse(block)
}

func parse(block map[string]interface{}) {

	tranB, _ := json.Marshal(block["transactions"])
	var trans []map[string]interface{}
	json.Unmarshal(tranB, &trans)

	for _, txObj := range trans {

		tx := txObj["hash"].(string)

		//RECEIPT
		result := eth.GetTransactionReceipt(tx)["result"].(map[string]interface{})

		//LOGS
		var logs []map[string]interface{}
		b, _ := json.Marshal(result["logs"])
		json.Unmarshal(b, &logs)

		if len(logs) > 0 {
			txObj["amount"] = logs[0]["data"].(string)
		}
		//

		Pool_TX[tx] = txObj

		txByte, _ := json.MarshalIndent(txObj, " ", "")
		redis.Session.Set(tx, txByte, 0)

		fmt.Println("TransactionS : ", tx, string(txByte))
	}
}
