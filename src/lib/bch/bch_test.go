package bch

import (
	"testing"
)

func Test_bitcash(t *testing.T) {

	Connect("simnet")

	GetBlockCount()
	//	blockHash, _ := GetBlockHash(blockHeight)
	//	GetBlock(blockHash)

	//GetNewAddress("bitcoin")
	//WalletPassphrase("123456", 10)
	//DumpPrivKey("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh")

	//ImportPrivKey("cQJYynSnzuUbNisDb7FsM2tpKi7Hu3HKtxegWohemwf8YU1EDduD", "bitcoin", true)
	//ImportAddress("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh")

	ValidateAddress("bchreg:qq4gxf9fm9mgemmmx8n9u879mf3c3qmzlc3kc929xx")

	//ListAccounts()
	//ListAddress()

	//GetBalanceAccount("bitcoin")
	//GetBalance("bchreg:qq4gxf9fm9mgemmmx8n9u879mf3c3qmzlc3kc929xx")
	//GetReceivedByAddress("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh")

	//SendFrom("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh", "2MssTcpVUccPxTbZRAPjygF5F5sr4DBZHbz", float64(0.001))
	//GetTransaction("0987d04f5ea00a54dcb606392c5bc04ccee798f8d60f44dfd789c6f0662402cd")
}
