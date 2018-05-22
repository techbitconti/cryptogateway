package eth

import (
	"fmt"
	//"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestNewAccount(t *testing.T) {

	privKey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privKey.PublicKey)

	prvHex := common.Bytes2Hex(crypto.FromECDSA(privKey))
	pubHex := common.ToHex(crypto.FromECDSAPub(&privKey.PublicKey))
	addrHex := strings.ToLower(addr.Hex())

	fmt.Println("prvKey Hex : ", prvHex)
	fmt.Println("pubKey : ", pubHex)
	fmt.Println("address : ", addrHex)
}

/*
func TestDeploy(t *testing.T) {

	Connect("http://localhost:8545")

	contracts := SolidityCompile("/Users/A/Desktop/nexle/mass/truffle/contracts/ERC20.sol")
	code := contracts["ERC20"]["code"]
	abi := contracts["ERC20"]["abi"]
	fmt.Println("contracts : ", contracts["ERC20"])
	fmt.Println("abi : ", abi)
	fmt.Println("code : ", code)
	SolidityDeploy("9384681aa3a5d9ce4c0536b8bc0904ef9050a8881b4230b700c3f42dfac952d8", abi, code)
}
*/

/*
func TestCall(t *testing.T) {

	Connect("http://localhost:8545")

	// check balance owner
	hex1, _ := SolidityCallRaw("0x3bacad5eec783560c6bad7bb4d8306a9751df882", "0xC1d2a735CEb9F680bA51aF887729aD71BE73E5fb", `balanceOf(address)`, "0x3bacad5eec783560c6bad7bb4d8306a9751df882")
	fmt.Println("check balance owner :", new(big.Int).SetBytes(hex1))

	// transfer token from owner to receiver
	SolidityTransactRaw("9384681aa3a5d9ce4c0536b8bc0904ef9050a8881b4230b700c3f42dfac952d8", "0xC1d2a735CEb9F680bA51aF887729aD71BE73E5fb", `transfer(address,uint256)`, nil, "0x466f3694a32e7e0f556c67542ce41a8461979016", big.NewInt(10))

	// check balance receiver
	hex2, _ := SolidityCallRaw("0x3bacad5eec783560c6bad7bb4d8306a9751df882", "0xC1d2a735CEb9F680bA51aF887729aD71BE73E5fb", `balanceOf(address)`, "0x466f3694a32e7e0f556c67542ce41a8461979016")
	fmt.Println("check balance receiver :", new(big.Int).SetBytes(hex2))

	// check balance owner
	hex3, _ := SolidityCallRaw("0x3bacad5eec783560c6bad7bb4d8306a9751df882", "0xC1d2a735CEb9F680bA51aF887729aD71BE73E5fb", `balanceOf(address)`, "0x3bacad5eec783560c6bad7bb4d8306a9751df882")
	fmt.Println("check balance owner :", new(big.Int).SetBytes(hex3))
}
*/
