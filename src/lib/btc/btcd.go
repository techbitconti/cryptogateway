// Copyright (c) 2014-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btc

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	//"time"
	"config"
	"encoding/json"
	"net/http"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/davecgh/go-spew/spew"
)

var Btcd *rpcclient.Client
var Chaincfg chaincfg.Params
var NotifyHandlers rpcclient.NotificationHandlers

var btcdHomeDir string = config.PATH_BTC

func Connect_btcd(net string) {

	if Btcd != nil {
		return
	}

	host := ""
	Chaincfg = chaincfg.SimNetParams

	switch net {
	case "mainnet":
		Chaincfg = chaincfg.MainNetParams
		host = "localhost:8332"
	case "testnet":
		Chaincfg = chaincfg.TestNet3Params
		host = "localhost:18332"
	case "simnet":
		Chaincfg = chaincfg.SimNetParams
		host = "localhost:18554"
	}

	NotifyHandlers = rpcclient.NotificationHandlers{
		OnFilteredBlockConnected: func(height int32, header *wire.BlockHeader, txns []*btcutil.Tx) {
			log.Printf("Block connected: %v (%d) %v", header.BlockHash(), height, header.Timestamp)
		},

		OnFilteredBlockDisconnected: func(height int32, header *wire.BlockHeader) {
			log.Printf("Block disconnected: %v (%d) %v", header.BlockHash(), height, header.Timestamp)
		},

		OnWalletLockState: func(locked bool) {
			log.Println("OnWalletLockState : ", locked)
		},

		OnAccountBalance: func(account string, balance btcutil.Amount, confirmed bool) {
			log.Println("OnAccountBalance : ", "account : ", account, " balance : ", balance, " confirmed : ", confirmed)
		},

		OnRecvTx: func(transaction *btcutil.Tx, details *btcjson.BlockDetails) {
			log.Println("OnRecvTx : ", transaction, " details : ", details)
		},
	}

	certs, err := ioutil.ReadFile(filepath.Join(btcdHomeDir, "rpc.cert"))
	if err != nil {
		log.Println("btcdHomeDir : ", err)
	}
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		Endpoint:     "ws",
		User:         "123",
		Pass:         "123",
		Certificates: certs,
	}
	client, err := rpcclient.New(connCfg, &NotifyHandlers)
	if err != nil {
		log.Println(err)
	}
	Btcd = client

	// Register for block connect and disconnect notifications.
	if err := NotifyBlocks(); err != nil {
		log.Println(err)
	}
	log.Println("NotifyBlocks: Registration Complete")

	unspent, err := client.ListUnspent()
	if err != nil {
		log.Println(err)
	}
	log.Println("Num unspent outputs (utxos): %d", len(unspent))
	if len(unspent) > 0 {
		log.Println("First utxo:\n%v", spew.Sdump(unspent[0]))
	}

	/*
		log.Println("Client shutdown in 10 seconds...")
		time.AfterFunc(time.Second*10, func() {
			log.Println("Client shutting down...")
			client.Shutdown()
			log.Println("Client shutdown complete.")
		})

		// Wait until the client either shuts down gracefully (or the user
		// terminates the process with Ctrl+C).
		client.WaitForShutdown()
	*/

}

func NotifyBlocks() error {

	return Btcd.NotifyBlocks()
}

func NotifyReceived(arr ...string) error {

	addresses := []btcutil.Address{}

	for _, v := range arr {
		a, _ := DecodeAddress(v)
		addresses = append(addresses, a)
	}

	return Btcd.NotifyReceived(addresses)
}

func GetBlockCount() int64 {

	// Get the current block count.
	blockCount, err := Btcd.GetBlockCount()
	if err != nil {
		log.Println(err)
		return 0
	}
	log.Println("Block count: %d", blockCount)

	return blockCount
}

func GetBlockHash(blockHeight int64) (blockHash *chainhash.Hash, err error) {

	blockHash, err = Btcd.GetBlockHash(blockHeight)
	if err != nil {
		log.Println("Error GetBlockHash : ", err)
		return
	}
	log.Println("GetBlockHash: %d", blockHash)

	return
}

func GetBlockHeader(blockHash *chainhash.Hash) (blockHeader *wire.BlockHeader, err error) {

	blockHeader, err = GetBlockHeader(blockHash)
	log.Println("GetBlockHeader: %d", blockHeader)
	return
}

func GetBlock(blockHash *chainhash.Hash) (block *wire.MsgBlock, err error) {

	block, err = Btcd.GetBlock(blockHash)
	log.Println("GetBlock : ", block)

	return
}

func GetBlockFromStr(blockHash string) (block *wire.MsgBlock, err error) {

	hx, _ := NewHashFromStr(blockHash)
	block, err = GetBlock(hx)

	return
}

func NewHashFromStr(hex string) (hash *chainhash.Hash, err error) {

	hash, err = chainhash.NewHashFromStr(hex)
	log.Println("NewHashFromStr : ", hash, err)
	return
}

// GetGenerate returns true if the server is set to mine, otherwise false.
func GetGenerate() (ok bool, err error) {

	ok, err = Btcd.GetGenerate()

	log.Println("GetGenerate : ", ok)

	return
}

// SetGenerate sets the server to generate coins (mine) or not.
func SetGenerate(enable bool, numCPUs int) (err error) {

	err = Btcd.SetGenerate(enable, numCPUs)
	log.Println("SetGenerate : ", err)

	return
}

func VerifyChainBlocks(checkLevel, numBlocks int32) (ok bool, err error) {

	ok, err = Btcd.VerifyChainBlocks(checkLevel, numBlocks)

	return
}

func InvalidateBlock(blockHash *chainhash.Hash) (err error) {

	err = Btcd.InvalidateBlock(blockHash)

	return
}

func WalletPassphrase(pass string, second int64) (bool, error) {

	//  WalletPassphrase
	err := Btcd.WalletPassphrase(pass, second)
	if err != nil {
		log.Println("WalletPassphrase", err)

		return false, err
	}
	log.Println("WalletPassphrase: ", pass)

	return true, nil
}

func CreateNewAccount(account string) (string, error) {

	//CreateNewAccount
	err := Btcd.CreateNewAccount(account)
	if err != nil {
		log.Println("Error CreateNewAccount", err)
		return account, err
	}
	log.Println("CreateNewAccount")

	return account, nil
}

func DecodeAddress(addr string) (address btcutil.Address, err error) {

	address, err = btcutil.DecodeAddress(addr, &Chaincfg)
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
	acc, err = Btcd.ValidateAddress(address)
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

	value, _ := btcutil.NewAmount(f)

	return value.ToBTC()
}

func ToSatoshi(amount string) float64 {
	f, ok := ValidateAmount(amount)
	if !ok {
		return float64(0)
	}

	value, _ := btcutil.NewAmount(f)

	return value.ToUnit(btcutil.AmountSatoshi)
}

func ListAccounts() (map[string]btcutil.Amount, error) {

	// ListAccounts
	list, err := Btcd.ListAccounts()
	if err != nil {
		log.Println("Error ListAccounts", err)

		return nil, err
	}
	log.Println("ListAccounts: ", list)

	return list, nil
}

func ListAddress() (list []btcutil.Address) {

	accounts, _ := ListAccounts()

	for acc, _ := range accounts {
		arr, _ := GetAddressesByAccount(acc)

		list = append(list, arr...)
	}

	log.Println("ListAddress: ", list)
	return
}

func GetBalanceAccount(account string) (amount btcutil.Amount, err error) {

	// GetBalance
	amount, err = Btcd.GetBalance(account)
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

func GetBalanceExplore(addr string) float64 {

	// GetBalance
	//account, _ := GetAccount(addr)
	//amount, err = GetBalanceAccount(account)

	api := "https://testnet.blockexplorer.com/api/addr/"
	url := api + addr
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error GetBalance", err)
		return 0
	}
	defer resp.Body.Close()

	result := map[string]interface{}{}
	body, rr := ioutil.ReadAll(resp.Body)
	if rr != nil {
		log.Println("Error GetBalance", rr)
		return 0
	}
	json.Unmarshal(body, &result)

	amount := result["balance"].(float64)
	log.Println("GetBalance ", addr, amount)

	return amount
}

func GetNewAddress(account string) (address btcutil.Address, err error) {

	list, _ := ListAccounts()
	if _, exist := list[account]; exist {
		log.Println("Error account exist ==>", account)
		return
	}

	// GetNewAddress
	address, err = Btcd.GetNewAddress(account)
	if err != nil {
		log.Println("Error GetNewAddress", err)
	}
	log.Println("GetNewAddress: ", address)

	return
}

// GetAccountAddress returns the current Bitcoin address for receiving payments
// to the specified account.
func GetAccountAddress(account string) (address btcutil.Address, err error) {

	address, err = Btcd.GetAccountAddress(account)
	if err != nil {
		log.Println("Error GetAccountAddress", err)
	}
	log.Println("GetAccountAddress: ", address)

	return
}

func GetAccount(addr string) (account string, err error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		log.Println("Error GetAccount", rr)
		return
	}

	// GetAccount
	account, err = Btcd.GetAccount(address)
	if err != nil {
		log.Println("Error GetAccount", err)
	}
	log.Println("GetAccountAddress: ", account, addr)

	return
}

func GetAddressesByAccount(account string) (address []btcutil.Address, err error) {

	// GetAddressesByAccount
	address, err = Btcd.GetAddressesByAccount(account)
	if err != nil {
		log.Println("Error GetAddressesByAccount", err)
	}
	log.Println("GetAddressesByAccount: ", account, address)

	return
}

func GetReceivedByAccount(account string) (amount btcutil.Amount, err error) {

	amount, err = Btcd.GetReceivedByAccount(account)
	if err != nil {
		log.Println("Error GetReceivedByAccount")
	}
	log.Println("GetReceivedByAccount : ", amount)

	return
}

func GetReceivedByAddress(addr string) (amount btcutil.Amount, err error) {

	address, _ := DecodeAddress(addr)

	amount, err = Btcd.GetReceivedByAddress(address)
	if err != nil {
		log.Println("Error GetReceivedByAddress")
	}
	log.Println("GetReceivedByAddress : ", amount)

	return
}

func ListReceivedByAccount() (btcj []btcjson.ListReceivedByAccountResult, err error) {

	btcj, err = Btcd.ListReceivedByAccount()
	if err != nil {
		log.Println("Error ListReceivedByAccount")
	}
	log.Println("ListReceivedByAccount :", btcj)

	return
}

func ListReceivedByAddress() (btcj []btcjson.ListReceivedByAddressResult, err error) {

	btcj, err = Btcd.ListReceivedByAddress()
	if err != nil {
		log.Println("Error ListReceivedByAddress")
	}
	log.Println("ListReceivedByAddress", btcj)

	return
}

// NOTE: This function requires to the wallet to be unlocked
func DumpPrivKey(addr string) (*btcutil.WIF, error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return nil, rr
	}

	// DumpPrivKey
	wif, err := Btcd.DumpPrivKey(address)
	if err != nil {
		log.Println("DumpPrivKey", err)
		return nil, err
	}
	log.Println("DumpPrivKey: ", wif)

	return wif, nil
}

func ImportPrivKey(prv, label string, rescan bool) error {

	wif, err := btcutil.DecodeWIF(prv)
	if err != nil {
		log.Println("Error ImportPrivKey")
	}
	log.Println("ImportPrivKey : ", wif)

	return Btcd.ImportPrivKeyRescan(wif, label, rescan)
}

func ImportAddress(addr string) error {

	err := Btcd.ImportAddress(addr)
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
	signature, err = Btcd.SignMessage(address, message)

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

	signed, err = Btcd.VerifyMessage(address, signature, message)

	log.Println("VerifyMessage: ", signed)

	return
}

// NOTE: This function requires to the wallet to be unlocked
func SendFrom(fromAddress string, toAddress string, value float64) (tx *chainhash.Hash, err error) {

	fromAccount, frr := GetAccount(fromAddress)
	if frr != nil {
		log.Println("Error fromAddress", frr)
		return
	}

	to, trr := DecodeAddress(toAddress)
	if trr != nil {
		log.Println("Error toAddress", trr)
		return
	}

	amount, vrr := btcutil.NewAmount(value)
	if vrr != nil {
		log.Println("Error NewAmount", vrr)
		return
	}

	tx, err = Btcd.SendFrom(fromAccount, to, amount)

	log.Println("SendFrom: ", tx)

	return
}

// NOTE: This function requires to the wallet to be unlocked.  See the
func SendMany(fromAccount string, amounts map[btcutil.Address]btcutil.Amount) (tx *chainhash.Hash, err error) {

	tx, err = Btcd.SendMany(fromAccount, amounts)
	if err != nil {
		log.Println("Error SendMany", err)
	}

	log.Println("SendMany: ", tx)

	return tx, nil

}

// NOTE: This function requires to the wallet to be unlocked
func SendToAddress(addr string, value float64) (tx *chainhash.Hash, err error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return
	}

	amount, vrr := btcutil.NewAmount(value)
	if vrr != nil {
		log.Println("Error NewAmount", vrr)
		return
	}

	tx, err = Btcd.SendToAddress(address, amount)

	log.Println("SendToAddress: ", tx)

	return
}

// NOTE: This function requires to the wallet to be unlocked
func SendToAddressComment(addr string, value float64, comment, commentTo string) (tx *chainhash.Hash, err error) {

	address, rr := DecodeAddress(addr)
	if rr != nil {
		log.Println("Error DecodeAddress", rr)
		return
	}

	amount, vrr := btcutil.NewAmount(value)
	if vrr != nil {
		log.Println("Error NewAmount", vrr)
		return
	}

	tx, err = Btcd.SendToAddressComment(address, amount, comment, commentTo)

	log.Println("SendToAddress: ", tx)

	return
}

// CreateRawTransaction returns a new transaction spending the provided inputs
// and sending to the provided addresses.
func CreateRawTransaction(inputs []btcjson.TransactionInput,
	amounts map[btcutil.Address]btcutil.Amount, lockTime *int64) (*wire.MsgTx, error) {

	return Btcd.CreateRawTransaction(inputs, amounts, lockTime)
}

// SendRawTransaction submits the encoded transaction to the server which will
// then relay it to the network.
func SendRawTransaction(tx *wire.MsgTx, allowHighFees bool) (*chainhash.Hash, error) {
	return Btcd.SendRawTransaction(tx, allowHighFees)
}

func SignRawTransaction(tx *wire.MsgTx) (*wire.MsgTx, bool, error) {
	return Btcd.SignRawTransaction(tx)
}

func GetRawTransaction(txHex string) (tx *btcutil.Tx, err error) {

	txHash, _ := NewHashFromStr(txHex)

	tx, err = Btcd.GetRawTransaction(txHash)
	log.Println("GetRawTransaction :", tx)

	return
}

func GetRawTransactionVerbose(txHex string) (btcj *btcjson.TxRawResult, err error) {

	txHash, _ := NewHashFromStr(txHex)
	btcj, err = Btcd.GetRawTransactionVerbose(txHash)

	b, _ := json.MarshalIndent(btcj, "", " ")

	log.Println("GetRawTransactionVerbose : ", string(b))

	return
}

func GetTransaction(txHex string) (btcj *btcjson.GetTransactionResult, err error) {

	txHash, _ := NewHashFromStr(txHex)
	btcj, err = Btcd.GetTransaction(txHash)
	log.Println("GetTransaction : ", btcj)

	return
}

func ListTransactions(account string) (btcj []btcjson.ListTransactionsResult, err error) {

	btcj, err = Btcd.ListTransactions(account)
	log.Println("ListTransactions : ", btcj)

	return
}
