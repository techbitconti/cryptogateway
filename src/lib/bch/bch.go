package bch

import (
	"encoding/json"
	"errors"
	"log"
	"strconv"

	"github.com/bchsuite/bchd/bchec"
	"github.com/bchsuite/bchd/bchjson"
	"github.com/bchsuite/bchd/chaincfg"
	"github.com/bchsuite/bchd/chaincfg/chainhash"
	"github.com/bchsuite/bchd/rpcclient"
	"github.com/bchsuite/bchd/wire"
	"github.com/bchsuite/bchutil"
)

var NET string

var Bch *rpcclient.Client
var Chaincfg chaincfg.Params
var NotifyHandlers rpcclient.NotificationHandlers

func Connect(net, host string) {

	if Bch != nil {
		return
	}

	NET = net

	switch net {
	case "mainnet":
		//host = "localhost:7332"
		host += ":7332"
		Chaincfg = chaincfg.MainNetParams
	case "testnet":
		//host = "localhost:17332"
		host += ":17332"
		Chaincfg = chaincfg.TestNet3Params
	case "simnet":
		//host = "localhost:17443"
		host += ":17443"
		Chaincfg = chaincfg.RegressionNetParams
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

	Bch = client

}

func GetBlockCount() int64 {

	// Get the current block count.
	blockCount, err := Bch.GetBlockCount()
	if err != nil {
		log.Println(err)
		return 0
	}
	//log.Println("Block count: %d", blockCount)

	return blockCount
}

func GetBlockHash(blockHeight int64) (blockHash *chainhash.Hash, err error) {

	blockHash, err = Bch.GetBlockHash(blockHeight)
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

	block, err = Bch.GetBlock(blockHash)
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
	tx, err := Bch.Generate(num)

	log.Println("Generate", err, tx)

	return tx, err
}

// GetGenerate returns true if the server is set to mine, otherwise false.
func GetGenerate() (ok bool, err error) {

	ok, err = Bch.GetGenerate()

	log.Println("GetGenerate : ", ok)

	return
}

// SetGenerate sets the server to generate coins (mine) or not.
func SetGenerate(enable bool, numCPUs int) (err error) {

	err = Bch.SetGenerate(enable, numCPUs)
	log.Println("SetGenerate : ", err)

	return
}

func VerifyChainBlocks(checkLevel, numBlocks int32) (ok bool, err error) {

	ok, err = Bch.VerifyChainBlocks(checkLevel, numBlocks)

	return
}

func InvalidateBlock(blockHash *chainhash.Hash) (err error) {

	err = Bch.InvalidateBlock(blockHash)

	return
}

func WalletPassphrase(pass string, second int64) (bool, error) {

	//  WalletPassphrase
	err := Bch.WalletPassphrase(pass, second)
	if err != nil {
		log.Println("WalletPassphrase", err)

		return false, err
	}
	log.Println("WalletPassphrase: ", pass)

	return true, nil
}

func CreateNewAccount(account string) (string, error) {

	//CreateNewAccount
	err := Bch.CreateNewAccount(account)
	if err != nil {
		log.Println("Error CreateNewAccount", err)
		return account, err
	}
	log.Println("CreateNewAccount")

	return account, nil
}

/*
func DecodeAddress(addr string) (address bchutil.Address, err error) {

	address, err = bchutil.DecodeAddress(addr, &Chaincfg)
	if err != nil {
		log.Println("Error DecodeAddress", err)
	}
	log.Println("DecodeAddress : ", address.EncodeAddress(), address.EncodeAddress())
	return
}
*/

func ValidateAddress(addr string) (acc *bchjson.ValidateAddressWalletResult, err error) {

	address, err := DecodeAddress(addr, &Chaincfg)
	if err != nil {
		log.Println("Error  DecodeAddress", err)
		return
	}

	// ValidateAddress
	acc, err = Bch.ValidateAddress(address)
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

func ToBCH(amount string) float64 {

	f, ok := ValidateAmount(amount)
	if !ok {
		return float64(0)
	}

	value, _ := bchutil.NewAmount(f)

	return value.ToBCH()
}

func ToSatoshi(amount string) float64 {
	f, ok := ValidateAmount(amount)
	if !ok {
		return float64(0)
	}

	value, _ := bchutil.NewAmount(f)

	return value.ToUnit(bchutil.AmountSatoshi)
}

func ListAccounts() (map[string]bchutil.Amount, error) {

	// ListAccounts
	list, err := Bch.ListAccounts()
	if err != nil {
		log.Println("Error ListAccounts", err)

		return nil, err
	}
	log.Println("ListAccounts: ", list)

	return list, nil
}

func ListAddress() (list []bchutil.Address) {

	accounts, _ := ListAccounts()

	for acc, _ := range accounts {

		arr, _ := GetAddressesByAccount(acc)

		list = append(list, arr...)
	}

	log.Println("ListAddress: ", list)
	return
}

func GetBalanceAccount(account string) (amount bchutil.Amount, err error) {

	// GetBalance
	amount, err = Bch.GetBalance(account)
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

	log.Println("GetBalanceX ", addr, amount.ToBCH())

	return amount.ToBCH()
}

func GetNewAddress(account string) (address *CashAddressPubKeyHash, err error) {

	list, _ := ListAccounts()
	if _, exist := list[account]; exist {
		err = errors.New("Error account exist ==>" + account)
		log.Println("Error account exist ==>", account)
		return
	}

	// GetNewAddress
	//	address, err = Bch.GetNewAddress(account)
	//	if err != nil {
	//		log.Println("Error GetNewAddress", err)
	//	}
	//	log.Println("GetNewAddress: ", address)

	priv, _ := bchec.NewPrivateKey(bchec.S256())

	pubHash := bchutil.Hash160(priv.PubKey().SerializeCompressed())
	addrHash, _ := bchutil.NewAddressPubKeyHash(pubHash, &Chaincfg)

	address, err = NewCashAddressPubKeyHash(addrHash.ScriptAddress(), &Chaincfg)

	wif, _ := bchutil.NewWIF(priv, &Chaincfg, true)
	ImportPrivKey(wif.String(), account, true)

	log.Println("GetNewAddress", "account : ", account, "   -- addrCash : ", address.String(), "   --- privKey : ", wif.String())

	return
}

// GetAccountAddress returns the current Bitcoin address for receiving payments
// to the specified account.
func GetAccountAddress(account string) (address bchutil.Address, err error) {

	address, err = Bch.GetAccountAddress(account)
	if err != nil {
		log.Println("Error GetAccountAddress", err)
	}
	log.Println("GetAccountAddress: ", address)

	return
}

func GetAccount(addr string) (account string, err error) {

	address, rr := DecodeAddress(addr, &Chaincfg)
	if rr != nil {
		err = rr
		log.Println("Error GetAccount", rr)
		return
	}

	// GetAccount
	account, err = Bch.GetAccount(address)
	if err != nil {
		log.Println("Error GetAccount", err)
	}
	log.Println("GetAccountAddress: ", account, addr)

	return
}

func GetAddressesByAccount(account string) (address []bchutil.Address, err error) {

	// GetAddressesByAccount
	address, err = Bch.GetAddressesByAccount(account)
	if err != nil {
		log.Println("Error GetAddressesByAccount", err)
	}
	log.Println("GetAddressesByAccount: ", account, address)

	return
}

func GetReceivedByAccount(account string) (amount bchutil.Amount, err error) {

	amount, err = Bch.GetReceivedByAccount(account)
	if err != nil {
		log.Println("Error GetReceivedByAccount")
	}
	log.Println("GetReceivedByAccount : ", amount)

	return
}

func GetReceivedByAddress(addr string) (amount bchutil.Amount, err error) {

	address, _ := DecodeAddress(addr, &Chaincfg)

	amount, err = Bch.GetReceivedByAddress(address)
	if err != nil {
		log.Println("Error GetReceivedByAddress")
	}
	log.Println("GetReceivedByAddress : ", amount)

	return
}

func ListReceivedByAccount() (btcj []bchjson.ListReceivedByAccountResult, err error) {

	btcj, err = Bch.ListReceivedByAccount()
	if err != nil {
		log.Println("Error ListReceivedByAccount")
	}
	log.Println("ListReceivedByAccount :", btcj)

	return
}

func ListReceivedByAddress() (btcj []bchjson.ListReceivedByAddressResult, err error) {

	btcj, err = Bch.ListReceivedByAddress()
	if err != nil {
		log.Println("Error ListReceivedByAddress")
	}
	log.Println("ListReceivedByAddress", btcj)

	return
}

// NOTE: This function requires to the wallet to be unlocked
func DumpPrivKey(addr string) (*bchutil.WIF, error) {

	address, rr := DecodeAddress(addr, &Chaincfg)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return nil, rr
	}

	// DumpPrivKey
	wif, err := Bch.DumpPrivKey(address)
	if err != nil {
		log.Println("DumpPrivKey", err)
		return nil, err
	}
	log.Println("DumpPrivKey: ", wif)

	return wif, nil
}

func ImportPrivKey(prv, label string, rescan bool) error {

	wif, err := bchutil.DecodeWIF(prv)
	if err != nil {
		log.Println("Error ImportPrivKey")
	}
	log.Println("ImportPrivKey : ", wif)

	return Bch.ImportPrivKeyRescan(wif, label, rescan)
}

func ImportAddress(addr string) error {

	err := Bch.ImportAddress(addr)
	if err != nil {
		log.Println("Error ImportAddress")
	}
	log.Println("ImportAddress", err)

	return err
}

// NOTE: This function requires to the wallet to be unlocked
func SignMessage(addr, message string) (signature string, err error) {

	address, rr := DecodeAddress(addr, &Chaincfg)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return
	}

	// SignMessage
	signature, err = Bch.SignMessage(address, message)

	log.Println("SignMessage: ", signature)

	return
}

// NOTE: This function requires to the wallet to be unlocked
func VerifyMessage(addr, signature, message string) (signed bool, err error) {

	address, rr := DecodeAddress(addr, &Chaincfg)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return
	}

	signed, err = Bch.VerifyMessage(address, signature, message)

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

	to, trr := DecodeAddress(toAddress, &Chaincfg)
	if trr != nil {
		err = trr
		log.Println("Error toAddress", trr)
		return
	}

	amount, vrr := bchutil.NewAmount(value)
	if vrr != nil {
		err = vrr
		log.Println("Error NewAmount", vrr)
		return
	}

	tx, err = Bch.SendFrom(fromAccount, to, amount)

	log.Println("SendFrom: ", tx)

	if NET == "simnet" {
		Generate(uint32(1))
	}

	return
}

// NOTE: This function requires to the wallet to be unlocked.  See the
func SendMany(fromAccount string, amounts map[bchutil.Address]bchutil.Amount) (tx *chainhash.Hash, err error) {

	tx, err = Bch.SendMany(fromAccount, amounts)
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

	address, rr := DecodeAddress(addr, &Chaincfg)
	if rr != nil {
		err = rr
		log.Println("Error DecodeAddress", rr)
		return
	}

	amount, vrr := bchutil.NewAmount(value)
	if vrr != nil {
		err = vrr
		log.Println("Error NewAmount", vrr)
		return
	}

	tx, err = Bch.SendToAddress(address, amount)

	log.Println("SendToAddress: ", tx)

	if NET == "simnet" {
		Generate(uint32(1))
	}

	return
}

// NOTE: This function requires to the wallet to be unlocked
func SendToAddressComment(addr string, value float64, comment, commentTo string) (tx *chainhash.Hash, err error) {

	address, rr := DecodeAddress(addr, &Chaincfg)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return
	}

	amount, vrr := bchutil.NewAmount(value)
	if vrr != nil {
		log.Println("Error NewAmount", vrr)
		return
	}

	tx, err = Bch.SendToAddressComment(address, amount, comment, commentTo)

	log.Println("SendToAddress: ", tx)

	return
}

// CreateRawTransaction returns a new transaction spending the provided inputs
// and sending to the provided addresses.
func CreateRawTransaction(inputs []bchjson.TransactionInput,
	amounts map[bchutil.Address]bchutil.Amount, lockTime *int64) (*wire.MsgTx, error) {

	return Bch.CreateRawTransaction(inputs, amounts, lockTime)
}

// SendRawTransaction submits the encoded transaction to the server which will
// then relay it to the network.
func SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	return Bch.SendRawTransaction(tx, allowHighFees)
}

func SignRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	return Bch.SignRawTransaction(tx)
}

func GetRawTransaction(txHex string) (tx *bchutil.Tx, err error) {

	txHash, _ := NewHashFromStr(txHex)

	tx, err = Bch.GetRawTransaction(txHash)
	log.Println("GetRawTransaction :", tx)

	return
}

func GetRawTransactionVerbose(txHex string) (btcj *bchjson.TxRawResult, err error) {

	txHash, _ := NewHashFromStr(txHex)
	btcj, err = Bch.GetRawTransactionVerbose(txHash)

	//	b, _ := json.MarshalIndent(btcj, "", " ")
	//	log.Println("GetRawTransactionVerbose : ", string(b))

	return
}

func GetTransaction(txHex string) (btcj *bchjson.GetTransactionResult, err error) {

	txHash, _ := NewHashFromStr(txHex)
	btcj, err = Bch.GetTransaction(txHash)
	log.Println("GetTransaction : ", btcj)

	return
}

func ListTransactions(account string) (btcj []bchjson.ListTransactionsResult, err error) {

	btcj, err = Bch.ListTransactions(account)
	log.Println("ListTransactions : ", btcj)

	return
}

func ListUnspent(addr string, minconf int64) ([]bchjson.ListUnspentResult, error) {
	utxo, err := Bch.ListUnspent()
	if err != nil {
		log.Println("Error ListUnspent : ", err)
	}

	unspent := make([]bchjson.ListUnspentResult, 0)
	for _, obj := range utxo {

		if obj.Address != addr {
			continue
		}

		if !obj.Spendable {
			continue
		}

		if obj.Confirmations < minconf {
			continue
		}

		unspent = append(unspent, obj)
	}

	b, _ := json.MarshalIndent(unspent, "", " ")

	log.Println("ListUnspent : ", string(b))

	return unspent, err
}
