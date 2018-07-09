package xlm

import (
	//"encoding/json"
	//"fmt"
	//"strconv"
	"testing"
)

func Test_Qurey(t *testing.T) {

	//  GetBalance("test", "GAYMTHMZCJYQCFLVRYZ4DMI4Q5B6VFFBJFQVGE2WJHIOPDC4UAG4TD3K", "native", "")
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

func Test_Tranfer(t *testing.T) {

	tx, _ := TxBuilder("test", "SBO3QHEVLPSI7SYFHCXFPCBWYN2WU6GL6Z7PW6FWSWY2GF67GJEK3PYJ", "GAYMTHMZCJYQCFLVRYZ4DMI4Q5B6VFFBJFQVGE2WJHIOPDC4UAG4TD3K", "9")
	txeB64 := TxSign(tx, "SBO3QHEVLPSI7SYFHCXFPCBWYN2WU6GL6Z7PW6FWSWY2GF67GJEK3PYJ")
	TxSubmit("test", txeB64)

}

/*
func Test_Parse(t *testing.T) {

	net := "public"

	block_embeded := LedgerAll(net, "", 1, ORDER_DESC)["_embedded"].(map[string]interface{})
	block_record, _ := json.Marshal(block_embeded["records"])
	var recordsBlock []map[string]interface{}
	json.Unmarshal(block_record, &recordsBlock)

	block_hash := recordsBlock[0]["hash"].(string)
	fmt.Println("block hash : ", block_hash)

	block_sequence := strconv.FormatFloat(recordsBlock[0]["sequence"].(float64), 'f', -1, 64)

	tx_embeded := TxForLedger(net, block_sequence, "", 200, ORDER_ASC)["_embedded"].(map[string]interface{})
	tx_record, _ := json.Marshal(tx_embeded["records"])
	var recordsTx []map[string]interface{}
	json.Unmarshal(tx_record, &recordsTx)

	for _, objTx := range recordsTx {

		txID := objTx["id"].(string)
		isTranfer := false

		payment_embeded := PaymentForTx(net, txID, "", 200, ORDER_DESC)["_embedded"].(map[string]interface{})
		payment_record, _ := json.Marshal(payment_embeded["records"])
		var recordsPay []map[string]interface{}
		json.Unmarshal(payment_record, &recordsPay)

		for _, objPay := range recordsPay {

			ttype := objPay["type"].(string)

			if ttype == "payment" {

				isTranfer = true

				if from, okF := objPay["from"]; okF {
					fmt.Println("from : ", from)
				}

				if to, okT := objPay["to"]; okT {
					fmt.Println("to : ", to)
				}

				if amount, okAm := objPay["amount"]; okAm {
					fmt.Println("amount : ", amount)
				}

				if asset_type, okAt := objPay["asset_type"]; okAt {
					fmt.Println("asset_type : ", asset_type)
				}

				if asset_code, okAc := objPay["asset_code"]; okAc {
					fmt.Println("asset_code : ", asset_code)
				}

				if asset_issuer, okAi := objPay["asset_issuer"]; okAi {
					fmt.Println("asset_issuer : ", asset_issuer)
				}
			}

			if ttype == "create_account" {

				isTranfer = true

				if funder, okF := objPay["funder"]; okF {
					fmt.Println("funder : ", funder)
				}

				if account, okT := objPay["account"]; okT {
					fmt.Println("account : ", account)
				}

				if starting_balance, okAm := objPay["starting_balance"]; okAm {
					fmt.Println("starting_balance : ", starting_balance)
				}
			}
		}

		if isTranfer {
			fee := objTx["fee_paid"].(float64)
			fmt.Println("fee_paid : ", fee)
		}
	}
}
*/
