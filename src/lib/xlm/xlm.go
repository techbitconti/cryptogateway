package xlm

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/agl/ed25519"
	"github.com/stellar/go/amount"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/hash"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/strkey"
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

func AmountParse(v string) (xdr.Int64, error) {
	return amount.Parse(v)
}

func AmountStringFromInt64(v int64) string {
	return amount.StringFromInt64(v)
}

func AmountString(v xdr.Int64) string {
	return amount.String(v)
}

func VerifyAmount(v string) bool {

	_, err := AmountParse(v)
	if err != nil {
		return false
	}

	return true
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

func KeyPairParse(addressOrSeed string) (keypair.KP, error) {
	return keypair.Parse(addressOrSeed)
}

func KeySignFull(kp *keypair.Full, input []byte) ([]byte, error) {
	return kp.Sign(input)
}

func KeySignFromAddress(kp *keypair.FromAddress, input []byte) ([]byte, error) {
	return kp.Sign(input)
}

func KeySignSeed(seed string, input []byte) ([]byte, error) {

	rawSeed, err := strkey.Decode(strkey.VersionByteSeed, seed)
	if err != nil {
		return []byte{}, err
	}

	reader := bytes.NewReader(rawSeed)
	_, priv, err := ed25519.GenerateKey(reader)
	if err != nil {
		return []byte{}, err
	}

	return xdr.Signature(ed25519.Sign(priv, input)[:]), nil
}

func VerifyAddress(addr string) bool {

	// prefix : G
	_, err := strkey.Decode(strkey.VersionByteAccountID, addr)
	if err != nil {
		fmt.Println("Error VerifyAddress", err)
		return false
	}

	return true
}

func VerifySeed(seed string) bool {

	// prefix : S
	_, err := strkey.Decode(strkey.VersionByteSeed, seed)
	if err != nil {
		fmt.Println("Error VerifySeed", err)
		return false
	}

	return true
}

func VerifyHashTx(tx string) bool {

	// prefix : T
	_, err := strkey.Decode(strkey.VersionByteHashTx, tx)
	if err != nil {
		fmt.Println("Error VerifyHashTx", err)
		return false
	}

	return true
}

func VerifyHashX(hx string) bool {

	// prefix : X
	_, err := strkey.Decode(strkey.VersionByteHashX, hx)
	if err != nil {
		fmt.Println("Error VerifyHashX", err)
		return false
	}

	return true
}

func AccountID(addr string) (accountID xdr.AccountId, err error) {

	err = accountID.SetAddress(addr)

	return
}

func AssestSetCredit(code, issuer string) (asset xdr.Asset, err error) {

	var accountID xdr.AccountId
	err = accountID.SetAddress(issuer)
	if err != nil {
		return
	}

	err = asset.SetCredit(code, accountID)
	if err != nil {
		return
	}

	return
}

func AssetSetNative() (asset xdr.Asset, err error) {

	err = asset.SetNative()

	return
}

func MemoNew(aType xdr.MemoType, value interface{}) (memo xdr.Memo, err error) {

	return xdr.NewMemo(aType, value)
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

	if !VerifySeed(fromSeed) {
		return
	}

	if !VerifyAddress(toAddress) {
		return
	}

	if !VerifyAmount(amount) {
		return
	}

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

func TxBuilderSign(seed string, tx *build.TransactionBuilder) (string, error) {

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

func TxBuilderSubmit(net, txeB64 string) {

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

func TxNew(from string) (tx *xdr.Transaction) {

	if !VerifyAddress(from) {
		return
	}

	var source xdr.AccountId
	err := source.SetAddress(from)
	if err != nil {
		return
	}

	tx = &xdr.Transaction{
		SourceAccount: source,
		SeqNum:        xdr.SequenceNumber(1),
		Operations:    []xdr.Operation{},
	}

	return
}

func TxAddFee(tx *xdr.Transaction, fee xdr.Uint32) {
	tx.Fee = fee
}

func TxAddMemo(tx *xdr.Transaction, memo xdr.Memo) {
	tx.Memo = memo
}

func TxAddPaymentOp(tx *xdr.Transaction, to, amount string, asset xdr.Asset) {

	var destination xdr.AccountId
	err := destination.SetAddress(to)
	if err != nil {
		return
	}

	lumens, err := AmountParse(amount)
	if err != nil {
		return
	}

	option := xdr.PaymentOp{
		Destination: destination,
		Asset:       asset,
		Amount:      lumens,
	}

	body, err := xdr.NewOperationBody(xdr.OperationTypePayment, option)
	if err != nil {
		return
	}

	operation := xdr.Operation{Body: body}

	tx.Operations = append(tx.Operations, operation)
}

func TxAddCreateAccountOp(tx *xdr.Transaction)      {}
func TxAddPathPaymentOp(tx *xdr.Transaction)        {}
func TxAddManageOfferOp(tx *xdr.Transaction)        {}
func TxAddCreatePassiveOfferOp(tx *xdr.Transaction) {}
func TxAddSetOptionsOp(tx *xdr.Transaction)         {}
func TxAddChangeTrustOp(tx *xdr.Transaction)        {}
func TxAddAllowTrustOp(tx *xdr.Transaction)         {}
func TxAddDestination(tx *xdr.Transaction)          {}
func TxAddManageDataOp(tx *xdr.Transaction)         {}

func TxEnvelopNew(tx xdr.Transaction) *xdr.TransactionEnvelope {

	txe := &xdr.TransactionEnvelope{
		Tx:         tx,
		Signatures: []xdr.DecoratedSignature{},
	}

	return txe
}

func TxAddSignature(seed string, tx *xdr.Transaction, txe *xdr.TransactionEnvelope) {

	if !VerifySeed(seed) {
		return
	}

	skp, err := KeyPairParse(seed)
	if err != nil {
		return
	}

	if skp.Address() != tx.SourceAccount.Address() {
		return
	}

	var txBytes bytes.Buffer
	_, err1 := xdr.Marshal(&txBytes, tx)
	if err1 != nil {
		return
	}

	txHash := hash.Hash(txBytes.Bytes())
	signature, err2 := skp.Sign(txHash[:])
	if err2 != nil {
		return
	}

	ds := xdr.DecoratedSignature{
		Hint:      skp.Hint(),
		Signature: xdr.Signature(signature[:]),
	}

	txe.Signatures = append(txe.Signatures, ds)
}

func TxEnvelopEncode(txe *xdr.TransactionEnvelope) (txeB64 string) {

	var txeBytes bytes.Buffer
	_, err := xdr.Marshal(&txeBytes, txe)
	if err != nil {
		return
	}
	txeB64 = base64.StdEncoding.EncodeToString(txeBytes.Bytes())

	fmt.Printf("tx base64: %s", txeB64)

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
