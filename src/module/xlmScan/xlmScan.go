package xlmScan

import (
	"bytes"
	"config"
	"encoding/json"
	//"fmt"
	"strconv"
	"time"

	"lib/xlm"
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

	net := config.XLM_NET

	block_embeded := xlm.LedgerAll(net, "", 1, xlm.ORDER_DESC)["_embedded"].(map[string]interface{})

	block_record, _ := json.Marshal(block_embeded["records"])
	var recordsBlock []map[string]interface{}
	json.Unmarshal(block_record, &recordsBlock)

	block_hash := recordsBlock[0]["hash"].(string)
	//fmt.Println("block hash : ", block_hash)

	blockOld, exist := MapBlock[block_hash]
	bytesOld, _ := json.Marshal(blockOld)
	bytesNew, _ := json.Marshal(recordsBlock[0])

	if !exist || !bytes.Equal(bytesOld, bytesNew) {

		MapBlock[block_hash] = recordsBlock[0]

		block_sequence := strconv.FormatFloat(recordsBlock[0]["sequence"].(float64), 'f', -1, 64)

		parse(net, block_sequence)
	}
}

func parse(net string, block_sequence string) {

	tx_embeded := xlm.TxForLedger(net, block_sequence, "", 200, xlm.ORDER_DESC)["_embedded"].(map[string]interface{})

	tx_record, _ := json.Marshal(tx_embeded["records"])
	var recordsTx []map[string]interface{}
	json.Unmarshal(tx_record, &recordsTx)

	for _, objTx := range recordsTx {

		txID := objTx["id"].(string)

		payment_embeded := xlm.PaymentForTx(net, txID, "", 200, xlm.ORDER_DESC)["_embedded"].(map[string]interface{})

		payment_record, _ := json.Marshal(payment_embeded["records"])
		var recordsPay []map[string]interface{}
		json.Unmarshal(payment_record, &recordsPay)

		for _, objPay := range recordsPay {

			ttype := objPay["type"].(string)

			if ttype == "payment" {

				from := objPay["from"].(string)
				//fmt.Println("from : ", from)

				to := objPay["to"].(string)
				//fmt.Println("to : ", to)

				amount := objPay["amount"].(string)
				//fmt.Println("amount : ", amount)

				//asset_type := objPay["asset_type"]
				//fmt.Println("asset_type : ", asset_type)

				//asset_code := objPay["asset_code"]
				//fmt.Println("asset_code : ", asset_code)

				//asset_issuer := objPay["asset_issuer"]
				//fmt.Println("asset_issuer : ", asset_issuer)

				fee := objTx["fee_paid"].(float64)
				//fmt.Println("fee_paid : ", fee)

				data := make(map[string]interface{})
				data["token"] = "XLM"
				data["transaction_id"] = txID
				data["transaction_fee"] = fee
				data["from_address"] = from
				data["to_address"] = to
				data["amount"] = amount

				if de, ok := dbScan.HMAP_DEPOSIT[data["to_address"].(string)]; ok {
					de.Notify(data)

					dbScan.Report_Fees("XLM", fee)
					dbScan.Report_Current("XLM")

					amountF, _ := strconv.ParseFloat(amount, 64)
					dbScan.Report_Deposit("XLM", amountF)

				}
			}

		}

	}

}
