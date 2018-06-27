package bch

import (
	"fmt"
	"testing"

	"github.com/bchsuite/bchd/bchec"
	"github.com/bchsuite/bchutil"
)

func Test_bitcash(t *testing.T) {

	Connect("simnet")

	//GetBlockCount()
	//	blockHash, _ := GetBlockHash(blockHeight)
	//	GetBlock(blockHash)

	//GetNewAddress("bitcoin")
	//WalletPassphrase("123456", 10)
	//DumpPrivKey("qq4gxf9fm9mgemmmx8n9u879mf3c3qmzlc3kc929xx")

	//ImportPrivKey("cQJYynSnzuUbNisDb7FsM2tpKi7Hu3HKtxegWohemwf8YU1EDduD", "bitcoin", true)
	//ImportAddress("qq4gxf9fm9mgemmmx8n9u879mf3c3qmzlc3kc929xx")

	priv, _ := bchec.NewPrivateKey(bchec.S256())

	pubHash := bchutil.Hash160(priv.PubKey().SerializeCompressed())
	addrHash, _ := bchutil.NewAddressPubKeyHash(pubHash, &Chaincfg)
	fmt.Println("NewAddressPubKeyHash : ", addrHash.String(), addrHash.EncodeAddress(), addrHash.ScriptAddress())
	DecodeAddress(addrHash.String())

	pub := priv.PubKey().SerializeCompressed()
	addrPub, _ := bchutil.NewAddressPubKey(pub, &Chaincfg)
	fmt.Println("NewAddressPubKey : ", addrPub.String(), addrPub.EncodeAddress(), addrPub.ScriptAddress())

	DecodeAddress(addrPub.String())

	//	wif1, _ := bchutil.NewWIF(priv, &Chaincfg, true)
	//	fmt.Println("wif1 : ", wif1.String(), wif1.SerializePubKey())

	//	fmt.Println("ImportPrivKey : ", ImportPrivKey(wif1.String(), "XLX111", true))

	//ListAccounts()
	//ListAddress()

	//GetBalanceAccount("bitcash")
	//GetBalance("bchreg:qq4gxf9fm9mgemmmx8n9u879mf3c3qmzlc3kc929xx")
	//GetReceivedByAddress("qq4gxf9fm9mgemmmx8n9u879mf3c3qmzlc3kc929xx")

	//SendFrom("qq4gxf9fm9mgemmmx8n9u879mf3c3qmzlc3kc929xx", "2MssTcpVUccPxTbZRAPjygF5F5sr4DBZHbz", float64(0.001))
	//GetTransaction("0987d04f5ea00a54dcb606392c5bc04ccee798f8d60f44dfd789c6f0662402cd")
}
