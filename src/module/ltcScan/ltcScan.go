package ltcScan

import (
	"bytes"
	"encoding/json"
	"lib/ltc"
	"time"

	"module/dbScan"

	"github.com/ltcsuite/ltcd/wire"
)

var MapBlock = map[string]interface{}{}

func Start() {

	update()
}

func update() {

	go func() {
		for {

			time.Sleep(10 * time.Second)

			getBlock()
		}
	}()
}

func getBlock() {

	blockHeight := ltc.GetBlockCount()
	blockHash, err1 := ltc.GetBlockHash(blockHeight)
	if err1 != nil {
		return
	}

	block, err2 := ltc.GetBlock(blockHash)
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
		result, err := ltc.GetTransaction(tx)
		if err != nil {
			return
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

				if de, ok := dbScan.HMAP_DEPOSIT[data["to_address"].(string)]; ok {
					de.Notify(data)

					dbScan.Report_Fees("LTC", result.Fee)
					dbScan.Report_Deposit("LTC", obj.Amount)
					dbScan.Report_Current("LTC")
				}
			}
		}
	}

}
