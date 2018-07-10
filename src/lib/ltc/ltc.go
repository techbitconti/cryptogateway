package ltc

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"config"

	"github.com/ltcsuite/ltcd/btcjson"
	"github.com/ltcsuite/ltcd/chaincfg"
	"github.com/ltcsuite/ltcd/chaincfg/chainhash"
	"github.com/ltcsuite/ltcd/rpcclient"
	"github.com/ltcsuite/ltcd/wire"
	"github.com/ltcsuite/ltcutil"
)

var NET string

var Ltcd *rpcclient.Client
var Chaincfg chaincfg.Params

var ltcdHomeDir string = config.PATH_LTC

func Connect(net string) {

	if Ltcd != nil {
		return
	}

	host := ""
	NET = net

	switch net {
	case "mainnet":
		host = "localhost:9332"
		Chaincfg = chaincfg.MainNetParams
	case "testnet":
		host = "localhost:19332"
		Chaincfg = chaincfg.TestNet4Params
	case "simnet":
		host = "localhost:19443"
		Chaincfg = chaincfg.SimNetParams
	}

	// Connect to local bitcoin core RPC server using HTTP POST mode.
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         "123",
		Pass:         "123",
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}
	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	client, err := rpcclient.New(connCfg, nil)
	if err != nil {
		log.Println(err)
	}
	//defer client.Shutdown()

	Ltcd = client

}

func GetBlockCount() int64 {

	// Get the current block count.
	blockCount, err := Ltcd.GetBlockCount()
	if err != nil {
		log.Println(err)
		return 0
	}
	//log.Println("Block count: %d", blockCount)

	return blockCount
}

func GetBlockHash(blockHeight int64) (blockHash *chainhash.Hash, err error) {

	blockHash, err = Ltcd.GetBlockHash(blockHeight)
	if err != nil {
		log.Println("Error GetBlockHash : ", err)
		return
	}
	//log.Println("GetBlockHash: %d", blockHash)

	return
}

func GetBlockHeader(blockHash *chainhash.Hash) (blockHeader *wire.BlockHeader, err error) {

	blockHeader, err = GetBlockHeader(blockHash)
	log.Println("GetBlockHeader: %d", blockHeader)
	return
}

func GetBlock(blockHash *chainhash.Hash) (block *wire.MsgBlock, err error) {

	block, err = Ltcd.GetBlock(blockHash)
	//log.Println("GetBlock : ", block)

	return
}

func GetBlockFromStr(blockHash string) (block *wire.MsgBlock, err error) {

	hx, _ := NewHashFromStr(blockHash)
	block, err = GetBlock(hx)

	return
}

func NewHashFromStr(hex string) (hash *chainhash.Hash, err error) {

	hash, err = chainhash.NewHashFromStr(hex)
	//log.Println("NewHashFromStr : ", hash, err)
	return
}

func Generate(num uint32) ([]*chainhash.Hash, error) {
	tx, err := Ltcd.Generate(num)

	log.Println("Generate", err, tx)

	return tx, err
}

// GetGenerate returns true if the server is set to mine, otherwise false.
func GetGenerate() (ok bool, err error) {

	ok, err = Ltcd.GetGenerate()

	log.Println("GetGenerate : ", ok)

	return
}

// SetGenerate sets the server to generate coins (mine) or not.
func SetGenerate(enable bool, numCPUs int) (err error) {

	err = Ltcd.SetGenerate(enable, numCPUs)
	log.Println("SetGenerate : ", err)

	return
}

func VerifyChainBlocks(checkLevel, numBlocks int32) (ok bool, err error) {

	ok, err = Ltcd.VerifyChainBlocks(checkLevel, numBlocks)

	return
}

func InvalidateBlock(blockHash *chainhash.Hash) (err error) {

	err = Ltcd.InvalidateBlock(blockHash)

	return
}

func WalletPassphrase(pass string, second int64) (bool, error) {

	//  WalletPassphrase
	err := Ltcd.WalletPassphrase(pass, second)
	if err != nil {
		log.Println("WalletPassphrase", err)

		return false, err
	}
	log.Println("WalletPassphrase: ", pass)

	return true, nil
}

func CreateNewAccount(account string) (string, error) {

	//CreateNewAccount
	err := Ltcd.CreateNewAccount(account)
	if err != nil {
		log.Println("Error CreateNewAccount", err)
		return account, err
	}
	log.Println("CreateNewAccount")

	return account, nil
}

func DecodeAddress(addr string) (address ltcutil.Address, err error) {

	address, err = ltcutil.DecodeAddress(addr, &Chaincfg)
	if err != nil {
		log.Println("Error DecodeAddress", err)
	}
	return
}

func ValidateAddress(addr string) (acc *btcjson.ValidateAddressWalletResult, err error) {

	address, err := DecodeAddress(addr)
	if err != nil {
		log.Println("Error  DecodeAddress", err)
		return
	}

	// ValidateAddress
	acc, err = Ltcd.ValidateAddress(address)
	if err != nil {
		log.Println("Error ValidateAddress", err)
		return
	}

	b, _ := json.MarshalIndent(acc, "", " ")
	log.Println(string(b))

	log.Println("ValidateAddress: ", acc)

	return
}

func ValidateAmount(amount string) (float64, bool) {

	f, err := strconv.ParseFloat(amount, 64)
	if err != nil || len(amount) > 9 {
		log.Println("ValidateAmount ", err)
		return float64(0), false
	}

	log.Println("ValidateAmount : ", f)

	return f, true
}

func ToBTC(amount string) float64 {

	f, ok := ValidateAmount(amount)
	if !ok {
		return float64(0)
	}

	value, _ := ltcutil.NewAmount(f)

	return value.ToBTC()
}

func ToSatoshi(amount string) float64 {
	f, ok := ValidateAmount(amount)
	if !ok {
		return float64(0)
	}

	value, _ := ltcutil.NewAmount(f)

	return value.ToUnit(ltcutil.AmountSatoshi)
}

func ListAccounts() (map[string]ltcutil.Amount, error) {

	// ListAccounts
	list, err := Ltcd.ListAccounts()
	if err != nil {
		log.Println("Error ListAccounts", err)

		return nil, err
	}
	//log.Println("ListAccounts: ", list)

	return list, nil
}

func ListAddress() (list []ltcutil.Address) {

	accounts, _ := ListAccounts()

	for acc, _ := range accounts {
		arr, _ := GetAddressesByAccount(acc)

		list = append(list, arr...)
	}

	log.Println("ListAddress: ", list)
	return
}

func GetBalanceAccount(account string) (amount ltcutil.Amount, err error) {

	// GetBalance
	amount, err = Ltcd.GetBalance(account)
	if err != nil {
		log.Println("Error GetBalanceAccount", err)
	}

	log.Println("GetBalanceAccount: ", amount)

	return
}

func GetBalance(addr string) float64 {

	account, err := GetAccount(addr)
	if err != nil {
		return 0
	}

	amount, rr := GetBalanceAccount(account)
	if rr != nil {
		return 0
	}

	log.Println("GetBalanceX ", addr, amount.ToBTC())

	return amount.ToBTC()
}

func GetNewAddress(account string) (address ltcutil.Address, err error) {

	list, _ := ListAccounts()
	if _, exist := list[account]; exist {
		err = errors.New("Error account exist ==>" + account)
		log.Println("Error account exist ==>", account)
		return
	}

	// GetNewAddress
	address, err = Ltcd.GetNewAddress(account)
	if err != nil {
		log.Println("Error GetNewAddress", err)
	}
	log.Println("GetNewAddress: ", address)

	return
}

// GetAccountAddress returns the current Bitcoin address for receiving payments
// to the specified account.
func GetAccountAddress(account string) (address ltcutil.Address, err error) {

	address, err = Ltcd.GetAccountAddress(account)
	if err != nil {
		log.Println("Error GetAccountAddress", err)
	}
	log.Println("GetAccountAddress: ", address)

	return
}

func GetAccount(addr string) (account string, err error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		err = rr
		log.Println("Error GetAccount", rr)
		return
	}

	// GetAccount
	account, err = Ltcd.GetAccount(address)
	if err != nil {
		log.Println("Error GetAccount", err)
	}
	log.Println("GetAccountAddress: ", account, addr)

	return
}

func GetAddressesByAccount(account string) (address []ltcutil.Address, err error) {

	// GetAddressesByAccount
	address, err = Ltcd.GetAddressesByAccount(account)
	if err != nil {
		log.Println("Error GetAddressesByAccount", err)
	}
	log.Println("GetAddressesByAccount: ", account, address)

	return
}

func GetReceivedByAccount(account string) (amount ltcutil.Amount, err error) {

	amount, err = Ltcd.GetReceivedByAccount(account)
	if err != nil {
		log.Println("Error GetReceivedByAccount")
	}
	log.Println("GetReceivedByAccount : ", amount)

	return
}

func GetReceivedByAddress(addr string) (amount ltcutil.Amount, err error) {

	address, _ := DecodeAddress(addr)

	amount, err = Ltcd.GetReceivedByAddress(address)
	if err != nil {
		log.Println("Error GetReceivedByAddress")
	}
	log.Println("GetReceivedByAddress : ", amount)

	return
}

func ListReceivedByAccount() (btcj []btcjson.ListReceivedByAccountResult, err error) {

	btcj, err = Ltcd.ListReceivedByAccount()
	if err != nil {
		log.Println("Error ListReceivedByAccount")
	}
	log.Println("ListReceivedByAccount :", btcj)

	return
}

func ListReceivedByAddress() (btcj []btcjson.ListReceivedByAddressResult, err error) {

	btcj, err = Ltcd.ListReceivedByAddress()
	if err != nil {
		log.Println("Error ListReceivedByAddress")
	}
	log.Println("ListReceivedByAddress", btcj)

	return
}

// NOTE: This function requires to the wallet to be unlocked
func DumpPrivKey(addr string) (*ltcutil.WIF, error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return nil, rr
	}

	// DumpPrivKey
	wif, err := Ltcd.DumpPrivKey(address)
	if err != nil {
		log.Println("DumpPrivKey", err)
		return nil, err
	}
	log.Println("DumpPrivKey: ", wif)

	return wif, nil
}

func ImportPrivKey(prv, label string, rescan bool) error {

	wif, err := ltcutil.DecodeWIF(prv)
	if err != nil {
		log.Println("Error ImportPrivKey")
	}
	log.Println("ImportPrivKey : ", wif)

	return Ltcd.ImportPrivKeyRescan(wif, label, rescan)
}

func ImportAddress(addr string) error {

	err := Ltcd.ImportAddress(addr)
	if err != nil {
		log.Println("Error ImportAddress")
	}
	log.Println("ImportAddress", err)

	return err
}

// NOTE: This function requires to the wallet to be unlocked
func SignMessage(addr, message string) (signature string, err error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return
	}

	// SignMessage
	signature, err = Ltcd.SignMessage(address, message)

	log.Println("SignMessage: ", signature)

	return
}

// NOTE: This function requires to the wallet to be unlocked
func VerifyMessage(addr, signature, message string) (signed bool, err error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return
	}

	signed, err = Ltcd.VerifyMessage(address, signature, message)

	log.Println("VerifyMessage: ", signed)

	return
}

// NOTE: This function requires to the wallet to be unlocked
func SendFrom(fromAddress string, toAddress string, value float64) (tx *chainhash.Hash, err error) {

	log.Println("SendFrom : ", fromAddress, " -- to : ", toAddress)

	fromAccount, frr := GetAccount(fromAddress)
	if frr != nil {
		err = frr
		log.Println("Error fromAddress", frr)
		return
	}

	to, trr := DecodeAddress(toAddress)
	if trr != nil {
		err = trr
		log.Println("Error toAddress", trr)
		return
	}

	amount, vrr := ltcutil.NewAmount(value)
	if vrr != nil {
		err = vrr
		log.Println("Error NewAmount", vrr)
		return
	}

	tx, err = Ltcd.SendFrom(fromAccount, to, amount)

	log.Println("SendFrom: ", tx)

	if NET == "simnet" {
		Generate(uint32(1))
	}

	return
}

// NOTE: This function requires to the wallet to be unlocked.  See the
func SendMany(fromAccount string, amounts map[ltcutil.Address]ltcutil.Amount) (tx *chainhash.Hash, err error) {

	tx, err = Ltcd.SendMany(fromAccount, amounts)
	if err != nil {
		log.Println("Error SendMany", err)
	}

	log.Println("SendMany: ", tx)

	if NET == "simnet" {
		Generate(uint32(1))
	}

	return tx, nil

}

// NOTE: This function requires to the wallet to be unlocked
func SendToAddress(addr string, value float64) (tx *chainhash.Hash, err error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		err = rr
		log.Println("Error DecodeAddress", rr)
		return
	}

	amount, vrr := ltcutil.NewAmount(value)
	if vrr != nil {
		err = vrr
		log.Println("Error NewAmount", vrr)
		return
	}

	tx, err = Ltcd.SendToAddress(address, amount)

	log.Println("SendToAddress: ", tx)

	if NET == "simnet" {
		Generate(uint32(1))
	}

	return
}

// NOTE: This function requires to the wallet to be unlocked
func SendToAddressComment(addr string, value float64, comment, commentTo string) (tx *chainhash.Hash, err error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return
	}

	amount, vrr := ltcutil.NewAmount(value)
	if vrr != nil {
		log.Println("Error NewAmount", vrr)
		return
	}

	tx, err = Ltcd.SendToAddressComment(address, amount, comment, commentTo)

	log.Println("SendToAddress: ", tx)

	return
}

// CreateRawTransaction returns a new transaction spending the provided inputs
// and sending to the provided addresses.
func CreateRawTransaction(inputs []btcjson.TransactionInput,
	amounts map[ltcutil.Address]ltcutil.Amount, lockTime *int64) (*wire.MsgTx, error) {

	return Ltcd.CreateRawTransaction(inputs, amounts, lockTime)
}

// SendRawTransaction submits the encoded transaction to the server which will
// then relay it to the network.
func SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	return Ltcd.SendRawTransaction(tx, allowHighFees)
}

func SignRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	return Ltcd.SignRawTransaction(tx)
}

func GetRawTransaction(txHex string) (tx *ltcutil.Tx, err error) {

	txHash, _ := NewHashFromStr(txHex)

	tx, err = Ltcd.GetRawTransaction(txHash)
	log.Println("GetRawTransaction :", tx)

	return
}

func GetRawTransactionVerbose(txHex string) (btcj *btcjson.TxRawResult, err error) {

	txHash, _ := NewHashFromStr(txHex)
	btcj, err = Ltcd.GetRawTransactionVerbose(txHash)

	//	b, _ := json.MarshalIndent(btcj, "", " ")
	//	log.Println("GetRawTransactionVerbose : ", string(b))

	return
}

func GetTransaction(txHex string) (btcj *btcjson.GetTransactionResult, err error) {

	txHash, _ := NewHashFromStr(txHex)
	btcj, err = Ltcd.GetTransaction(txHash)
	log.Println("GetTransaction : ", btcj)

	return
}

func ListTransactions(account string) (btcj []btcjson.ListTransactionsResult, err error) {

	btcj, err = Ltcd.ListTransactions(account)
	log.Println("ListTransactions : ", btcj)

	return
}
