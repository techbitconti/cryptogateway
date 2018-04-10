package btc

import (
	"testing"
)

func Test_bitcoind(t *testing.T) {

	Connect_bitcoind()

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

	//GetRawTransactionVerbose("c6b8619d339fb253cebeda1dae13312c8167c4d32c90897c85a220cc3d3240ec")

	SendFrom("2NAKhJLCi6yTM6oLjyG2U3sZJAdbcSMhgjh", "2MssTcpVUccPxTbZRAPjygF5F5sr4DBZHbz", float64(0.002))

}
