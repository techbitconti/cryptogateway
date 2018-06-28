package bch

import (
	"fmt"
	"testing"

	"github.com/bchsuite/bchd/bchec"
	"github.com/bchsuite/bchutil"
)

func Test_bitcash(t *testing.T) {

	Connect("simnet")

	priv, _ := bchec.NewPrivateKey(bchec.S256())

	pubHash := bchutil.Hash160(priv.PubKey().SerializeCompressed())
	addrHash, _ := bchutil.NewAddressPubKeyHash(pubHash, &Chaincfg)
	//fmt.Println("NewAddressPubKeyHash : ", addrHash.String(), addrHash.EncodeAddress())

	addrCash, err := NewCashAddressPubKeyHash(addrHash.ScriptAddress(), &Chaincfg)
	fmt.Println(err, addrCash)

	DecodeAddress(addrCash.EncodeAddress(), &Chaincfg)

	//	wif, _ := bchutil.NewWIF(priv, &Chaincfg, true)
	//	fmt.Println("wif : ", wif.String(), wif.SerializePubKey())
	//	fmt.Println("ImportPrivKey : ", ImportPrivKey(wif.String(), "XLX111", true))
}
