package bchScan

import (
	"bytes"
	"encoding/json"
	"lib/bch"
	"time"

	"module/dbScan"

	"github.com/bchsuite/bchd/wire"
)

var MapBlock = map[string]interface{}{}

func Start() {

	update()
}

func update() {

	go func() {
		for {

			getBlock()

			time.Sleep(10 * time.Second)

		}
	}()
}

func getBlock() {

	blockHeight := bch.GetBlockCount()
	blockHash, err1 := bch.GetBlockHash(blockHeight)
	if err1 != nil {
		return
	}

	block, err2 := bch.GetBlock(blockHash)
	if err2 != nil {
		return
	}

	blockOld, exist := MapBlock[blockHash.String()]
	bytesOld, _ := json.Marshal(blockOld)
	bytesNew, _ := json.Marshal(block)

	if !exist || !bytes.Equal(bytesOld, bytesNew) {

		MapBlock[blockHash.String()] = block

		parse(block)
	}

}

func parse(block *wire.MsgBlock) {

	hashList, _ := block.TxHashes()

	for _, txObj := range hashList {

		tx := txObj.String()
		result, err := bch.GetTransaction(tx)
		if err != nil {
			return
		}

		data := make(map[string]interface{})
		data["token"] = "BCH"
		data["transaction_id"] = result.TxID
		data["transaction_fee"] = result.Fee

		for _, obj := range result.Details {

			if obj.Category == "send" {
				data["from_address"] = obj.Address

			} else if obj.Category == "receive" {

				data["to_address"] = obj.Address
				data["amount"] = obj.Amount

				if de, ok := dbScan.HMAP_DEPOSIT[data["to_address"].(string)]; ok {
					de.Notify(data)

					dbScan.Report_Fees("BCH", result.Fee)
					dbScan.Report_Deposit("BCH", obj.Amount)
					dbScan.Report_Current("BCH")
				}
			}
		}
	}

}
