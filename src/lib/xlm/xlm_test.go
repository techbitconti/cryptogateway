package xlm

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

func Test_XLM(t *testing.T) {

	// 	GetBalance("test", "GCHKKQ5VWJBRQZHNMODO5BWYZKPNM2HDSJ26T4O644CNEQBYK7IXATKM", "credit_alphanum12", "nCntGameCoin")
	//	AccountDetails("test", "GCHKKQ5VWJBRQZHNMODO5BWYZKPNM2HDSJ26T4O644CNEQBYK7IXATKM")
	//	AssetCodeIssuer("test", "EUR", "GDZQZ6YRVUKC7AJVWNT5IKNSLVJTFEKAQ35DVH3YZUSA4RNY4BKX4Q6D", "", 200, ORDER_DESC)
	//	TxAll("test", "", 2, ORDER_DESC)
	//	TxForAccount("test", "9823770", "", 200, ORDER_ASC)
	//	TxByHash("test", "eedef59536c572c938e8934107b72a70291a420e384357fd33515196bb2c6bc0")
	//	TxForLedger("public", "18762210", "", 200, ORDER_ASC)

	//	LedgerByID("public", "18762210")
	// 	LedgerAll("public", "", 1, ORDER_DESC)

	// Operation
	//	OperationsForLedger("public", "18762210", "", 200, ORDER_DESC)
	//	OperationsForTx("public", "de668bb82ff67c96e33eaaf13af530ee46cf408c00339ecfe8086a3d97f9dbf7", "", 200, ORDER_DESC)

	// payment
	//	PaymentForLedger("test", "18760618", "", 200, ORDER_DESC)
	//	PaymentForTx("public", "3a024e2c0c4534fa246d0b756098213a7c0f7f919cdbbc6b7d93bc89723c541f", "", 200, ORDER_DESC)

	// offer

	// order book

	// trade

}

func Test_Parse(t *testing.T) {

	net := "public"

	block_embeded := LedgerAll(net, "", 1, ORDER_DESC)["_embedded"].(map[string]interface{})

	block_record, _ := json.Marshal(block_embeded["records"])
	var recordsBlock []map[string]interface{}
	json.Unmarshal(block_record, &recordsBlock)

	block_sequence := strconv.FormatFloat(recordsBlock[0]["sequence"].(float64), 'f', -1, 64)

	tx_embeded := TxForLedger(net, block_sequence, "", 200, ORDER_ASC)["_embedded"].(map[string]interface{})
	tx_record, _ := json.Marshal(tx_embeded["records"])
	var recordsTx []map[string]interface{}
	json.Unmarshal(tx_record, &recordsTx)

	for _, v := range recordsTx {
		txID := v["id"].(string)
		fmt.Println(txID)
		PaymentForTx(net, txID, "", 200, ORDER_DESC)

		// [] native
	}

}
