package ltc

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func Test_Litecoin(t *testing.T) {

	Connect("simnet")

	//GetBlockCount()
	//	blockHash, _ := GetBlockHash(0)
	//	GetBlock(blockHash)

	//GetNewAddress("litecoin")
	//GetAddressesByAccount("litecoin")
	//WalletPassphrase("123456", 10)
	//DumpPrivKey("mugRPuZajwD7UdAohNcwzgDG6HnuZj1cni")

	//ImportPrivKey("cQJYynSnzuUbNisDb7FsM2tpKi7Hu3HKtxegWohemwf8YU1EDduD", "litecoin", true)
	//ImportAddress("mugRPuZajwD7UdAohNcwzgDG6HnuZj1cni")

	//ListAccounts()
	//ListAddress()

	//GetBalanceAccount("litecoin")
	//GetBalance("mugRPuZajwD7UdAohNcwzgDG6HnuZj1cni")
	//GetReceivedByAddress("mugRPuZajwD7UdAohNcwzgDG6HnuZj1cni")

	//SendFrom("mgJVwHGgFjSvsRXWQ5iyUksWL6GoxbbdRs", "mq9ayFf2aZ7VmsFCYuvNGrSKP7TzEe3Ewf", float64(100))
	//GetTransaction("fdb4ab403ada197fda5a56d620833e3ecc6a8bbf402f586bed49d0dc16969e13")

	addr, priv := genAddressLTC()
	fmt.Println(addr, priv)
}

func genAddressLTC() (string, string) {

	utc := time.Now().Unix()
	decode := strconv.FormatInt(utc, 10)
	encode := base64.StdEncoding.EncodeToString([]byte(decode))

	address, err1 := GetNewAddress(encode)
	if err1 != nil {
		return "", ""
	}
	privKey, err2 := DumpPrivKey(address.String())
	if err2 != nil {
		return "", ""
	}

	return address.String(), privKey.String()
}
