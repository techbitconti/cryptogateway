package eth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	Abi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	Solc "github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"golang.org/x/net/context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

var Url string

var client *ethclient.Client

var txList map[string][]uint64

func Connect(url string) {

	if client != nil {
		return
	}

	Url = url

	// connect to rpc
	rpc, err := ethclient.Dial(url)
	if err != nil {
		fmt.Println("Error ethClient ", err)
		return
	}

	client = rpc

	txList = make(map[string][]uint64)

}

func call_RPC(method string, paramsIn ...interface{}) map[string]interface{} {

	fmt.Println("method : ", method, "  params : ", paramsIn)

	values := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": paramsIn, "id": 1}
	jsonStr, _ := json.Marshal(values)

	resp, err := http.Post(Url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("Error web3 call ", err)
		return nil
	}
	defer resp.Body.Close()

	result := map[string]interface{}{}
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(method, string(body))

	json.Unmarshal(body, &result)

	return result
}

func NewAccount(path, pass string) (string, error) {

	// Generate a new random account and a funded simulator
	key, err := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	fmt.Println("prvKey Hex : ", common.Bytes2Hex(crypto.FromECDSA(key)))
	fmt.Println("pubKey : ", common.ToHex(crypto.FromECDSAPub(&key.PublicKey)), "  addr : ", addr.Hex())

	ks := keystore.NewKeyStore(path, keystore.StandardScryptN, keystore.StandardScryptP)
	ks.ImportECDSA(key, pass)

	fmt.Println("ETH Accounts : ", ks.Accounts())

	return strings.ToLower(addr.Hex()), err
}

func IsHexAddress(address string) bool {

	return common.IsHexAddress(strings.ToLower(address))
}

func IsHexContract(addr string) bool {

	_, codeHex := CodeAt(addr)
	if codeHex == "0x0" {
		return false
	}

	return true
}

func ValidateAmount(amount string) bool {

	f, err := strconv.ParseFloat(amount, 64)
	if err != nil || len(amount) > 19 {
		fmt.Println("ValidateAmount ", err)
		return false
	}

	fmt.Println("ValidateAmount : ", f)

	return true
}

func ValidateToken(amount string) bool {

	f, err := strconv.ParseInt(amount, 0, 64)
	if err != nil || len(amount) > 19 {
		fmt.Println("ValidateToken ", err)
		return false
	}

	fmt.Println("ValidateToken : ", f)

	return true
}

func GetCoinbase() string {
	return call_RPC("eth_coinbase")["result"].(string)
}

func GetBlockNumber() string {
	return call_RPC("eth_blockNumber")["result"].(string)
}

func GetBlockByNumber(number interface{}) map[string]interface{} {
	return call_RPC("eth_getBlockByNumber", number, true)["result"].(map[string]interface{})
}

func GetAccounts() (arr []string) {
	result := call_RPC("eth_accounts")["result"]

	data, err := json.Marshal(result)
	if err != nil {
		fmt.Println("Error GetAccounts")
		return
	}

	json.Unmarshal(data, &arr)

	return
}

func GetCodeAt(addr string) interface{} {
	return call_RPC("eth_getCode", addr)["result"]
}

func GetTransactionByHash(tx interface{}) map[string]interface{} {
	return call_RPC("eth_getTransactionByHash", tx)
}

func GetFilterChanges(number interface{}) map[string]interface{} {
	return call_RPC("eth_getFilterChanges", number)
}

func GetTransactions(number interface{}) []interface{} {
	txs := GetBlockByNumber(number)["result"].(map[string]interface{})["transactions"].([]interface{})
	return txs
}

func GetTransactionCount(addr, block interface{}) interface{} {
	return call_RPC("eth_getTransactionCount", addr)["result"].(string)
}

func GetTransactionReceipt(tx interface{}) map[string]interface{} {
	return call_RPC("eth_getTransactionReceipt", tx)
}

func GetLogs(tx interface{}) []interface{} {
	logs := GetTransactionReceipt(tx)["result"].(map[string]interface{})["logs"].([]interface{})
	return logs
}

func GetSyncing() map[string]interface{} {
	return call_RPC("eth_syncing")
}

func GetBalance(addr string) *big.Int {

	result := call_RPC("eth_getBalance", addr, "latest")

	if _, ok := result["result"]; !ok {
		return big.NewInt(int64(0))
	}

	str := result["result"].(string)

	balance, _ := strconv.ParseInt(str, 0, 64)

	return big.NewInt(balance)
}

func UnlockAccount(addr, pass string, sec uint64) {
	call_RPC("personal_unlockAccount", addr, pass, sec)
}

func SendTransaction(message map[string]interface{}) string {
	return call_RPC("eth_sendTransaction", message)["result"].(string)
}

func SendRawTransaction(message interface{}) interface{} {
	return call_RPC("eth_sendRawTransaction", message)
}

func Sign(addr, data interface{}) interface{} {
	return call_RPC("eth_sign", addr, data)["result"].(string)
}

func Web3_sha3(hs interface{}) interface{} {
	return call_RPC("web3_sha3", hs)["result"].(string)
}

func Sha3FromHex(s string) string {
	return crypto.Keccak256Hash(common.FromHex(s)).Hex()
}

func Sha3FromEvent(s string) string {
	return crypto.Keccak256Hash([]byte(s)).Hex()
}

func AddressFromEvent(s string) string {
	data := common.FromHex(s)
	hex := common.ToHex(data[12:])
	return hex
}

func SuggestGasPrice() (gasPrice *big.Int, err error) {

	if client == nil {
		return
	}

	gasPrice, err = client.SuggestGasPrice(context.Background())

	return
}

func EstimateGas(_from, _to string, value *big.Int, data []byte) (gasLimit uint64, err error) {

	if client == nil {
		return
	}

	var from common.Address = common.HexToAddress(_from)
	var to common.Address = common.HexToAddress(_to)

	msg := ethereum.CallMsg{From: from, To: &to, Value: value, Data: data}

	gasLimit, err = client.EstimateGas(context.Background(), msg)

	return
}

func GetByteCode(addr, method string, amount *big.Int, params ...interface{}) []byte {

	event := Sha3FromEvent(method)

	sig := common.FromHex(event)[:4]

	data := make([]byte, 0)
	data = append(data, sig...)

	for _, value := range params {

		r := reflect.ValueOf(value)

		if r.Kind() == reflect.String {
			v := common.HexToAddress(value.(string)).Bytes()
			data = append(data, common.BytesToHash(v).Bytes()...)
		} else {
			v := common.BigToHash(value.(*big.Int)).Bytes()
			data = append(data, v...)
		}
	}

	input := common.ToHex(data)

	fmt.Println("method : ", method)
	fmt.Println("event : ", event)
	fmt.Println("sig : ", common.ToHex(sig))
	fmt.Println("value : ", params)
	fmt.Println("data : ", data)
	fmt.Println("input : ", input)

	return data
}

func SolidityCompile(path string) (mContracts map[string]map[string]string) {

	if _, err := exec.LookPath("solc"); err != nil {
		fmt.Println("exec solc", err)
		return
	}
	fmt.Println("SolidityCompile : ", path)

	mContracts = make(map[string]map[string]string)

	data, _ := ioutil.ReadFile(path)
	source := string(data)
	fmt.Println("source : ", source)

	contracts, crr := Solc.CompileSolidityString("", source)

	fmt.Println("compile contracts : ", contracts, crr)

	for stdin, c := range contracts {
		//c, _ := contracts["<stdin>:Test"]
		name := strings.Split(stdin, "<stdin>:")[1]
		code := c.Code
		abi, _ := json.Marshal(c.Info.AbiDefinition)

		fmt.Println("name : ", name)
		fmt.Println("code : ", code)
		fmt.Println("abi : ", string(abi))

		mCon := map[string]string{
			"code": code,
			"abi":  string(abi),
		}
		mContracts[name] = mCon
	}
	return
}

func SolidityDeploy(prvKey, abi, codeHex string, params ...interface{}) (string, string) {

	if client == nil {
		fmt.Println("Error SolidityDeploy Client")
		return "", ""
	}

	parsed, eabi := Abi.JSON(strings.NewReader(abi))
	if eabi != nil {
		fmt.Println("Error SolidityDeploy Abi", eabi)
		return "", ""
	}

	fmt.Println("abi : ", parsed)

	key, _ := crypto.HexToECDSA(prvKey)
	opts := bind.NewKeyedTransactor(key)
	fmt.Println("params........", len(params), params)
	addr, tx, contract, err := bind.DeployContract(opts, parsed, common.FromHex(codeHex), client, params...)
	if err != nil {
		fmt.Println("Error SolidityDeploy Bind", err)
		return "", ""
	}
	fmt.Println("contract : ", contract)
	fmt.Println("address : ", addr.Hex())
	fmt.Println("tx : ", tx)

	return addr.Hex(), tx.Hash().Hex()
}

func SolidityCall(result interface{}, addr, abi, method string, params ...interface{}) interface{} {

	if client == nil {
		fmt.Println("Err Client")
		return ""
	}

	parsed, eabi := Abi.JSON(strings.NewReader(abi))
	if eabi != nil {
		fmt.Println("Err eabi")
		return ""
	}

	fmt.Println("abi : ", parsed)

	contract := bind.NewBoundContract(common.HexToAddress(addr), parsed, client, nil)
	erc := contract.Call(nil, &result, method, params...)
	if erc != nil {
		fmt.Println("Err Call", erc)
		return ""
	}
	fmt.Println("Call : ", result)

	return result
}

//`balanceOf(address)`
func SolidityCallRaw(_from, _to, method string, params ...interface{}) (hexutil.Bytes, error) {

	if client == nil {
		fmt.Println("Err Client")
		return nil, nil
	}

	event := Sha3FromEvent(method)

	sig := common.FromHex(event)[:4]

	data := make([]byte, 0)
	data = append(data, sig...)

	for _, value := range params {

		r := reflect.ValueOf(value)

		if r.Kind() == reflect.String {
			v := common.HexToAddress(value.(string)).Bytes()
			data = append(data, common.BytesToHash(v).Bytes()...)
		} else {
			v := common.BigToHash(value.(*big.Int)).Bytes()
			data = append(data, v...)
		}
	}

	//	input := common.ToHex(data)

	//	fmt.Println("method : ", method)
	//	fmt.Println("event : ", event)
	//	fmt.Println("sig : ", common.ToHex(sig))
	//	fmt.Println("value : ", params)
	//	fmt.Println("data : ", data)
	//	fmt.Println("input : ", input)
	// Input //

	var from common.Address = common.HexToAddress(_from)
	var to common.Address = common.HexToAddress(_to)

	msg := ethereum.CallMsg{From: from, To: &to, Data: data}

	hex, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		fmt.Println("Error CallContract")
	}

	fmt.Println("CallContract : ", hex)

	return hex, err
}

func SolidityTransact(prvKey, addr, abi, method string, params ...interface{}) string {

	if client == nil {
		return ""
	}

	parsed, eabi := Abi.JSON(strings.NewReader(abi))
	if eabi != nil {
		return ""
	}

	key, _ := crypto.HexToECDSA(prvKey)
	addr0 := crypto.PubkeyToAddress(key.PublicKey)
	opts := bind.NewKeyedTransactor(key)

	// Nonce //
	nonce, ern := client.PendingNonceAt(context.Background(), addr0)
	if ern != nil {
		n, _ := fmt.Printf("%v", ern)
		fmt.Println(n)
	}
	fmt.Println("nonce : ", nonce)

	if len(txList[addr0.Hex()]) <= 0 {
		txList[addr0.Hex()] = make([]uint64, 0)
	} else {
		max := len(txList[addr0.Hex()])
		nc := txList[addr0.Hex()][max-1] + 1
		if nc >= nonce {
			nonce = nc
		}
	}
	txList[addr0.Hex()] = append(txList[addr0.Hex()], nonce)

	opts.Nonce = big.NewInt(int64(nonce))

	contract := bind.NewBoundContract(common.HexToAddress(addr), parsed, nil, client)
	tx, err := contract.Transact(opts, method, params...)
	if err != nil {
		fmt.Println("Error Transact", err)
		return ""
	}
	fmt.Println("Transact", tx.Hash().Hex(), "tx : ", tx)

	return tx.Hash().Hex()
}

//SolRaw(prvKey, addressContract, `transfer(address,uint256)`, `0x5e531bb813994f27f65d2d5f7dc7c51dbe5406f7`, big.NewInt(100))
func SolidityTransactRaw(prvKey, addr, method string, amount *big.Int, params ...interface{}) string {

	if client == nil {
		return ""
	}

	key0, _ := crypto.HexToECDSA(prvKey)
	addr0 := crypto.PubkeyToAddress(key0.PublicKey)
	opts := bind.NewKeyedTransactor(key0)

	// Input //
	//method := `print(uint256)`
	//value := big.NewInt(11)

	event := Sha3FromEvent(method)

	sig := common.FromHex(event)[:4]

	data := make([]byte, 0)
	data = append(data, sig...)

	for _, value := range params {

		r := reflect.ValueOf(value)

		if r.Kind() == reflect.String {
			v := common.HexToAddress(value.(string)).Bytes()
			data = append(data, common.BytesToHash(v).Bytes()...)
		} else {
			v := common.BigToHash(value.(*big.Int)).Bytes()
			data = append(data, v...)
		}
	}

	input := common.ToHex(data)

	//	fmt.Println("method : ", method)
	//	fmt.Println("event : ", event)
	//	fmt.Println("sig : ", common.ToHex(sig))
	//	fmt.Println("value : ", params)
	//	fmt.Println("data : ", data)
	fmt.Println("input : ", input)
	// Input //

	// gasLimit //
	var aC common.Address = common.HexToAddress(addr)
	msg := ethereum.CallMsg{From: opts.From, To: &aC, Value: amount, Data: data}
	gasLimit, erg := client.EstimateGas(context.Background(), msg)
	if erg != nil {
		fmt.Println("Error GasLimit")
		return ""
	}

	// Nonce //
	nonce, ern := client.PendingNonceAt(context.Background(), addr0)
	if ern != nil {
		fmt.Println("Error Nonce", ern)
	}
	fmt.Println("nonce : ", nonce)

	if len(txList[addr0.Hex()]) <= 0 {
		txList[addr0.Hex()] = make([]uint64, 0)
	} else {
		max := len(txList[addr0.Hex()])
		nc := txList[addr0.Hex()][max-1] + 1
		if nc >= nonce {
			nonce = nc
		}
	}
	txList[addr0.Hex()] = append(txList[addr0.Hex()], nonce)

	// NewTransaction
	tx := types.NewTransaction(nonce, common.HexToAddress(addr), amount, gasLimit, nil, data)
	//fmt.Println("tx un signed : ", tx)

	signer := types.HomesteadSigner{}
	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key0)
	if err != nil {
		return ""
	}
	rawTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		return ""
	}

	// SendTraction to Contract //
	ers := client.SendTransaction(context.Background(), rawTx)
	if ers != nil {
		fmt.Println("Error SendTraction to Contract")
		return ""
	}

	fmt.Println("rawTx : ", rawTx, rawTx.Hash().Hex())

	return rawTx.Hash().Hex()
}

func SendTransactionRaw(prvKey, to string, value *big.Int, data []byte) string {

	if client == nil {
		return ""
	}

	key0, _ := crypto.HexToECDSA(prvKey)
	addr0 := crypto.PubkeyToAddress(key0.PublicKey)
	opts := bind.NewKeyedTransactor(key0)

	// gasLimit //
	var aC common.Address = common.HexToAddress(to)
	msg := ethereum.CallMsg{From: opts.From, To: &aC, Value: value, Data: data}
	gasLimit, erg := client.EstimateGas(context.Background(), msg)
	if erg != nil {
		return ""
	}
	fmt.Println("msg : ", msg)

	// Nonce //
	nonce, ern := client.PendingNonceAt(context.Background(), addr0)
	if ern != nil {
		n, _ := fmt.Printf("%v", ern)
		fmt.Println(n)
	}
	fmt.Println("nonce : ", nonce)

	if len(txList[addr0.Hex()]) <= 0 {
		txList[addr0.Hex()] = make([]uint64, 0)
	} else {
		max := len(txList[addr0.Hex()])
		nc := txList[addr0.Hex()][max-1] + 1
		if nc >= nonce {
			nonce = nc
		}
	}
	txList[addr0.Hex()] = append(txList[addr0.Hex()], nonce)

	// NewTransaction
	tx := types.NewTransaction(nonce, common.HexToAddress(to), value, gasLimit, nil, data)
	//fmt.Println("tx un signed : ", tx)

	signer := types.HomesteadSigner{}
	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key0)
	if err != nil {
		return ""
	}
	rawTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		return ""
	}
	fmt.Println("rawTx : ", rawTx, rawTx.Hash().Hex())

	// SendTraction to Contract //
	ers := client.SendTransaction(context.Background(), rawTx)
	if ers != nil {
		return ""
	}

	return rawTx.Hash().Hex()
}

func CodeAt(addr string) ([]byte, string) {

	if client == nil {
		return []byte{}, ""
	}

	var aC common.Address = common.HexToAddress(addr)
	code, err := client.CodeAt(context.Background(), aC, nil)
	if err != nil {
		return []byte{}, ""
	}

	return code, common.ToHex(code)
}
