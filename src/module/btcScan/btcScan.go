package btcScan

import (
	"bytes"
	"encoding/json"
	"lib/btc"
	"time"

	"module/dbScan"

	"github.com/btcsuite/btcd/wire"
)

var MapBlock = map[string]interface{}{}

func Start() {

	update()
}

func update() {

	go func() {
		for {

			getBlock()

			time.Sleep(10000 * time.Millisecond)

		}
	}()
}

func getBlock() {

	blockHeight := btc.GetBlockCount()
	blockHash, err1 := btc.GetBlockHash(blockHeight)
	if err1 != nil {
		return
	}

	block, err2 := btc.GetBlock(blockHash)
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
		result, err := btc.GetTransaction(tx)
		if err != nil {
			return
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

				if de, ok := dbScan.HMAP_DEPOSIT[data["to_address"].(string)]; ok {
					de.Notify(data)

					dbScan.Report_Fees("BTC", result.Fee)
					dbScan.Report_Deposit("BTC", obj.Amount)
					dbScan.Report_Current("BTC")
				}
			}
		}
	}

}
