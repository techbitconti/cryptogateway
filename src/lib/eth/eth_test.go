package eth

import (
	"math/big"
	"testing"
)

func Test_SendPriv(t *testing.T) {

	Connect("http://localhost:8545")

	SendTransactionRaw("47de15108b35169c4aff4826d5c413fe117e361a900325f6d3df1f0e04cbd706", "0x26b20cf94cd11f45242ce7214c2f3fb0612a874a", big.NewInt(1), []byte{})

}
