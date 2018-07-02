package xlm

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/stellar/go/amount"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/hash"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/xdr"
)

var NET_PUBLIC = "public"
var NET_TEST = "test"

func HorizonNetwork(net string) (network *horizon.Client) {

	switch net {
	case NET_PUBLIC:
		network = horizon.DefaultPublicNetClient
	case NET_TEST:
		network = horizon.DefaultTestNetClient
	}

	return
}

func AmountMustParse(v string) xdr.Int64 {
	return amount.MustParse(v)
}

func AmountParse(v string) (xdr.Int64, error) {
	return amount.Parse(v)
}

func AmountStringFromInt64(v int64) string {
	return amount.StringFromInt64(v)
}

func AmountString(v xdr.Int64) string {
	return amount.String(v)
}

func KeyPairRandom() (full *keypair.Full, err error) {

	full, err = keypair.Random()
	if err != nil {
		fmt.Println("Error KeyPairRandom", err)
	}

	fmt.Println("KeyPairRandom", full.Seed(), full.Address())

	return
}

func KeyPairFromPassphrase(passphrase string) (full *keypair.Full, err error) {

	rawSeed := hash.Hash([]byte(passphrase))

	full, err = keypair.FromRawSeed(rawSeed)
	if err != nil {
		fmt.Println("Error KeyPairFromPassphrase", err)
	}

	fmt.Println("KeyPairFromPassphrase", full.Seed(), full.Address())

	return
}

func KeyPairFromRaw(rawSeed [32]byte) (full *keypair.Full, err error) {

	full, err = keypair.FromRawSeed(rawSeed)
	if err != nil {
		fmt.Println("Error KeyPairFromSeed", err)
	}

	fmt.Println("KeyPairFromSeed", full.Seed(), full.Address())

	return
}

func KeypairParse(addressOrSeed string) (keypair.KP, error) {
	return keypair.Parse(addressOrSeed)
}

func AssestsCreate() {

}

func AssetsAll() {

}

func LedgerAll() {

}

func LedgerByID() {

}

func OfferForAccount() {

}

func OperationsAll() {

}

func OperationsByID() {

}

func OperationsForAccount() {

}

func OperationsForLedger() {

}

func OperationsForTx() {

}

func OrderBookDetails() {

}

func FindPaymentPath() {

}

func TradeAggregations() {

}

func TradeAll() {

}

func TradeForAccount() {

}

func PaymentAll() {

}

func PaymenForAccount() {

}

func PaymentForLedger() {

}

func PaymentForTx() {

}

func TxBuilder(net, fromSeed, toAddress, amount string) (tx *build.TransactionBuilder, err error) {

	network := HorizonNetwork(net)

	tx, err = build.Transaction(
		build.SourceAccount{AddressOrSeed: fromSeed},
		build.AutoSequence{SequenceProvider: network},
		build.Payment(
			build.Destination{AddressOrSeed: toAddress},
			build.NativeAmount{Amount: amount},
		),
	)
	if err != nil {
		fmt.Println("Error TxBuilder", err)
	}

	fmt.Println("TxBuilder", tx)

	return
}

func TxSign(seed string, tx *build.TransactionBuilder) (string, error) {

	txe, err := tx.Sign(seed)
	if err != nil {
		fmt.Println("Error TxSign", err)
		return "", err
	}

	txeB64, err := txe.Base64()
	if err != nil {
		fmt.Println("Error Base64", err)
		return "", err
	}

	fmt.Printf("tx base64: %s", txeB64)

	return txeB64, nil
}

func TxSubmit(net, txeB64 string) {

	resp, err := HorizonNetwork(net).SubmitTransaction(txeB64)
	if err != nil {
		panic(err)
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)
}

func TxDecode(data string) (tx xdr.TransactionEnvelope) {

	rawr := strings.NewReader(data)
	b64r := base64.NewDecoder(base64.StdEncoding, rawr)

	bytesRead, err := xdr.Unmarshal(b64r, &tx)

	fmt.Printf("read %d bytes\n", bytesRead)

	if err != nil {
		fmt.Println("Error txDecode", err)
	}

	fmt.Printf("This tx has %d operations\n", len(tx.Tx.Operations))

	return
}

func TxAll() {

}

func TxByHash() {

}

func TxForAccount() {

}

func TxForLedger() {

}
