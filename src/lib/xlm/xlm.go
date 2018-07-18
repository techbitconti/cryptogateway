package xlm

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/agl/ed25519"
	"github.com/stellar/go/amount"
	"github.com/stellar/go/build"
	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/go/hash"
	"github.com/stellar/go/keypair"
	networkStellar "github.com/stellar/go/network"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

var ORDER_ASC = `asc`
var ORDER_DESC = `desc`

var ASSET_TYPE_NATIVE = `native`
var ASSET_TYPE_CREDIT_ALPHANUM4 = `credit_alphanum4`
var ASSET_TYPE_CREDIT_ALPHANUM12 = `credit_alphanum12`

func HorizonNetwork(net string) (network *horizon.Client) {

	switch net {
	case "public":
		network = horizon.DefaultPublicNetClient
	case "test":
		network = horizon.DefaultTestNetClient
	}

	return
}

func HorizonPassPhrase(net string) (pass string) {

	switch net {
	case "public":
		pass = networkStellar.PublicNetworkPassphrase
	case "test":
		pass = networkStellar.TestNetworkPassphrase
	}

	return
}

/*-------------------------------Amount--------------------------------------*/

func AmountParse(v string) (xdr.Int64, error) { // stroops
	return amount.Parse(v)
}

func AmountStringFromInt64(v int64) string {
	return amount.StringFromInt64(v)
}

func AmountString(v xdr.Int64) string {
	return amount.String(v)
}

func ToLumens(v string) float64 {

	stroops, err := amount.Parse(v)
	if err != nil {
		return 0
	}

	return float64(stroops) / math.Pow10(7)
}

func ToStroops(lumen float64) float64 {
	return lumen * math.Pow10(7)
}

/*-------------------------------Verify--------------------------------------*/

func VerifyAmount(v string) bool {

	_, err := AmountParse(v)
	if err != nil {
		return false
	}

	return true
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

func VerifyTx(net, txHash string) bool {
	tx := TxByHash(net, txHash)
	if tx.ID != "" && tx.Ledger != 0 {
		return true
	}

	return false
}

/*-------------------------------KeyPair--------------------------------------*/

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

/*-------------------------------Builder--------------------------------------*/

func NewAccountID(addr string) (accountID xdr.AccountId, err error) {

	err = accountID.SetAddress(addr)

	return
}

func NewAssest(asset_type, asset_code, asset_issuer string) (asset xdr.Asset, err error) {

	if asset_type != ASSET_TYPE_NATIVE && asset_type !=
		ASSET_TYPE_CREDIT_ALPHANUM4 && asset_type !=
		ASSET_TYPE_CREDIT_ALPHANUM12 {

		err = errors.New("Error asset_type")
		return
	}

	var accountID xdr.AccountId
	err = accountID.SetAddress(asset_issuer)
	if err != nil {
		return
	}

	err = asset.SetCredit(asset_code, accountID)
	if err != nil {
		return
	}

	switch asset_type {
	case ASSET_TYPE_NATIVE:
		asset.Type = xdr.AssetTypeAssetTypeNative
	case ASSET_TYPE_CREDIT_ALPHANUM4:
		asset.Type = xdr.AssetTypeAssetTypeCreditAlphanum4
	case ASSET_TYPE_CREDIT_ALPHANUM12:
		asset.Type = xdr.AssetTypeAssetTypeCreditAlphanum12
	}

	return
}

func NewNative() (asset xdr.Asset, err error) {

	err = asset.SetNative()

	return
}

func NewMemo(aType xdr.MemoType, value interface{}) (memo xdr.Memo, err error) {

	return xdr.NewMemo(aType, value)
}

func NewOption(muts ...interface{}) build.SetOptionsBuilder {

	return build.SetOptions(muts)
}

func SetOption_Inflation(option *build.SetOptionsBuilder, to string) {

	var destination xdr.AccountId
	err := destination.SetAddress(to)
	if err != nil {
		return
	}

	option.Mutate(build.InflationDest(to))
}

func SetOption_SetAuthRequired(option *build.SetOptionsBuilder) {
	option.Mutate(build.SetAuthRequired())
}

func SetOption_SetAuthRevocable(option *build.SetOptionsBuilder) {
	option.Mutate(build.SetAuthRevocable())
}

func SetOption_SetAuthImmutable(option *build.SetOptionsBuilder) {
	option.Mutate(build.SetAuthImmutable())
}

func SetOption_ClearAuthRequired(option *build.SetOptionsBuilder) {
	option.Mutate(build.ClearAuthRequired())
}

func SetOption_ClearAuthRevocable(option *build.SetOptionsBuilder) {
	option.Mutate(build.ClearAuthRevocable())
}

func SetOption_ClearAuthImmutable(option *build.SetOptionsBuilder) {
	option.Mutate(build.ClearAuthImmutable())
}

func SetOption_MasterWeight(option *build.SetOptionsBuilder, weight uint32) {
	option.Mutate(build.MasterWeight(weight))
}

func SetOption_Threshold(option *build.SetOptionsBuilder, low, medium, high uint32) {
	option.Mutate(build.SetThresholds(low, medium, high))
}

func SetOption_AddSigner(option *build.SetOptionsBuilder, address string, weight uint32) {

	if !VerifyAddress(address) {
		return
	}

	option.Mutate(build.AddSigner(address, weight))
}

func SetOption_HomeDomain(option *build.SetOptionsBuilder, domain string) {
	option.Mutate(build.HomeDomain(domain))
}

func NewPayment(net, fromSeed, toAddress, amount string) (tx *build.TransactionBuilder, err error) {

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
	pass := HorizonPassPhrase(net)

	tx, err = build.Transaction(
		build.SourceAccount{AddressOrSeed: fromSeed},
		build.AutoSequence{SequenceProvider: network},
		build.Payment(
			build.Destination{AddressOrSeed: toAddress},
			build.NativeAmount{Amount: amount},
		),
	)
	tx.NetworkPassphrase = pass
	if err != nil {
		fmt.Println("Error TxBuilder", err)
		return
	}
	fmt.Println("TxBuilder", tx)

	return
}

func TxSign(tx *build.TransactionBuilder, fromSeed string) string {
	txe, err := tx.Sign(fromSeed)
	if err != nil {
		fmt.Println("Error TxSign", err)
		return ""
	}

	txeB64, err := txe.Base64()
	if err != nil {
		fmt.Println("Error Base64", err)
		return ""
	}
	fmt.Println("tx base64: %s", txeB64)

	return txeB64
}

func TxSubmit(net, txeB64 string) string {

	resp, err := HorizonNetwork(net).SubmitTransaction(txeB64)
	if err != nil {
		fmt.Println("Error TxSubmit", err)
	}

	fmt.Println("transaction posted in ledger:", resp.Ledger)

	return resp.Hash
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

func TxBuilder(sourceAccount string) (tx *xdr.Transaction) {

	if !VerifyAddress(sourceAccount) {
		return
	}

	var source xdr.AccountId
	err := source.SetAddress(sourceAccount)
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

func TxAddTimeBound(tx *xdr.Transaction, timebound xdr.TimeBounds) {
	tx.TimeBounds = &timebound
}

func TxAddPaymentOp(tx *xdr.Transaction, to, amount string, asset xdr.Asset) {

	if !VerifyAddress(to) {
		return
	}

	var destination xdr.AccountId
	err := destination.SetAddress(to)
	if err != nil {
		return
	}

	lumens, err := AmountParse(amount)
	if err != nil {
		return
	}

	op := xdr.PaymentOp{
		Destination: destination,
		Asset:       asset,
		Amount:      lumens,
	}

	body, err := xdr.NewOperationBody(xdr.OperationTypePayment, op)
	if err != nil {
		return
	}

	operation := xdr.Operation{Body: body}

	tx.Operations = append(tx.Operations, operation)
}

func TxAddCreateAccountOp(tx *xdr.Transaction, to, balance string) {

	if !VerifyAddress(to) {
		return
	}

	var destination xdr.AccountId
	err := destination.SetAddress(to)
	if err != nil {
		return
	}

	lumens, err := AmountParse(balance)
	if err != nil {
		return
	}

	op := xdr.CreateAccountOp{
		Destination:     destination,
		StartingBalance: lumens,
	}

	body, err := xdr.NewOperationBody(xdr.OperationTypeCreateAccount, op)
	if err != nil {
		return
	}

	operation := xdr.Operation{Body: body}

	tx.Operations = append(tx.Operations, operation)
}

func TxAddPathPaymentOp(tx *xdr.Transaction) {}

func TxAddManageOfferOp(tx *xdr.Transaction, selling, buying xdr.Asset, amount xdr.Int64, price xdr.Price, offerID xdr.Uint64) {

	op := xdr.ManageOfferOp{
		Selling: selling,
		Buying:  buying,
		Amount:  amount,
		Price:   price,
		OfferId: offerID,
	}

	body, err := xdr.NewOperationBody(xdr.OperationTypeManageOffer, op)
	if err != nil {
		return
	}

	operation := xdr.Operation{Body: body}

	tx.Operations = append(tx.Operations, operation)
}

func TxAddCreatePassiveOfferOp(tx *xdr.Transaction, selling, buying xdr.Asset, amount xdr.Int64, price xdr.Price) {

	op := xdr.CreatePassiveOfferOp{
		Selling: selling,
		Buying:  buying,
		Amount:  amount,
		Price:   price,
	}

	body, err := xdr.NewOperationBody(xdr.OperationTypeCreatePassiveOffer, op)
	if err != nil {
		return
	}

	operation := xdr.Operation{Body: body}

	tx.Operations = append(tx.Operations, operation)
}

func TxAddSetOptionsOp(tx *xdr.Transaction, op *xdr.SetOptionsOp) {

	body, err := xdr.NewOperationBody(xdr.OperationTypeSetOptions, op)
	if err != nil {
		return
	}

	operation := xdr.Operation{Body: body}

	tx.Operations = append(tx.Operations, operation)
}

func TxAddChangeTrustOp(tx *xdr.Transaction, asset_type, asset_code, asset_issuer, limit string) {

	if limit == "" {
		limit = string(build.MaxLimit)
	}

	tokens, errL := AmountParse(limit)
	if errL != nil {
		return
	}

	ass, errA := NewAssest(asset_type, asset_code, asset_issuer)
	if errA != nil {
		return
	}

	op := xdr.ChangeTrustOp{Line: ass, Limit: tokens}

	body, err := xdr.NewOperationBody(xdr.OperationTypeChangeTrust, op)
	if err != nil {
		return
	}

	operation := xdr.Operation{Body: body}

	tx.Operations = append(tx.Operations, operation)
}

func TxAddAllowTrustOp(tx *xdr.Transaction, trustor, asset_code string, asset_type xdr.AssetType, authorize bool) {

	if !VerifyAddress(trustor) {
		return
	}

	var destination xdr.AccountId
	err := destination.SetAddress(trustor)
	if err != nil {
		return
	}

	ass, errA := xdr.NewAllowTrustOpAsset(asset_type, asset_code)
	if errA != nil {
		return
	}

	op := xdr.AllowTrustOp{Trustor: destination, Asset: ass, Authorize: authorize}

	body, err := xdr.NewOperationBody(xdr.OperationTypeAllowTrust, op)
	if err != nil {
		return
	}

	operation := xdr.Operation{Body: body}

	tx.Operations = append(tx.Operations, operation)
}

func TxAddAccountMergeOp(tx *xdr.Transaction, to string) {

	if !VerifyAddress(to) {
		return
	}

	accountMergeBuilder := build.AccountMerge(build.Destination{to})

	tx.Operations = append(tx.Operations, accountMergeBuilder.O)
}

func TxAddInflationOp(tx *xdr.Transaction, to string, lumens xdr.Int64) {

	if !VerifyAddress(to) {
		return
	}

	var destination xdr.AccountId
	err := destination.SetAddress(to)
	if err != nil {
		return
	}

	op := xdr.InflationPayout{Destination: destination, Amount: lumens}

	body, err := xdr.NewOperationBody(xdr.OperationTypeInflation, op)
	if err != nil {
		return
	}

	operation := xdr.Operation{Body: body}

	tx.Operations = append(tx.Operations, operation)
}

func TxAddManageDataOp(tx *xdr.Transaction, key string, value []byte) {

	dataBuilder := build.SetData(key, value)

	tx.Operations = append(tx.Operations, dataBuilder.O)
}

func TxEnvelopBuilder(tx xdr.Transaction) *xdr.TransactionEnvelope {

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

/*-------------------------------Query--------------------------------------*/

func cursor_limit_order(cursor string, limit uint, order string) string {

	curs := "?cursor=" + cursor
	lim := "&limit=" + strconv.Itoa(int(limit))
	ordr := "&order=" + order

	return curs + lim + ordr
}

func call(url string) ([]byte, bool) {

	//fmt.Println(".....url......", url)

	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	var problem horizon.Problem
	json.Unmarshal(body, &problem)
	if problem.Type != "" {
		fmt.Println("Error call", string(body))
		return []byte{}, false
	}

	return body, true
}

func FriendBot(addr string) (result interface{}) {

	url := "https://horizon-testnet.stellar.org/friendbot?addr=" + addr

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......FriendBot........")
	fmt.Println(string(body))

	return
}

func AccountDetails(net, id string) (result horizon.Account) {

	url := HorizonNetwork(net).URL + "/accounts/" + id

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	//fmt.Println(".......AccountDetails........")
	//fmt.Println(string(body))

	return
}

func GetBalance(net, id, asset_type, asset_code string) (balance float64) {

	info := AccountDetails(net, id)

	for _, v := range info.Balances {

		if asset_type == v.Type && asset_code == v.Code {
			balance, _ = strconv.ParseFloat(v.Balance, 64)
			break
		}
	}

	fmt.Println(".........GetBalance........", balance)

	return
}

func AssetCodeIssuer(net, code, issuer string, cursor string, limit uint, order string) (result map[string]interface{}) {

	asset_code := "&asset_code=" + code
	asset_issuer := "&asset_issuer=" + issuer

	url := HorizonNetwork(net).URL + "/assets" + cursor_limit_order(cursor, limit, order) + asset_code + asset_issuer

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......AssetCodeIssuer........")
	fmt.Println(string(body))

	return
}

func LedgerAll(net string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/ledgers" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	//fmt.Println(".......LedgerAll........")
	//fmt.Println(string(body))

	return
}

func LedgerByID(net, id string) (result horizon.Ledger) {

	url := HorizonNetwork(net).URL + "/ledgers/" + id

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......LedgerID........")
	fmt.Println(string(body))

	return

}

func OfferForAccount(net, id string, cursor string, limit uint, order string) (result horizon.OffersPage) {

	url := HorizonNetwork(net).URL + "/accounts/" + id + "/offers" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OfferForAccount........")
	fmt.Println(string(body))

	return
}

func OperationsAll(net string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/operations" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OperationsAll........")
	fmt.Println(string(body))

	return
}

func OperationsByID(net, id string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/operations/" + id

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OperationsByID........")
	fmt.Println(string(body))

	return
}

func OperationsForAccount(net, id string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/accounts/" + id + "/operations" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OperationsForAccount........")
	fmt.Println(string(body))

	return
}

func OperationsForLedger(net, id string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/ledgers/" + id + "/operations" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OperationsForLedger........")
	fmt.Println(string(body))

	return
}

func OperationsForTx(net, txHash string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/transactions/" + txHash + "/operations" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OperationsForTx........")
	fmt.Println(string(body))

	return
}

func OrderBookDetails(net, selling_asset_type, selling_asset_code, selling_asset_issuer,
	buying_asset_type, buying_asset_code, buying_asset_issuer, limit string) (result horizon.OrderBookSummary) {

	selling_asset_type = "?selling_asset_type=" + selling_asset_type
	selling_asset_code = "&selling_asset_code=" + selling_asset_code
	selling_asset_issuer = "&selling_asset_issuer=" + selling_asset_issuer

	buying_asset_type = "&buying_asset_type=" + buying_asset_type
	buying_asset_code = "&buying_asset_code=" + buying_asset_code
	buying_asset_issuer = "&buying_asset_issuer=" + buying_asset_issuer

	limit = "&limit=" + limit

	params := selling_asset_type + selling_asset_code + selling_asset_issuer + buying_asset_type + buying_asset_code + buying_asset_issuer + limit

	url := HorizonNetwork(net).URL + "/order_book" + params

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OrderBookDetails........")
	fmt.Println(string(body))

	return
}

func FindPaymentPath(net, destination_account, destination_asset_type, destination_asset_code, destination_asset_issuer, destination_amount, source_account string) (result map[string]interface{}) {

	destination_account = "?destination_account=" + destination_account
	destination_asset_type = "&destination_asset_type=" + destination_asset_type
	destination_asset_code = "&destination_asset_code=" + destination_asset_code
	destination_asset_issuer = "&destination_asset_issuer=" + destination_asset_issuer
	destination_amount = "&destination_amount=" + destination_amount
	source_account = "&destination_amount=" + destination_amount

	params := destination_account + destination_asset_type + destination_asset_code + destination_asset_issuer + destination_amount + source_account

	url := HorizonNetwork(net).URL + "/paths" + params

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OrderBookDetails........")
	fmt.Println(string(body))

	return
}

func TradeAggregations(net, base_asset_type, base_asset_code, base_asset_issuer,
	counter_asset_type, counter_asset_code, counter_asset_issuer,
	order string,
	limit, start_time, end_time, resolution uint64) (result map[string]interface{}) {

	base_asset_type = "?base_asset_type=" + base_asset_type
	base_asset_code = "&base_asset_code=" + base_asset_code
	base_asset_issuer = "&base_asset_issuer=" + base_asset_issuer

	counter_asset_type = "&counter_asset_type=" + counter_asset_type
	counter_asset_code = "&counter_asset_code=" + counter_asset_code
	counter_asset_issuer = "&counter_asset_issuer=" + counter_asset_issuer
	order = "&order=" + order

	limitStr := "&limit=" + strconv.FormatUint(limit, 10)
	start_timeStr := "&start_time=" + strconv.FormatUint(start_time, 10)
	end_timeStr := "&end_time=" + strconv.FormatUint(end_time, 10)
	resolutionStr := "&resolution=" + strconv.FormatUint(resolution, 10)

	params := base_asset_type + base_asset_code + base_asset_issuer +
		counter_asset_type + counter_asset_code + counter_asset_issuer +
		order + limitStr + start_timeStr + end_timeStr + resolutionStr

	url := HorizonNetwork(net).URL + "/trade_aggregations" + params

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OrderBookDetails........")
	fmt.Println(string(body))

	return
}

func TradeAll(net, base_asset_type, base_asset_code, base_asset_issuer,
	counter_asset_type, counter_asset_code, counter_asset_issuer,
	offer_id, cursor, order, limit string) (result map[string]interface{}) {

	base_asset_type = "?base_asset_type=" + base_asset_type
	base_asset_code = "&base_asset_code=" + base_asset_code
	base_asset_issuer = "&base_asset_issuer" + base_asset_issuer
	counter_asset_type = "&counter_asset_type=" + counter_asset_type
	counter_asset_code = "&counter_asset_code=" + counter_asset_code
	counter_asset_issuer = "&counter_asset_issuer=" + counter_asset_issuer
	offer_id = "&offer_id=" + offer_id
	cursor = "&cursor=" + cursor
	order = "&order=" + order
	limit = "&limit=" + limit

	params := base_asset_type + base_asset_code + base_asset_issuer +
		counter_asset_type + counter_asset_code + counter_asset_issuer +
		offer_id + cursor + order + limit

	url := HorizonNetwork(net).URL + "/trades" + params

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......TradeAll........")
	fmt.Println(string(body))

	return

}

func TradeForAccount(net, id string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/accounts/" + id + "/trades" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......TradeForAccount........")
	fmt.Println(string(body))

	return
}

func PaymentAll(net string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/payments" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......PaymentAll........")
	fmt.Println(string(body))

	return
}

func PaymenForAccount(net, id string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/accounts/" + id + "/payments" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......OperationsForAccount........")
	fmt.Println(string(body))

	return
}

func PaymentForLedger(net, id string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/ledgers/" + id + "/payments" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......PaymentForLedger........")
	fmt.Println(string(body))

	return
}

func PaymentForTx(net, txHash string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/transactions/" + txHash + "/payments" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	//fmt.Println(".......PaymentForTx........")
	//fmt.Println(string(body))

	return
}

func TxAll(net string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/transactions" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......TxAll........")
	fmt.Println(string(body))

	return
}

func TxByHash(net, txHash string) (result horizon.Transaction) {

	url := HorizonNetwork(net).URL + "/transactions/" + txHash

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......TxByHash........")
	fmt.Println(string(body))

	return
}

func TxForAccount(net, id string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/accounts/" + id + "/transactions" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	fmt.Println(".......TxForAccount........")
	fmt.Println(string(body))

	return

}

func TxForLedger(net, id string, cursor string, limit uint, order string) (result map[string]interface{}) {

	url := HorizonNetwork(net).URL + "/ledgers/" + id + "/transactions" + cursor_limit_order(cursor, limit, order)

	body, ok := call(url)
	if !ok {
		return
	}
	json.Unmarshal(body, &result)

	//fmt.Println(".......TxForLedger........")
	//fmt.Println(string(body))

	return
}
