package eth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	solc "github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	rpc "github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/net/context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var prvKey string = "6969e3d0ec8d0af65e0cd6bbd94cbd8a50418ef3b34e35041450d4d7640d9f23"
var addrKey string = "0xa1f2d3d790b05f11c1c718bf5cae0deff0d5e9ee"

func post() {
	values := map[string]interface{}{"jsonrpc": "2.0", "method": "eth_coinbase", "params": []interface{}{}, "id": 1}
	jsonStr, _ := json.Marshal(values)

	resp, err := http.Post("http://127.0.0.1:8545", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}

func request() {

	url := "http://127.0.0.1:8545"
	//fmt.Println("URL:>", url)

	values := map[string]interface{}{"jsonrpc": "2.0", "method": "eth_gasPrice", "params": []interface{}{}, "id": 1}
	jsonStr, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}

func call(method string, paramsIn ...interface{}) {

	fmt.Println("method : ", method, "  params : ", paramsIn)

	values := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": paramsIn, "id": 1}
	jsonStr, _ := json.Marshal(values)

	resp, err := http.Post("http://127.0.0.1:8545", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("response Headers:", resp.Header)

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func Test(t *testing.T) {

}

func TestKeyStore(t *testing.T) {
	ks := keystore.NewKeyStore("/home/trung/ethereum_private/keystore", keystore.StandardScryptN, keystore.StandardScryptP)

	//	n := len(ks.Accounts())
	//	for i := 0; i < n; i++ {
	//		if i > 0 {
	//			ks.Delete(ks.Accounts()[i], "123456")
	//		}
	//	}

	//	list := ks.Accounts()
	//	for _, acc := range list {
	//		fmt.Println("acc : ", acc.Address.Hex())
	//	}

	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	fmt.Println("prvKey Hex : ", common.Bytes2Hex(crypto.FromECDSA(key)))
	fmt.Println("pubKey : ", common.ToHex(crypto.FromECDSAPub(&key.PublicKey)), "  addr : ", addr.Hex())

	ks.ImportECDSA(key, "123456")

	listN := ks.Accounts()
	for _, acc := range listN {
		fmt.Println("acc : ", acc.Address.Hex())
	}

}

func TestWeb3(t *testing.T) {

	//call("eth_blockNumber")
	//call("eth_getBlockByNumber", "0x02", true)

	//call("eth_getBalance","0x09ca5f6bf1299e6a8e8e910b7e45ec8cc0ac0556")
	//call("eth_sendTransaction", map[string]string{"from": "0xd712a8124134d050bc3799f7056ffaef7591ce15", "to": "0xcf7aa8b8ff2c37b831d0ff5b3e03b4b1d428231d", "value": "1000"})

	//call("eth_getTransactionCount", "0xd712a8124134d050bc3799f7056ffaef7591ce15", "pending")
	call("eth_getTransactionReceipt", "0xb19c616275a48be3dd853f89a0f50c43325a7e8fc1296ecc9fa5a93e3f2a8793")
	//call("eth_getTransactionByHash", "0x119a567eb01c67420ce5fd99fd1a7ab6c2138062c26ae628b2b43d8ac49f1a76")

	//	b := hexutil.Big(*big.NewInt(1000000))
	//	value := b.String()
	//	msg := map[string]interface{}{"from": "0x911fd1328d3782d1f7237d4f4c61d3fa69610225", "to": "0xa1f2d3d790b05f11c1c718bf5cae0deff0d5e9ee", "value": value}
	//	.SendTransaction(msg)
	//	call("eth_getBalance", "0xa1f2d3d790b05f11c1c718bf5cae0deff0d5e9ee", "latest")
	//	fmt.Println(.GetBalance("0xa1f2d3d790b05f11c1c718bf5cae0deff0d5e9ee"))

	//	param := map[string]interface{}{}
	//	param["topics"] = []string{"0x015d831579c6c98a55e2fbb92250d99250365ad639a8181ab603ad6e1ac6802c"}
	//	call("eth_newFilter", param)
	//	call("eth_getFilterChanges", "0x06")
	//	call("eth_getLogs", param)

}

func TestSha3(t *testing.T) {

	//	s := "0xb5485962447535b482100b41cbea458c624379cf136495834495190587684580970546719416274516772976153406203899607085165482491682"
	//	fmt.Println(web3.Sha3FromHex(s))

	LogJoinOK := "LogJoinOK(uint256,bytes32,uint256)"
	fmt.Println("LogJoinOK : ", Sha3FromEvent(LogJoinOK))

	LogJoinFail := "LogJoinFail(uint256,bytes32,uint256,uint256)"
	fmt.Println("LogJoinFail :", Sha3FromEvent(LogJoinFail))

	LogSplitOK := "LogSplitOK(uint256,address,uint256)"
	fmt.Println("LogSplitOK : ", Sha3FromEvent(LogSplitOK))

	LogSplitFail := "LogSplitFail(uint256,address,uint256)"
	fmt.Println("LogSplitFail : ", Sha3FromEvent(LogSplitFail))

	LogPublishOK := "LogPublishOK(uint256,address,bytes)"
	fmt.Println("LogPublishOK : ", Sha3FromEvent(LogPublishOK))

	LogFailCampaign := "LogFailCampaign(uint256,uint256)"
	fmt.Println("LogFailCampaign : ", Sha3FromEvent(LogFailCampaign))

	LogStart := "LogStart(uint256,uint256,uint256,uint256,uint256,uint256)"
	fmt.Println("LogStart : ", Sha3FromEvent(LogStart))

	//	fmt.Println("AddressToHex : ", AddressFromEvent("0x0000000000000000000000007e5b63af5d65d6b72bfc5c00f7195a0515e4b0ed"))

	//	print := `print(uint256)`
	//	value := big.NewInt(11)

	//	event := web3.Sha3FromEvent(print)
	//	signature := common.FromHex(event)[:4]
	//	data := append(signature, common.BigToHash(value).Bytes()...)

	//	fmt.Println("event : ", event)
	//	fmt.Println("signature : ", common.ToHex(signature))
	//	fmt.Println("input : ", value)
	//	fmt.Println("data : ", common.ToHex(data))
}

var testSource string = `
	pragma solidity ^0.4.8;
	contract Test{
		
		event LogMultiply(uint256 a, uint256 b, uint256 c);
		event LogPrint(uint256 p);
		
		function Test(){
			
		}
		
	   	function multiply(uint256 a, uint256 b) returns(uint256 c) {
			c = a * b;
			LogMultiply(a, b, c);		
	       	return c;
	   	}
		
		function print(uint256 p) returns(uint256){
			LogPrint(p);
			return p;
		}
	}`

func TestNonce(t *testing.T) {

	b := hexutil.Big(*big.NewInt(11))
	value := b.String()
	fmt.Println("value : ", value)

	call("personal_unlockAccount", "0xd0e602df00041ae48fe4c1856f73b208b6f554d0", "123456", uint64(86400))
	call("eth_sendTransaction", map[string]string{"from": "0xd0e602df00041ae48fe4c1856f73b208b6f554d0", "to": "0xf5da8903ce27e622a7c2af267cd9ad1a97eb8270", "value": value})
	call("eth_getTransactionCount", "0xd0e602df00041ae48fe4c1856f73b208b6f554d0", "pending")

	//connect to testrpc
	client, ss := rpc.Dial("http://127.0.0.1:8545")
	if ss != nil {
		panic(ss)
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress("0x911fd1328d3782d1f7237d4f4c61d3fa69610225"))
	if err != nil {
		n, _ := fmt.Printf("%v", err)
		fmt.Println(n)
	}
	fmt.Println("nonce : ", nonce)
}

func TestDeploy(t *testing.T) {

	//code :  0x6060604052600a600281905560038190556004819055600555601e6006556007805460ff19169055341561002f57fe5b5b60008054600160a060020a03191633600160a060020a03161790555b5b610a6a8061005c6000396000f300606060405236156100675763ffffffff60e060020a60003504166307da68f5811461006957806314d8085c1461007b5780635d1e2d1b146100995780639d3f08f7146100ba578063a45b27b31461014a578063be9a6555146101ac578063e5009bb6146101be575bfe5b341561007157fe5b6100796101ce565b005b341561008357fe5b6100796004356024356044356064356101f7565b005b34156100a157fe5b610079600160a060020a0360043516602435610212565b005b34156100c257fe5b6100ca610324565b604080516020808252835181830152835191928392908301918501908083838215610110575b80518252602083111561011057601f1990920191602091820191016100f0565b505050905090810190601f16801561013c5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b341561015257fe5b60408051602060046024803582810135601f8101859004850286018501909652858552610079958335600160a060020a031695939460449493929092019181908401838280828437509496506103e095505050505050565b005b34156101b457fe5b610079610502565b005b6100796004356024356106af565b005b60005433600160a060020a039081169116146101e957610000565b6007805460ff191690555b5b565b60028490556003839055600482905560058190555b50505050565b60005433600160a060020a0390811691161461022d57610000565b6001546000908152600860209081526040808320600160a060020a03861684526004019091529020541561026f5761026a600154838360016108dd565b61031f565b604051600160a060020a0383169082156108fc029083906000818181858888f1935050505015156102ae5761026a600154838360026108dd565b61031f565b600180546000908152600860209081526040808320600160a060020a0387168085526004909101835292819020859055925483519081529081019190915280820183905290517fbf70f2bc0f0d0ad1228883287eb77d54c1218b36126bd6584f9628876cab93489181900360600190a15b5b5050565b61032c61094c565b60018054600090815260086020908152604080832033600160a060020a0316845260030182529182902080548351600295821615610100026000190190911694909404601f810183900483028501830190935282845291908301828280156103d55780601f106103aa576101008083540402835291602001916103d5565b820191906000526020600020905b8154815290600101906020018083116103b857829003601f168201915b505050505090505b90565b60005433600160a060020a039081169116146103fb57610000565b6001546000908152600860209081526040808320600160a060020a0386168452600301825290912082516104319284019061095e565b507f988614b62bfb00d31e7d37f5e3f2bba5130522883d3004e0de625522f6a5f89660015483836040518084815260200183600160a060020a0316600160a060020a03168152602001806020018281038252838181518152602001915080519060200190808383600083146104c1575b8051825260208311156104c157601f1990920191602091820191016104a1565b505050905090810190601f1680156104ed5780820380516001836020036101000a031916815260200191505b5094505050505060405180910390a15b5b5050565b61050a6109dd565b60005433600160a060020a0390811691161461052557610000565b60075460ff161561053557610000565b6001600081548092919060010191905055504381600001818152505060025481600001510181602001818152505060035481602001510181604001818152505060045481604001510181606001818152505060055481606001510181608001818152505060008160a00181815250506001600760006101000a81548160ff02191690831515021790555080600860006001548152602001908152602001600020600082015181600501556020820151816006015560408201518160070155606082015181600801556080820151816009015560a082015181600a015560c082015181600b0160006101000a81548160ff0219169083151502179055509050507f1c7c8f7872ca176b9b1d7dd41add9f3ae61564f9b851efe5a227525040d3724d6001548260000151836020015184604001518560600151866080015160405180878152602001868152602001858152602001848152602001838152602001828152602001965050505050505060405180910390a15b5b5b50565b6000341161070857600180546040805191825260208201859052818101849052606082019290925290517fd53480493ed6e0b0ad596a84e1d779c8ed9a953bc1e0268c8ef30c378b0096139181900360800190a161031f565b6001546000908152600860209081526040808320858452909152902054600160a060020a03161561078157600154604080519182526020820184905281810183905260026060830152517fd53480493ed6e0b0ad596a84e1d779c8ed9a953bc1e0268c8ef30c378b0096139181900360800190a161031f565b600180546000908152600860209081526040808320600160a060020a03331684529093019052205434018190111561080157600154604080519182526020820184905281810183905260036060830152517fd53480493ed6e0b0ad596a84e1d779c8ed9a953bc1e0268c8ef30c378b0096139181900360800190a161031f565b60015460005260086020525b600180546000908152600860208181526040808420600a018054860190558454845280842087855282528084208054600160a060020a03331673ffffffffffffffffffffffffffffffffffffffff19909116811790915585548552838352818520818652860183528185208054340190558554855292825280842092845260029092018152918190208490559154825190815290810184905280820183905290517fa02f15a5453c2fd5348a53745e46efb4691ca0ebbdb3fa609d2fd122e3601d189181900360600190a15b5050565b600084815260086020908152604091829020600b01805460ff191660011790558151868152600160a060020a0386169181019190915280820183905290517f9d96d39003ab484c7c3dd9c378d33438dfc94956cba331602a752382351d38399181900360600190a15b50505050565b60408051602081019091526000815290565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061099f57805160ff19168380011785556109cc565b828001600101855582156109cc579182015b828111156109cc5782518255916020019190600101906109b1565b5b506109d9929150610a1d565b5090565b60e0604051908101604052806000815260200160008152602001600081526020016000815260200160008152602001600081526020016000151581525090565b6103dd91905b808211156109d95760008155600101610a23565b5090565b905600a165627a7a723058203d84385b15206806ab5cc5988e38179a84e8a56b6c1b43519ad4e05cf4d80d5b0029
	//abi :  [{"constant":false,"inputs":[],"name":"stop","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"_block_Join","type":"uint256"},{"name":"_block_Publish","type":"uint256"},{"name":"_block_Review","type":"uint256"},{"name":"_block_Split","type":"uint256"}],"name":"setConfigBlock","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"amount","type":"uint256"}],"name":"split","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"review","outputs":[{"name":"","type":"bytes"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"from","type":"address"},{"name":"data","type":"bytes"}],"name":"publish","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"start","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"k","type":"bytes32"},{"name":"b","type":"uint256"}],"name":"join","outputs":[],"payable":true,"type":"function"},{"inputs":[],"payable":false,"type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"roundID","type":"uint256"},{"indexed":false,"name":"block_Start","type":"uint256"},{"indexed":false,"name":"block_Join","type":"uint256"},{"indexed":false,"name":"block_Publish","type":"uint256"},{"indexed":false,"name":"block_Review","type":"uint256"},{"indexed":false,"name":"block_Split","type":"uint256"}],"name":"LogStart","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"roundID","type":"uint256"},{"indexed":false,"name":"K","type":"bytes32"},{"indexed":false,"name":"B","type":"uint256"}],"name":"LogJoinOK","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"roundID","type":"uint256"},{"indexed":false,"name":"K","type":"bytes32"},{"indexed":false,"name":"B","type":"uint256"},{"indexed":false,"name":"err","type":"uint256"}],"name":"LogJoinFail","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"roundID","type":"uint256"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"amount","type":"uint256"}],"name":"LogSplitOK","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"roundID","type":"uint256"},{"indexed":false,"name":"to","type":"address"},{"indexed":false,"name":"err","type":"uint256"}],"name":"LogSplitFail","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"roundID","type":"uint256"},{"indexed":false,"name":"from","type":"address"},{"indexed":false,"name":"data","type":"bytes"}],"name":"LogPublishOK","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"roundID","type":"uint256"},{"indexed":false,"name":"err","type":"uint256"}],"name":"LogFailCampaign","type":"event"}]
	//address := "0x60322b54143ac440328960dc6d4e82a06324afe1"

	//	url := "http://127.0.0.1:8545"
	//	prv := prvKey
	//	Connect(url)
	//	code, abi := SolidityCompile("/home/trung/Desktop/Solarium/JoinSplit/SmartContract/contracts/JoinSplit.sol")
	//	SolidityDeploy(prv, abi, code)

	//web3.SolidityTransact(prv, "0x60322b54143ac440328960dc6d4e82a06324afe1", abi, "start")
}

func TestSol(t *testing.T) {

	//addr := "0xa27f7898de645079de55ed74f1509cfc01df1027"
	//.SolidityTransact(prv, addr, abi, "print", big.NewInt(22))
	//	var result *big.Int
	//	.SolidityCall(&result, addr, abi, "print", big.NewInt(22))

	// connect to testrpc
	client, ss := rpc.Dial("http://127.0.0.1:8545")
	if ss != nil {
		fmt.Println(ss, client)
		return
	}

	// generate
	key0, _ := crypto.HexToECDSA(prvKey)
	//addr0 := crypto.PubkeyToAddress(key0.PublicKey)

	// compile contract from text source
	path := "/home/trung/Desktop/Solarium/JoinSplit/SmartContract/contracts/"
	name := "Test.sol"
	data, _ := ioutil.ReadFile(path + name)
	source := string(data)
	fmt.Println("source : ", source)

	contracts, _ := solc.CompileSolidityString("", source)
	c, _ := contracts["<stdin>:Test"]
	code := c.Code
	definition, _ := json.Marshal(c.Info.AbiDefinition)
	parsed, _ := abi.JSON(strings.NewReader(string(definition)))
	fmt.Println("code : ", code)
	fmt.Println("definition : ", string(definition))
	fmt.Println("abi : ", parsed)

	opts := bind.NewKeyedTransactor(key0)
	fmt.Println("opts : ", opts)

	// Deploy //
	addr, txC, bound, err := bind.DeployContract(opts, parsed, common.FromHex(code), client)
	//backend.Commit()
	if err != nil {
		panic(err)
	}
	fmt.Println("contract : ", bound)
	fmt.Println("address : ", addr.Hex())
	fmt.Println("txC : ", txC)

	// SendTraction to Contract //
	/*
		contract := bind.NewBoundContract(common.HexToAddress("0x84767e0d2ff992dcddbb266fa2131501e58305e2"), parsed, client, client)
		txs, err := contract.Transact(opts, "print", big.NewInt(12))
		if err != nil {
			panic(err)
		}
		fmt.Println("Transact", txs.Hash().Hex(), "txs : ", txs)

		var result *big.Int
		erc := contract.Call(nil, &result, "print", big.NewInt(12))
		if erc != nil {
			panic(erc)
		}
		fmt.Println("Call : ", result)
	*/

}

//SolRaw(prvKey, addressContract, `transfer(address,uint256)`, `0x5e531bb813994f27f65d2d5f7dc7c51dbe5406f7`, big.NewInt(100))
func testSolRaw(prvKey, addr, method string, params ...interface{}) {

	// connect to testrpc
	client, ss := rpc.Dial("http://127.0.0.1:8545")
	if ss != nil {
		fmt.Println(ss, client)
		return
	}

	key0, _ := crypto.HexToECDSA(prvKey)
	addr0 := crypto.PubkeyToAddress(key0.PublicKey)
	opts := bind.NewKeyedTransactor(key0)

	// Input //
	//method := `print(uint256)`
	//value := big.NewInt(11)

	event := Sha3FromEvent(method)
	sig := common.FromHex(event)[:4]
	//data := append(sig, common.BigToHash(value).Bytes()...)
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
	// Input //

	// gasLimit //
	var aC common.Address = common.HexToAddress(addr)
	msg := ethereum.CallMsg{From: opts.From, To: &aC, Value: nil, Data: data}
	gasLimit, erg := client.EstimateGas(context.Background(), msg)
	if erg != nil {
		panic(erg)
	}
	// gasLimit //

	// Nonce //
	nonce, ern := client.PendingNonceAt(context.Background(), addr0)
	if ern != nil {
		n, _ := fmt.Printf("%v", ern)
		fmt.Println(n)
	}
	fmt.Println("nonce : ", nonce)
	// Nonce //

	// SendTraction to Contract //
	tx := types.NewTransaction(nonce, common.HexToAddress(addr), nil, gasLimit, nil, data)
	fmt.Println("tx un signed : ", tx)
	signer := types.HomesteadSigner{}
	signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key0)
	if err != nil {
		panic(err)
	}
	rawTx, err := tx.WithSignature(signer, signature)
	if err != nil {
		panic(err)
	}
	fmt.Println("rawTx : ", rawTx, rawTx.Hash().Hex())

	ers := client.SendTransaction(context.Background(), rawTx)
	if ers != nil {
		panic(ers)
	}
}

func TestContract(t *testing.T) {

	//	eth.Connect("http://localhost:8545")

	//	contracts := eth.SolidityCompile("/Users/A/Desktop/nexle/mass/truffle/contracts/TestSource.sol")
	//	code := contracts["Test"]["code"]
	//	abi := contracts["Test"]["abi"]

	//	eth.SolidityDeploy(config.ETH_SIM.PrivKey, abi, code)

	//	eth.SolidityTransactRaw(config.ETH_SIM.PrivKey, config.ERC20_SIM.Address, `transfer(address,uint256)`, nil, "0x34e58de83ae76f96eecb8765890669a9784d641d", big.NewInt(10000000))

	//	hex, _ := eth.SolidityCallRaw(config.ETH_SIM.Address, config.ERC20_SIM.Address, `balanceOf(address)`, "0x34e58de83ae76f96eecb8765890669a9784d641d")
	//	fmt.Println(common.BytesToHash(hex).Big(), new(big.Int).SetBytes(hex))

}
