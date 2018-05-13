package btc

import (
	"testing"
)

func Test_bitcoind(t *testing.T) {

	Connect_bitcoind("simnet")

	GetBlockCount()
	//	blockHash, _ := GetBlockHash(blockHeight)
	//	GetBlock(blockHash)

	//GetNewAddress("bitcoin")
	//WalletPassphrase("123456", 10)
	//DumpPrivKey("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh")

	//ImportPrivKey("cQJYynSnzuUbNisDb7FsM2tpKi7Hu3HKtxegWohemwf8YU1EDduD", "bitcoin", true)
	//ImportAddress("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh")

	//ListAccounts()
	//ListAddress()

	//GetBalanceAccount("bitcoin")
	//GetBalanceExplore("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh")
	//GetBalance("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh")
	//GetReceivedByAddress("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh")

	//SendFrom("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh", "2MssTcpVUccPxTbZRAPjygF5F5sr4DBZHbz", float64(0.001))
	//GetRawTransactionVerbose("0987d04f5ea00a54dcb606392c5bc04ccee798f8d60f44dfd789c6f0662402cd")
}
