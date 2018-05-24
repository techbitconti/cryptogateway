package etherScan

import (
	"encoding/json"
	"fmt"
	"lib/eth"
	"strconv"
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

func getBlock() {

	numHash := eth.GetBlockNumber()
	numInt, _ := strconv.ParseInt(numHash, 0, 64)

	block := eth.GetBlockByNumber(numHash)
	bBytes, _ := json.MarshalIndent(block, " ", "")

	fmt.Println("..............................................................")
	fmt.Println("blockNumber: ", numHash, numInt)
	fmt.Println(string(bBytes))
	fmt.Println("..............................................................")

	parse(block)
}

func parse(block map[string]interface{}) {

	var trans []map[string]interface{}

	tranB, _ := json.Marshal(block["transactions"])
	json.Unmarshal(tranB, &trans)

	for _, txObj := range trans {

		tx := txObj["hash"].(string)
		b, _ := json.MarshalIndent(txObj, "", " ")
		fmt.Println(tx, string(b))

	}
}
